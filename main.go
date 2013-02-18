package main

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var optAddr = flag.String("addr", ":8000", "address to listen on")
var optPath = flag.String("path", ".", "path of the static files to serve")
var optStaticDir = flag.String("public", "public", "name of the 'public' directory")

var root string
var tree *Dir

func main() {
	flag.Parse()

	var err error

	if root, err = filepath.Abs(*optPath);  err != nil {
		log.Fatalf("Could not get the absolute path of %v. %v", *optPath, err)
	}

	setupStaticDir()
	setupSignals()

  var errs []error
  if tree, errs = ReadTree(root); len(errs) > 0 {
    for _, err = range errs {
      log.Printf("%v\n", err)
    }

    log.Fatalf("Found %v errors", len(errs))
  }

	http.HandleFunc("/", renderPage)

	if err = http.ListenAndServe(*optAddr, nil); err != nil {
		log.Fatalf("Could not serve static files at path %v. %v", root, err)
	}
}

func setupStaticDir() {
	static := "/" + *optStaticDir + "/"
	public := filepath.Join(root, *optStaticDir)

	h := http.StripPrefix(static, http.FileServer(http.Dir(public)))

	http.Handle(static, h)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(public, "favicon.ico"))
	})
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(public, "robots.txt"))
	})
}

func setupSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)

	go func() {
		<-c
		//tree = parseTree()
	}()
}

func renderPage(w http.ResponseWriter, r *http.Request) {
  
	write(w, r, page.Content, page.Hash)
}

func renderIndex(w http.ResponseWriter, r *http.Request, pIdx []PageIndex) {
  b := bytes.Buffer{}

  fmt.Fprintln(b, "<h2>Blog</h2>")
  for i := 0; i < len(pIdx); i++ {
    fmt.Fprintf(b, "<h3>%v</h3>\n<p>%v</p>\n", pIdx[i].Name, pIdx[i].Summary)
  }

  b = baseLayout.Render(b)

}

func write(w http.ResponseWriter, r *http.Request, value []byte, hash []byte) {
  if len(hash) > 0 {
    strHash := fmt.Sprintf("%x", hash)

    if etag := r.Header.Get("If-None-Match"); strings.EqualFold(etag, strHash) {
      w.WriteHeader(http.StatusNotModified)
      return
    }

    w.Header().Set("Etag", strHash)
  }

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		b := new(bytes.Buffer)
		gz := gzip.NewWriter(b)
		gz.Write(tmpl[0])
		gz.Write(value)
		gz.Write(tmpl[1])
		gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Length", strconv.Itoa(b.Len()))

		b.WriteTo(w)
	} else {
		w.Header().Set("Content-Length", strconv.Itoa(len(tmpl[0])+len(value)+len(tmpl[1])))
		w.Write(tmpl[0])
		w.Write(value)
		w.Write(tmpl[1])
	}
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	errFmt := "<h2>Oops! We've hit a bit of a problem...</h2><p>%v</p>"

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	b := new(bytes.Buffer)
	b.Write(tmpl[0])

	if pErr, ok := err.(*PageError); ok {
		w.WriteHeader(pErr.StatusCode)
		fmt.Fprintf(b, errFmt, pErr.Message)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(b, errFmt, "Page not available")
	}

	b.Write(tmpl[1])
	b.WriteTo(w)
}


