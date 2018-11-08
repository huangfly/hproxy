// main project main.go
package main

import (
	"log"
	_ "os"

	_ "github.com/huangfly/hproxy/balance"
	"github.com/huangfly/hproxy/proxy"
)

func main() {
	proxy := proxy.NewProxySvr()
	log.Fatal(proxy.ListenAndServe())

}
