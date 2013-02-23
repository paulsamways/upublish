package main

import (
  "net/http"
  "fmt"
  "sort"
  "bytes"
)

func renderIndex(d *Dir, w http.ResponseWriter, r *http.Request) {
  p := make(Pages, len(d.Meta.Pages))
  copy(p, d.Meta.Pages)
  sort.Sort(ByDate{p})

  b := &bytes.Buffer{}
  for _, v := range p {
    fmt.Fprintf(b, "<h3><a href=\"%v\">%v</a></h3><p>%v</p><p>%v</p>", v.Path, v.Title, v.Date, v.Summary)
  }

  write(w, r, 200, b.Bytes(), nil, d.Layout)
}

// Sorters

type Pages []*Page

func (p Pages) Len() int { return len(p) }
func (p Pages) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type ByDate struct {
  Pages
}
func (p ByDate) Less(i, j int) bool {
  return p.Pages[i].Date.Unix() > p.Pages[j].Date.Unix()
}
