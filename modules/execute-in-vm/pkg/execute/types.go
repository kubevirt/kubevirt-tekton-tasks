package execute

import "time"

type RemoteExecutor interface {
	Init(ipAddress string) error
	TestConnection() bool
	RemoteExecute(timeout time.Duration) error
}
