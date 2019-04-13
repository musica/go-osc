package osc

import (
	"encoding/binary"
	"math"
	"time"
	"unicode/utf8"
)

type Arguments struct {
	typeTags   []byte
	dataBinary []byte
}

func NewArguments() *Arguments {
	return &Arguments{}
}

func (args *Arguments) addTypeTag(typeTags ...byte) {
	args.typeTags = append(args.typeTags, typeTags...)
}

func (args *Arguments) addDataBinary(dataBinary ...byte) {
	args.dataBinary = append(args.dataBinary, dataBinary...)
}

func (args *Arguments) Clear() *Arguments {
	args.typeTags = nil
	args.dataBinary = nil
	return args
}

func (args *Arguments) TypeTags() []byte {
	if args.typeTags == nil {
		return nil
	}
	return copyBytes(args.typeTags)
}

func (args *Arguments) DataBinary() []byte {
	if args.dataBinary == nil {
		return nil
	}
	return copyBytes(args.dataBinary)
}

func (args *Arguments) AddInt32(v int32) *Arguments {
	args.addTypeTag('i')
	dataBinary := make([]byte, 4)
	binary.BigEndian.PutUint32(dataBinary, uint32(v))
	args.addDataBinary(dataBinary...)
	return args
}

func (args *Arguments) AddFloat32(v float32) *Arguments {
	args.addTypeTag('f')
	dataBinary := make([]byte, 4)
	binary.BigEndian.PutUint32(dataBinary, math.Float32bits(v))
	args.addDataBinary(dataBinary...)
	return args
}

func (args *Arguments) AddString(v string) *Arguments {
	args.addTypeTag('s')
	args.addDataBinary(createOSCString(v)...)
	return args
}

func (args *Arguments) AddBlob(v []byte) *Arguments {
	args.addTypeTag('b')
	args.addDataBinary(createOSCBlob(v)...)
	return args
}

func (args *Arguments) AddInt64(v int64) *Arguments {
	args.addTypeTag('h')
	dataBinary := make([]byte, 8)
	binary.BigEndian.PutUint64(dataBinary, uint64(v))
	args.addDataBinary(dataBinary...)
	return args
}

func (args *Arguments) AddTimeTag(v time.Time) *Arguments {
	args.addTypeTag('t')
	args.addDataBinary(timeToTimeTag(v)...)
	return args
}

func (args *Arguments) AddFloat64(v float64) *Arguments {
	args.addTypeTag('d')
	dataBinary := make([]byte, 8)
	binary.BigEndian.PutUint64(dataBinary, math.Float64bits(v))
	args.addDataBinary(dataBinary...)
	return args
}

func (args *Arguments) AddSymbol(v string) *Arguments {
	args.addTypeTag('S')
	args.addDataBinary(createOSCString(v)...)
	return args
}

func (args *Arguments) AddASCIICharacter(v rune) *Arguments {
	args.addTypeTag('c')
	dataBinary := make([]byte, 4)
	utf8.EncodeRune(dataBinary, v)
	args.addDataBinary(dataBinary...)
	return args
}

func (args *Arguments) AddRGBAColor(r uint8, g uint8, b uint8, a uint8) *Arguments {
	args.addTypeTag('r')
	args.addDataBinary(r, g, b, a)
	return args
}

func (args *Arguments) AddMIDI(portID byte, statusByte byte, data1 byte, data2 byte) *Arguments {
	args.addTypeTag('m')
	args.addDataBinary(portID, statusByte, data1, data2)
	return args
}

func (args *Arguments) AddBoolean(v bool) *Arguments {
	if v {
		args.addTypeTag('T')
	} else {
		args.addTypeTag('F')
	}
	return args
}

func (args *Arguments) AddNil() *Arguments {
	args.addTypeTag('N')
	return args
}

func (args *Arguments) AddImpulse() *Arguments {
	args.addTypeTag('I')
	return args
}

func (args *Arguments) AddArray(v *Arguments) *Arguments {
	args.addTypeTag('[')
	args.addTypeTag(v.typeTags...)
	args.addTypeTag(']')
	args.addDataBinary(v.dataBinary...)
	return args
}

func (args *Arguments) MarshalBinary() ([]byte, error) {
	return args.DataBinary(), nil
}
