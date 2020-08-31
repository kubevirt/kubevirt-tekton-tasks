package execute

type RemoteExecutor interface {
	Init(ipAddress string) error
	TestConnection() bool
	RemoteExecute() error
}
