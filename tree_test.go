package main

import (
  "testing"
)

func TestReadMetaFile(t *testing.T) {
  m, err := readMetaFile("testdata", "meta.json")
  if err != nil {
    t.Error(err)
  }

  if len(m.Pages) != 1 {
    t.Fatalf("Page count was not 1.")
  }

  p := m.Pages[0]

  if p.Path != "testdata/test" {
    t.Fatalf("Path was not 'testdata/test', got '%v'.", p.Path)
  }

  if p.Title != "Test Page" {
    t.Fatalf("Page title was not 'Test Page', got '%v'.", p.Title)
  }

  if p.Summary != "This is a test page" {
    t.Fatalf("Page summary was not 'This is a test page', got '%v'.", p.Summary)
  }
}
