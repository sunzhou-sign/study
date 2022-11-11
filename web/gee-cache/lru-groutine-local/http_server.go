package lru_groutine_local

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"web/gee-cache/consistenthash"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self     string
	basePath string

	sync.Mutex
	peers       *consistenthash.Map    // 根据具体的 key 选择节点
	httpGetters map[string]*httpGetter // 映射远程节点与对应的 httpGetter, httpGetter 与远程节点的地址 baseURL 有关
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (h *HTTPPool) Log(str string) {
	log.Printf("[Server %s] %s", h.self, str)
}

func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	h.Log(r.Method + " " + r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

func (h *HTTPPool) Set(peers ...string) {
	h.Lock()
	defer h.Unlock()

	h.peers = consistenthash.New(defaultReplicas, nil)
	h.peers.Add(peers...)

	h.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseURL: peer + h.basePath}
	}
}

func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.Lock()
	defer h.Unlock()

	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		h.Log(fmt.Sprintf("Pick peer %s", peer))
		return h.httpGetters[peer], true
	}

	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)
