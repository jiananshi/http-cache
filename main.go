package main

import (
  "net/http"
  "io/ioutil"
  "fmt"
  "log"
  "time"
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

    ifModified := r.Header.Get("if-modified-since")
    w.Header().Set("last-modified", lastModified)
    w.Header().Set("cache-control", "public, max-age=5")

    if ifModified == lastModified {
      lastModified = time.Now().Format(http.TimeFormat)
      w.WriteHeader(304)
    } else {
      fmt.Fprintf(w, string(data))
    }
  })

  log.Fatal(http.ListenAndServe(":8088", nil))
}

