# Overview

A Filesystem implementation backed by SFTP.

# Example

```
package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	config := &ssh.ClientConfig{
		User: "foo",
		Auth: []ssh.AuthMethod{
			ssh.Password("pass"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	key := "test"

	sshCli, err := ssh.Dial("tcp", "127.0.0.1:2222", config)
	if err != nil {
		log.Fatal(err)
	}
	cli, err := sftp.NewClient(sshCli)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	sftp := NewSftpBacked(cli, "upload")
	if err := sftp.Put(ctx, key, bytes.NewReader([]byte("bar"))); err != nil {
		log.Fatal(err)
	}

	reader, err := sftp.Get(ctx, key)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	if string(data) != "bar" {
		log.Fatalf("results do not match! expected %s, but got %s", "bar", string(data))
	}
	if err := sftp.Delete(ctx, key); err != nil {
		log.Fatal(err)
	}
}
```