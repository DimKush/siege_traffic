package stream

import "github.com/google/uuid"

type ResponseHttp struct {
	UID          uuid.UUID
	protocolType string
}

func (h *ResponseHttp) Parse(body []byte) {

}

func NewHttpResponse(uid uuid.UUID) HttpBody {
	return &ResponseHttp{UID: uid}
}
