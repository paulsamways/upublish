package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var optAddr = flag.String("addr", ":8000", "address to listen on")
var optPath = flag.String("path", ".", "path of the static files to serve")
var optStaticDir = flag.String("public", ".public", "path of the 'public' directory")

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

  var ok bool
  if tree, ok = readTree(); !ok {
    log.Fatalf("Exiting...")
  }

	http.HandleFunc("/", renderPage)

	if err = http.ListenAndServe(*optAddr, nil); err != nil {
		log.Fatalf("Could not serve static files at path %v. %v", root, err)
	}
}

func setupStaticDir() {
	public := filepath.Join(root, *optStaticDir)

	h := http.StripPrefix("/public/", http.FileServer(http.Dir(public)))

	http.Handle("/public/", h)
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

    log.Print("\nReloading...")
    var d *Dir
    var ok bool

    if d, ok = readTree(); !ok {
      log.Println("Reload unsuccessful")
    } else {
      //lock
      tree = d
    }
	}()
}

func readTree() (*Dir, bool) {
  var dir *Dir
  var errs []error
  if dir, errs = ReadTree(root); len(errs) > 0 {
    for _, err := range errs {
      log.Printf("%v\n", err)
    }

    log.Printf("Found %v errors\n", len(errs))
    return nil, false
  }

  return dir, true
}

func renderPage(w http.ResponseWriter, r *http.Request) {
  p, file := filepath.Split(r.URL.Path)

  if file == "" {
    file = "index"
  }

  dir := tree.FindByPath(p)

  if dir == nil {
    write(w, r, []byte("Path Not found"), nil, tree.Layout)
    return
  }

  cf, ok := dir.Files[file]
  if !ok {
    write(w, r, []byte("Page not found"), nil, dir.Layout)
    return
  }

  write(w, r, cf.Content, cf.Hash, dir.Layout)
}

type nCloseWriter struct {
  io.Writer
}
func (w nCloseWriter) Close() error {
  return nil
}
func NoOpCloseWriter(w io.Writer) io.WriteCloser {
  return nCloseWriter{w}
}

func write(w http.ResponseWriter, r *http.Request, value []byte, hash []byte, layout *LayoutFile) {
  if len(hash) > 0 {
    h := make([]byte, 16)
    copy(h, hash)

    if layout != nil {
      for i := 0; i < 16; i++ {
        h[i] ^= layout.Hash[i]
      }
    }

    strHash := fmt.Sprintf("%x", h)

    if etag := r.Header.Get("If-None-Match"); strings.EqualFold(etag, strHash) {
      w.WriteHeader(http.StatusNotModified)
      return
    }

    w.Header().Set("Etag", strHash)
  }

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

  b := &bytes.Buffer{}
  var writer io.WriteCloser = NoOpCloseWriter(b)
  if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
    writer = gzip.NewWriter(b)
  }

  if layout != nil {
    writer.Write(layout.Pre)
    writer.Write(value)
    writer.Write(layout.Post)
  } else {
    writer.Write(value)
  }

  writer.Close()

	w.Header().Set("Content-Length", strconv.Itoa(b.Len()))
  b.WriteTo(w)
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Fprintf(w, "<h2>Oops! We've hit a bit of a problem...</h2><p>%v</p>", err)
}
