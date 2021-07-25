package interfaces

type HttpProcessor interface {
	Process() (msg string, err error)
}
