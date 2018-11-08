// consishash
package balance

import (
	"crypto/sha1"
	_ "log"
	"sort"
	"strconv"
	"sync"
)

type HashRing []uint32

func (this HashRing) Len() int {
	return len(this)
}
func (this HashRing) Less(i, j int) bool {
	return this[i] < this[j]
}
func (this HashRing) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this HashRing) Sort() {
	sort.Sort(this)
}

type HashServer struct {
	VirtualNodes map[uint32]string
	Nodes        map[string]int
	Ring         HashRing
	Lock         *sync.RWMutex
}

//生成一个hash算法的负载均衡器
func NewHashServer() *HashServer {
	return &HashServer{
		VirtualNodes: make(map[uint32]string),
		Nodes:        make(map[string]int),
		Lock:         &sync.RWMutex{},
	}
}

//向负载均衡器里添加多个代理服务器节点
func (this *HashServer) AddNodes(args ...NodeServer) {
	for _, v := range args {
		this.addNode(v.Ip, v.Wight)
	}
}

func (this *HashServer) addNode(ip string, virualweight int) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	this.Nodes[ip] = virualweight
	this.BuildRing()
}

//删除节点
func (this *HashServer) DelNode(ip string) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	delete(this.Nodes, ip)
	this.BuildRing()
}

//根据ip计算hash值 找到对应的服务器
func (this *HashServer) LoadBalance(ip string) string {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	if this.Empty() {
		return ""
	}
	hashfn := sha1.New()
	hashfn.Write([]byte(ip))
	hashval := this.Byte2Uint32(hashfn.Sum(nil)[6:10])
	index := sort.Search(len(this.Ring), func(i int) bool { return this.Ring[i] >= hashval })
	if index == len(this.Ring) {
		index = 0
	}
	return this.VirtualNodes[this.Ring[index]]
}

//构建哈希环
func (this *HashServer) BuildRing() {
	for key, val := range this.Nodes {
		for i := 0; i < val; i++ {
			hashfn := sha1.New()
			hashfn.Write([]byte(key + strconv.Itoa(i)))
			hashval := this.Byte2Uint32(hashfn.Sum(nil)[6:10])
			this.VirtualNodes[hashval] = key
			this.Ring = append(this.Ring, hashval)
		}
	}
	this.Ring.Sort()
}

func (this *HashServer) Empty() bool {
	return len(this.Ring) == 0
}

func (this *HashServer) Byte2Uint32(src []byte) uint32 {
	if len(src) < 4 {
		return 0
	}
	dst := (uint32(src[3]) << 24) | (uint32(src[2]) << 16) | (uint32(src[1]) << 8) | (uint32(src[0]))
	return dst
}
