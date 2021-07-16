package actor

type ActorMgr struct {
	actorList map[string]IActor
}

func (this *ActorMgr) Init() {
	this.actorList = make(map[string]IActor)
}

func (this *ActorMgr) AddActor(name string, actor IActor)  {
	this.actorList[name] = actor
}

func (this *ActorMgr) GetActor(name string) IActor {
	actor, exist := this.actorList[name]
	if exist{
		return actor
	}
	return nil
}

func (this *ActorMgr) GetAllAcotr() map[string]IActor {
	return this.actorList
}

func (this *ActorMgr) SendMsgToActor(name string, clentId uint32, targetId uint64, msgId interface{}, msg interface{})  {
	actor, found := this.actorList[name]
	if found {
		actor.SendMsg(clentId, targetId, msgId, msg)
	}
}

func (this *ActorMgr) SendToActor(name string, clentId uint32, targetId uint64, buf []byte)  {
	actor, found := this.actorList[name]
	if found {
		actor.Send(clentId, targetId, buf)
	}
}

var GActorMgr = &ActorMgr{}