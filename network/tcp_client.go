package network

import (
	"github.com/yzyshine/cercis/log"
	"github.com/yzyshine/cercis/utils"
	"io"
	"net"
)

type TCPClient struct {
	clientId   uint32
	remoteAddr string
	state      uint8

	server      *TCPServer
	conn        *net.TCPConn
	connectType uint8

	sendChan chan []byte
}

func (this *TCPClient) Init() {
	if this.state == CONNECT_CLIENT {
		this.sendChan = make(chan []byte, SEND_CHAN_SIZE)
	}
}

func (this *TCPClient) Start() {
	if this.server == nil {
		return
	}

	this.state = SS_ACCEPT
	this.conn.SetNoDelay(true)

	go this.Run()
}

func (this *TCPClient) Close() {
	if this.conn != nil {
		this.conn.Close()
	}

	if this.server != nil {
		this.server.DelClient(this.clientId)
	}

	this.Clear()
}

func (this *TCPClient) Clear() {
	this.state = SS_SHUT_DOWN
	this.clientId = 0
	this.remoteAddr = ""

	this.server = nil
	this.conn = nil
}

func (this *TCPClient) Run() {
	loop := func() bool {
		defer func() {
			if err := recover(); err != nil {
				utils.TraceCode(err)
			}
		}()

		if this.state == SS_SHUT_DOWN || this.conn == nil {
			return false
		}

		buf, err := this.server.msgParser.Read(this)
		if err == io.EOF {
			log.Info("client: %s close\n", this.remoteAddr)
			return false
		}

		if err != nil {
			return false
		}

		this.server.OnReciveMsg(this.clientId, buf)

		return true
	}

	for {
		if !loop() {
			break
		}
	}

	this.Close()
}

func (this *TCPClient) IsShutDown() bool  {
	return this.state == SS_SHUT_DOWN
}

func (this *TCPClient) SendLoop() {
	for {
		select {
		case buff := <-this.sendChan:
			if buff == nil { //信道关闭
				return
			} else {
				this.DoSend(buff)
			}
		}
	}
}

func (this *TCPClient) Send(buf []byte) {
	defer func() {
		if err := recover(); err != nil{
			utils.TraceCode(err)
		}
	}()

	if this.state == SS_SHUT_DOWN {
		return
	}

	if this.connectType == CONNECT_CLIENT {
		this.sendChan <- buf
	} else {
		this.DoSend(buf)
	}
}

func (this *TCPClient) DoSend(buf []byte) {
	if this.conn == nil {
		return
	}

	this.conn.Write(buf)
}
