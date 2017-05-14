package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"github.com/gorilla/mux"
	"net/url"
	"log"
	"os"
	"io"
)

const (
	defaultHandlerName = "localhost.com"
	defaultHeaderKey = "Hostname"
)

func debugFunc(w http.ResponseWriter, r *http.Request, )  {
	dump, _ := httputil.DumpRequest(r, true)
	fmt.Fprintln(w, string(dump))
	fmt.Println(string(dump), "\n", "------------------")
	return
}

func DumpDebugHandler(f http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		fmt.Println(string(dump), "\n", "------------------")
		f.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func HostnameHeaderHandler(f http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("ERROR: %s : %s\n", defaultHandlerName, err)
		}
		w.Header().Set(defaultHeaderKey, hostname)
		f.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func main(){
	bm := mux.NewRouter()
	bm.HandleFunc("/{ep:[a-z/0-9]+}", debugFunc)
	us := "http://localhost:8000/"
	u, _ := url.Parse(us)
	proxy := httputil.NewSingleHostReverseProxy(u)
	pp :=  DumpDebugHandler(HostnameHeaderHandler(proxy))
	err := http.ListenAndServe(":8081", pp)
	//m  := mux.NewRouter()
	//m.HandleFunc("/{ep:[a-z/0-9]+}", handler(proxy))
	//err := http.ListenAndServe(":8081", m)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = mux.Vars(r)["rest"]
		p.ServeHTTP(w, r)
	}
}