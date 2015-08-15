package sss

type LeakyBuf struct {
	bufSize  int
	freeList chan []byte
}

const (
	leakyBufSize = 4049
	maxNBuf      = 2048
)

var leakyBuf = NewLeakyBuf(maxNBuf, leakyBufSize)

//新建缓存
func NewLeakyBuf(n, bufSize int) *LeakyBuf {
	return &LeakyBuf{
		bufSize:  bufSize,
		freeList: make(chan []byte, n),
	}
}

//获取一个缓存空间
func (lb *LeakyBuf) Get() (b []byte) {
	select {
	case b = <-lb.freeList:
	default:
		b = make([]byte, lb.bufSize)
	}
	return
}

//释放一个缓存空间
func (lb *LeakyBuf) Put(b []byte) {
	if len(b) != lb.bufSize {
		panic("invalid buffer size that's put into leaky buffer")
	}
	select {
	case lb.freeList <- b:
	default:
	}
	return
}
