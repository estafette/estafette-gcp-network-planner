# estafette-gcp-network-planner

Both a CLI and a library to plan your Google Cloud Platform network ranges to avoid overlapping ranges in separate projects, so they can easily be peered at a later stage.

## Usage

To install this cli on your mac run:

```bash
brew install estafette/stable/gcp-network-planner
```

Then run `estafette help` to see what commands are available.

### > estafette manifest validate

The check whether your .estafette.yaml manifest file is valid runL

```bash
gcp-network-planner suggest --filter labels.environment:dev
```

## Development

For local development when running `go build .` the generated binary can be used with

```bash
./gcp-network-planner help
```

Development versions get released as a pre-release version on github for each commit and have their brew formular updates in a development tap repository. You can install it via

```bash
brew install estafette/dev/gcp-network-planner-dev
```

And then use it with

```bash
gcp-network-planner-dev help
```

## Releases

Official releases are performed by creating and pushing a branch with the same version as specified in `version.semver.releaseBranch`.