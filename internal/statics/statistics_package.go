package statics

import (
	"fmt"
	"time"
)

type PackageStatistics struct {
	eventDate   time.Time
	processName string
	processPid  uint32
	ipSrc       string
	portSrc     uint32
	ipDst       string
	portDst     uint32
	packageSize int
	packageType string
}

func (p PackageStatistics) String() string {
	return fmt.Sprintf("%v | package type: %s process: %s - %d sender: %s:%v => reciver: %s:%v package size: %d",
		p.eventDate, p.packageType, p.processName, p.processPid, p.ipSrc, p.portSrc, p.ipDst, p.portDst, p.packageSize)
}
