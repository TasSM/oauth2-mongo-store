# Mongo Storage for [OAuth 2.0](https://github.com/go-oauth2/oauth2)

This implementation of the go-oauth2 store uses the newer, officially maintained MongoDB driver for Go.
Deliberately, this implementation stores data in a de-normalised way which maybe better fits common usage patterns for MongoDB collections.

## Install

``` bash
$ go get -u -v gopkg.in/go-oauth2/mongo.v3
```

## Usage

``` go
package main

import (
	// stuff
)

func main() {
	manager := manage.NewDefaultManager()

	// use mongodb token store
	manager.MapTokenStorage(
		mongo.NewTokenStore(mongo.NewConfig(
			"mongodb://127.0.0.1:27017",
			"oauth2",
		)),
	)
	// ...
}
```

## MIT License

```
Copyright (c) 2022 Tasman Mayers
```