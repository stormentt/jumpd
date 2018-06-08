package auth

import (
	"fmt"

	"git.secrecy.rocks/stormentt/jumpd/users"
	"golang.org/x/crypto/ssh"
)

func Password(conn ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	user := users.Get(conn.User())
	if user == nil {
		return nil, fmt.Errorf("user doesn't exist: %s\n", conn.User())
	}
	if !user.CheckPassword(string(pass)) {
		return nil, fmt.Errorf("bad password for %s\n", conn.User())
	}

	return &ssh.Permissions{
		Extensions: map[string]string{
			"user": user.Name,
		},
	}, nil
}
