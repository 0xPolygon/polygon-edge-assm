package ports

import (
	"Trapesys/polygon-edge-assm/internal/adapters/core"
)

type ILambdaAPIPort interface {
	SetConfig(config core.Core)
	SetNodes(config core.Core)
	GetConfig() core.Config
	GetNodes() core.Nodes
}
