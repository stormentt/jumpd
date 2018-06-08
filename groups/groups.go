package groups

import "git.secrecy.rocks/stormentt/jumpd/servers"

var Groups GroupList

type Group struct {
	Name       string
	Access     []string
	ServerList servers.ServerList
}

type GroupList []Group

func Defaults() {
	Groups.Defaults()
}

func Get(name string) *Group {
	return Groups.Get(name)
}

func (gl GroupList) Defaults() {
	for i := range gl {
		group := &gl[i]
		group.Defaults()
	}
}

func (gl GroupList) Get(name string) *Group {
	for _, group := range gl {
		if group.Name == name {
			return &group
		}
	}

	return nil
}

func (g *Group) Defaults() {
	for _, s := range g.Access {
		server := servers.Get(s)
		if server != nil {
			g.ServerList = append(g.ServerList, *server)
		}
	}
}

func (g Group) HasAccess(serverName string) bool {
	for _, server := range g.Access {
		if server == serverName {
			return true
		}
	}

	return false
}
