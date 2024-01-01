package curlfuzz

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func (c *Config) gen_fuzzing_form(new_raw, req *http.Request, bodysave bytes.Buffer) {
	forms := create_index_form(req, bodysave)
	for l, form := range forms {
		content, _ := strconv.ParseInt(string(form), 10, 64)
		new_raw.ContentLength = content
		new_raw.Body = io.NopCloser(strings.NewReader(string(form)))
		dump, _ := httputil.DumpRequest(new_raw, true)
		newDumpRequest := new_dump(dump)
		fmt.Println(newDumpRequest)

		c.create_template(*new_raw, newDumpRequest, fmt.Sprintf("post-form-%v", l))
	}

}

func create_index_form(req *http.Request, bodysave bytes.Buffer) []string {
	var indexed []string
	log.Println("Creating body json index for fuzzing")
	indexed = snipper_form(req, bodysave)

	return indexed
}

func snipper_form(req *http.Request, bodysave bytes.Buffer) []string {
	var array_string []string
	form_body, _ := url.ParseQuery(bodysave.String())
	fmt.Println(form_body)

	var param []string
	for k, _ := range form_body {
		param = append(param, k)
	}

	for _, parv := range param {
		count := 1
		var query_raw string
		for k, i := range form_body {
			if k == parv {
				log.Println("Indexing object param ", k)
				query_raw += fmt.Sprintf("%s=%s&", k, fmt.Sprintf("§%v§", "path"))
			} else {
				log.Println("Indexing object param ", k)
				t := &url.URL{Path: i[0]}
				query_raw += fmt.Sprintf("%s=%s&", k, t.String())
			}
			count = count + 1
		}
		array_string = append(array_string, query_raw)
	}

	return array_string
}
