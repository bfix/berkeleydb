package berkeleydb

import (
	"testing"
)

const TestFilename = "test_db.db"

func TestNewDB(t *testing.T) {

	_, err := NewDB()

	if err != nil {
		t.Errorf("Expected error code 0, got %d", err)
	}
}
func TestOpen(t *testing.T) {
	db, err := NewDB()

	if err != nil {
		t.Errorf("Unexpected failure of CreateDB")
	}

	err = db.Open(TestFilename, "", DbBtree, DbCreate)

	if err != nil {
		t.Errorf("Could not open test_db.db. Error code %s", err)
	}

	flags, err := db.Flags()
	if err != nil {
		t.Errorf("Could not get Flags: %s", err)
	}
	if flags != DbCreate {
		t.Errorf("Expected flag to match DB_CREATE, got %d", flags)
	}

	err = db.Close()
	if err != nil {
		t.Errorf("Could not close file %s: %s", TestFilename, err)
		return
	}

}

func openDB() (*Db, error) {
	db, err := NewDB()

	if err != nil {
		return nil, err
	}

	err = db.Open(TestFilename, "", DbBtree, DbCreate)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func closeDB(db *Db) error {
	return db.Close()
}

func TestPut(t *testing.T) {
	db, err := openDB()
	defer closeDB(db)

	err = db.Put("key", "value")
	if err != nil {
		t.Error("Expected clean Put.", err)
	}
}

func TestGet(t *testing.T) {
	db, err := openDB()
	defer closeDB(db)

	err = db.Put("key", "value")
	if err != nil {
		t.Error("Expected clean Put: ", err)
	}

	val, err := db.Get("key")
	if err != nil {
		t.Error("Unexpected error in Get: ", err)
		return
	}

	if val != "value" {
		t.Error("Expected 'value', got ", val)
	}
}

func TestDelete(t *testing.T) {
	db, err := openDB()
	defer closeDB(db)

	err = db.Put("key", "value")
	if err != nil {
		t.Error("Expected clean Put: ", err)
	}

	err = db.Delete("key")
	if err != nil {
		t.Error("Expected a clean delete, got ", err)
	}

	err = db.Delete("nosuchkey")
	if err == nil {
		t.Error("Expected error, got ", err)
	}
}

func TestRemove(t *testing.T) {
	db, _ := NewDB()

	err := db.Remove(TestFilename)
	if err != nil {
		t.Errorf("Could not delete %s. Expected 0, got %s", TestFilename, err)
	}
}

func TestRename(t *testing.T) {
	db, _ := NewDB()
	db.Open(TestFilename, "", DbHash, DbCreate)
	db.Close()

	db, _ = NewDB()

	newname := "foo_" + TestFilename
	err := db.Rename(TestFilename, newname)
	if err != nil {
		t.Errorf("Could not rename %s to %s", TestFilename, newname)
	}

	db, _ = NewDB()
	err = db.Remove(newname)
	if err != nil {
		t.Errorf("Could not remove %s", newname)
	}
}
