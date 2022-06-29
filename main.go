package main

import (
	"flag"
	"fmt"
	"log"
	"multifile/utils"
	"net/http"
	"os"
)

const home = "static"

func main() {
	//IP := flag.String("IP", "127.0.0.1", "ip address of server, default to be 127.0.0.1")
	SSL := flag.Bool("SSL", false, "open SSL or not, default is true")
	Port := flag.Uint("Port", 0, "port to be used, default is 443")
	flag.Parse()
	_, crtErr := os.Stat("resources/tls.crt")
	_, keyErr := os.Stat("resources/tls.key")
	root := utils.NewRoot(os.DirFS(home), "404.html", "index.html")
	//router := mux.NewRouter()
	//router.PathPrefix("/cache/").HandlerFunc(root.CacheHandler)
	//router.HandleFunc("/", root.ServeHttp)
	withGzipped := utils.Gzip(root)
	portString := fmt.Sprintf(":%d", *Port)
	if *SSL && (crtErr == nil && keyErr == nil) {
		if *Port == 0 {
			portString = fmt.Sprintf(":%d", 443)
		}
		go func() {
			http.ListenAndServe(":80", http.HandlerFunc(utils.Redirect))
		}()
		log.Fatal(http.ListenAndServeTLS(portString, "resources/tls.crt", "resources/tls.key", withGzipped))
	} else {
		if *Port == 0 {
			portString = fmt.Sprintf(":%d", 8080)
		}
		log.Fatal(http.ListenAndServe(portString, withGzipped))
	}
}
