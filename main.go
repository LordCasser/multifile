package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"multifile/utils"
	"net/http"
	"os"
)

//go:embed static
var static embed.FS

func main() {
	IP := flag.String("IP", "127.0.0.1", "ip address of server, default to be 127.0.0.1")
	SSL := flag.Bool("SSL", true, "open SSL or not, default is true")
	Port := flag.Uint("Port", 443, "port to be used, default is 443")
	flag.Parse()
	utils.CreateCert(*IP)
	_, crtErr := os.Stat("resources/https.crt")
	_, keyErr := os.Stat("resources/https.key")

	if *SSL && (crtErr == nil && keyErr == nil) {
		portString := fmt.Sprintf(":%d", *Port)
		fsys, _ := fs.Sub(static, "static")
		handle := utils.New(fsys, "404.html")
		withGzipped := utils.Gzip(handle)

		http.Handle("/", handle)
		_ = http.ListenAndServe(":8080", nil)
		log.Fatal(http.ListenAndServeTLS(portString, "resources/https.crt", "resources/https.key", nil))
	} else {
		portString := fmt.Sprintf(":%d", *Port)
		log.Fatal(http.ListenAndServe(portString, nil))
	}

}

/*
	withGzipped := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/home.html" {
			w.WriteHeader(http.StatusOK)
			w.Write(homeFile)
		} else if _, err := fs.Stat(root, r.URL.Path[1:]); err == nil {
			fileServer.ServeHTTP(w, r)
		} else {
			_, filename := filepath.Split(r.URL.Path)
			if filepath.Ext(r.URL.Path) == ".js" || filepath.Ext(r.URL.Path) == ".css" {
				if r.Header.Get("if-none-match") == filename {
					w.WriteHeader(http.StatusNotModified)
					return
				}
				w.Header().Set("etag", filename)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(indexFile)
		}
	}))
*/
