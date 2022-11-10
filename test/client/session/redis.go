package gin_session

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"sync"
	"time"
)

//redis版的sessiondata的数据结构
type RedisSD struct {
	ID      string
	Data    map[string]interface{}
	rwLock  sync.RWMutex  // 读写锁，锁的是上面的Data
	expired int           // 过期时间
	client  *redis.Client // redis连接池
}

func (r *RedisSD) Get(key string) (value interface{}, err error) {
	// 获取读锁
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	value, ok := r.Data[key]
	if !ok {
		err = fmt.Errorf("invalid Key")
		return
	}
	return
}

func (r *RedisSD) Set(key string, value interface{}) {
	// 获取写锁
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	r.Data[key] = value
}

func (r *RedisSD) Del(key string) {
	// 删除key对应的键值对
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	delete(r.Data, key)
}

func (r *RedisSD) Save() {
	//将最新的sessiondata保存到redis中
	value, err := json.Marshal(r.Data)
	if err != nil {
		fmt.Printf("redis 序列化sessiondata失败 err=%v\n", err)
		return
	}
	//入库
	r.client.Set(r.ID, value, time.Duration(r.expired)*time.Second) //注意这里要用time.Duration转换一下
}

func (r *RedisSD) GetID() string { //为了拿到接口的ID数据
	return r.ID
}

//大仓库
type RedisMgr struct {
	Session map[string]SessionData //存储所有的session的一个大切片
	rwLock  sync.RWMutex
	client  *redis.Client //redis连接池
}

//NewRedisMgr  redis版初始化session仓库,构造函数
func NewRedisMgr() Mgr {
	//返回一个对象实例
	return &RedisMgr{
		Session: make(map[string]SessionData, 1024),
	}

}

//RedisMgr初始化
func (r *RedisMgr) Init(addr string, option ...string) { //这里的option...代表不定参数，参数个数不确定
	//	初始化redis连接池
	var (
		passwd string
		db     string
	)
	if len(option) == 1 {
		passwd = option[0]
	} else if len(option) == 2 {
		passwd = option[0]
		db = option[1]
	}
	//转换一下db数据类型，输入为string，需要转成int
	dbValue, err := strconv.Atoi(db)
	if err != nil {
		dbValue = 0 //如果转换失败，geidb一个默认值
	}
	r.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,  // no password set
		DB:       dbValue, // use default DB
	})

	_, err = r.client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

//GetSessionData 根据传进来的SessionID找到对应Session
func (r *RedisMgr) GetSessionData(sessionId string) (sd SessionData, err error) {

	//1.r.sesion已经从redis中拿到数据
	if r.Session == nil {
		err = r.LoadFromRedis(sessionId)
		if err != nil {
			return nil, err
		}
	}
	//2.r.session[sessionID]拿到sessionData
	// 获取读锁
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	sd, ok := r.Session[sessionId]
	if !ok {
		err = fmt.Errorf("无效的sessionId")
		return
	}
	return
}

//CreatSession 创建一个session记录
func (r *RedisMgr) CreatSession() (sd SessionData) {
	//1. 构造一个sessionID
	uuidObj := uuid.NewV4()
	//2.创建一个sessionData
	sd = NewRedisSessionData(uuidObj.String(), r.client) //从连接池中拿去一个client连接传给小红方块
	//3.创建对应关系
	r.Session[sd.GetID()] = sd
	//返回
	return
}

//NewRedisSessionData  的构造函数,用于构造sessiondata小仓库，小红块
func NewRedisSessionData(id string, client *redis.Client) SessionData {
	return &RedisSD{
		ID:     id,
		Data:   make(map[string]interface{}, 8),
		client: client,
	}
}

//加载数据库里的数据
func (r *RedisMgr) LoadFromRedis(sessionID string) (err error) {
	//1.连接redis
	//2.根据sessioniD拿到数据
	value, err := r.client.Get(sessionID).Result()
	if err != nil {
		//redis中wusessioinid对应的sessiondata
		fmt.Errorf("连接数据库失败")
		return
	}
	//3.反序列化成 r.session
	err = json.Unmarshal([]byte(value), &r.Session)
	if err != nil {
		//反序列化失败
		fmt.Println("连接数据库失败")
		return
	}
	return
}
