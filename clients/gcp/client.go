package gcp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2/google"
	crmv1 "google.golang.org/api/cloudresourcemanager/v1"
	computev1 "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	iamv1 "google.golang.org/api/iam/v1"
)

var (
	// ErrAPIForbidden is returned when the api returns a 401
	ErrAPIForbidden = wrapError{msg: "The api is not allowed for the current service account"}

	// ErrAPINotEnabled is returned when an api is not enabled
	ErrAPINotEnabled = wrapError{msg: "The api is not enabled"}

	// ErrUnknownProjectID is returned when an api throws 'googleapi: Error 400: Unknown project id: 0, invalid'
	ErrUnknownProjectID = wrapError{msg: "The project id is unknown"}

	// ErrProjectNotFound is returned when an api throws 'googleapi: Error 404: The requested project was not found., notFound'
	ErrProjectNotFound = wrapError{msg: "The project is not found"}

	// ErrEntityNotFound is returned when pubsub topics return html with a 404
	ErrEntityNotFound = wrapError{msg: "Entity is not found"}

	// ErrEntityNotActive is returned when cloud sql instance is not running and its databases cannot be fetched
	ErrEntityNotActive = wrapError{msg: "Entity is not runactivening"}
)

//go:generate mockgen -package=gcp -destination ./mock.go -source=client.go
type Client interface {
	GetProjectByLabels(ctx context.Context, filters []string) (projects []*crmv1.Project, err error)
	GetProjectNetworks(ctx context.Context, projects []*crmv1.Project) (networks []*computev1.Network, err error)
	GetProjectSubnetworks(ctx context.Context, projects []*crmv1.Project) (subnetworks []*computev1.Subnetwork, err error)
}

