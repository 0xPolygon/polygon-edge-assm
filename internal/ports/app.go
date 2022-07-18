package ports

import (
	"Trapesys/polygon-edge-assm/internal/adapters/core"
)

type IApp interface {
	LambdaHandler(request core.Core) (string, error)
}
