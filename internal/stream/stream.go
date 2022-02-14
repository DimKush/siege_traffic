package stream

import (
	"bufio"
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"io"
	"log"
	"net/http"
)

type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
	request        HttpBody
	response       HttpBody
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)

	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {
			log.Println("Error reading stream", h.net, h.transport, ":", err)
		} else {
			bodyBytesReq := tcpreader.DiscardBytesToEOF(req.Body)
			req.Body.Close()
			log.Println("Received request from stream", h.net, h.transport, ":", req, "with", bodyBytesReq, "bytes in request body")
		}

		bufResp := bufio.NewReader(&h.r)
		resp, err := http.ReadResponse(bufResp, req)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {
			log.Println("Error reading stream", h.net, h.transport, ":", err)
		} else {
			bodyBytesResp := tcpreader.DiscardBytesToEOF(req.Body)
			resp.Body.Close()
			log.Println("Received response from stream", h.net, h.transport, ":", resp, "with", bodyBytesResp, "bytes in request body")
		}
	}
}
