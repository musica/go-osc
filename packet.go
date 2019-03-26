package osc

import "encoding"

type Packet interface {
	encoding.BinaryMarshaler
	internal()
}
