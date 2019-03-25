package osc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"math"
	"strings"
	"unicode/utf8"
)

type Message struct {
	address   string
	arguments []*Argument
}

func (msg *Message) internal() {}

func NewMessage() *Message {
	return &Message{}
}

func (msg *Message) Clear() *Message {
	msg.address = ""
	msg.arguments = nil

	return msg
}

func (msg *Message) Address() string {
	return msg.address
}

func (msg *Message) SetAddress(address string) {
	msg.address = address
}

func (msg *Message) Arguments() []*Argument {
	return msg.arguments
}

func (msg *Message) AddArguments(args ...*Argument) *Message {
	msg.arguments = append(msg.arguments, args...)
	return msg
}

func (msg *Message) MarshalBinary() ([]byte, error) {
	typeTag := strings.Builder{}
	dataBinary := strings.Builder{}

	for _, arg := range msg.arguments {
		typeTag.WriteString(arg.typeTag)
		argBinary, _ := arg.MarshalBinary()
		dataBinary.Write(argBinary)
	}

	messageBinary := strings.Builder{}
	messageBinary.WriteString(createOSCString(msg.address))
	messageBinary.WriteString(createOSCString("," + typeTag.String()))
	messageBinary.WriteString(dataBinary.String())

	return []byte(messageBinary.String()), nil
}

func (msg *Message) UnmarshalBinary(data []byte) error {
	if data == nil || len(data) == 0 {
		return errInvalidData
	}

	dataReader := bufio.NewReader(bytes.NewReader(data))

	address, err := getOSCString(dataReader)
	if err != nil {
		return err
	}
	msg.address = string(address)

	typeTags, err := getOSCString(dataReader)
	if err != nil {
		return err
	}

	if !checkBracketsBalance(typeTags) {
		return errInvalidData
	}

	typeTagsReader := bufio.NewReader(bytes.NewReader(typeTags))
	ch, err := typeTagsReader.ReadByte()
	if err != nil {
		return err
	}
	if ch != ',' {
		return errInvalidData
	}

	arguments, err := msg.unmarshalArguments(typeTagsReader, dataReader)
	if err != nil {
		return err
	}
	msg.arguments = append(msg.arguments, arguments...)

	return nil
}

func (msg *Message) unmarshalArguments(typeTags *bufio.Reader, data *bufio.Reader) ([]*Argument, error) {
	arguments := make([]*Argument, 0)

	for {
		if typeTags.Buffered() == 0 {
			break
		}

		tt, err := typeTags.ReadByte()
		if err != nil {
			return nil, err
		}

		switch tt {
		case 'i':
			v, err := getOSCData(data, 4)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetInt32(int32(binary.BigEndian.Uint32(v)))
			arguments = append(arguments, arg)

		case 'f':
			v, err := getOSCData(data, 4)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetFloat32(math.Float32frombits(binary.BigEndian.Uint32(v)))
			arguments = append(arguments, arg)

		case 's':
			v, err := getOSCString(data)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetString(string(v))
			arguments = append(arguments, arg)

		case 'b':
			v, err := getOSCBlob(data)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetBlob(v)
			arguments = append(arguments, arg)

		case 'h':
			v, err := getOSCData(data, 8)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetInt64(int64(binary.BigEndian.Uint64(v)))
			arguments = append(arguments, arg)

		case 't':
			v, err := getOSCData(data, 8)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetTimeTag(timeTagToTime(v))
			arguments = append(arguments, arg)

		case 'd':
			v, err := getOSCData(data, 8)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetFloat64(math.Float64frombits(binary.BigEndian.Uint64(v)))
			arguments = append(arguments, arg)

		case 'S':
			v, err := getOSCString(data)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetSymbol(string(v))
			arguments = append(arguments, arg)

		case 'c':
			v, err := getOSCData(data, 4)
			if err != nil {
				return nil, err
			}

			ch, _ := utf8.DecodeRune(v)
			arg := NewArgument().SetASCIICharacter(ch)
			arguments = append(arguments, arg)

		case 'r':
			v, err := getOSCData(data, 4)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetRGBAColor(v[0], v[1], v[2], v[3])
			arguments = append(arguments, arg)

		case 'm':
			v, err := getOSCData(data, 4)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetMIDI(v[0], v[1], v[2], v[3])
			arguments = append(arguments, arg)

		case 'T':
			arg := NewArgument().SetBoolean(true)
			arguments = append(arguments, arg)

		case 'F':
			arg := NewArgument().SetBoolean(false)
			arguments = append(arguments, arg)

		case 'N':
			arg := NewArgument().SetNil()
			arguments = append(arguments, arg)

		case 'I':
			arg := NewArgument().SetImpulse()
			arguments = append(arguments, arg)

		case '[':
			argArray, err := msg.unmarshalArguments(typeTags, data)
			if err != nil {
				return nil, err
			}

			arg := NewArgument().SetArray(argArray...)
			arguments = append(arguments, arg)

		case ']':
			return arguments, nil

		default:
			return nil, errInvalidData
		}
	}

	return arguments, nil
}
