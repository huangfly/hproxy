package balance

import (
	_ "log"
	"sync"
)

type Server struct {
	Ip           string
	Wight        int
	CurrentWitht int
}

type WightServer struct {
	ServerBuf []Server
	Lock      *sync.Mutex
}

func (this *WightServer) getMaxWight() int {
	var maxWight int
	var maxIndex int
	sum := 0
	this.Lock.Lock()
	defer this.Lock.Unlock()
	for index, svr := range this.ServerBuf {
		if svr.CurrentWitht > maxWight {
			maxWight = svr.CurrentWitht
			maxIndex = index
		}
		sum += svr.CurrentWitht
	}

	this.ServerBuf[maxIndex].CurrentWitht -= sum
	return maxIndex
}

func (this *WightServer) addWight() {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	for index, svr := range this.ServerBuf {
		this.ServerBuf[index].CurrentWitht += svr.Wight
	}
}

func (this *WightServer) LoadBalance() string {
	this.addWight()
	index := this.getMaxWight()
	return this.ServerBuf[index].Ip
}

func NewWightServer(args ...Server) *WightServer {
	wightsvr := &WightServer{
		ServerBuf: make([]Server, 0),
		Lock:      &sync.Mutex{},
	}
	for _, arg := range args {
		wightsvr.ServerBuf = append(wightsvr.ServerBuf, arg)
	}
	return wightsvr
}
