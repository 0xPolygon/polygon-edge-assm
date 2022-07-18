package ports

import "Trapesys/polygon-edge-assm/internal/adapters/core"

// ICore is the core interface
type ICore interface {
	GetConfig() *core.Config
	GetNodesInfo() *core.Nodes
	GetCoreJSON() string
	GetCore() *core.Core
}
