package service

import (
	"log"
	"urls/internal/repo"
)

type CreateUrlJob struct {
	Long  string
	Short string
}

type WriteExecutor struct {
	repo    repo.UrlRepo
	JobChan chan CreateUrlJob
	Cancel  chan bool
}

func NewWriteExecutor() *WriteExecutor {
	return &WriteExecutor{
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
				err := e.repo.CreateUrl(job.Short, job.Long)
				if err != nil {
					log.Printf("failed to save url in database")
				}
			case <-e.Cancel:
				return
			}
		}
	}()

	return e
}
