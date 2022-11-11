package lru_groutine_local

import (
	"fmt"
	"log"
	"sync"
)

// 设计了一个回调函数(callback)，在缓存不存在时，调用这个函数，得到源数据
type Getter interface {
	Get(key string) ([]byte, error)
}

// 函数接口类型
type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

/* 核心代码 */
// 一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name。
//比如可以创建三个 Group，缓存学生的成绩命名为 scores，缓存学生信息的命名为 info，缓存学生课程的命名为 courses
type Group struct {
	name      string
	getter    Getter // 第二个属性是 getter Getter，即缓存未命中时获取源数据的回调(callback)。
	mainCache cache  // 第三个属性是 mainCache cache，即一开始实现的并发缓存。

	peers PeerPicker
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group) // 保存不同命名空间缓存
)

// 构建函数 NewGroup 用来实例化 Group，并且将 group 存储在全局变量 groups 中
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g

	return g
}

// GetGroup 用来特定名称的 Group，这里使用了只读锁 RLock()，因为不涉及任何冲突变量的写操作。
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

// 获取缓存信息
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 命中缓存返回
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	// 未命中，查找数据源
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	// 修改 load 方法，使用 PickPeer() 方法选择节点，若非本机节点，则调用 getFromPeer() 从远程获取。若是本机节点或失败，则回退到 getLocally()
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[GeeCache] Failed to get from peer ", err)
		}
	}

	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用数据源查询函数
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	// 数据源查询到，加入缓存
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// 将 实现了 PeerPicker 接口的 HTTPPool 注入到 Group 中
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// 使用实现了 PeerGetter 接口的 httpGetter 从访问远程节点，获取缓存值
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}
