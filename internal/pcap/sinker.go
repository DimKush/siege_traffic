package pcap

import "sync"

type Sinker struct {
	mut         sync.Mutex
	utilization Utilization
}

func (c *Sinker) Fetch(seg Segment) {
	c.mut.Lock()
	defer c.mut.Unlock()

	if _, ok := c.utilization[seg.Connection]; !ok {
		c.utilization[seg.Connection] = &ConnectionInfo{
			Interface: seg.Interface,
		}
	}

	switch seg.Direction {
	case DirectionUpload:
		c.utilization[seg.Connection].UploadBytes += seg.DataLen
		c.utilization[seg.Connection].UploadPackets += 1

	case DirectionDownload:
		c.utilization[seg.Connection].DownloadBytes += seg.DataLen
		c.utilization[seg.Connection].DownloadPackets += 1
	}
}

func (c *Sinker) GetUtilization() Utilization {
	c.mut.Lock()
	defer c.mut.Unlock()

	utilization := c.utilization
	c.utilization = make(Utilization)
	return utilization
}

func NewSinker() *Sinker {
	return &Sinker{utilization: make(Utilization)}
}
