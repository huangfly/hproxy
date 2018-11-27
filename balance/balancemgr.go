// balancemgr.go
package balance

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//定义负载均衡接口interface
type BalanceMgr interface {
	LoadBalance(ip string) string
}

//定义json配置文件格式
type JsonConfig struct {
	Method      string
	NodeServers []NodeServer
}

//定义单台服务器,ip,权重，配置文件中不需要配置Currentwitht
//这玩意hash算法不需要，只有加权轮询算法需要，但是为了通用性
//而且这个字段定义在结构体不定义在配置文件，解析不影响，初始化为0
type NodeServer struct {
	Ip           string
	Wight        int
	CurrentWitht int
}

var BalanceInitErr = errors.New("not found the specific balance method")
var conf JsonConfig

func printCurrentPath() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	strings.Replace(dir, "\\", "/", -1)
	log.Println("current path is ", dir)
}

func init() {
	data, err := ioutil.ReadFile("../conf/config.json")
	if err != nil {
		log.Println("read config.json file error, ", err.Error())
		printCurrentPath()
		os.Exit(-1)
	}
	err = json.Unmarshal([]byte(data), &conf)
	if err != nil {
		log.Println("json unmarshal error, ", err.Error())
		os.Exit(-1)
	}
	log.Println(conf)
}

func GetBalanceInstance() (BalanceMgr, error) {
	switch conf.Method {
	case "hash":
		mgr := NewHashServer()
		mgr.AddNodes(conf.NodeServers...)
		return mgr, nil
	case "round":
		mgr := NewWightServer(conf.NodeServers...)
		return mgr, nil
	default:
	}
	return nil, BalanceInitErr
}
