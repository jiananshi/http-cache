package main

import (
  "net/http"
  "io"
  "io/ioutil"
  "fmt"
  "log"
  "time"
  "crypto/md5"
  "encoding/hex"
)

func main() {
  lastModified := time.Now().Format(http.TimeFormat)

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    data, err := ioutil.ReadFile("index.html")
    if err != nil {
      log.Fatal(err)
    }
    w.Header().Set("cache-control", "no-cache")
    fmt.Fprintf(w, string(data))
  })

  http.HandleFunc("/test.js", func(w http.ResponseWriter, r *http.Request) {
    data, err := ioutil.ReadFile("test.js")

    if err != nil {
      log.Fatal(err)
    }

    h := md5.New()
    io.WriteString(h, string(data))
    etag := hex.EncodeToString(h.Sum(nil))

    ifModified := r.Header.Get("if-modified-since")
    ifNoneMatch := r.Header.Get("if-none-match");

    w.Header().Set("last-modified", lastModified)
    w.Header().Set("cache-control", "public, max-age=5")
    w.Header().Set("etag", etag)

    if ifNoneMatch == etag {
      w.WriteHeader(304)
    } else if ifModified == lastModified {
      lastModified = time.Now().Format(http.TimeFormat)
      w.WriteHeader(304)
    } else {
      fmt.Fprintf(w, string(data))
    }
  })

  log.Fatal(http.ListenAndServe(":8088", nil))
}

