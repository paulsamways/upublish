## &micro;Publish

### About

&micro;Publish is a web publishing engine written in [Go](http://golang.org).

Features include:

- File-based storage, making it very easy to get started. No database required.
- Folder-inherited layouts, so you can control the look and feel of your site.
- Content pages are written in Markdown.

### Usage

By default, &micro;Publish expects to be executed in the directory that you wish 
to serve.

``` bash
$ upublish
```


#### Command Line Options

Name      | Default       | Description
----------|---------------|---------------------------------------------------
-addr     | :8080         | The address to listen for incoming connections on
-path     | ./            | Path to serve files from
-public   | .public       | Directory where the static files are located


#### Layouts
The only file that is required by &micro;Publish is a single layout file, 
layout.html, which must contain the token "{{content}}".

Example:

``` HTML
<!DOCTYPE html>
<html>
  <head>
    <title>Example</title>
  </head>
  <body>
    <h1>Example &micro;Publish Template</h1>
    <div class="content">
      {{content}}
    </div>
  </body>
</html>
```

A layout file will be used for any page rendered in the current directory,
or any sub-directory recursively. When a layout file is created in a
sub-directory, the layout will be rendered within the section defined by
the parent layout.

Example:


/layout.html

``` HTML
<div class="main">
  {{content}}
</div>
```

/articles/layout.html

``` HTML
<h1>Articles</h1>
<div class="article">
  {{content}}
</div>
```

/articles/abc.md

``` MarkDown
  ABC
```

GET: /articles/abc

``` HTML
<div class="main">
  <h1>Articles</h1>
  <div class="article">
    ABC
  </div>
</div>
```

#### Content Files

Content in &micro;Publish is written using Markdown and stored as .md files.
The location of the files relative to the root determines the URL that will 
be used to access the pages.

``` Bash
$ upublish -path="/srv/http/mysite"
```

/srv/http/mysite/projects/xyz/about.md -> /projects/xyz/about


When a request is made which does not specify a file, &micro;Publish will 
attempt to serve an index.md file.

#### Reloading Pages

&micro;Publish caches all content pages and layouts when the server starts,
due to this you will need to either restart the process when updating content,
or send the &micro;Publish process a USR1 signal. This will force all files 
to be reloaded and the cache updated accordingly.

``` Bash
  # refresh *all* upublish processes
  killall -USR1 upublish

  # refresh the upublish process with PID 1234
  kill -USR1 1234

  # refresh upublish whenver something changes; requires inotify-tools
  while inotifywait -r -e modify -e create -e delete .; do 
    killall -USR1 upublish; 
  done
```

#### Hosting

&micro;Publish has been written to run as a standalone process. The easiest
way to start the process on a Linux machine is to use Systemd. See my article
on running [Go web servers using Systemd](/articles/go-systemd) for further information.

### Getting &micro;Publish

The source can be found at https://github.com/PaulSamways/upublish.

Alternatively, &micro;Publish can be fetched via 'go get':

``` Bash
$ go get github.com/PaulSamways/upublish
```
