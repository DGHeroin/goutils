package main

import (
    "fmt"
    "goutils/task"
    "sync/atomic"
    "time"
)

func main() {
    t := task.New(30)
    cc := t.AddWorker()
    time.AfterFunc(time.Nanosecond, func() {
        close(cc)
    })
    var n int32
    for i := 0; i < 1000; i++ {
        t.Run(func(params ...interface{}) {
            atomic.AddInt32(&n, 1)
            // log.Println(params...)
        }, i)
    }

    t.Close()
    fmt.Println("total run:", n)
}
