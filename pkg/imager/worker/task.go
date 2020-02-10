package worker

import (
	"sync"
)

type Task interface {
	Run(wg *sync.WaitGroup)
}

type TaskImpl struct {
	Err   error
	Image interface{}
	f     func() (interface{}, error)
}

func NewTask(f func() (interface{}, error)) Task {
	return &TaskImpl{f: f}
}

func (t *TaskImpl) Run(wg *sync.WaitGroup) {
	t.Image, t.Err = t.f()
	wg.Done()
}
