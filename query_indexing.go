package curlfuzz

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func snipper_query(req *http.Request) []string {
	var array_string []string
	url_query := req.URL.Query()

	var param []string
	for k, _ := range url_query {
		param = append(param, k)
	}
	for _, parv := range param {
		count := 1
		var query_raw string
		for k, i := range url_query {
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
