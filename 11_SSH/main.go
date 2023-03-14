package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/wonderivan/logger"
)

type SSHClient struct {
	c *ssh.Client
}

func NewSSHClient(host string, username string, passwd string, timeout time.Duration) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		Timeout:         timeout * time.Second,
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(passwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//创建ssh连接
	sshClient, err := ssh.Dial("tcp", host, config)
	if err != nil {
		logger.Error("ssh Dial Failed : %s\n", err)
		return nil, err
	}
	return &SSHClient{c: sshClient}, nil

}

func (s *SSHClient) Exec(cmd string) ([]byte, error) {
	session, err := s.c.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return session.CombinedOutput(cmd)
}

func (s *SSHClient) Close() {
	s.c.Close()
}

func main() {
	logger.SetLogger(`{"Console": {"level": "INFO", "color" : true}}`)
	ssh, err := NewSSHClient("192.168.1.11:22", "root", "123456", 20)
	if err != nil {
		logger.Error("NewSSHClient Faield %s\n", err)
		return
	}
	fmt.Println("create NewSSHClient success")
	defer ssh.Close()
	out, err := ssh.Exec("cat /home/123")
	if err != nil {
		logger.Warn("Exec Faield %s\n", err)
	}
	fmt.Println("11111111111111111111111111111111")
	fmt.Println(string(out))
	out, err = ssh.Exec("cd /tmp")
	if err != nil {
		logger.Warn("Exec Faield %s\n", err)
	}
	fmt.Println("222222222222222222222222222222222")
	fmt.Println(string(out))

}
