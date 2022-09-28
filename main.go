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
const resources = "resources"

func FileCheckCreate(filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(filename, os.FileMode(0755))
			if err != nil {
				log.Panicf("make %s dir error\n%s\n", filename, err)
				return
			}
		} else {
			log.Panicf("create %s dir filed\n%s\n", filename, err)
			return
		}
	}
}

func main() {
	SSL := flag.Bool("SSL", false, "open SSL or not, default is true")
	Port := flag.Uint("Port", 0, "port to be used, default is 443")
	init := flag.Bool("init", false, "init the environment")

	flag.Parse()
	if *init && (Port != nil || SSL != nil) {
		log.Println("[X]Init don't need other agreements")
		flag.PrintDefaults()
		return
	} else {
		FileCheckCreate(home)
		FileCheckCreate(resources)
	}

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
