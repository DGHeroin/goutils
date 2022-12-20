package rorate

import (
    "sync"
    "sync/atomic"
    "time"
)

type (
    RotateObject struct {
        mu       sync.Mutex
        buf      []interface{}
        MaxSize  int
        onRotate func(arr []interface{}, t0, t1 time.Time, num int)
        onAdd    func(interface{}) bool
        t0       time.Time
        c        int32
        wg       *sync.WaitGroup
    }
)

func NewRotateObject() *RotateObject {
    return &RotateObject{
        MaxSize: 1000,
        wg:      &sync.WaitGroup{},
    }
}
func (b *RotateObject) Add(p interface{}) {
    var isOverLimit bool

    b.mu.Lock()
    isOverLimit = len(b.buf) >= b.MaxSize
    b.mu.Unlock()

    if isOverLimit {
        b.rotate()
    }
    if b.onAdd != nil {
        if !b.onAdd(p) {
            return
        }
    }

    num := atomic.AddInt32(&b.c, 1)
    if num == 1 {
        b.t0 = time.Now()
    }

    b.mu.Lock()
    defer b.mu.Unlock()
    b.buf = append(b.buf, p)
}
func (b *RotateObject) rotate() {
    b.mu.Lock()
    defer b.mu.Unlock()
    if b.onRotate != nil && len(b.buf) > 0 {
        var slice = b.buf
        b.wg.Add(1)
        num := len(slice)
        t0 := b.t0
        t1 := time.Now()
        go func() {
            defer b.wg.Done()
            b.onRotate(slice, t0, t1, num)
        }()
    }
    b.buf = []interface{}{}
    atomic.StoreInt32(&b.c, 0)
}
func (b *RotateObject) Flush() {
    b.rotate()
}
func (b *RotateObject) Close() {
    b.rotate()
    b.wg.Wait()
}
func (b *RotateObject) Current() []interface{} {
    b.mu.Lock()
    defer b.mu.Unlock()
    return b.buf
}
func (b *RotateObject) OnRotate(fn func(object []interface{}, t0, t1 time.Time, num int)) {
    b.onRotate = fn
}
func (b *RotateObject) OnWrite(fn func(interface{}) bool) {
    b.onAdd = fn
}
