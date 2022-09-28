package utils

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const cachePath = "cache"

type Root struct {
	fs         fs.FS
	fileServer http.Handler
	fallback   string
	Index      string
}

func NewRoot(fsys fs.FS, fallback string, index string) *Root {
	if index == "" {
		return &Root{
			fs:         fsys,
			fileServer: http.FileServer(http.FS(fsys)),
			fallback:   fallback,
			Index:      "index.html",
		}
	} else {
		return &Root{
			fs:         fsys,
			fileServer: http.FileServer(http.FS(fsys)),
			fallback:   fallback,
			Index:      index,
		}
	}
}

func (t *Root) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	//log.Println("url: ", path)
	sub := strings.Split(path, "/")
	//log.Println(sub, len(sub))
	if len(sub) == 2 && sub[1] != "" && filepath.Ext(sub[1]) == "" {
		path += "/" + t.Index
	} else {
		if path[len(path)-1] == '/' {
			path += t.Index
		}
	}
	path = path[1:]
	//log.Println("resolved path: ", path)
	//index
	if path == t.Index {
		f, err := t.fs.Open(t.Index)
		if err != nil {
			t.nofound(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		f.Close()
		return
	}

	if _, err := fs.Stat(t.fs, path); err == nil {
		//_, filename := filepath.Split(path)
		if r.Header.Get("if-none-match") == path {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("etag", path)
		//log.Println(path)
		f, err := t.fs.Open(path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		f.Close()
	} else {
		//log.Println("error")
		w.WriteHeader(http.StatusNotFound)
	}

}

func (t *Root) nofound(w http.ResponseWriter) {
	if t.fallback == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		f, err := t.fs.Open(t.fallback)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		//w.WriteHeader(http.StatusNotFound)
		io.Copy(w, f)
		f.Close()
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

}

func Redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("nofound to: %s", target)
	http.Redirect(w, req, target,
		// see comments below and consider the codes 308, 302, or 301
		http.StatusTemporaryRedirect)
}
