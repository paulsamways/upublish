package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

var root string
var cache map[string]*Page
var tmplBefore, tmplAfter string

func main() {
  parseOptions()

	var err error
	root, err = filepath.Abs(*optPath)

	if err != nil {
		log.Fatalf("Could not get the absolute path of %v. %v", *optPath, err)
	}

	cache = make(map[string]*Page)

	setupStatic()
	setupTemplate()

	http.HandleFunc("/", renderer)

	err = http.ListenAndServe(*optAddr, nil)

	if err != nil {
		log.Fatalf("Could not serve static files at path %v. %v", root, err)
	}
}

func setupStatic() {
  static := "/"+*optStaticDir+"/"
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


func setupTemplate() {
	path := filepath.Join(root, *optTemplate)
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal("Could not read template: ", err)
	}

	parts := strings.Split(string(bytes), "{{content}}")

	if len(parts) != 2 {
		log.Fatal("Template was not in a valid format")
	}

	tmplBefore = parts[0]
	tmplAfter = parts[1]
}

func renderer(w http.ResponseWriter, r *http.Request) {
	p, file := path.Split(r.URL.Path)

	if p == "" {
		p = *optHomeDir
	}

	if file == "" {
		file = *optDefault
	}

  abs := path.Join(root, p, file) + "." + *optExt

  var page *Page
  var ok bool
  var err error

  if page, ok = cache[abs]; !ok {
    page, err = GetPage(abs)

    if err != nil {
      errFmt := "<h2>Oops! We've hit a bit of a problem...</h2><p>%v</p>"

      if pErr, ok := err.(*PageError); ok {
        w.WriteHeader(pErr.StatusCode)
        write(w, fmt.Sprintf(errFmt, pErr.Message))
      } else {
        w.WriteHeader(http.StatusInternalServerError)
        write(w, fmt.Sprintf(errFmt, "Page not available"))
      }

      log.Printf("[%v] %v", abs, err)

      return
    }

    if !*optDebug {
      cache[abs] = page
    }
  }

  if etag := r.Header.Get("If-None-Match"); strings.EqualFold(etag, page.Hash) {
    w.WriteHeader(http.StatusNotModified)
    return
  }

  w.Header().Add("Etag", page.Hash)
  write(w, page.Content)
}

func write(w http.ResponseWriter, value string) {
  fmt.Fprintf(w, tmplBefore+value+tmplAfter)
}
