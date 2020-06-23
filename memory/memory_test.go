package memory

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestMemoryFilesystem(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	fileStr := "test"
	resourceKey := "sys/test"
	if err := store.Put(ctx, resourceKey, strings.NewReader(fileStr)); err != nil {
		t.Fatal(err)
	}

	rc, err := store.Get(ctx, resourceKey)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	newStr := buf.String()
	if newStr != fileStr {
		t.Fatalf("results do not match! expected %s but got %s", fileStr, newStr)
	}

	if err := store.Delete(ctx, resourceKey); err != nil {
		t.Fatal(err)
	}
}

func TestGetKeyNotFound(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	resourceKey := "notfound"
	if _, err := store.Get(ctx, resourceKey); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestDeleteKeyNotFound(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	resourceKey := "notfound"
	if err := store.Delete(ctx, resourceKey); err == nil {
		t.Fatal("expected error but got nil")
	}
}
