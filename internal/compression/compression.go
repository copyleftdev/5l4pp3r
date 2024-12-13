package compression

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
)

type Compressor interface {
	Compress([]byte) ([]byte, error)
}

func NewCompressor(algorithm string, level int) (Compressor, error) {
	switch algorithm {
	case "zlib":
		return &ZlibCompressor{level: level}, nil
	case "gzip":
		return &GzipCompressor{level: level}, nil
	default:
		return nil, fmt.Errorf("unsupported compression algorithm: %s", algorithm)
	}
}

type ZlibCompressor struct {
	level int
}

func (z *ZlibCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, z.level)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type GzipCompressor struct {
	level int
}

func (g *GzipCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, g.level)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
