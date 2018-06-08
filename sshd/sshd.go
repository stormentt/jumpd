package sshd

import (
	"fmt"
	"net"
	"strconv"

	"git.secrecy.rocks/stormentt/jumpd/auth"
	"git.secrecy.rocks/stormentt/jumpd/proxy"
	"git.secrecy.rocks/stormentt/jumpd/servers"
	"git.secrecy.rocks/stormentt/jumpd/users"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func Start() {
	hostKey, err := GetHostKey()
	if err != nil {
		log.Fatal(err)
	}

	config := &ssh.ServerConfig{
		PasswordCallback: auth.Password,
	}

	config.AddHostKey(hostKey)

	host := viper.GetString("config.host")
	port := viper.GetString("config.port")
	address := fmt.Sprintf("%s:%s", host, port)

	socket, err := net.Listen("tcp", address)
	if err != nil {
		log.WithFields(log.Fields{
			"address": address,
			"error":   err,
		}).Fatal("unable to listen for connections")
	}

	for {
		nConn, err := socket.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("unable to accept connection")
		} else {
			go HandleSSH(nConn, config)
		}
	}
}

func HandleSSH(nConn net.Conn, config *ssh.ServerConfig) {
	conn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to complete handshake")
		return
	}

	go ssh.DiscardRequests(reqs)

	user := users.Get(conn.Permissions.Extensions["user"])
	log.WithFields(log.Fields{
		"user": conn.Permissions.Extensions["user"],
	}).Info("authorized")

	for newCh := range chans {
		ch, _, heldReqs := acceptNewChannel(newCh)
		defer ch.Close()
		if ch == nil {
			break
		}

		var server *servers.Server
		if len(user.Default) != 0 {
			server = servers.Get(user.Default)
		}

		if server == nil {
			server = interact(ch, user)
		}

		log.WithFields(log.Fields{
			"server": server,
		}).Info("auth")

		t := terminal.NewTerminal(ch, "> ")

		toCl, err := server.Connect(t, user.Name)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("unable to connect to backend")
			break
		}
		proxy.Proxy(ch, heldReqs, toCl)

		fmt.Fprintf(ch, "\r\n")
	}
}

func acceptNewChannel(newCh ssh.NewChannel) (ssh.Channel, <-chan *ssh.Request, <-chan *ssh.Request) {
	if newCh.ChannelType() != "session" {
		newCh.Reject(ssh.UnknownChannelType, "unknown channel type")
		return nil, nil, nil
	}

	ch, reqs, err := newCh.Accept()
	if err != nil {
		log.Fatalf("Could not accept channel: %v", err)
	}

	heldReqs := make(chan *ssh.Request, 5)
	go replyToReqs(reqs, heldReqs)

	return ch, reqs, heldReqs
}

func interact(ch ssh.Channel, user *users.User) (server *servers.Server) {
	t := terminal.NewTerminal(ch, "> ")
	fmt.Fprintf(ch, "Hello, %s\r\n", user.Name)
	for i := 0; i < 3; i++ {
		printServerList(ch, user)
		server = selectServer(ch, t)

		// TODO: access check
		if server == nil {
			fmt.Fprintf(ch, "invalid server\r\n")
			continue
		}

		fmt.Fprintf(ch, "connecting to %s...\r\n", server.Name)
		return server
	}

	return nil
}

func selectServer(ch ssh.Channel, t *terminal.Terminal) *servers.Server {
	input, err := t.ReadLine()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to read line")
		fmt.Fprintf(ch, "unable to read line: %s\r\n", err)
		return nil
	}

	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Fprintf(ch, "error: %s\r\n", err)
		return nil
	}

	return servers.Get(id)
}

func printServerList(ch ssh.Channel, user *users.User) {
	fmt.Fprintf(ch, "Servers:\r\n")
	for _, server := range user.ServerList {
		fmt.Fprintf(ch, "  (%d) %s\r\n", server.ID, server.Name)
	}
}

func replyToReqs(inreqs <-chan *ssh.Request, heldReqs chan<- *ssh.Request) {
	for req := range inreqs {
		if req.Type == "pty-req" || req.Type == "shell" {
			req.Reply(true, nil)
			req.WantReply = false
		}

		heldReqs <- req
	}
}
