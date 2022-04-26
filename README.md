# CouchDB Storage for [OAuth 2.0](https://github.com/go-oauth2/oauth2)

## Install

``` bash
$ go get github.com/DataChassis/go-oauth2-couchdb
```

## Usage

``` go
package main

import (
	couchdb "github.com/DataChassis/go-oauth2-couchdb"
	"github.com/go-oauth2/oauth2/manage"
)

func main() {
	manager := manage.NewDefaultManager()

	// use couchdb token store
	tokenConfig := couchdb.NewConfig("http://localhost:5984", "oauth2-tokens", "username", "password")
	manager.MapTokenStorage(couchdb.NewTokenStore(tokenConfig))

	// use couchdb client store
	clientConfig := couchdb.NewConfig("http://localhost:5984", "oauth2-clients", "username", "password")
	manager.MapClientStorage(couchdb.NewClientStore(clientConfig))

	// ...
}
```

## MIT License

```
Copyright (c) 2022 Data Chassis Ltd
Portions Copyright (c) 2016 Lyric
```
