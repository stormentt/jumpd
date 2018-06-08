package servers

import (
	"fmt"
	"os"
)

var Servers ServerList

type Server struct {
	ID         int
	Host       string
	Name       string
	Port       int
	RemoteUser string
}
type ServerList []Server

func Get(by interface{}) *Server {
	return Servers.Get(by)
}

func Defaults() {
	Servers.Defaults()
}

func (s *Server) GetPath() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s ServerList) Defaults() {
	idList := make(map[int]struct{})
	for i := range s {
		server := &s[i]

		if server.ID == 0 {
			server.ID = i + 1
		}

		_, ok := idList[server.ID]
		if ok {
			fmt.Printf("dupe ID %d\n", server.ID)
			os.Exit(1)
		}
		idList[server.ID] = struct{}{}

		if len(server.Host) == 0 {
			server.Host = "localhost"
		}

		if server.Port == 0 {
			server.Port = 22
		}

		if len(server.Name) == 0 {
			server.Name = server.GetPath()
		}
	}
}

func (s ServerList) Get(by interface{}) *Server {
	switch by.(type) {
	case string:
		return s.GetByName(by.(string))
	case int:
		return s.GetByID(by.(int))
	default:
		return nil
	}
}

func (s ServerList) GetByName(name string) *Server {
	for _, server := range s {
		if server.Name == name {
			return &server
		}
	}

	return nil
}

func (s ServerList) GetByID(id int) *Server {
	for _, server := range s {
		if server.ID == id {
			return &server
		}
	}

	return nil
}

func (s ServerList) GetSlice(names []string) ServerList {
	var new ServerList
	for _, name := range names {
		got := s.Get(name)
		if got != nil {
			new = append(new, *got)
		}
	}
	return new
}
