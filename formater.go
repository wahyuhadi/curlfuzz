package curlfuzz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type BodySave struct {
	Body bytes.Buffer
}

func (c *Config) Format(req *http.Request) *http.Request {
	var b bytes.Buffer

	if req.Body != nil {
		b.ReadFrom(req.Body)
	}

	var bodysave BodySave
	bodysave.Body = b
	newBody := io.NopCloser(&b)

	//Modified Request@!
	if req.Method != "OPTIONS" && req.Method != "HEAD" {
		var new_raw http.Request
		new_raw.Method = req.Method
		new_raw.Proto = req.Proto
		new_raw.ProtoMajor = req.ProtoMajor
		new_raw.ProtoMinor = req.ProtoMinor
		// new_raw.Host = "{{ Hostname }}"
		new_raw.Header = req.Header
		new_raw.URL = req.URL
		lenQuery := len(req.URL.Query())

		//Generate fuzzing in query param
		if lenQuery != 0 && req.Method == "GET" {
			c.gen_fuzzing_query(&new_raw, req)
		}

		uri, _ := url.Parse(req.URL.String())
		new_raw.URL = uri
		log.Println(uri)

		if req.Method == "POST" || req.Method == "PATCH" || req.Method == "PUT" {
			type_body := check_body_type(bodysave.Body.Bytes(), req)
			log.Println(type_body)
			if type_body == "json" {
				// json body
				c.gen_fuzzing_jbody(new_raw, *req, bodysave.Body)
			}

			if type_body == "form" {
				c.gen_fuzzing_form(&new_raw, req, bodysave.Body)
			}
		}
	}

	req.Body = newBody
	return req
}
func check_body_type(s []byte, req *http.Request) string {
	log.Println("Validating request body")
	var js map[string]interface{}
	if json.Unmarshal(s, &js) == nil {
		log.Println("Body type is json")
		return "json"
	}

	if strings.Contains(req.Header.Get("Content-Type"), "form") {
		log.Println("Body type is post form ")
		return "form"
	}
	log.Println("undifined body format")
	return "undifined"
}
func (c *Config) gen_fuzzing_query(new_raw, req *http.Request) {
	uri, _ := url.Parse(req.URL.String())
	new_raw.URL = uri
	array_query := create_index_query(req)
	for l, query := range array_query {
		uri.RawQuery = query
		dump, _ := httputil.DumpRequest(new_raw, true)
		newDumpRequest := new_dump(dump)
		c.create_template(*new_raw, newDumpRequest, fmt.Sprintf("query-%v", l))
	}
}

func create_index_query(req *http.Request) []string {
	log.Println("Creating query param index for fuzzing")
	var indexed []string
	indexed = snipper_query(req)
	return indexed
}

// dumping to http request format data
func new_dump(dump []byte) []byte {
	newDumpRequest := []byte{}
	for index, tmpDump := range bytes.Split(dump, []byte("\r\n")) {
		if index == 0 {
			newDumpRequest = append(newDumpRequest, append(tmpDump, "\r\n"...)...)
			continue
		}
		appendSpace := append([]byte("        "), tmpDump...)
		newDumpRequest = append(newDumpRequest, append(appendSpace, "\r\n"...)...)
	}
	return newDumpRequest
}

func snipper_bjson(req *http.Request, body map[string]interface{}) []map[string]interface{} {
	var bodyo []string
	var objects = []map[string]interface{}{}

	for k, _ := range body {
		bodyo = append(bodyo, k)
	}
	for _, object := range bodyo {
		var new_body = make(map[string]interface{})
		for k, i := range body {
			if object == k {
				log.Println("Indexing object body ", k)
				body := fmt.Sprintf("§%v§", "path")
				new_body[k] = body
			} else {
				log.Println("Indexing objec body ", k)
				new_body[k] = i
			}
		}
		objects = append(objects, new_body)
	}

	return objects
}
