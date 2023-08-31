package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var (
	port        int64
	service     string
	servicePort string
	destServer  string
)

func init() {
	flag.Int64Var(&port, "port", 9999, "port to listen to")
	flag.StringVar(&service, "service", "", "service which is running on localhost")
	flag.StringVar(&servicePort, "service.port", "8080", "port of the service which is running on localhost")
	flag.Parse()

	var envs map[string]string
	envs, err := godotenv.Read(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	destServer = envs["DEST_SERVER"]
	if destServer == "" {
		log.Fatal("Ensure you have a valid .env file with DEST_SERVER variable")
	}
}

// Serve a reverse proxy for a given url
func serveReverseProxy(res http.ResponseWriter, req *http.Request, target string) { // parse the parsedUrl
	parsedUrl, _ := url.Parse(target)
	proxy := &httputil.ReverseProxy{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			TLSNextProto:    make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
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
		fmt.Println(destServer)
		flag.Usage()
		os.Exit(1)
	}

	router := mux.NewRouter()

	router.PathPrefix(fmt.Sprintf("/%s/", service)).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newpath := strings.TrimPrefix(r.URL.String(), "/"+service)
		newpath = "http://localhost:" + servicePort + "/" + strings.TrimPrefix(newpath, "/")
		log.Printf("%s -> %s\n", r.URL.String(), newpath)
		serveReverseProxy(w, r, newpath)
	}))

	router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newpath := destServer + ":9999" + r.URL.String()
		log.Printf("%s -> %s\n", r.URL.String(), newpath)
		serveReverseProxy(w, r, newpath)
	}))

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("localhost:%d", port),
	}

	fmt.Printf("Redirecting /%s/* requests to http://localhost:%s/*\n", service, servicePort)
	fmt.Printf("Running dev gateway on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServeTLS("./server.crt", "./server.key"))
}
