// balancemgr.go
package balance

import (
	"errors"
	_ "log"
)

var BalanceInitErr = errors.New("not found the specific balance method")

type BalanceMgr interface {
	LoadBalance(ip string) string
}

func GetBalanceInstance(method string) (BalanceMgr, error) {
	switch method {
	case "hash":
		mgr := NewHashServer()
		mgr.AddNode("192.168.5.112:9090", 2)
		return mgr, nil
	case "round":
		s := NodeServer{"192.168.5.112:8088", 1, 0}

		mgr := NewWightServer(s)
		return mgr, nil
	default:
	}
	return nil, BalanceInitErr
}
