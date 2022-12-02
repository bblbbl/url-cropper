package service

import (
	"context"
	"urls/internal/repo"
	"urls/pkg/etc"
)

type CreateUrlJob struct {
	Long string
	Hash string
}

type WriteExecutor struct {
	ctx     context.Context
	repo    repo.UrlRepo
	JobChan chan CreateUrlJob
	Cancel  chan bool
}

func NewWriteExecutor(ctx context.Context) *WriteExecutor {
	return &WriteExecutor{
		ctx:     ctx,
		repo:    repo.NewMysqlUrlRepo(),
		JobChan: make(chan CreateUrlJob, 10),
		Cancel:  make(chan bool),
	}
}

func (e *WriteExecutor) Start() *WriteExecutor {
	go func() {
		for {
			select {
			case job := <-e.JobChan:
				if err := e.repo.CreateUrl(repo.NewUrl(job.Hash, job.Long)); err != nil {
					etc.GetLogger().Error("failed to save url: " + job.Long)
				}
			case <-e.ctx.Done():
				return
			}
		}
	}()

	return e
}
