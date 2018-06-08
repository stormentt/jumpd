package servers

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func (s *Server) Connect(t *terminal.Terminal, username string) (*ssh.Client, error) {
	if len(s.RemoteUser) != 0 {
		username = s.RemoteUser
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PasswordCallback(func() (secret string, err error) {
				return t.ReadPassword(fmt.Sprintf("%s@%s password: ", username, s.Name))
			}),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote_addr net.Addr, key ssh.PublicKey) error {
			//TODO: host key checking
			return nil
		}),
	}

	return ssh.Dial("tcp", s.GetPath(), config)
}
