package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	md "github.com/russross/blackfriday"
)

var optAddr = flag.String("addr", ":8000", "address to listen on")
var optPath = flag.String("path", ".", "path of the static files to serve")
var optPublic = flag.String("public", "public", "name of the 'public' directory")
var optTemplate = flag.String("tmpl", "template.html", "template to use")
var optHomeDir = flag.String("home", "", "home directory")
var optDefault = flag.String("default", "index", "default file to render")
var optNoCache = flag.Bool("noCache", true, "sends the Cache-Control headers to the client to prevent caching")

var root string
var cache map[string]string
var tmplBefore, tmplAfter string

func main() {
	flag.Parse()

	var err error
	root, err = filepath.Abs(*optPath)

	if err != nil {
		log.Fatalf("Could not get the absolute path of %v. %v", *optPath, err)
	}

	log.Printf("Rendering files in path %v.\n", root)

  cache = make(map[string]string)

	setupPublic()
	setupTemplate()

	http.HandleFunc("/", renderer)

	err = http.ListenAndServe(*optAddr, nil)

	if err != nil {
		log.Fatalf("Could not serve static files at path %v. %v", root, err)
	}
}

func setupPublic() {
	public := filepath.Join(root, *optPublic)

	log.Printf("Public folder is %v.\n", public)

	h := http.StripPrefix("/public/", http.FileServer(http.Dir(public)))

	if *optNoCache {
		h = CachePreventionHandler(h)
	}

	http.Handle("/public/", h)
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

  page, err := getPage(p, file)

  if err != nil {
		fmt.Fprint(w, tmplBefore+"Page doesn't exist!"+tmplAfter)
		return
	}

	r.Body.Close()
	fmt.Fprintf(w, tmplBefore+page+tmplAfter)
}

func getPage(p, file string) (string, error) {
  fp := path.Join(root, p, file)

  if v, ok := cache[fp]; ok {
    return v, nil
  }

	ext := filepath.Ext(fp)
	showMd := ext == ".md"

	if ext == "" {
		fp += ".md"
	}

  bytes, err := ioutil.ReadFile(fp)

  if err != nil {
    return "", err
  }

	if showMd {
		return string(bytes), nil
	}

  html := string(md.MarkdownCommon(bytes))

  cache[fp] = html
  return html, nil
}

func CachePreventionHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
		h.ServeHTTP(w, r)
	})
}
