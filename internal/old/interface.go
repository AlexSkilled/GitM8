package old

type Processor interface {
	Process(payload []byte) (msg string, skip bool, err error)
}

type Notifier interface {
	Notify(payload string) error
}

type UpdateWatcher interface {
	Polling()
}

type Configuration interface {
	GetBool(string) bool
	GetInt(string) int
	GetInt32(string) int32
	GetInt64(string) int64
	GetString(string) string
}
