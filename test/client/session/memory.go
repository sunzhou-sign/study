package gin_session

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"sync"
)

//内存版的session服务

//SessionData支持的操作

//type MemSD struct {
//	ID string
//	Data map[string]interface{}
//	rwLock sync.RWMutex // 读写锁，锁的是上面的Data
//	// 过期时间
//}

//memory的sessiondata
type MemSD struct {
	ID     string
	Data   map[string]interface{}
	rwLock sync.RWMutex // 读写锁，锁的是上面的Data
	// 过期时间
}

// Get 根据key获取值
func (m *MemSD) Get(key string) (value interface{}, err error) {
	// 获取读锁
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	value, ok := m.Data[key]
	if !ok {
		err = fmt.Errorf("invalid Key")
		return
	}
	return
}

// Set 根据key获取值
func (m *MemSD) Set(key string, value interface{}) {
	// 获取写锁
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	m.Data[key] = value
}

// Del 删除Key对应的键值对
func (m *MemSD) Del(key string) {
	// 删除key对应的键值对
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	delete(m.Data, key)
}

//Save方法，被动设置的，因为要照顾redis版的接口
func (m *MemSD) Save() {
	return
}

//GetID 为了拿到接口的ID数据
func (m *MemSD) GetID() string {
	return m.ID
}

//管理全局的session
type MemoryMgr struct {
	Session map[string]SessionData //存储所有的session的一个大切片
	rwLock  sync.RWMutex           //读写锁，用于读多写少的情况，读锁可以重复的加，写锁互斥
}

//内存版初始化session仓库
func NewMemory() Mgr {
	return &MemoryMgr{
		Session: make(map[string]SessionData, 1024),
	}
}

//init方法
func (m *MemoryMgr) Init(addr string, option ...string) {
	//这里创建Init方法纯属妥协，其实memory版的并不需要初始化，前面NewMemory已经把活干完了
	//这里只是为了满足接口的定义，因为redis里需要这个方法取去连接数据库
	return
}

//GetSessionData 根据传进来的SessionID找到对应Session
func (m *MemoryMgr) GetSessionData(sessionId string) (sd SessionData, err error) {
	// 获取读锁
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	sd, ok := m.Session[sessionId]
	if !ok {
		err = fmt.Errorf("无效的sessionId")
		return
	}
	return
}

//CreatSession 创建一个session记录
func (m *MemoryMgr) CreatSession() (sd SessionData) {
	//1. 构造一个sessionID
	uuidObj := uuid.NewV4()
	//2.创建一个sessionData
	sd = NewMemorySessionData(uuidObj.String())
	//3.创建对应关系
	m.Session[sd.GetID()] = sd
	//返回
	return
}

//NewRedisSessionData  的构造函数,用于构造sessiondata小仓库，小红块
func NewMemorySessionData(id string) SessionData {
	return &MemSD{
		ID:   id,
		Data: make(map[string]interface{}, 8),
	}
}
