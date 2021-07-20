package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
)

var (
	taskDone     = "DONE"
	taskAlloc    = "ALLOCATED"
	taskNotAlloc = "NOTALLOCATED"
)

type Coordinator struct {
	// Your definitions here.
	R          int //nReduce value
	mapTask    []string
	reduceTask []string
	mapDone    bool
	reduceDone bool
	files      []string
	mu         sync.Mutex
}

// Your code here -- RPC handlers for the worker to call.
func (c *Coordinator) Task(args *RequestTaskArgs, reply *RequestTaskReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	noTask := true
	if !c.mapDone {
		for i := range c.mapTask {
			if c.mapTask[i] == taskNotAlloc {
				reply.TaskType = MapTaskType
				reply.TaskArg = append(reply.TaskArg, c.files[i])
				c.mapTask[i] = taskAlloc
				noTask = false
				//TODO:启动定时器
				break
			}
		}
		if noTask {
			reply.TaskType = 0
		}
		return nil
	}
	if !c.reduceDone {
		for i := range c.reduceTask {
			if c.reduceTask[i] == taskNotAlloc {
				reply.TaskType = ReduceTaskType
				reply.TaskArg = append(reply.TaskArg, strconv.Itoa(i))
				c.reduceTask[i] = taskAlloc
				//TODO:启动定时器
				break
			}
		}
		if noTask {
			reply.TaskType = 0
		}
		return nil
	}
	//all tasks finished
	reply.ShouldExit = true
	return nil
}

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{
		R:          nReduce,
		mapTask:    make([]bool, nReduce),
		reduceTask: make([]bool, nReduce),
		mapDone:    false,
		reduceDone: false,
		files:      files,
	}

	// Your code here.

	c.server()
	return &c
}
