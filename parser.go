package curlfuzz

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func (r *Request) ToJson(format bool) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	if format {
		encoder.SetIndent("", "  ")
	}
	_ = encoder.Encode(r)
	return string(buffer.Bytes())
}

func (r *Request) ToRequest() http.Request {
	var bodys io.Reader
	// make sure body is valid format
	// if body is "" from Request Struct so body is nil
	if r.Body == "" {
		bodys = nil
	} else {
		bodys = strings.NewReader(r.Body)
	}
	req, _ := http.NewRequest(r.Method, r.Url, bodys)
	// adding header to req
	for head, val := range r.Header {
		req.Header.Add(head, val)
	}
	return *req
}
