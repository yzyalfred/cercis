package network

import (
	"errors"
	"fmt"
	"github.com/yzyshine/cercis/log"
	"net"
	"sync"
	"sync/atomic"
)

type TCPServer struct {
	remoteAddr string
	state      uint8
	maxConnNum int
	listener   *net.TCPListener

	connectType uint8

	idSeed     uint32
	clientList map[uint32]*TCPClient
	clientLock sync.RWMutex

	msgParser *TCPMsgParser

	reciveMsgHandler func(uint32, []byte)
}

func (this *TCPServer) Init(ip string, port int) {
	this.remoteAddr = fmt.Sprintf("%s:%d", ip, port)
	this.maxConnNum = MAX_CONN_NUM

	this.idSeed = 0
	this.clientList = make(map[uint32]*TCPClient)
}

func (this *TCPServer) Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", this.remoteAddr)
	if err != nil {
		log.Fatal("%v", err)
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("%v", err)
	}

	log.Info("Start Server: %s", this.remoteAddr)

	this.listener = ln
	this.state = SS_ACCEPT

	go this.Run()
}

func (this *TCPServer) SetConnectType(connectType uint8) {
	this.connectType = connectType
}

func (this *TCPServer) Run() {
	for {
		tcpConn, err := this.listener.AcceptTCP()
		if err != nil {
			return
		}

		fmt.Printf("客户端：%s已连接！\n", tcpConn.RemoteAddr().String())

		this.handleConn(tcpConn, tcpConn.RemoteAddr().String())
	}
}

func (this *TCPServer) AssignClientId() uint32 {
	return atomic.AddUint32(&this.idSeed, 1)
}

func (this *TCPServer) GetClientById(clientId uint32) *TCPClient {
	this.clientLock.RLock()
	client, found := this.clientList[clientId]
	this.clientLock.RUnlock()
	if found {
		return client
	}
	return nil
}

func (this *TCPServer) AddClinet(tcpConn *net.TCPConn, addr string, connectType uint8) *TCPClient {
	client := &TCPClient{}
	client.remoteAddr = addr
	client.clientId = this.AssignClientId()
	client.server = this
	client.conn = tcpConn
	client.connectType = connectType
	this.clientLock.Lock()
	this.clientList[client.clientId] = client
	this.clientLock.Unlock()
	client.Start()
	return client
}

func (this *TCPServer) DelClient(clientId uint32) {
	this.clientLock.Lock()
	delete(this.clientList, clientId)
	this.clientLock.Unlock()
}

func (this *TCPServer) handleConn(tcpConn *net.TCPConn, addr string) bool {
	if tcpConn == nil {
		return false
	}

	pClient := this.AddClinet(tcpConn, addr, this.connectType)
	if pClient == nil {
		return false
	}

	return true
}

func (this *TCPServer) SendMsg(clientId uint32, args ...interface{}) error {
	client := this.GetClientById(clientId)
	if client != nil {
		return this.msgParser.Write(client, args...)
	}
	return errors.New("cleintId not found")
}

func (this *TCPServer) OnReciveMsg(clienId uint32, buf []byte) {
	this.reciveMsgHandler(clienId, buf)
}