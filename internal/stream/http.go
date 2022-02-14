package stream

type HttpBody interface {
	Parse([]byte)
}
