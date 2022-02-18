package sniffer

import (
	"github.com/DimKush/siege_traffic/internal/dns"
	"github.com/DimKush/siege_traffic/internal/option"
	"github.com/DimKush/siege_traffic/internal/pcap"
	"time"
)

var (
	// ProcessPID - pid of the process
	ProcessPID uint32
)

type Sniffer struct {
	options     *option.Options
	dnsResolver *dns.DNSResolver
	client      *pcap.Client
}

func NewSniffer(opts *option.Options) *Sniffer {
	dnsResolver := dns.NewDnsResolver()
	//pcapClient, err := pcap.NewPcapClient(opts.deviceName, opts.processPID, dnsResolver.Lookup)

	return &Sniffer{
		options:     opts,
		dnsResolver: dnsResolver,
		client:      pcap.NewPcapClient(opts.FilterClientIP, opts.FilterServerIP, opts.DeviceName, dnsResolver.Lookup),
	}
}

func (s *Sniffer) Start() {

	for {
		ticker := time.Tick(time.Duration(10 * time.Second))

		select {

		case <-ticker:
			s.Refresh()
		}
	}
}

func (s *Sniffer) Refresh() {
	s.client.DropToLog()
}

// logging
//pidSocket
//s.statistics.Put(statics.Statistic{OpenSockets: openSockets, Utilization: utilization})
//s.ui.viewer.Render(s.statsManager.GetStats())
