package pcap

import (
	"github.com/google/gopacket/pcap"
)

type (
	Utilization map[Connection]*ConnectionInfo
	direction   uint8
)

const (
	ProtoUDP        string    = "udp"
	DirectionUpload direction = iota
	DirectionDownload
)

type RemoteConnect struct {
	RemoteIP   string
	RemotePort uint32
}

type LocalConnect struct {
	LocalIP       string
	LocalPort     uint32
	LocalProtocol string
}

type Connection struct {
	Local         LocalConnect
	Remote        RemoteConnect
	DirectionType direction
}

type ConnectionInfo struct {
	Interface       string
	UploadPackets   int
	DownloadPackets int
	UploadBytes     int
	DownloadBytes   int
}

type Segment struct {
	Interface  string
	DataLen    int
	Connection Connection
	Direction  direction
}

/*
func (p Process) String() string {
	return fmt.Sprintf("pid : %d Name : %s", p.PID, p.Name)
}
*/

func GetAllDevices() ([]pcap.Interface, error) {
	return pcap.FindAllDevs()
}
