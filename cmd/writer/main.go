package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"urls/internal/messaging"
	"urls/internal/repo"
	"urls/pkg/database"
	"urls/pkg/etc"
)

const (
	buffCap      = 1000
	flushBuffSec = 30
)

type urlBuffer struct {
	mu   sync.Mutex
	buff []repo.Url
}

func main() {
	etc.InitLogger()

	ctx, cancelFunc := context.WithCancel(context.Background())

	consumer := messaging.NewKafkaUrlConsumer(ctx)
	urlRepo := repo.NewMysqlUrlRepo()

	go consumer.ListenMessages()

	buff := &urlBuffer{
		buff: make([]repo.Url, 0, buffCap),
	}

	go func(buff *urlBuffer) {
		for url := range consumer.Messages() {
			buff.buff = append(buff.buff, url)
		}
	}(buff)

	go func(buff *urlBuffer, ctx context.Context) {
		flushTime := time.Now().Add(flushBuffSec * time.Second).Unix()
		for {
			select {
			case <-ctx.Done(): // todo: flush buffer
				break
			default:
				buff.mu.Lock()
				if len(buff.buff) >= buffCap || time.Now().Unix() >= flushTime {
					err := urlRepo.BatchCreateUrl(buff.buff)
					if err != nil {
						etc.GetLogger().Error("failed to batch create urls")
					}

					buff.buff = make([]repo.Url, 0, buffCap)
					flushTime = time.Now().Add(flushBuffSec * time.Second).Unix()
					buff.mu.Unlock()
				} else {
					buff.mu.Unlock()
					time.Sleep(2 * time.Second)
				}
			}
		}
	}(buff, ctx)

	etc.GetLogger().Info("urls writer started")

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-term
	terminate(cancelFunc)
}

func terminate(cancelFunc context.CancelFunc) {
	cancelFunc()
	time.Sleep(2 * time.Second)

	_ = database.GetConnection().Close()
}
