package worker

import "sync"

type Pool struct {
	Tasks []TaskImpl

	concurrency int
	tasksChan   chan TaskImpl
	wg          sync.WaitGroup
}

func NewPool(tasks []TaskImpl, concurrency int) *Pool {
	return &Pool{
		Tasks:       tasks,
		concurrency: concurrency,
		tasksChan:   make(chan TaskImpl),
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}

	p.wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.tasksChan <- task
	}

	close(p.tasksChan)

	p.wg.Wait()
}

func (p *Pool) work() {
	for task := range p.tasksChan {
		task.Run(&p.wg)
	}
}
