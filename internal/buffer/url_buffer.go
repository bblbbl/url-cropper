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

type UrlBuffer struct {
	mu        sync.Mutex
	urlRepo   repo.UrlRepo
	buff      []repo.Url
	flushTime int64
}

func NewUrlBuffer(r repo.UrlRepo) *UrlBuffer {
	b := &UrlBuffer{
		buff:    make([]repo.Url, 0, buffCap),
		urlRepo: r,
	}

	b.updateFlushTime()
	return b
}

func (b *UrlBuffer) Flush(force bool) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.buff) >= buffCap || time.Now().Unix() >= b.flushTime || force {
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

func (b *UrlBuffer) updateFlushTime() {
	b.flushTime = time.Now().Add(flushBuffSec * time.Second).Unix()
}
