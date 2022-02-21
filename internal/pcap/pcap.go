package pcap

import (
	"errors"
	"fmt"
	"github.com/DimKush/siege_traffic/internal/dns"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/rs/zerolog/log"
)

type handler struct {
	device     string
	pcapHandle *pcap.Handle
}

type Client struct {
	ipAddresses, devicesPrefix                     []string
	handlers                                       []*handler
	bpfFilter                                      string
	disableDNSResolve                              bool
	wg                                             sync.WaitGroup
	lookup                                         dns.Lookup
	clientPackageSizes                             []int
	serverPackageSizes                             []int
	badClientPackagesCount, badServerPackagesCount int
	clientIP, serverIP                             string
	clientPort, serverPort                         uint32
}

func NewPcapClient(clientIP, serverIP string, device string, lookup dns.Lookup) *Client {
	client := &Client{
		handlers: make([]*handler, 0),
		lookup:   lookup,
		clientIP: clientIP,
		serverIP: serverIP,
	}

	if err := client.isDeviceAvaible(device); err != nil {
		fmt.Println("Device is not avaible")
		os.Exit(1)
	}

	for _, h := range client.handlers {
		go client.listen(h)
	}

	return client
}

func (c *Client) isDeviceAvaible(device string) error {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Println("No devices avaible")
		os.Exit(1)
	}

	fmt.Println("Current", device)
	for _, d := range devices {

		if d.Name != device {
			continue
		}
		h, err := c.getHandler(d.Name, c.bpfFilter)
		if err != nil {
			continue
		}

		c.handlers = append(c.handlers, &handler{
			device:     d.Name,
			pcapHandle: h,
		})

		for _, address := range d.Addresses {
			c.ipAddresses = append(c.ipAddresses, address.IP.String())
		}
	}

	if len(c.handlers) == 0 {
		return errors.New("No handler found")
	}

	return nil
}

func (c *Client) getHandler(device, filter string) (*pcap.Handle, error) {
	handle, err := pcap.OpenLive(device, 65535, false, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	if c.bpfFilter != "" {
		if err := handle.SetBPFFilter(filter); err != nil {
			handle.Close()
			return nil, err
		}
	}

	return handle, nil
}

func (c *Client) parsePacket(device string, packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return
	}

	ipv4pkg := ipLayer.(*layers.IPv4)
	if ipv4pkg == nil {
		return
	}

	srcIP := ipv4pkg.SrcIP.String()
	dstIP := ipv4pkg.DstIP.String()

	var dataLen int

	if ipv4pkg.Protocol.String() == "UDP" {

		if srcIP == c.serverIP || srcIP == c.clientIP || dstIP == c.serverIP || dstIP == c.clientIP {
			udpLayer := packet.Layer(layers.LayerTypeUDP)
			udpPkg, ok := udpLayer.(*layers.UDP)

			dataLen = len(udpPkg.Contents) + len(udpPkg.Payload)
			if ipv4pkg.SrcIP.String() == c.clientIP {
				c.clientPort = c.parsePort(udpPkg.SrcPort.String())
				c.serverPort = c.parsePort(udpPkg.DstPort.String())

				if ok {
					c.clientPackageSizes = append(c.clientPackageSizes, dataLen)
				} else {
					c.badClientPackagesCount++
				}
			} else if ipv4pkg.SrcIP.String() == c.serverIP {
				c.serverPort = c.parsePort(udpPkg.SrcPort.String())
				c.clientPort = c.parsePort(udpPkg.DstPort.String())

				if ok {
					c.serverPackageSizes = append(c.clientPackageSizes, dataLen)
				} else {
					c.badServerPackagesCount++
				}
			}

		}
	}
}

func (c *Client) listen(ph *handler) {
	c.wg.Add(1)
	defer c.wg.Done()

	packetSource := gopacket.NewPacketSource(ph.pcapHandle, ph.pcapHandle.LinkType())
	packetSource.Lazy = true
	packetSource.NoCopy = true
	packetSource.DecodeStreamsAsDatagrams = true
	for {
		select {
		case packet, ok := <-packetSource.Packets():
			if !ok {
				return
			}

			c.parsePacket(ph.device, packet)
		}
	}
}

func (c *Client) Close() {
	for _, h := range c.handlers {
		h.pcapHandle.Close()
	}
	c.wg.Wait()
}

func (c *Client) parsePort(s string) uint32 {
	idx := strings.Index(s, "(")
	if idx == -1 {
		i, _ := strconv.Atoi(s)
		return uint32(i)
	}

	i, _ := strconv.Atoi(s[:idx])
	return uint32(i)
}

func (c *Client) DropToLog() {
	rw := &sync.RWMutex{}
	rw.Lock()

	sumClPkg := 0
	maxClPkSz := 0
	for _, p := range c.clientPackageSizes {
		sumClPkg += p
		if maxClPkSz < p {
			maxClPkSz = p
		}
	}

	avgClPkgSz := sumClPkg / len(c.clientPackageSizes)

	log.Info().Msgf("sender: %s:%d => receiver: %s:%d sent bytes: %db, average UDP package size : %db, the biggest package size %db count of packages: %d, count of the bad packages: %d",
		c.clientIP, c.clientPort, c.serverIP, c.serverPort, sumClPkg, avgClPkgSz, maxClPkSz, len(c.clientPackageSizes), c.badClientPackagesCount,
	)

	sumSrvPkg := 0
	maxSrvPkSz := 0
	for _, p := range c.serverPackageSizes {
		sumSrvPkg += p
		if maxSrvPkSz < p {
			maxSrvPkSz = p
		}
	}

	avgSrvPkgSz := sumSrvPkg / len(c.serverPackageSizes)

	log.Info().Msgf("sender: %s:%d => receiver: %s:%d sent bytes: %db average UDP package size : %db the biggest package size %db, count of packages: %d, count of the bad packages: %d",
		c.serverIP, c.serverPort, c.clientIP, c.clientPort, sumSrvPkg, avgSrvPkgSz, maxSrvPkSz, len(c.serverPackageSizes), c.badServerPackagesCount,
	)

	c.badServerPackagesCount = 0
	c.badClientPackagesCount = 0
	c.serverPackageSizes = []int{}
	c.clientPackageSizes = []int{}

	rw.Unlock()
}
