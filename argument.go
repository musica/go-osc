package osc

import (
	"encoding/binary"
	"math"
	"strings"
	"time"
	"unicode/utf8"
)

type Argument struct {
	typeTag    string
	data       interface{}
	dataBinary []byte
}

func NewArgument() *Argument {
	return &Argument{}
}

func (arg *Argument) Clear() *Argument {
	arg.typeTag = ""
	arg.data = nil
	arg.dataBinary = nil
	return arg
}

func (arg *Argument) TypeTag() string {
	return arg.typeTag
}

func (arg *Argument) Data() interface{} {
	return arg.data
}

func (arg *Argument) DataBinary() []byte {
	if arg.dataBinary == nil {
		return nil
	}

	dataBinary := make([]byte, len(arg.dataBinary))
	copy(dataBinary, arg.dataBinary)

	return dataBinary
}

func (arg *Argument) SetInt32(v int32) *Argument {
	arg.typeTag = "i"
	arg.data = v

	arg.dataBinary = make([]byte, 4)
	binary.BigEndian.PutUint32(arg.dataBinary, uint32(v))

	return arg
}

func (arg *Argument) SetFloat32(v float32) *Argument {
	arg.typeTag = "f"
	arg.data = v

	arg.dataBinary = make([]byte, 4)
	binary.BigEndian.PutUint32(arg.dataBinary, math.Float32bits(v))

	return arg
}

func (arg *Argument) SetString(v string) *Argument {
	arg.typeTag = "s"
	arg.data = v

	arg.dataBinary = []byte(createOSCString(v))

	return arg
}

func (arg *Argument) SetBlob(v []byte) *Argument {
	arg.typeTag = "b"
	arg.data = v

	arg.dataBinary = createOSCBlob(v)

	return arg
}

func (arg *Argument) SetInt64(v int64) *Argument {
	arg.typeTag = "h"
	arg.data = v

	arg.dataBinary = make([]byte, 8)
	binary.BigEndian.PutUint64(arg.dataBinary, uint64(v))

	return arg
}

func (arg *Argument) SetTimeTag(v time.Time) *Argument {
	arg.typeTag = "t"
	arg.data = v

	arg.dataBinary = timeToTimeTag(v)

	return arg
}

func (arg *Argument) SetFloat64(v float64) *Argument {
	arg.typeTag = "d"
	arg.data = v

	arg.dataBinary = make([]byte, 8)
	binary.BigEndian.PutUint64(arg.dataBinary, math.Float64bits(v))

	return arg
}

func (arg *Argument) SetSymbol(v string) *Argument {
	arg.typeTag = "S"
	arg.data = v

	arg.dataBinary = []byte(createOSCString(v))

	return arg
}

func (arg *Argument) SetASCIICharacter(v rune) *Argument {
	arg.typeTag = "c"
	arg.data = v

	arg.dataBinary = make([]byte, 4)
	utf8.EncodeRune(arg.dataBinary, v)

	return arg
}

func (arg *Argument) SetRGBAColor(r uint8, g uint8, b uint8, a uint8) *Argument {
	arg.typeTag = "r"
	arg.data = []uint8{r, g, b, a}

	arg.dataBinary = []byte{r, g, b, a}

	return arg
}

func (arg *Argument) SetMIDI(portID byte, statusByte byte, data1 byte, data2 byte) *Argument {
	arg.typeTag = "m"
	arg.data = []byte{portID, statusByte, data1, data2}

	arg.dataBinary = []byte{portID, statusByte, data1, data2}

	return arg
}

func (arg *Argument) SetBoolean(v bool) *Argument {
	if v {
		arg.typeTag = "T"
	} else {
		arg.typeTag = "F"
	}
	arg.data = v

	arg.dataBinary = nil

	return arg
}

func (arg *Argument) SetNil() *Argument {
	arg.typeTag = "N"
	arg.data = nil

	arg.dataBinary = nil

	return arg
}

func (arg *Argument) SetImpulse() *Argument {
	arg.typeTag = "I"
	arg.data = nil

	arg.dataBinary = nil

	return arg
}

func (arg *Argument) SetArray(v ...*Argument) *Argument {
	arg.data = v

	typeTag := strings.Builder{}
	dataBinary := strings.Builder{}

	typeTag.WriteByte('[')

	for _, x := range v {
		typeTag.WriteString(x.typeTag)
		dataBinary.Write(x.dataBinary)
	}

	typeTag.WriteByte(']')

	arg.typeTag = typeTag.String()
	arg.dataBinary = []byte(dataBinary.String())

	return arg
}

func (arg *Argument) MarshalBinary() ([]byte, error) {
	return arg.dataBinary, nil
}
