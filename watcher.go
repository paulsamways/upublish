package main

import (
  "log"
  "path/filepath"

  "github.com/howeyc/fsnotify"
)

var watcher *fsnotify.Watcher

func setupWatchers() error {
  var err error
  watcher, err = fsnotify.NewWatcher()

  if err != nil {
    log.Println("Watcher could not be initialised")
    return err
  }

  go func() {
    for {
      select {
      case ev := <-watcher.Event:
        invalidate(ev.Name)
      case err := <-watcher.Error:
        log.Println(err)
      }
    }
  }()

  err = watcher.Watch(root)
  if err != nil {
    log.Println("Could not watch path")
    return err
  }

  return nil
}

func invalidate(file string) {
  if file == filepath.Join(root, *optTemplate) {
    setupTemplate()
	  cache = make(map[string]*Page)
  } else {
    delete(cache, file)
  }
}
