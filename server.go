package main

import (
  "fmt"
  "log"
  "encoding/json"
  "net/http"
)

func helpHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome to the mock proxy-admin tool!\n")
  fmt.Fprintf(w, "Try one of these routes:\n")
  fmt.Fprintf(w, "→ /get-proxies: returns a list of proxies\n")
  fmt.Fprintf(w, "→ /get-db: returns a list of databases\n")
  fmt.Fprintf(w, "→ /get-proxy-conf: returns the proxy global configuration\n")
}

func getProxies() []*Proxy {
  return []*Proxy{NewProxy("proxy1", "127.0.0.1", "50051"),
                  NewProxy("proxy2", "127.0.0.1", "50052")}
}

func getDatabases() []*Database {
  return []*Database{NewDatabase("db1", "127.0.0.1", "8443"),
                  NewDatabase("db2", "127.0.0.1", "8444"),
                  NewDatabase("db3", "127.0.0.1", "5444")}
}

func getProxiesHandler(w http.ResponseWriter, r *http.Request) {
  js, _ := json.Marshal(getProxies())
  fmt.Fprintf(w, "%s\n", js)
}

func getDatabasesHandler(w http.ResponseWriter, r *http.Request) {
  js, _ := json.Marshal(getDatabases())
  fmt.Fprintf(w, "%s\n", js)
}

func main() {
  http.HandleFunc("/help", helpHandler)
  http.HandleFunc("/get-proxies", getProxiesHandler)
  http.HandleFunc("/get-databases", getDatabasesHandler)
  log.Fatal(http.ListenAndServe(":8081", nil))
}
