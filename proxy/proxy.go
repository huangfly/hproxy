// proxy project proxy.go
package proxy

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/huangfly/hproxy/balance"
)

var proxyHeaders = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

type ProxySvr struct {
	Trans *http.Transport
}

func NewProxySvr() *http.Server {
	return &http.Server{
		Addr:    ":8989",
		Handler: &ProxySvr{Trans: &http.Transport{Proxy: http.ProxyFromEnvironment, DisableKeepAlives: true}},
	}
}

//重写ServerHttp接口
func (this *ProxySvr) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println("start proxy")
	mgr, _ := balance.GetBalanceInstance("round")
	ip := mgr.LoadBalance(req.Host)
	req.Host = ip
	req.URL.Host = ip
	req.URL.Scheme = "http"
	switch req.Method {
	case "CONNECT":
		ProxyHttpsHandler(rw, req)
	case "GET":

	default:
		this.ProxyHttpHandler(rw, req)
	}
}

//转发http
func (this *ProxySvr) ProxyHttpHandler(rw http.ResponseWriter, req *http.Request) {

	this.DelHeads(req)
	addXForwardIpToHead(req)

	res, err := this.Trans.RoundTrip(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	this.RewriteHead(rw.Header(), res.Header)
	rw.WriteHeader(res.StatusCode)
	_, err = io.Copy(rw, res.Body)
	if err != nil {
		if err != io.EOF {
			return
		}
	}
}

func (this *ProxySvr) ProxyHttpsHandler(rw http.ResponseWriter, req *http.Request) {
	hijack, _ := rw.(http.Hijacker)
	clientConn, _, err := hijack.Hijack()
	if err != nil {
		log.Println("https hijack error : ", err.Error())
		return
	}

	proxyConn, err := net.Dial("tcp", req.URL.Host)
	if err != nil {
		log.Println("https dial server error : ", err.Error())
		return
	}
	clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	go func() {
		_, err := io.Copy(proxyConn, clientConn)
		if err != nil && err != io.EOF {
			log.Println("proxy to server failed : ", err.Error())
			return
		}
		proxyConn.Close()
		clientConn.Close()
	}()
	go func() {
		_, err := io.Copy(clientConn, proxyConn)
		if err != nil && err != io.EOF {
			log.Println("proxy to server failed : ", err.Error())
			return
		}
		clientConn.Close()
		proxyConn.Close()
	}()
}

//添加X-Forwarded-For
func addXForwardIpToHead(req *http.Request) {
	if ip, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if proxyip, ok := req.Header["X-Forwarded-For"]; ok {
			ip = strings.Join(proxyip, ", ") + ", " + ip
		}
		req.Header.Set("X-Forwarded-For", ip)
	}
}

//重写头部
func (this *ProxySvr) RewriteHead(dst, src http.Header) {
	for dstkey, _ := range dst {
		dst.Del(dstkey)
	}
	for srckey, srcv := range src {
		for _, v := range srcv {
			dst.Add(srckey, v)
		}
	}
}

//删除hop-to-hop头部
func (this *ProxySvr) DelHeads(req *http.Request) {
	for _, head := range proxyHeaders {
		if req.Header.Get(head) != "" {
			req.Header.Del(head)
		}
	}
}
