[![GoDoc](https://godoc.org/github.com/bfix/berkeleydb?status.svg)](https://godoc.org/github.com/bfix/berkeleydb)
[![Travis](https://api.travis-ci.org/bfix/berkeleydb.svg?branch=master)](https://travis-ci.org/bfix/berkeleydb)


### BerkeleyDB Bindings

This package provides BerkeleyDB wrappers for the C library using `cgo`.

To build, you will need a relatively recent version of BerkeleyDB.



### Example
```go
package main

import (
        "fmt"

        "github.com/jsimonetti/berkeleydb"
)

func main() {
        var err error

        db, err := berkeleydb.NewDB()
        if err != nil {
                fmt.Printf("Unexpected failure of CreateDB %s\n", err)
        }

        err = db.Open("./test.db", "", berkeleydb.DbHash, berkeleydb.DbCreate)
        if err != nil {
                fmt.Printf("Could not open test_db.db. Error code %s", err)
                return
        }
        defer db.Close()

        err = db.Put([]byte("key"), byte[]("value"))
        if err != nil {
                fmt.Printf("Expected clean Put: %s\n", err)
        }

        value, err := db.Get([]byte("key"))
        if err != nil {
                fmt.Printf("Unexpected error in Get: %s\n", err)
                return
        }
        fmt.Printf("value: %s\n", string(value))
}
```
