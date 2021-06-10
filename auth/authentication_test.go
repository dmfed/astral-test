package auth

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

var testData = []byte(`
user1 pass1
user2 pass2
`)

func Test_newCK(t *testing.T) {
	fmt.Println("Testing CredetinalsKeeper")
	ck, _ := newCK(bytes.NewReader(testData))
	testAuthenticator(ck, t)
}

func Test_New(t *testing.T) {
	filename := "test.txt"
	if err := os.WriteFile(filename, testData, 0666); err != nil {
		fmt.Println("error creating test file")
		t.SkipNow()
	}
	defer os.Remove(filename)
	ck, _ := New(filename)
	testAuthenticator(ck, t)
}

func testAuthenticator(ck Authenticator, t *testing.T) {
	if !ck.CredentialsAreValid("user1", "pass1") || !ck.CredentialsAreValid("user2", "pass2") {
		fmt.Println("valid passwords do not check out")
		t.Fail()
	}
	if ck.CredentialsAreValid("user3", "some") {
		fmt.Println("invalid password checks out")
		t.Fail()
	}
}
