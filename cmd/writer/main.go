package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urls/internal/buffer"
	"urls/internal/messaging"
	"urls/internal/repo"
	"urls/pkg/database"
	"urls/pkg/etc"
)

func main() {
	etc.InitLogger()

	ctx, cancelFunc := context.WithCancel(context.Background())

	consumer := messaging.NewKafkaUrlConsumer(ctx)

	go consumer.ListenMessages()

	buff := buffer.NewUrlBuffer(repo.NewMysqlUrlWriteRepo())

	go func(buff *buffer.UrlBuffer) {
		for url := range consumer.Messages() {
			buff.Append(url)
		}
	}(buff)

	go func(buff *buffer.UrlBuffer, ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				buff.Flush(true)
				buff.DoneFlush() <- struct{}{}
				break
			default:
				if flushed := buff.Flush(false); !flushed {
					time.Sleep(2 * time.Second)
				}
			}
		}
	}(buff, ctx)

	etc.GetLogger().Info("urls writer started")

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-term
	terminate(cancelFunc, buff.DoneFlush())
}

func terminate(cancelFunc context.CancelFunc, flushChan buffer.FlushChan) {
	cancelFunc()
	<-flushChan

	database.CloseMysqlWriteConnection()
}
