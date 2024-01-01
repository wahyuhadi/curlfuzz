package curlfuzz

func CurlToNuc(c, key string) {
	config := Config{
		Key:  key,
		Curl: c,
	}
	curl, _ := Parse(c)

	http_req := curl.ToRequest()
	config.Format(&http_req)
}
