package proxy

import "bytes"

type RequestReader struct {
	*bytes.Buffer
}

func (m RequestReader) Close() error { return nil }
