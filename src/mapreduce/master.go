package mapreduce

import "container/list"
import "fmt"


type WorkerInfo struct {
	address string
	// You can add definitions here.
}


// Clean up all workers by sending a Shutdown RPC to each one of them Collect
// the number of jobs each work has performed.
func (mr *MapReduce) KillWorkers() *list.List {
	l := list.New()
	for _, w := range mr.Workers {
		DPrintf("DoWork: shutdown %s\n", w.address)
		args := &ShutdownArgs{}
		var reply ShutdownReply
		ok := call(w.address, "Worker.Shutdown", args, &reply)
		if ok == false {
			fmt.Printf("DoWork: RPC %s shutdown error\n", w.address)
		} else {
			l.PushBack(reply.Njobs)
		}
	}
	return l
}

func AssignJobsToWorkers(mr *MapReduce, doneChannel chan int, job JobType, nJobs int, nJobsOther int) {
	for i :=0; i < nJobs; i++ {
		go func(jobNum int) {
			worker := <-mr.registerChannel

			args := &DoJobArgs{mr.file, job, jobNum, nJobsOther}
			var reply DoJobReply

			call
		}(i)
	}
}

func (mr *MapReduce) RunMaster() *list.List {
	mapDoneChannel, reduceDoneChannel := make(chan int, mr.nMap), make(chan int, mr.nReduce)
	AssignJobsToWorkers(mr, mapDoneChannel, Map, mr.nMap, mr.nReduce)

	for i := 0; i < mr.nMap; i++ {
		<-mapDoneChannel
	}

	AssignJobsToWorkers(mr, reduceDoneChannel, Reduce, mr.nReduce, mr.nMap)

	for i := 0; i < mr.nMap; i++ {
		<-reduceDoneChannel
	}
	
	// Your code here
	return mr.KillWorkers()
}
