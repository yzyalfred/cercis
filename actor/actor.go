package actor

import (
	"github.com/yzyshine/cercis/processor"
	"github.com/yzyshine/cercis/utils"
	"github.com/yzyshine/cercis/utils/mpsc"
	"sync/atomic"
	"time"
)

type CallIO struct {
	ClienId  uint32
	TargetId uint64

	Buff []byte
}

type msgCallFuncType func(uint32, uint64, interface{}, interface{})

type Actor struct {
	isStart bool

	eventChan chan uint8

	// msg
	msgChan     chan bool
	msgFlag     int32
	msgQuene    *mpsc.Queue
	msgCallFunc msgCallFuncType

	// timer
	timer         *time.Ticker
	timerCallFunc func()

	// Processor
	Processor processor.IProcessor
}

func NewActor() *Actor {
	actor := &Actor{}
	actor.Init()
	return actor
}

func (this *Actor) Init() {
	this.isStart = false

	this.eventChan = make(chan uint8)

	this.msgChan = make(chan bool)
	this.msgFlag = 0
	this.msgQuene = mpsc.New()

	// no timer
	this.timer = time.NewTicker(1<<63 - 1)
	this.timerCallFunc = nil

	// pb Processor default
	this.Processor = processor.NewPBProcessor()
}

func (this *Actor) Start() {
	if !this.isStart {
		go this.run()
		this.isStart = true
	}
}

func (this *Actor) Stop() {
	this.eventChan <- ACTOR_EVENT_CLOSE
}

func (this *Actor) SendMsg(clientId uint32, targetId uint64, msgId interface{}, msg interface{}) {
	buf, err := this.Processor.Marshal(msgId, msg)
	if err == nil {
		this.Send(clientId, targetId, buf)
	}
}

func (this *Actor) Send(clentId uint32, targetId uint64, buf []byte) {
	this.msgQuene.Push(CallIO{
		ClienId:  clentId,
		TargetId: targetId,
		Buff:     buf,
	})
	if atomic.CompareAndSwapInt32(&this.msgFlag, 0, 1) {
		this.msgChan <- true
	}
}

func (this *Actor) RegisterMsgCallFunc(f msgCallFuncType) {
	this.msgCallFunc = f
}

func (this *Actor) RegisterTimer(duration time.Duration, callFunc func()) {
	this.timer.Stop()
	this.timer = time.NewTicker(duration)
	this.timerCallFunc = callFunc
}

func (this *Actor) SetProcessor(processorType uint8) {
	if processorType == processor.PROCESSOR_TYPE_PB {
		this.Processor = processor.NewPBProcessor()
	}
}

func (this *Actor) clear() {
	this.isStart = false

	if this.timer != nil {
		this.timer.Stop()
		this.timerCallFunc = nil
	}
}

func (this *Actor) run() {
	for {
		if !this.loop() {
			break
		}
	}

	this.clear()
}

func (this *Actor) loop() bool {
	defer func() {
		if err := recover(); err != nil {
			utils.TraceCode()
		}
	}()
	select {
	case <-this.msgChan:
		this.consumeMsg()
	case eventId := <-this.eventChan:
		if eventId == ACTOR_EVENT_CLOSE {
			return false
		}
	case <-this.timer.C:
		this.timerCallFunc()
	}

	return true
}

func (this *Actor) RegisterMsg(msgId interface{}, msg interface{}) {
	this.Processor.Register(msgId, msg)
}

func (this *Actor) consumeMsg() {
	for data := this.msgQuene.Pop(); data != nil; data = this.msgQuene.Pop() {
		this.handleMsg(data.(CallIO))
	}
	atomic.StoreInt32(&this.msgFlag, 0)
}

func (this *Actor) handleMsg(callIo CallIO) {
	msgId, msg, err := this.Processor.Unmarshal(callIo.Buff)
	if err == nil {
		this.msgCallFunc(callIo.ClienId, callIo.TargetId, msgId, msg)
	}
}
