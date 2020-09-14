package network

type Type string

const (
	TypeNode    Type = "node"
	TypePod     Type = "pod"
	TypeService Type = "service"
	TypeMaster  Type = "master"
	TypeOther   Type = "other"

	TypeUnknown Type = ""
)
