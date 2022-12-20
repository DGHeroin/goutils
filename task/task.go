package task

import (
    "sync"
    "sync/atomic"
)

type (
    Task struct {
        Run    func(v ...interface{})
        Params []interface{}
    }
    Pool struct {
        wg  sync.WaitGroup
        num int64 // workers num
        ch  chan *Task
        sync.Mutex
        isClose bool
    }
)

func New(bufferSize uint64) *Pool {
    p := &Pool{
        ch: make(chan *Task, bufferSize),
    }
    go p.run()
    return p
}
func (p *Pool) run() (closeCh chan struct{}) {
    atomic.AddInt64(&p.num, 1)
    closeCh = make(chan struct{})
    go func() {
        defer func() {
            atomic.AddInt64(&p.num, -1)
        }()
        for {
            select {
            case <-closeCh:
                return
            case task := <-p.ch:
                p.fire(task)
            }
        }
    }()
    return
}
func (p *Pool) fire(task *Task) {
    if task == nil {
        return
    }
    defer func() {
        p.wg.Done()
        recover()
    }()
    task.Run(task.Params...)
}
func (p *Pool) commitTask(task *Task) {
    p.Lock()
    defer p.Unlock()
    if task == nil || p.isClose {
        return
    }
    p.wg.Add(1)
    p.ch <- task
}
func (p *Pool) Run(fn func(params ...interface{}), params ...interface{}) {
    task := &Task{
        Run:    fn,
        Params: params,
    }
    p.commitTask(task)
}
func (p *Pool) Close() {
    p.isClose = true
    p.wg.Wait()
    close(p.ch)
}
func (p *Pool) AddWorker() chan struct{} {
    return p.run()
}
