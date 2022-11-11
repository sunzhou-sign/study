package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"web/gee"
	lru_groutine_local "web/gee-cache/lru-groutine-local"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

// web框架测试
func main01() {
	r := gee.New()
	r.Use(gee.Logger()) // global midlleware
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, fmt.Sprintf("hello %s, you're at %s\n", c.Param("name"), c.Path))
		})
	}

	r.Run(":9999")
}

// 分布式缓存测试
// 缓存雪崩：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。常因为缓存服务器宕机，或缓存设置了相同的过期时间引起。
func main02() {
	// 模拟数据源
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	// 数据源获取函数
	var geeter lru_groutine_local.Getter = lru_groutine_local.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	})

	// 新建缓存
	lru_groutine_local.NewGroup("scores", 2<<10, geeter)

	addr := "localhost:9999"
	// 拦截所以请求
	httpCache := lru_groutine_local.NewHTTPPool(addr)
	log.Fatalln(http.ListenAndServe(addr, httpCache))
}

func createGroup() *lru_groutine_local.Group {
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	return lru_groutine_local.NewGroup("scores", 2<<10, lru_groutine_local.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *lru_groutine_local.Group) {
	peers := lru_groutine_local.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *lru_groutine_local.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

// 分布式节点
func main03() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], addrs, gee)
}
