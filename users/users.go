package users

import (
	"crypto/hmac"
	"fmt"
	"os"

	"git.secrecy.rocks/stormentt/jumpd/groups"
	"git.secrecy.rocks/stormentt/jumpd/servers"
)

var Users UserList

type User struct {
	Name    string
	Groups  []string
	Pass    string
	Access  []string
	Default string

	GroupList   groups.GroupList
	ServerList  servers.ServerList
	AutoForward string
	Aliases     []string
}

type UserList []User

func Get(name string) *User {
	return Users.Get(name)
}

func Defaults() {
	Users.Defaults()
}

func (ul UserList) Defaults() {
	for i := range ul {
		user := &ul[i]
		user.Defaults()
	}
}

func (ul UserList) Get(name string) *User {
	for _, user := range ul {
		if user.Name == name {
			return &user
		}

		for _, alias := range user.Aliases {
			if alias == name {
				return &user
			}
		}
	}

	return nil
}

func (u User) HasAccess(serverName string) bool {
	for _, server := range u.ServerList {
		if server.Name == serverName {
			return true
		}
	}

	for _, group := range u.GroupList {
		if group.HasAccess(serverName) {
			return true
		}
	}

	return false
}

func (u User) CheckPassword(pass string) bool {
	return hmac.Equal([]byte(u.Pass), []byte(pass))
}

func (u *User) Defaults() {
	if len(u.Name) == 0 {
		fmt.Println("all users need names")
		os.Exit(1)
	}

	if len(u.Pass) == 0 {
		fmt.Println("all users need passwords")
		os.Exit(1)
	}

	for _, g := range u.Groups {
		group := groups.Get(g)
		if group != nil {
			u.GroupList = append(u.GroupList, *group)
		}
	}

	for _, s := range u.Access {
		server := servers.Get(s)
		if server != nil {
			u.ServerList = append(u.ServerList, *server)
		}
	}
}
