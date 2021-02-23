package berkeleydb

import (
	"os"
	"testing"
)

const TestDirectory = "./TEST_ENV"

func TestNewEnvironment(t *testing.T) {
	_, err := os.Stat(TestDirectory)
	if err != nil && os.IsNotExist(err) {
		e := os.Mkdir(TestDirectory, os.ModeDir|os.ModePerm)
		if e != nil {
			t.Fatalf("Failed to create directory: %s", e)
		}
	}

	_, err = NewEnvironment()

	if err != nil {
		t.Errorf("Expected environment, got %s", err)
	}

}

func TestOpenEnvironment(t *testing.T) {
	env, _ := NewEnvironment()
	err := env.Open(TestDirectory, DbCreate|DbInitMpool, 0)
	if err != nil {
		t.Errorf("Expected to open DB, got %s", err)
	}

	err = env.Close()
	if err != nil {
		t.Errorf("Expected to close DB, got %s", err)
	}
}

func TestOpenDBInEnvironment(t *testing.T) {
	env, _ := NewEnvironment()
	err := env.Open(TestDirectory, DbCreate|DbInitMpool, 0755)
	if err != nil {
		t.Error("Expected to open DB, got ", err)
		return
	}

	// Now create, open, and close a DB
	db, err := NewDBInEnvironment(env)
	if err != nil {
		t.Error("Expected to create new DB: ", err)
	}

	err = db.Open(TestFilename, "", DbBtree, DbCreate)
	if err != nil {
		t.Error("Expected to open DB, got ", err)
	}

	// Test that the DB file was actually created.
	_, err = os.Stat(TestDirectory + "/" + TestFilename)
	if err != nil {
		t.Error("Expected to stat .db, got ", err)
	}

	err = db.Close()
	if err != nil {
		t.Error("Expected to close the DB, got ", err)
	}

	err = env.Close()
	if err != nil {
		t.Errorf("Expected to close DB, got %s", err)
	}
}
func TestTeardown(t *testing.T) {
	err := os.RemoveAll(TestDirectory)
	if err != nil {
		t.Fatalf("Expected to remove fixtures, got %s", err)
	}
}
