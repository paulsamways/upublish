package main

import (
	"testing"
	"time"
)

func TestReadMetaFile(t *testing.T) {
	m, err := readMetaFile("testdata", "meta.json")
	if err != nil {
		t.Error(err)
	}

	if len(m.Pages) != 1 {
		t.Fatalf("Page count was not 1.")
	}

	p := *m.Pages[0]
	pExp := Page{
		"testdata/test",
		"Test Page",
		"This is a test page",
		[]string{"One", "Two", "Three"},
		time.Date(2013, 02, 22, 14, 22, 0, 0, time.UTC),
	}

	const x = "%v was not decoded correctly; expected %v but got %v."

	if p.Path != pExp.Path {
		t.Fatalf(x, "Path", pExp.Path, p.Path)
	}
	if p.Title != pExp.Title {
		t.Fatalf(x, "Title", pExp.Title, p.Title)
	}
	if p.Summary != pExp.Summary {
		t.Fatalf(x, "Summary", pExp.Summary, p.Summary)
	}
	if p.Date != pExp.Date {
		t.Fatalf(x, "Date", pExp.Date, p.Date)
	}
  if len(p.Tags) != len(pExp.Tags) {
    t.Fatalf(x, "Tags", pExp.Tags, p.Tags)
  }
}
