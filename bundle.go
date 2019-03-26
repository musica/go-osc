package osc

import (
	"encoding/binary"
	"strings"
	"time"
)

type Bundle struct {
	timeTag  time.Time
	elements []Packet
}

func (bnd *Bundle) internal() {}

func NewBundle() *Bundle {
	return &Bundle{}
}

func (bnd *Bundle) Clear() *Bundle {
	bnd.timeTag = time.Time{}
	bnd.elements = nil
	return bnd
}

func (bnd *Bundle) TimeTag() time.Time {
	return bnd.timeTag
}

func (bnd *Bundle) SetTimeTag(v time.Time) *Bundle {
	bnd.timeTag = v
	return bnd
}

func (bnd *Bundle) Elements() []Packet {
	return bnd.elements
}

func (bnd *Bundle) AddElements(elements ...Packet) *Bundle {
	bnd.elements = append(bnd.elements, elements...)
	return bnd
}

func (bnd *Bundle) MarshalBinary() ([]byte, error) {
	binaryData := strings.Builder{}
	binaryData.Write(createOSCString(bundleIdentifier))
	binaryData.Write(timeToTimeTag(bnd.timeTag))

	for _, packet := range bnd.elements {
		packetBinary, _ := packet.MarshalBinary()

		packetLength := make([]byte, 4)
		binary.BigEndian.PutUint32(packetLength, uint32(len(packetBinary)))

		binaryData.Write(packetLength)
		binaryData.Write(packetBinary)
	}

	return []byte(binaryData.String()), nil
}
