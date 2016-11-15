package flv

import "io"
import "errors"
import "fmt"
import "encoding/binary"

type FlvHeader struct {
	// Flv Header (3+6 bytes)
	// Signature [3]byte = 'F', 'L', 'V'
	Version    int    // U8
	Flags      int    // U8
	HeaderSize uint32 // U32
}

func (f *FlvHeader) Parse(r io.Reader) (err error) {
	// SIGNATURE VERSION FLAGS HeaderSize
	//     3        1      1       2
	buf := make([]byte, 9)
	if _, err = io.ReadAtLeast(r, buf, 9); err != nil {
		return
	}

	if buf[0] != 'F' || buf[1] != 'L' || buf[2] != 'V' {
		err = errors.New("Not a valid FLV file")
		return
	}

	if buf[3] != '\x01' {
		err = fmt.Errorf("Unsupported FLV version: %d", int(buf[3]))
		return
	}

	f.Version = int(buf[3])
	f.Flags = int(buf[4])
	f.HeaderSize = binary.BigEndian.Uint32(buf[5:])
	return
}

type FlvTag struct {
	PrevTagSize uint32
	// Tag Header (11 bytes)
	Type        int // 1-byte, U8
	DataSize    int // 3-byte, U24
	Timestamp   int // 3-byte, U24
	TimestampEx int // 1-byte, U8, in case you need to extend the precision of Timestamp to 32bit
	StreamID    int // 3-byte, U24, should always be 0
	// Tag Data
	Data []byte
}

func (t *FlvTag) Parse(r io.Reader) (err error) {
	buf := make([]byte, 15)
	if _, err = io.ReadAtLeast(r, buf, 15); err != nil {
		return
	}

	t.PrevTagSize = binary.BigEndian.Uint32(buf[0:4])
	t.Type = int(buf[4])
	t.DataSize = readU24(buf[5:8])
	t.Timestamp = readU24(buf[8:11])
	t.TimestampEx = int(buf[11])
	t.StreamID = readU24(buf[12:15])
	return
}

func (t *FlvTag) ReadData(r io.Reader) (err error) {
	if len(t.Data) < t.DataSize {
		t.Data = make([]byte, t.DataSize)
	}
	_, err = io.ReadAtLeast(r, t.Data, t.DataSize)
	return
}

func readU24(buf []byte) int {
	return int(buf[2]) | int(buf[1])<<8 | int(buf[0])<<16
}