// NewClient returns a new gcp.Client
func NewClient(ctx context.Context, concurrency int) (Client, error) {

	// use service account to authenticate against gcp apis
	googleClient, err := google.DefaultClient(ctx, iamv1.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	computev1Service, err := computev1.New(googleClient)
	if err != nil {
		return nil, err
	}

	crmv1Service, err := crmv1.New(googleClient)
	if err != nil {
		return nil, err
	}

	return &client{
		computev1Service: computev1Service,
		crmv1Service:     crmv1Service,

		concurrency: concurrency,
	}, nil
}

type client struct {
	computev1Service *computev1.Service
	crmv1Service     *crmv1.Service

	concurrency int
}

func (c *client) GetProjectByLabels(ctx context.Context, filters []string) (projects []*crmv1.Project, err error) {

	filters = append(filters, "lifecycleState:ACTIVE")

	log.Info().Msgf("Retrieving projects for filters %v...", filters)

	projects = make([]*crmv1.Project, 0)

	nextPageToken := ""

	for {
		// retrieving projects matching labels (by page)
		var resp *crmv1.ListProjectsResponse
		err = c.substituteErrorsWithPredefinedErrors(foundation.Retry(func() error {
			listCall := c.crmv1Service.Projects.List()
			if nextPageToken != "" {
				listCall.PageToken(nextPageToken)
			}

			// set filter
			listCall.Filter(strings.Join(filters, " "))

			listCall.Context(ctx)
			resp, err = listCall.Do()
			if err != nil {
				return err
			}

			return nil
		}, c.getRetryOptions()...))
		if err != nil {
			return projects, fmt.Errorf("Can't get projects by filters %v: %w", filters, err)
		}

		projects = append(projects, resp.Projects...)

		if resp.NextPageToken == "" {
			break
		}
		nextPageToken = resp.NextPageToken
	}

	log.Debug().Msgf("Retrieved %v projects for filters %v", len(projects), filters)

	return
}

func (c *client) getProjectNetworks(ctx context.Context, projectID string) (networks []*computev1.Network, err error) {
	if projectID == "" {
		return nil, fmt.Errorf("GetProjectNetworks argument projectID is empty")
	}

	log.Debug().Msgf("Retrieving networks for project %v...", projectID)

	nextPageToken := ""
	for {
		var resp *computev1.NetworkList
		err = c.substituteErrorsWithPredefinedErrors(foundation.Retry(func() error {
			listCall := c.computev1Service.Networks.List(projectID)
			if nextPageToken != "" {
				listCall.PageToken(nextPageToken)
			}
			listCall.Context(ctx)
			resp, err = listCall.Do()
			if err != nil {
				return err
			}
			return nil
		}, c.getRetryOptions()...))
		if err != nil && !errors.Is(err, ErrAPIForbidden) {
			return networks, fmt.Errorf("Can't get project networks for project id %v: %w", projectID, err)
		}
		if err != nil && errors.Is(err, ErrAPIForbidden) {
			return networks, nil
		}
		networks = append(networks, resp.Items...)

		if resp.NextPageToken == "" {
			break
		}
		nextPageToken = resp.NextPageToken
	}

	log.Debug().Msgf("Retrieved %v networks for project %v", len(networks), projectID)

	return
}

func (c *client) GetProjectNetworks(ctx context.Context, projects []*crmv1.Project) (networks []*computev1.Network, err error) {

	// http://jmoiron.net/blog/limiting-concurrency-in-go/
	semaphore := make(chan bool, c.concurrency)
	cancelled := false

	resultChannel := make(chan struct {
		Networks []*computev1.Network
		Err      error
	}, len(projects))

	for _, p := range projects {
		select {
		// try to fill semaphore up to it's full size otherwise wait for a routine to finish
		case semaphore <- true:

			go func(ctx context.Context, p *crmv1.Project) {
				// lower semaphore once the routine's finished, making room for another one to start
				defer func() { <-semaphore }()

				networks, err := c.getProjectNetworks(ctx, p.ProjectId)

				resultChannel <- struct {
					Networks []*computev1.Network
					Err      error
				}{networks, err}
			}(ctx, p)

		case <-ctx.Done():
			log.Info().Msg("User has canceled execution, stopping retrieval of networks...")
			cancelled = true
		}
		if cancelled {
			log.Info().Msg("User has canceled execution, waiting for pending retrieval of networks to finish...")
			break
		}
	}

	// try to fill semaphore up to it's full size which only succeeds if all routines have finished or execution has been canceled
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- true
	}

	if cancelled {
		log.Info().Msg("User has canceled execution, checking retrieved networks...")
	}

	// check for errors and aggregate all networks
	close(resultChannel)
	for r := range resultChannel {
		if r.Err != nil {
			err = r.Err
			return
		}
		networks = append(networks, r.Networks...)
	}

	return
}

func (c *client) getProjectSubnetworks(ctx context.Context, projectID string) (subnetworks []*computev1.Subnetwork, err error) {

	log.Info().Msgf("Retrieving subnetworks for project %v...", projectID)

	nextPageToken := ""
	for {
		var resp *computev1.SubnetworkAggregatedList
		err = c.substituteErrorsWithPredefinedErrors(foundation.Retry(func() error {

			listCall := c.computev1Service.Subnetworks.AggregatedList(projectID)
			if nextPageToken != "" {
				listCall.PageToken(nextPageToken)
			}
			listCall.Context(ctx)
			resp, err = listCall.Do()
			if err != nil {
				return err
			}
			return nil
		}, c.getRetryOptions()...))
		if err != nil && !errors.Is(err, ErrAPIForbidden) {
			return subnetworks, fmt.Errorf("Can't get project subnetworks for project id %v: %w", projectID, err)
		}
		if err != nil && errors.Is(err, ErrAPIForbidden) {
			return subnetworks, nil
		}

		for _, v := range resp.Items {
			if v.Subnetworks != nil && len(v.Subnetworks) > 0 {
				subnetworks = append(subnetworks, v.Subnetworks...)
			}
		}

		if resp.NextPageToken == "" {
			break
		}
		nextPageToken = resp.NextPageToken
	}

	log.Debug().Msgf("Retrieved %v subnetworks for project %v", len(subnetworks), projectID)

	return
}

