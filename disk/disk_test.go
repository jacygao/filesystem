package disk

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func genKey() (string, error) {
	b := make([]byte, 4)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func TestFilesystemService(t *testing.T) {
	ctx := context.Background()
	// Generate a key.
	keygen, err := genKey()
	if err != nil {
		t.Fatal(err)
	}

	// Create a multi-level deep resource key.
	key := fmt.Sprintf("1/2/%s", keygen)
	path := "testdata"
	ds, err := NewDiskStore(path)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	if err := ds.Put(ctx, key, strings.NewReader("bar")); err != nil {
		t.Fatal(err)
	}
	reader, err := ds.Get(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "bar" {
		t.Fatalf("results do not match! expected %s, but got %s", "bar", string(data))
	}

	fileInfo, err := ds.List(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(fileInfo) != 1 {
		t.Fatal("expected one directory to be returned")
	}

	fileInfo, err = ds.List(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo[0].Name() != "2" {
		t.Fatal("expected to see directory 2")
	}

	fileInfo, err = ds.List(ctx, "1/2")
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo[0].Name() != keygen {
		t.Fatalf("expected to see file keygen %s", keygen)
	}

	if err := ds.Delete(ctx, key); err != nil {
		t.Fatal(err)
	}
}
