package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

type transport struct {
	http.RoundTripper
}

func debugPrint(dumps string) {
	fmt.Println("------- Debug start --------")
	fmt.Print(dumps)
	fmt.Println("------- Debug end --------")
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	dumps, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	debugPrint(string(dumps))

	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	resp.Body = body
	_body, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	debugPrint(string(_body))
	return resp, nil
}

func main() {
	var dst string
	var port int
	flag.StringVar(&dst, "dst", "http://localhost:8080", "proxy destination")
	flag.IntVar(&port, "port", 8081, "listen port")
	flag.Parse()
	u, err := url.Parse(dst)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// see: http://stackoverflow.com/questions/31535569/golang-how-to-read-response-body-of-reverseproxy
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &transport{http.DefaultTransport}
	http.Handle("/", proxy)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), proxy)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
