package formats

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type DecodingResult struct {
	DataSize int64
}

type parserFunc func(*os.File) (*DecodingResult, error)
type validatorFunc func(*os.File) (bool, error)

type Parser struct {
	format   string
	Parse    parserFunc
	Validate validatorFunc
}

var extensionList = map[string]string{
	"jpeg": "JPEG",
	"jpg":  "JPEG",
	"jpe":  "JPEG",
	"png":  "PNG",
}

var parserList = map[string]Parser{
	"JPEG": {"JPEG", ParseJpeg, IsJpeg},
	"PNG":  {"PNG", ParsePNG, IsPNG},
}

func (p *Parser) Format() string {
	return p.format
}

// GetParser returns the function to parse the given file
// The parser is determined from file extension only and
// parsing will fail if file contents mismatch it.
func GetParser(fileName string) (*Parser, error) {
	ext := strings.ToLower(path.Ext(fileName))
	if ext != "" && ext[0] == '.' {
		ext = ext[1:]
	}
	format, ok := extensionList[ext]
	if !ok {
		return nil, fmt.Errorf("unrecognized extension '%s'")
	}
	parser, _ := parserList[format]
	return &parser, nil
}
