package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	md "github.com/russross/blackfriday"
)

type Page struct {
	Path    string
	Content []byte
	Hash    string
}

type PageError struct {
	StatusCode int
	Message    string
}

func (p PageError) Error() string {
	return fmt.Sprintf("HTTP %v: %v", p.StatusCode, p.Message)
}

func GetPage(absPath string) (*Page, error) {
	bytes, err := ioutil.ReadFile(absPath)

	if err != nil {

		if os.IsNotExist(err) {
			return nil, &PageError{http.StatusNotFound, "The page does not exist."}
		}

		return nil, err
	}

	page := &Page{}
	page.Path = absPath
	page.Content = md.MarkdownCommon(bytes)

	hTmp := make([]byte, 16)
	h := hash(bytes)

	for i := 0; i < 16; i++ {
		hTmp[i] = h[i] ^ tmplHash[i]
	}

	page.Hash = fmt.Sprintf("%x", hTmp)

	return page, nil
}
