package httptool

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
)

// CustomHTTPRequest Send http request with custom headers
func CustomHTTPRequest(method string, addr string, data []byte, client *http.Client, updHeader func(req *http.Request) *http.Request) ([]byte, error) {
	var req *http.Request
	var err error

	if len(data) > 0 {
		requestData := bytes.NewBuffer(data)
		req, err = http.NewRequest(method, addr, requestData)
	} else {
		req, err = http.NewRequest(method, addr, nil)
	}

	if err != nil {
		return []byte{}, err
	}

	req = updHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Body, _ = gzip.NewReader(resp.Body)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

// GetHTTPRequest send GET http request with custom headers
func GetHTTPRequest(addr string, client *http.Client, updHeader func(req *http.Request) *http.Request) ([]byte, error) {
	return CustomHTTPRequest("GET", addr, []byte{}, client, updHeader)
}

// PostFormHTTPRequest send POST (from) http request with custom headers
func PostFormHTTPRequest(addr string, data []byte, client *http.Client, updHeader func(req *http.Request) *http.Request) ([]byte, error) {
	return CustomHTTPRequest("POST", addr, data, client, func(req *http.Request) *http.Request {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		updHeader(req)

		return req
	})
}

// SetFirefoxHeaders return http request with firefox headers
func SetFirefoxHeaders(req *http.Request) *http.Request {
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:64.0) Gecko/20100101 Firefox/64.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "ru")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")

	return req
}

// SetChromeHeaders return http request with google chrome headers
func SetChromeHeaders(req *http.Request) *http.Request {
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("Accept-Encoding", "gzip, deflate")

	return req
}

// GetFirefoxRequest return GET request with firefox headers
func GetFirefoxRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return req, err
	}

	return SetFirefoxHeaders(req), err
}

// PostFirefoxRequest return POST request with firefox headers
func PostFirefoxRequest(url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return req, err
	}

	return SetFirefoxHeaders(req), err
}

// GetChromeRequest return GET request with google chrome headers
func GetChromeRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return req, err
	}

	return SetChromeHeaders(req), err
}

// PostChromeRequest return POST request with google chrome headers
func PostChromeRequest(url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return req, err
	}

	return SetChromeHeaders(req), err
}
