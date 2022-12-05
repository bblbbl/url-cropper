package buffer

import (
	"sync"
	"time"
	"urls/internal/repo"
	"urls/pkg/etc"
)

type FlushChan chan struct{}

type UrlBuffer struct {
	mu           sync.Mutex
	urlRepo      repo.UrlWriteRepo
	buff         []repo.Url
	flushTime    int64
	flushChan    FlushChan
	buffCap      int
	buffFlushSec int
}

func NewUrlBuffer(r repo.UrlWriteRepo, buffCap, buffFlushSec int) *UrlBuffer {
	b := &UrlBuffer{
		buff:         make([]repo.Url, 0, buffCap),
		flushChan:    make(chan struct{}),
		urlRepo:      r,
		buffCap:      buffCap,
		buffFlushSec: buffFlushSec,
	}

	b.updateFlushTime()
	return b
}

func (b *UrlBuffer) Flush(force bool) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if (len(b.buff) >= b.buffCap || time.Now().Unix() >= b.flushTime || force) && len(b.buff) > 0 {
		err := b.urlRepo.BatchCreateUrl(b.buff)
		if err != nil {
			etc.GetLogger().Error("failed to batch create urls")
		}

		b.buff = make([]repo.Url, 0, b.buffCap)
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
	b.flushTime = time.Now().Add(time.Duration(b.buffFlushSec) * time.Second).Unix()
}
