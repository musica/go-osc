package osc

import (
	"bufio"
	"bytes"
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
	binaryData.WriteString(createOSCString(bundleIdentifier))
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

func (bnd *Bundle) UnmarshalBinary(data []byte) error {
	if data == nil || len(data) == 0 {
		return errInvalidData
	}

	dataReader := bufio.NewReader(bytes.NewReader(data))

	v, err := getOSCString(dataReader)
	if err != nil || string(v) != bundleIdentifier {
		return err
	}

	timeTag, err := getOSCData(dataReader, 8)
	if err != nil {
		return err
	}
	bnd.timeTag = timeTagToTime(timeTag)

	elements, err := bnd.unmarshalElements(dataReader)
	if err != nil {
		return err
	}
	bnd.elements = elements

	return nil
}

func (bnd *Bundle) unmarshalElements(data *bufio.Reader) ([]Packet, error) {
	packets := make([]Packet, 0)

	for {
		if data.Buffered() == 0 {
			break
		}

		elementSizeBytes, err := getOSCData(data, 4)
		if err != nil {
			return nil, err
		}

		elementSize := int32(binary.BigEndian.Uint32(elementSizeBytes))
		if elementSize == 0 {
			continue
		}

		elementContents, err := getOSCData(data, int(elementSize))
		if err != nil {
			return nil, err
		}

		if bytes.HasPrefix(elementContents, []byte(bundleIdentifier)) {
			bundle := NewBundle()
			if err := bundle.UnmarshalBinary(elementContents); err != nil {
				return nil, err
			}
			packets = append(packets, bundle)
		} else {
			message := NewMessage()
			if err := message.UnmarshalBinary(elementContents); err != nil {
				return nil, err
			}
			packets = append(packets, message)
		}
	}

	return packets, nil
}
