package main

import (
  "crypto/md5"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"

  md "github.com/russross/blackfriday"
)

type Page struct {
  Path string
  Content []byte
  Hash string
}

type PageError struct {
  StatusCode int
  Message string
}

func (p PageError) Error() string {
  return fmt.Sprintf("HTTP %v: %v", p.StatusCode, p.Message)
}

func GetPage(absPath string) (*Page, error) {
	bytes, err := ioutil.ReadFile(absPath)

	if err != nil {

    if os.IsNotExist(err) {
      return nil, &PageError {http.StatusNotFound, "The page does not exist."}
    }

		return nil, err
  }

  page := &Page{}
  page.Path = absPath
  page.Hash = hash(bytes)
	page.Content = md.MarkdownCommon(bytes)

	return page, nil
}

func hash(value []byte) string {
  h := md5.New()
  h.Write(value)

  return fmt.Sprintf("%x", h.Sum(nil))
}