func (c *client) GetProjectSubnetworks(ctx context.Context, projects []*crmv1.Project) (subnetworks []*computev1.Subnetwork, err error) {

	// http://jmoiron.net/blog/limiting-concurrency-in-go/
	semaphore := make(chan bool, c.concurrency)
	cancelled := false

	resultChannel := make(chan struct {
		Subnetworks []*computev1.Subnetwork
		Err         error
	}, len(projects))

	for _, p := range projects {
		select {
		// try to fill semaphore up to it's full size otherwise wait for a routine to finish
		case semaphore <- true:
			go func(ctx context.Context, p *crmv1.Project) {
				// lower semaphore once the routine's finished, making room for another one to start
				defer func() { <-semaphore }()

				subnetworks, err := c.getProjectSubnetworks(ctx, p.ProjectId)

				resultChannel <- struct {
					Subnetworks []*computev1.Subnetwork
					Err         error
				}{subnetworks, err}
			}(ctx, p)

		case <-ctx.Done():
			log.Info().Msg("User has canceled execution, stopping retrieval of subnetworks...")
			cancelled = true
		}
		if cancelled {
			log.Info().Msg("User has canceled execution, waiting for pending retrieval of subnetworks to finish...")
			break
		}
	}

	// try to fill semaphore up to it's full size which only succeeds if all routines have finished or execution has been canceled
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- true
	}

	if cancelled {
		log.Info().Msg("User has canceled execution, checking retrieved subnetworks...")
	}

	// check for errors and aggregate all subnetworks
	close(resultChannel)
	for r := range resultChannel {
		if r.Err != nil {
			err = r.Err
			return
		}
		subnetworks = append(subnetworks, r.Subnetworks...)
	}

	return
}

func (c *client) isRetryableErrorCustomOption() foundation.RetryOption {
	return func(c *foundation.RetryConfig) {
		c.IsRetryableError = func(err error) bool {
			switch e := err.(type) {
			case *googleapi.Error:
				// Retry on 429 and 5xx, according to
				// https://cloud.google.com/storage/docs/exponential-backoff.
				return e.Code == http.StatusTooManyRequests || (e.Code >= 500 && e.Code < 600)
			case *url.Error:
				// Retry socket-level errors ECONNREFUSED and ENETUNREACH (from syscall).
				// Unfortunately the error type is unexported, so we resort to string
				// matching.
				retriable := []string{"connection refused", "connection reset"}
				for _, s := range retriable {
					if strings.Contains(e.Error(), s) {
						return true
					}
				}
				return false
			case interface{ Temporary() bool }:
				return e.Temporary()
			default:
				return false
			}
		}
	}
}

func (c *client) substituteErrorsWithPredefinedErrors(err error) error {
	if err == nil {
		return nil
	}

	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusForbidden {
		return ErrAPIForbidden.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusBadRequest && err.Error() == "googleapi: Error 400: Unknown project id: 0, invalid" {
		return ErrUnknownProjectID.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusBadRequest && strings.HasSuffix(err.Error(), "has not enabled BigQuery., invalid") {
		return ErrAPINotEnabled.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusBadRequest && strings.HasSuffix(err.Error(), "Invalid request: Invalid request since instance is not running., invalid") {
		return ErrEntityNotActive.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusNotFound && err.Error() == "googleapi: Error 404: The requested project was not found., notFound" {
		return ErrProjectNotFound.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusNotFound {
		return ErrEntityNotFound.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusNoContent {
		return ErrEntityNotFound.wrap(err)
	}

	return err
}

func (c *client) getRetryOptions() []foundation.RetryOption {
	return []foundation.RetryOption{
		c.isRetryableErrorCustomOption(),
		foundation.LastErrorOnly(true),
		foundation.Attempts(5),
		foundation.DelayMillisecond(1000),
	}
}
