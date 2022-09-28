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
	SSL := flag.Bool("SSL", false, "open SSL or not, default is true")
	Port := flag.Uint("Port", 0, "port to be used, default is 443")
	flag.Parse()
	_, crtErr := os.Stat("resources/tls.crt")
	_, keyErr := os.Stat("resources/tls.key")
	root := utils.NewRoot(os.DirFS(home), "404.html", "index.html")

	withGzipped := utils.Gzip(root)
	if *SSL && (crtErr == nil && keyErr == nil) {
		if *Port == 0 {
			*Port = 443
		}
		if *Port == 443 {
			go func() {
				err := http.ListenAndServe(":80", http.HandlerFunc(utils.Redirect))
				if err != nil {
					log.Println(err)
				}
			}()
		}
		portString := fmt.Sprintf(":%d", *Port)

		log.Fatal(http.ListenAndServeTLS(portString, "resources/https.crt", "resources/https.key", nil))
	} else {
		if *Port == 0 {
			*Port = 80
		}
		portString := fmt.Sprintf(":%d", *Port)
		log.Fatal(http.ListenAndServe(portString, withGzipped))
	}
}
