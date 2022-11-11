package lru_groutine_local

// ByteView 只有一个数据成员，b []byte，b 将会存储真实的缓存值。选择 byte 类型是为了能够支持任意的数据类型的存储，例如字符串、图片等
type ByteView struct {
	b []byte
}

func (b ByteView) Len() int {
	return len(b.b)
}

// b 是只读的，使用 ByteSlice() 方法返回一个拷贝，防止缓存值被外部程序修改
func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.b)
}

func (b ByteView) String() string {
	return string(b.b)
}

func cloneBytes(dataB []byte) []byte {
	c := make([]byte, len(dataB))
	copy(c, dataB)
	return c
}
