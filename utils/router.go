package utils

import (
	"io"
	"io/fs"
	"net/http"
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
	if path[len(path)-1] == '/' {
		path += t.Index
	}
	path = path[1:]
	//log.Println(path)
	//index
	if path == t.Index {
		f, err := t.fs.Open(t.Index)
		if err != nil {
			t.redirect(w)
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

	//resources
	//if _, err := fs.Stat(t.fs, path); err == nil {
	//	t.fileServer.ServeHTTP(w, r)
	//} else {
	//	_, filename := filepath.Split(path)
	//	if filepath.Ext(path) != "" {
	//		if r.Header.Get("if-none-match") == filename {
	//			w.WriteHeader(http.StatusNotModified)
	//			return
	//		}
	//		w.WriteHeader(http.StatusNotFound)
	//		return
	//	}
	//	if t.fallback == "" {
	//		w.WriteHeader(http.StatusNotFound)
	//	} else {
	//		f, err := t.fs.Open(t.fallback)
	//		if err != nil {
	//			if err == fs.ErrNotExist {
	//				w.WriteHeader(http.StatusNotFound)
	//			} else {
	//				w.WriteHeader(http.StatusInternalServerError)
	//				w.Write([]byte(err.Error()))
	//			}
	//		}
	//		w.WriteHeader(http.StatusOK)
	//		io.Copy(w, f)
	//		f.Close()
	//	}
	//}
}

func (t *Root) redirect(w http.ResponseWriter) {
	if t.fallback == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		f, err := t.fs.Open(t.fallback)
		if err != nil {
			if err == fs.ErrNotExist {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		//w.WriteHeader(http.StatusNotFound)
		io.Copy(w, f)
		f.Close()
	}
}

//func (t *Root) CacheHandler(w http.ResponseWriter, r *http.Request) {
//
//	path := r.URL.Path
//	path = path[1:]
//	_, filename := filepath.Split(path)
//	if filepath.Ext(r.URL.Path) == ".js" || filepath.Ext(r.URL.Path) == ".css" {
//		if r.Header.Get("if-none-match") == filename {
//			w.WriteHeader(http.StatusNotModified)
//			return
//		}
//		w.Header().Set("etag", filename)
//		//log.Println(cachePath + "/" + filename)
//		f, err := t.fs.Open(cachePath + "/" + filename)
//		//log.Println("err", err)
//		if err != nil {
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//
//		w.WriteHeader(http.StatusOK)
//		io.Copy(w, f)
//		f.Close()
//	}
//}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

}
