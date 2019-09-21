//
// common.go
// Copyright (C) 2019 flo <flo@knightknight>
//
// Distributed under terms of the MIT license.
//

package main

type Proxy struct {
  Name string
  IP string
  Port string
}

func NewProxy(name string, ip string, port string) *Proxy {
  proxy := new(Proxy)
  proxy.Name = name
  proxy.IP = ip
  proxy.Port = port
  return proxy
}

type Database struct {
  Name string
  IP string
  Port string
}

func NewDatabase(name string, ip string, port string) *Database {
  database := new(Database)
  database.Name = name
  database.IP = ip
  database.Port = port
  return database
}

