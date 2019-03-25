package osc

import (
	"bufio"
	"encoding/binary"
	"errors"
	"strings"
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

func createOSCString(str string) string {
	oscString := strings.Builder{}
	oscString.WriteString(str + "\x00")

	for i := 0; i < getPaddingLength(len(str)+1, 4); i++ {
		oscString.WriteByte('\x00')
	}

	return oscString.String()
}

func createOSCBlob(blob []byte) []byte {
	blobLength := make([]byte, 4)
	binary.BigEndian.PutUint32(blobLength, uint32(len(blob)))

	oscBlob := strings.Builder{}
	oscBlob.Write(blobLength)
	oscBlob.Write(blob)

	for i := 0; i < getPaddingLength(len(blob), 4); i++ {
		oscBlob.WriteByte('\x00')
	}

	return []byte(oscBlob.String())
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

func checkBracketsBalance(data []byte) bool {
	bracketCount := 0

	for _, ch := range data {
		switch ch {
		case '[':
			bracketCount++
		case ']':
			bracketCount--
			if bracketCount < 0 {
				return false
			}
		}
	}

	return bracketCount == 0
}

func getOSCData(data *bufio.Reader, nBytes int) ([]byte, error) {
	if nBytes > data.Buffered() {
		return nil, errInvalidData
	}

	v := make([]byte, nBytes)
	if _, err := data.Read(v); err != nil {
		return nil, err
	}

	return v, nil
}

func getOSCString(data *bufio.Reader) ([]byte, error) {
	oscString, err := data.ReadBytes('\x00')
	if err != nil {
		return nil, err
	}

	oscString = oscString[:len(oscString)-1]

	paddingLength := getPaddingLength(len(oscString)+1, 4)
	if paddingLength > data.Buffered() {
		return nil, errInvalidData
	}

	if _, err := data.Discard(paddingLength); err != nil {
		return nil, err
	}

	return oscString, nil
}

func getOSCBlob(data *bufio.Reader) ([]byte, error) {
	oscBlobSizeBytes, err := getOSCData(data, 4)
	if err != nil {
		return nil, err
	}

	oscBlobSize := int32(binary.BigEndian.Uint32(oscBlobSizeBytes))
	oscBlob, err := getOSCData(data, int(oscBlobSize))
	if err != nil {
		return nil, err
	}

	paddingLength := getPaddingLength(len(oscBlob), 4)
	if paddingLength > data.Buffered() {
		return nil, errInvalidData
	}

	if _, err := data.Discard(paddingLength); err != nil {
		return nil, err
	}

	return oscBlob, nil
}
