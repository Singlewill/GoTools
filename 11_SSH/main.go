package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/wonderivan/logger"
)

func dialSSH(addr string, username string, passwd string, timeout int) {
	//state := statIDLE
	config := &ssh.ClientConfig{
		Timeout:         time.Duration(timeout) * time.Second,
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.Password(passwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//创建ssh连接
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		logger.Error("ssh Dial Failed : %s\n", err)
		return
	}
	defer sshClient.Close()
	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		logger.Error("ssh NewSession Failed : %s\n", err)
		return
	}
	defer session.Close()

	

	//直接cd到/tmp目录
	_, err = session.CombinedOutput("cd /tmp")
	if err != nil {
		logger.Error("ssh ConbinedOutout Failed : %s\n", err)
		return
	}
	//cat /bin/echo判断架构
	combo, err := session.CombinedOutput("cat /bin/echo")
	if err != nil {
		logger.Error("ssh ConbinedOutout Failed : %s\n", err)
		return
	}

	fmt.Println(string(combo))
}

func main() {
	logger.SetLogger(`{"Console": {"level": "INFO", "color" : true}}`)
	dialSSH("192.168.1.6:22", "ll", "kalo", 10)

}
