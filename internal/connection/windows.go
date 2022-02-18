package connection

import (
	"fmt"
	"github.com/DimKush/siege_traffic/internal/pcap"
	"github.com/shirou/gopsutil/net"
	"log"
	"sync"
)

type Socket struct {
	pid                 uint32
	wg                  sync.WaitGroup
	connectionsDescribe map[string]*pcap.Connection
}

func (ps *Socket) checkSocketConnections(protocol string) error {
	connections, err := net.Connections(protocol)
	if err != nil {
		return err
	}

	for _, connection := range connections {
		if connection.Pid != int32(ps.pid) {
			continue
		}

		local := pcap.LocalConnect{
			LocalIP:       connection.Laddr.IP,
			LocalPort:     connection.Laddr.Port,
			LocalProtocol: protocol,
		}

		remote := pcap.RemoteConnect{
			RemoteIP:   connection.Raddr.IP,
			RemotePort: connection.Raddr.Port,
		}

		connMeta := &pcap.Connection{Local: local, Remote: remote}

		_, ok := ps.connectionsDescribe[remote.RemoteIP]
		if !ok {
			ps.connectionsDescribe[remote.RemoteIP] = connMeta
		}

		fmt.Println(local, remote)
	}

	return nil
}

func NewSocket(pid uint32) *Socket {
	return &Socket{pid: pid, connectionsDescribe: map[string]*pcap.Connection{}}
}

func (p *Socket) RefreshSocket() (*Socket, error) {
	err := p.checkSocketConnections(pcap.ProtoUDP)
	if err != nil {
		log.Fatal("Socket ", err)
		return nil, err
	}

	return p, nil
}
