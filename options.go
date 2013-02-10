package main

import (
  "flag"
)

var optAddr = flag.String("addr", ":8000", "address to listen on")
var optPath = flag.String("path", ".", "path of the static files to serve")
var optStaticDir = flag.String("public", "public", "name of the 'public' directory")
var optTemplate = flag.String("tmpl", "template.html", "template to use")
var optHomeDir = flag.String("home", "", "home directory")
var optDefault = flag.String("default", "index", "default file to render")
var optDebug = flag.Bool("debug", false, "enables debug mode")
var optExt = flag.String("ext", "md", "extension of the markdown files")

func parseOptions() {
	flag.Parse()
}
