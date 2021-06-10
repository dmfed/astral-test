package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var testElements = []byte(`
[
	{
		"payload": "first element"
	},
	{
		"payload": "second element"
	},
	{
		"payload": "third element"
	},
	{
		"payload": "fourth element"
	}
]`)

func Test_StorageSQLite(t *testing.T) {
	fmt.Println("Testing SQLite Storage methods")
	testDBfile := "./test.db"
	if _, err := os.Stat(testDBfile); err == nil {
		os.Remove(testDBfile)
	}
	st, err := OpenSQLiteStorage(testDBfile)
	if err != nil {
		fmt.Printf("error creating %s: %v", testDBfile, err)
		t.Fail()
	}
	defer os.Remove(testDBfile) // cleanup
	defer st.Close()            // close Sqlite connection
	testStorageMethods(t, st)
}

func getTestElements() ([]Element, error) {
	var e []Element
	err := json.Unmarshal(testElements, &e)
	return e, err
}

func testStorageMethods(t *testing.T, st Storage) {
	defer st.Close()
	elements, err := getTestElements()
	if err != nil {
		fmt.Println("test data fails to unmarshal. stopping.")
		t.FailNow()
	}

	// Testing Put()
	var ids []ID
	for _, e := range elements {
		id, err := st.Put(e)
		if err != nil || id == BadID {
			fmt.Println("storage fails to put:", err)
			t.FailNow()
		} else {
			ids = append(ids, id)
		}
	}

	// Testing Get()
	received, err := st.Get()
	if err != nil {
		fmt.Println("storage fails to get all:", err)
		t.FailNow()
	}
	for i, e := range received {
		if e.Payload != elements[i].Payload || e.ID != ids[i] {
			fmt.Println("recevied does not match original")
			t.FailNow()
		}
	}

	// Testing Get() by ID
	for i, id := range ids {
		e, err := st.Get(id)
		if err != nil {
			fmt.Println("error getting existing ID:", err)
			t.Fail()
		}
		if e[0].Payload != elements[i].Payload || e[0].ID != ids[i] {
			fmt.Println("Storage returns incorrect element by ID")
			t.Fail()
		}
	}

	// Testing Get() with incorrect ID
	if e, err := st.Get(BadID); err == nil {
		fmt.Println("Get() called with incorrect ID and returned no error")
		fmt.Println(e)
		t.Fail()
	}

	// Testing Upd()
	if err := st.Upd(received[0].ID, received[1]); err != nil {
		fmt.Println("failed to update existing ID:", err)
		t.Fail()
	}
	updated, err := st.Get(received[0].ID)
	if err != nil {
		fmt.Println("failed to get update element by id", err)
		t.Fail()
	}
	if updated[0].Payload != received[1].Payload {
		fmt.Println("updated element differs from original")
		fmt.Println("updated:", updated[0])
		fmt.Println("expected:", received[1])
		t.Fail()
	}
	toDelete, _ := st.Get()

	// Testing Upd() with incorrect ID
	if err := st.Upd(BadID, received[0]); err == nil {
		fmt.Println("Upd() with incorrect ID returns no error")
		t.Fail()
	}

	// Testing Del() with incorrect ID
	if err := st.Del(BadID); err == nil {
		fmt.Println("Del() with incorrect ID returns no error:", err)
		t.Fail()
	}

	// Testing Del()
	for _, e := range toDelete {
		if err := st.Del(e.ID); err != nil {
			fmt.Println("error deleting element by ID:", err)
			t.Fail()
		}
	}
	if remaining, _ := st.Get(); len(remaining) > 0 {
		fmt.Println("elements remaining in Storage after delete:", err)
		t.Fail()
	}
}
