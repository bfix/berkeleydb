package berkeleydb

// #cgo LDFLAGS: -ldb
// #include <db.h>
// #include <stdlib.h>
// #include "bdb.h"
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

const version string = "0.0.1"

// Flags for opening a database or environment.
const (
	DbCreate   = C.DB_CREATE
	DbExcl     = C.DB_EXCL
	DbRdOnly   = C.DB_RDONLY
	DbTruncate = C.DB_TRUNCATE

	// DbInitMpool is used in environment only
	DbInitMpool = C.DB_INIT_MPOOL
)

// Database types.
const (
	DbBtree   = C.DB_BTREE
	DbHash    = C.DB_HASH
	DbRecno   = C.DB_RECNO
	DbQueue   = C.DB_QUEUE
	DbUnknown = C.DB_UNKNOWN
)

// Cursor modes.
const (
	CrsFirst = C.DB_FIRST
	CrsNext  = C.DB_NEXT
	CrsPrev  = C.DB_PREV
	CrsLast  = C.DB_LAST
)

// Db is the structure that holds the database connection
type Db struct {
	db *C.DB
}

// Cursor holds the current cursor position
type Cursor struct {
	dbc *C.DBC
}

// Errno is the error number
type Errno int

// NewDB initialises a new bdb connection
func NewDB() (*Db, error) {
	var db *C.DB
	err := C.db_create(&db, nil, 0)

	if err > 0 {
		return nil, createError(err)
	}

	return &Db{db}, nil
}

// NewDBInEnvironment initialises a new bdb connection in an environment
func NewDBInEnvironment(env *Environment) (*Db, error) {
	var db *C.DB
	err := C.db_create(&db, env.environ, 0)

	if err > 0 {
		return nil, createError(err)
	}

	return &Db{db}, nil
}

// OpenWithTxn opens the database in transaction mode (transactions are not yet supported by all
// funtions)
func (handle *Db) OpenWithTxn(filename string, txn *C.DB_TXN, dbtype C.DBTYPE, flags C.u_int32_t) error {
	db := handle.db
	file := C.CString(filename)
	defer C.free(unsafe.Pointer(file))

	ret := C.go_db_open(db, txn, file, nil, dbtype, flags, 0)

	return createError(ret)
}

// Open a database file
func (handle *Db) Open(filename, dbname string, dbtype C.DBTYPE, flags C.u_int32_t) error {
	file := C.CString(filename)
	defer C.free(unsafe.Pointer(file))
	var dbn *C.char = nil
	if len(dbname) > 0 {
		dbn = C.CString(dbname)
		defer C.free(unsafe.Pointer(dbn))
	}

	ret := C.go_db_open(handle.db, nil, file, dbn, dbtype, flags, 0)

	return createError(ret)
}

// Close the database file
func (handle *Db) Close() error {
	ret := C.go_db_close(handle.db, 0)

	return createError(ret)
}

// Flags returns the flags of the database connection
func (handle *Db) Flags() (C.u_int32_t, error) {
	var flags C.u_int32_t

	ret := C.go_db_get_open_flags(handle.db, &flags)

	return flags, createError(ret)
}

// Remove the database
func (handle *Db) Remove(filename string) error {
	file := C.CString(filename)
	defer C.free(unsafe.Pointer(file))

	ret := C.go_db_remove(handle.db, file)

	return createError(ret)
}

// Rename the database filename
func (handle *Db) Rename(oldname, newname string) error {
	oname := C.CString(oldname)
	defer C.free(unsafe.Pointer(oname))
	nname := C.CString(newname)
	defer C.free(unsafe.Pointer(nname))

	ret := C.go_db_rename(handle.db, oname, nname)

	return createError(ret)
}

func prepDbt(in []byte) *C.DBT {
	dbt := new(C.DBT)
	dbt.size = C.uint(len(in))
	dbt.data = unsafe.Pointer(&in[0])
	return dbt
}

// Put a key/value pair into the database
func (handle *Db) Put(key, value []byte) error {
	dbKey := prepDbt(key)
	dbVal := prepDbt(value)
	ret := C.go_db_put(handle.db, dbKey, dbVal, C.uint(0))
	if ret > 0 {
		return createError(ret)
	}
	return nil
}

// Get a value from the database by key
func (handle *Db) Get(key []byte) ([]byte, error) {
	dbKey := prepDbt(key)
	dbVal := new(C.DBT)
	ret := C.go_db_get(handle.db, dbKey, dbVal)
	data := C.GoBytes(unsafe.Pointer(dbVal.data), C.int(dbVal.size))
	return data, createError(ret)
}

// Delete a value from the database by key
func (handle *Db) Delete(key []byte) error {
	dbKey := prepDbt(key)
	ret := C.go_db_del(handle.db, dbKey)
	return createError(ret)
}

//Cursor returns a handle for the database cursor
func (handle *Db) Cursor() (*Cursor, error) {
	var dbc *C.DBC

	err := C.go_db_cursor(handle.db, &dbc)

	if err > 0 {
		return nil, createError(err)
	}

	return &Cursor{dbc}, nil
}

// Get [FIRST|NEXT|PREV|LAST] record from cursor and return the key/value pair
func (cursor *Cursor) Get(mode int) ([]byte, []byte, error) {
	key := new(C.DBT)
	value := new(C.DBT)
	ret := C.go_cursor_get(cursor.dbc, key, value, C.int(mode))
	dKey := C.GoBytes(unsafe.Pointer(key.data), C.int(key.size))
	dVal := C.GoBytes(unsafe.Pointer(value.data), C.int(value.size))
	return dKey, dVal, createError(ret)
}

// UTILITY FUNCTIONS

// Version returns the version of the database and binding
func Version() string {
	libVersion := C.GoString(C.db_full_version(nil, nil, nil, nil, nil))

	tpl := "%s (Go bindings v%s)"
	return fmt.Sprintf(tpl, libVersion, version)
}

// DBError contains the database Error
type DBError struct {
	Code    int
	Message string
}

func createError(code C.int) error {
	if code == 0 {
		return nil
	}
	msg := C.db_strerror(code)
	e := DBError{int(code), C.GoString(msg)}
	return errors.New(e.Error())
}

// Error return the string representation of the error
func (e *DBError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
