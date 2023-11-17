package logger

// LogKey for setting keys.
type LogKey uint8

// Log keys.
//
//go:generate stringer -output=stringer.LogKey.go -type=LogKey -linecomment
const (
	_           LogKey = iota
	Version            // version
	PanicReason        // panic_reason
	URL                // url
	Error              // err
	Reason             // reason
	Host               // host
	Port               // port
	Module             // module
	Environment        // environment
	Stack              // stack
	TaskID             // task_id
	TaskKind           // task_kind
)
