package ssh

import (
	"bytes"
	"errors"
	"io"
	"net"

	"golang.org/x/crypto/ssh"

	"github.com/FreifunkBremen/freifunkmanager/lib/log"
)

type SSHResultHandler func([]byte, error)

type SSHResultStringHandler func(string, error)

func SSHResultToString(result []byte) string {
	if len(result) > 0 {
		result = result[:len(result)-1]
	}
	return string(result)
}

func SSHResultToStringHandler(handler SSHResultStringHandler) SSHResultHandler {
	return func(result []byte, err error) {
		handler(SSHResultToString(result), err)
	}
}

func (m *Manager) RunEverywhere(cmd string, handler SSHResultHandler) {
	for host, client := range m.clients {
		result, err := m.run(host, client, cmd)
		handler(result, err)
	}
}

func (m *Manager) RunOn(addr net.TCPAddr, cmd string) ([]byte, error) {
	client := m.ConnectTo(addr)
	if client != nil {
		return m.run(addr.IP.String(), client, cmd)
	}
	return nil, errors.New("no connection for runOn")
}

func (m *Manager) run(host string, client *ssh.Client, cmd string) ([]byte, error) {
	session, err := client.NewSession()
	defer session.Close()

	if err != nil {
		log.Log.Warnf("can not create session on %s: %s", host, err)
		delete(m.clients, host)
		return nil, err
	}
	stdout, err := session.StdoutPipe()
	buffer := &bytes.Buffer{}
	go io.Copy(buffer, stdout)
	if err != nil {
		log.Log.Warnf("can not create pipe for run on %s: %s", host, err)
		delete(m.clients, host)
		return nil, err
	}
	err = session.Run(cmd)
	if err != nil {
		log.Log.Warnf("could not run %s on %s: %s", cmd, host, err)
		delete(m.clients, host)
		return nil, err
	}
	var result []byte
	for {
		b, err := buffer.ReadByte()
		if err != nil {
			break
		}
		result = append(result, b)
	}
	return result, nil
}
