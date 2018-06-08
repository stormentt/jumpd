package proxy

import (
	"io"

	"golang.org/x/crypto/ssh"
)

func Proxy(fromCh ssh.Channel, fromReqs <-chan *ssh.Request, to *ssh.Client) {
	defer to.Close()
	toCh, toReqs, err := to.OpenChannel("session", nil)
	if err != nil {
		panic(err)
	}

	proxy(fromCh, fromReqs, toCh, toReqs)
}

func copyCh(out, in ssh.Channel, closer chan<- struct{}) {
	io.Copy(out, in)
	closer <- struct{}{}
}

func proxy(fromCh ssh.Channel, fromReqs <-chan *ssh.Request, toCh ssh.Channel, toReqs <-chan *ssh.Request) {
	defer fromCh.Close()
	defer toCh.Close()

	closer := make(chan struct{}, 1)
	go copyCh(fromCh, toCh, closer)
	go copyCh(toCh, fromCh, closer)

	for {
		select {
		case fReq := <-fromReqs:
			if fReq == nil {
				return
			}

			r, err := toCh.SendRequest(fReq.Type, fReq.WantReply, fReq.Payload)
			if err != nil {
				return
			}

			fReq.Reply(r, nil)
		case tReq := <-toReqs:
			if tReq == nil {
				return
			}

			r, err := fromCh.SendRequest(tReq.Type, tReq.WantReply, tReq.Payload)
			if err != nil {
				return
			}

			tReq.Reply(r, nil)
		case <-closer:
			return
		}
	}
}
