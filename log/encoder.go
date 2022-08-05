package log

type Encoder interface {
	Encode(*Entry) ([]byte, error)
}
