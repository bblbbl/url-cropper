package buffer

import (
	"sync"
	"time"
	"urls/internal/repo"
	"urls/pkg/etc"
)

const (
	buffCap      = 1000
	flushBuffSec = 30
)

type FlushChan chan struct{}

type UrlBuffer struct {
	mu        sync.Mutex
	urlRepo   repo.UrlWriteRepo
	buff      []repo.Url
	flushTime int64
	flushChan FlushChan
}

func NewUrlBuffer(r repo.UrlWriteRepo) *UrlBuffer {
	b := &UrlBuffer{
		buff:      make([]repo.Url, 0, buffCap),
		flushChan: make(chan struct{}),
		urlRepo:   r,
	}

	b.updateFlushTime()
	return b
}

func (b *UrlBuffer) Flush(force bool) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if (len(b.buff) >= buffCap || time.Now().Unix() >= b.flushTime || force) && len(b.buff) > 0 {
		err := b.urlRepo.BatchCreateUrl(b.buff)
		if err != nil {
			etc.GetLogger().Error("failed to batch create urls")
		}

		b.buff = make([]repo.Url, 0, buffCap)
		b.updateFlushTime()

		return true
	}

	return false
}

func (b *UrlBuffer) Append(url repo.Url) {
	b.buff = append(b.buff, url)
}

func (b *UrlBuffer) DoneFlush() FlushChan {
	return b.flushChan
}

func (b *UrlBuffer) updateFlushTime() {
	b.flushTime = time.Now().Add(flushBuffSec * time.Second).Unix()
}
