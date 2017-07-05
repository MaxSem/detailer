package formats

import "os"
import "fmt"
import "encoding/binary"
import "io"

const pngHeader = "\x89PNG\r\n\x1a\n"

func IsPNG(f *os.File) (bool, error) {
	buf := make([]byte, len(pngHeader))

	num, err := f.Read(buf)
	if err != nil {
		return false, err
	}

	if num != len(pngHeader) {
		return false, nil
	}

	return string(buf) == pngHeader, nil
}

func ParsePNG(f *os.File) (*DecodingResult, error) {
	isPng, err := IsPNG(f)
	if err != nil {
		return nil, err
	}
	if !isPng {
		return nil, fmt.Errorf("not a PNG file")
	}
	buf := make([]byte, 8)

	for {
		numRead, err := f.Read(buf[:8]) // length and chunk type fields
		if err == io.EOF {
			return nil, fmt.Errorf("unexpected end of file")
		}
		if err != nil {
			return nil, err
		}
		if numRead < 8 {
			return nil, fmt.Errorf("not enough data to read next chunk")
		}

		length := binary.BigEndian.Uint32(buf)
		chunkType := string(buf[4:8])
		pos, err := f.Seek(int64(length+4 /* CRC field */), io.SeekCurrent)
		if err != nil {
			return nil, err
		}
		if chunkType == "IEND" {
			return &DecodingResult{pos}, nil
		}
	}
}
