package osc

import (
	"encoding/binary"
	"errors"
	"time"
)

const (
	secondsFrom1900To1970 = 2208988800
	bundleIdentifier      = "#bundle"
)

var (
	errInvalidData = errors.New("invalid data")
)

func getPaddingLength(len int, multipleOf int) int {
	return (multipleOf - (len % multipleOf)) % multipleOf
}

func createOSCString(data string) []byte {
	paddingLength := getPaddingLength(len(data)+1, 4) + 1
	oscString := make([]byte, len(data)+paddingLength)
	copy(oscString, []byte(data))
	return oscString
}

func createOSCBlob(data []byte) []byte {
	paddingLength := getPaddingLength(len(data), 4)
	oscBloc := make([]byte, 4+len(data)+paddingLength)
	binary.BigEndian.PutUint32(oscBloc, uint32(len(data)))
	copy(oscBloc[4:], data)
	return oscBloc
}

func timeToTimeTag(v time.Time) []byte {
	msb32 := uint64((secondsFrom1900To1970 + v.Unix()) << 32)
	lsb32 := uint64(v.Nanosecond())

	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, msb32+lsb32)

	return data
}

func timeTagToTime(v []byte) time.Time {
	seconds := binary.BigEndian.Uint32(v[0:4]) - secondsFrom1900To1970
	nanoseconds := binary.BigEndian.Uint32(v[4:8])

	return time.Unix(int64(seconds), int64(nanoseconds))
}

func copyBytes(data []byte) []byte {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	return dataCopy
}
