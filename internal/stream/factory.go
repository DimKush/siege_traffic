package stream

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/google/uuid"
)

type HttpStreamFactory struct{}

func (st *HttpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	uidStream := uuid.New()

	stream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
		request:   NewHttpRequest(uidStream),
		response:  NewHttpResponse(uidStream),
	}

	go stream.run()

	return &stream.r
}
