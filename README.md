strongswan vici golang client
=============================
[![Build Status](https://travis-ci.org/bronze1man/goStrongswanVici.svg)](https://travis-ci.org/bronze1man/goStrongswanVici)
[![GoDoc](https://godoc.org/github.com/bronze1man/goStrongswanVici?status.svg)](https://godoc.org/github.com/bronze1man/goStrongswanVici)
[![docs examples](https://sourcegraph.com/api/repos/github.com/bronze1man/goStrongswanVici/badges/docs-examples.png)](https://sourcegraph.com/github.com/bronze1man/goStrongswanVici)
[![Total views](https://sourcegraph.com/api/repos/github.com/bronze1man/goStrongswanVici/counters/views.png)](https://sourcegraph.com/github.com/bronze1man/goStrongswanVici)

a golang implement of strongswan vici plugin client.

### document
http://godoc.org/github.com/bronze1man/goStrongswanVici
https://github.com/strongswan/strongswan/tree/master/src/libcharon/plugins/vici

### Implemented command list
* version()
* list-sas()
* terminate()
If you need some commands, but it is not here .you can implement yourself, and send a pull request to this project.

### example
```go
package main

import (
	"fmt"
	"github.com/bronze1man/goStrongswanVici"
)

func main(){
    // create a client.
	client, err := goStrongswanVici.NewClientConnFromDefaultSocket()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// get strongswan version
	v, err := client.Version()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", v)

	// get all conns info from strongswan
	connInfo, err := client.ListAllVpnConnInfo()
	if err != nil {
		panic(err)
	}
	fmt.Printf("found %d connections. \n", len(connInfo))

	// kill all conns in strongswan.
	for _, info := range connInfo {
		fmt.Printf("kill connection id %s\n", info.Uniqueid)
		err = client.Terminate(&goStrongswanVici.TerminateRequest{
			Ike_id: info.Uniqueid,
		})
		if err != nil {
			panic(err)
		}
	}
}
```
