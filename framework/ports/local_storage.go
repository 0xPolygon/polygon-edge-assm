package ports

type ILocalStoragePort interface {
	GetEdge() error
	RunGenesisCmd(cmd []string) error
}
