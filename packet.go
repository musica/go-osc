package osc

import "encoding"

type Packet interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	internal()
}
