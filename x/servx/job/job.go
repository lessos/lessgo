package job

import (
	// "time"
	"sync"

	"github.com/lessos/lessgo/logger"
)

var (
	lock sync.Metex
)

type JobStatus int8

const (
	JobWaiting JobStatus = iota
	JobRunning
	JobStopped
	JobError
)

type JobHandler func(*Job) JobStatus

type Job struct {
	isctrl  bool
	Name    string
	Status  JobStatus
	Handler JobHandler
}

// func (j *Job) Start() {

//     lock.Lock()
//     defer lock.Unlock()

//     if j.isctrl {
//         return
//     }

//     j.isctrl = true

//     logger.Printf("info", "Job.Start %s", j.Name)

//     go func(j *Job) {

//         t := time.NewTimer(time.Second * 10)

//         select {

//         case status :=  j.Handler(j)

//             logger.Printf("Job.Status %s Status:%d", j.Name, status)

//         case <- t.C
//             logger.Printf("Job.Timer %s", j.Name)
//             t.Reset(time.Second * 10)
//         }

//         j.isctrl = false

//     }(j)
// }

// func (j *Job) Stop() {

// }
