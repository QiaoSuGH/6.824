package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

var CTask = "Coordinator.Task"

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}

const (
	MapTaskType    = 1
	ReduceTaskType = 2
)

type RequestTaskArgs struct {
	TaskType int
}

type RequestTaskReply struct {
	ShouldExit bool //if all tasks finished, worker shoud exits
	TaskType   int  //if == 0 -- not allocated
	TaskArg    []string
}
