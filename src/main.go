package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var (
	port    int64
	service string
)

func init() {
	flag.Int64Var(&port, "port", 9999, "port to listen to")
	flag.StringVar(&service, "service", "", "service which is running on localhost")
	flag.Parse()
}

// Serve a reverse proxy for a given url
func serveReverseProxy(res http.ResponseWriter, req *http.Request, target string) { // parse the parsedUrl
	parsedUrl, _ := url.Parse(target)
	proxy := &httputil.ReverseProxy{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Director: func(req *http.Request) {
			req.URL = parsedUrl
		},
	}
	proxy.ServeHTTP(res, req)
}

func main() {
	if service == "" {
		fmt.Println("Please provide a service name (e.g. postgresql)")
		flag.Usage()
		log.Fatal()
	}

	router := mux.NewRouter()

	router.PathPrefix(fmt.Sprintf("/%s/", service)).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newpath := strings.TrimPrefix(r.URL.String(), "/"+service)
		newpath = "http://localhost:8080/" + strings.TrimPrefix(newpath, "/")
		log.Printf("%s -> %s\n", r.URL.String(), newpath)
		serveReverseProxy(w, r, newpath)
	}))

	router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newpath := "secretpath:9999" + r.URL.String()
		log.Printf("%s -> %s\n", r.URL.String(), newpath)
		serveReverseProxy(w, r, newpath)
	}))

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
	}

	fmt.Printf("Redirecting /%s/* requests to http://localhost:8080/*\n", service)
	fmt.Printf("Running dev gateway on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
