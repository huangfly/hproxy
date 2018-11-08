# hproxy
一个http的反向代理服务器，coding。

hproxy 提供http反向代理的功能，同时支持两种负载均衡算法分别是加权轮询以及hash算法，通过配置文件config.json配置

##配置文件

Method可以配置为round 或 hash

{
    
	"Method": "round", 
    
	"NodeServers": [
    
    {

	"Ip": "192.168.5.112:9090", 

	"Wight": 5
 
 }, 
 
 {

 "Ip": "192.168.5.113:9090", 

 "Wight": 5

 }, 

 {

 "Ip": "192.168.5.115:9090", 

 "Wight": 5

 }

 ]
}

##使用示例
import (
	"log"

	"github.com/huangfly/hproxy/proxy"
)

func main() {
	proxy := proxy.NewProxySvr(8088)
	log.Fatal(proxy.ListenAndServe())

}