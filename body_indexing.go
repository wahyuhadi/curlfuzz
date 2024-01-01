package curlfuzz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

// json body
func (c *Config) gen_fuzzing_jbody(new_raw, req http.Request, bodysave bytes.Buffer) {
	var body map[string]interface{}
	json.NewDecoder(&bodysave).Decode(&body)
	objs := create_index_bjson(&req, body)
	for l, o := range objs {
		fuzzedJSON, _ := json.Marshal(o)
		content, _ := strconv.ParseInt(string(fuzzedJSON), 10, 64)
		new_raw.Body = io.NopCloser(strings.NewReader(string(fuzzedJSON)))
		new_raw.ContentLength = content
		dump, _ := httputil.DumpRequest(&new_raw, true)
		newDumpRequest := new_dump(dump)
		c.create_template(new_raw, newDumpRequest, fmt.Sprintf("body-%v", l))
	}
}

func create_index_bjson(req *http.Request, body map[string]interface{}) []map[string]interface{} {

	var indexed = []map[string]interface{}{}
	log.Println("Creating body json index for fuzzing")
	indexed = snipper_bjson(req, body)

	return indexed
}
