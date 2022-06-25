package utils

import (
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
)

type OneFile struct {
	fs         fs.FS
	fileServer http.Handler
	fallback   string
	Index      string
}

func New(fsys fs.FS, fallback string) *OneFile {
	return &OneFile{
		fs:         fsys,
		fileServer: http.FileServer(http.FS(fsys)),
		fallback:   fallback,
		Index:      "index.html",
	}
}

func (o *OneFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path[len(path)-1] == '/' {
		path += o.Index
	}
	path = path[1:]
	//withGzipped := utils.Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path == "/home.html" {
	//	w.WriteHeader(http.StatusOK)
	//	w.Write(homeFile)
	//	} else if _, err := fs.Stat(root, r.URL.Path[1:]); err == nil {
	//		fileServer.ServeHTTP(w, r)
	//	} else {
	//		_, filename := filepath.Split(r.URL.Path)
	//		if filepath.Ext(r.URL.Path) == ".js" || filepath.Ext(r.URL.Path) == ".css" {
	//			if r.Header.Get("if-none-match") == filename {
	//				w.WriteHeader(http.StatusNotModified)
	//				return
	//			}
	//			w.Header().Set("etag", filename)
	//		}
	//		w.WriteHeader(http.StatusOK)
	//		w.Write(indexFile)
	//	}
	//}))
	if path == o.Index {
		f, err := o.fs.Open(o.Index)
		if err != nil {
			if err == fs.ErrNotExist {
				//w.WriteHeader(http.StatusNotFound)
				if o.fallback == "" {
					w.WriteHeader(http.StatusNotFound)
				} else {
					f, err := o.fs.Open(o.fallback)
					if err != nil {
						if err == fs.ErrNotExist {
							w.WriteHeader(http.StatusNotFound)
						} else {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte(err.Error()))
						}
					}
					//w.WriteHeader(http.StatusOK)
					w.WriteHeader(http.StatusNotFound)
					io.Copy(w, f)
					f.Close()
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
		}
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
		f.Close()
	}

	if _, err := fs.Stat(o.fs, path); err == nil {

		o.fileServer.ServeHTTP(w, r)
	} else {
		_, filename := filepath.Split(path)
		if filepath.Ext(path) != "" {
			if r.Header.Get("if-none-match") == filename {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if o.fallback == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			f, err := o.fs.Open(o.fallback)
			if err != nil {
				if err == fs.ErrNotExist {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
				}
			}
			w.WriteHeader(http.StatusOK)
			io.Copy(w, f)
			f.Close()
		}

	}
}
