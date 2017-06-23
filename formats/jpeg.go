// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// JPEG is defined in ITU-T T.81: http://www.w3.org/Graphics/JPEG/itu-t81.pdf.
package formats

import "os"
import "math"
import "encoding/binary"

import "fmt"
import "encoding/hex"

const (
	sequenceStart = 0xff
	jfif          = 0x4A464946
)

const (
	sof0Marker = 0xc0 // Start Of Frame (Baseline Sequential).
	sof1Marker = 0xc1 // Start Of Frame (Extended Sequential).
	sof2Marker = 0xc2 // Start Of Frame (Progressive).
	dhtMarker  = 0xc4 // Define Huffman Table.
	rst0Marker = 0xd0 // ReSTart (0).
	rst7Marker = 0xd7 // ReSTart (7).
	soiMarker  = 0xd8 // Start Of Image.
	eoiMarker  = 0xd9 // End Of Image.
	sosMarker  = 0xda // Start Of Scan.
	dqtMarker  = 0xdb // Define Quantization Table.
	driMarker  = 0xdd // Define Restart Interval.
	comMarker  = 0xfe // COMment.
	// "APPlication specific" markers aren't part of the JPEG spec per se,
	// but in practice, their use is described at
	// http://www.sno.phy.queensu.ca/~phil/exiftool/TagNames/JPEG.html
	app0Marker  = 0xe0
	app14Marker = 0xee
	app15Marker = 0xef
)

// Finds SOS (Start of Scan) index in the buffer.
// Assumes that the buffer is large enough to hold all the segments
func findSos(buffer []byte) (int, error) {
	if !testJpegHeader(buffer[:]) {
		return -1, fmt.Errorf("not a JPEG file")
	}
	pos := 2 // Skip initial 0xFF, SOI

	for {
		if len(buffer) < pos+2 {
			return -1, fmt.Errorf("not enough space in the buffer for another segment")
		}
		segmentID := buffer[pos+1]
		if buffer[pos] != sequenceStart || segmentID == 0 {
			return -1, fmt.Errorf("segment start not found at the expected position (%d), surroundings: %s",
				pos, hex.EncodeToString(buffer[pos-4:pos+4]))
		}
		if segmentID == sosMarker {
			return pos + 2, nil
		}
		if len(buffer) < pos+4 {
			return -1, fmt.Errorf("not enough space in the buffer for segment length")
		}
		segmentLen := int(binary.BigEndian.Uint16(buffer[pos+2 : pos+4]))
		fmt.Printf("Segment ID=%2X, length=%d\n", segmentID, segmentLen)

		pos += 2 /* id */ + segmentLen
		if len(buffer) <= pos {
			return -1, fmt.Errorf("not enough space in the buffer for section %2Xd", segmentID)
		}
	}
}

func ParseJpeg(f *os.File) (*DecodingResult, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	bufferSize := int(math.Min(float64(stat.Size()), float64(BufferSize)))
	buffer := make([]byte, bufferSize)

	var sizeSoFar int64
	bytesRead, err := f.Read(buffer[:])
	if err != nil {
		return nil, err
	}

	buffer = buffer[:bytesRead]

	pos, err := findSos(buffer)

	for {
		if err != nil {
			return nil, err
		}

		buffer = buffer[:bytesRead]

		for ; pos < len(buffer)-2; pos++ {
			if buffer[pos] == sequenceStart && buffer[pos+1] == eoiMarker {
				return &DecodingResult{sizeSoFar + int64(pos) + 2}, nil
			}
		}
		buffer[0] = buffer[len(buffer)-1]
		sizeSoFar += int64(len(buffer)) - 1
		bytesRead, err = f.Read(buffer[1:])
		bytesRead++
		pos = 0
	}

	return nil, fmt.Errorf("EOS marker not found")
}

func testJpegHeader(buffer []byte) bool {
	return buffer[0] == sequenceStart && buffer[1] == soiMarker
}

func IsJpeg(f *os.File) (bool, error) {
	buffer := make([]byte, 16)
	num, err := f.Read(buffer[:])
	if err != nil {
		return false, err
	}

	if num < 16 {
		return false, nil
	}

	return testJpegHeader(buffer), nil
}
