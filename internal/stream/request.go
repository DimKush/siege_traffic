package stream

import "github.com/google/uuid"

type RequestHttp struct {
	UID          uuid.UUID
	protocolType string
}

func (h *RequestHttp) Parse(body []byte) {

}

func NewHttpRequest(uid uuid.UUID) HttpBody {
	return &RequestHttp{UID: uid}
}
