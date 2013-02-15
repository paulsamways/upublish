package Layout

import (
  "bytes"
  "io/ioutil"
  "path/filepath"
  "strings"
)

type File struct {
  path string

  Bytes []byte
  Hash []byte
  Offset int

  Sub []*File
}

var Filename string
var Token []byte = []byte("{{CONTENT}}")

func (f *File) Matching(path string) *File {
  if ok, _ := strings.HasPrefix(path, f.path); ok {
    for _, v := range f.Sub {
      if v.Matching(path) {
        return v
      }
    }
    return f
  }
  return nil
}

func Graph(path string) (*File, error) {
  f, err := read(filepath.Join(path, Filename))
  if err != nil {
    return err
  }


}

func read(path string) (*Layout, error) {
  var b []byte
  var err error

  if b, err = ioutil.ReadFile(path); err != nil {
    return nil, err
  }

  bs := bytes.Split(b, Token)

  f := &File{path}
  f.Bytes = make([]byte, len(b) - len(Token))
  f.Hash = hash(b)
  f.Offset = bytes.Index(b, Token)

  return f, nil
}
