package main

import (
  "bytes"
  "crypto/md5"
  "path/filepath"
  "fmt"
  "io/ioutil"
  "os"
  "strings"

	md "github.com/russross/blackfriday"
)

var LayoutFilename = "layout.html"
var MetaFilename = "meta.json"

type Dir struct {
  Name string

  Layout *LayoutFile
  Meta *MetaFile
  Files map[string]*ContentFile

  Directories map[string]*Dir
}

type ContentFile struct {
  Name string

  Content []byte
  Hash []byte
}

type LayoutFile struct {
  Pre, Post []byte
  Hash []byte
}

type MetaFile struct {
}

func ReadTree(base string) (*Dir, []error) {
  errors := make([]error, 0)

  var parse func(current string, parentLayout *LayoutFile) *Dir

  parse = func(current string, parentLayout *LayoutFile) *Dir {
    dir := &Dir{}
    dir.Name = filepath.Base(current)
    dir.Layout = parentLayout
    dir.Files = make(map[string]*ContentFile, 0)

    fd, err := os.Open(current);
    if err != nil {
      errors = append(errors, fmt.Errorf("Failed to open directory '%v': %v", current, err))
      return nil
    }

    files, err := fd.Readdir(0)
    if err != nil {
      errors = append(errors, fmt.Errorf("Failed to enumerate directory '%v': %v", current, err))
    }

    subdirs := make([]string, 0)

    for _, file := range files {
      n := file.Name()

      if n[0] == '.' {
        continue
      }

      if file.IsDir() {
        subdirs = append(subdirs, filepath.Join(current, n))
        continue
      }

      var err error

      switch {
        case filepath.Ext(n) == ".md":
          var c *ContentFile
          if c, err = readContentFile(current, n); err != nil {
            errors = append(errors, fmt.Errorf("Failed to read content file '%v': %v",
              filepath.Join(current, n), err))
          }

          dir.Files[c.Name] = c
        case n == "layout.html":
          if dir.Layout, err = readLayoutFile(current, n, parentLayout); err != nil {
            errors = append(errors, fmt.Errorf("Failed to read layout file '%v': %v",
              filepath.Join(current, n), err))
          }
        case n == "meta.json":
          if dir.Meta, err = readMetaFile(current, n); err != nil {
            errors = append(errors, fmt.Errorf("Failed to read meta file '%v': %v",
              filepath.Join(current, n), err))
          }
      }
    }

    if len(subdirs) > 0 {
      dir.Directories = make(map[string]*Dir)

      for _, subdir := range subdirs {
        n := filepath.Base(subdir)
        dir.Directories[n] = parse(subdir, dir.Layout)
      }
    }

    return dir
  }

  return parse(base, nil), errors
}

func readContentFile(dir, name string) (*ContentFile, error) {
	b, err := ioutil.ReadFile(filepath.Join(dir, name))

  if err != nil {
    return nil, err
  }

	cf := &ContentFile{}
  cf.Name = name[:len(name)-3]
  cf.Content = md.MarkdownCommon(b)
  cf.Hash = hash(cf.Content)

  return cf, nil
}
func readLayoutFile(dir, name string, parent *LayoutFile) (*LayoutFile, error) {
  b, err := ioutil.ReadFile(filepath.Join(dir, name))

  if err != nil {
    return nil, err
  }

  spl := bytes.Split(b, []byte("{{content}}"))

  if len(spl) != 2 {
    return nil, fmt.Errorf("{{content}} token not found")
  }

  lf := &LayoutFile{}
  lf.Pre, lf.Post = spl[0], spl[1]
  lf.Hash = hash(b)

  if parent != nil {
    x := make([]byte, len(lf.Pre) + len(parent.Pre))
    copy(x, parent.Pre)
    copy(x[len(parent.Pre):], lf.Pre)

    lf.Pre = x
    lf.Post = append(lf.Post, parent.Post...)

    for i := 0; i < 16; i++ {
      lf.Hash[i] ^= parent.Hash[i]
    }
  }

  return lf, nil
}
func readMetaFile(dir, name string) (*MetaFile, err or) {
  return nil, nil
}

func hash(value []byte) []byte {
	h := md5.New()
	h.Write(value)
  return h.Sum(nil)
}

func (d *Dir) FindByPath(path string) *Dir {
  dir, _ := filepath.Split(path)
  dir = strings.Trim(dir, string(os.PathSeparator))

  if len(dir) == 0 {
    return d
  }

  dirSpl := strings.Split(dir, string(os.PathSeparator))
  match := d

  for _, n := range dirSpl {
    v, ok := match.Directories[n]
    if !ok {
      return nil
    }
    match = v
  }

  return match
}
