// main project main.go
package main

import (
	"log"

	"github.com/huangfly/hproxy/proxy"
)

func main() {
	proxy := proxy.NewProxySvr(8088)
	log.Fatal(proxy.ListenAndServe())

}
