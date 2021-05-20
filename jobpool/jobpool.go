package jobpool

import (
	"container/list"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
)

//** TYPES

type (
	// queueJob is a control structure for queuing jobs.
	queueJob struct {
		Jobber                   // The object to execute the job routine against.
		priority      bool       // If the job needs to be placed on the priority queue.
		resultChannel chan error // Used to inform the queue operation is complete.
	}

	// dequeueJob is a control structure for dequeue jobs.
	dequeueJob struct {
		ResultChannel chan *queueJob // Used to return the queued job to be processed.
	}

	// JobPool maintains queues and Go routines for processing jobs.
	JobPool struct {
		priorityJobQueue     *list.List       // The priority job queue.
		normalJobQueue       *list.List       // The normal job queue.
		queueChannel         chan *queueJob   // Channel allows the thread safe placement of jobs into the queue.
		dequeueChannel       chan *dequeueJob // Channel allows the thread safe removal of jobs from the queue.
		shutdownQueueChannel chan string      // Channel used to shutdown the queue routine.
		jobChannel           chan string      // Channel to signal to a job routine to process a job.
		shutdownJobChannel   chan struct{}    // Channel used to shutdown the job routines.
		shutdownWaitGroup    sync.WaitGroup   // The WaitGroup for shutting down existing routines.
		queuedJobs           int32            // The number of pending jobs in queued.
		activeRoutines       int32            // The number of routines active.
		queueCapacity        int32            // The max number of jobs we can store in the queue.
	}
)

//** INTERFACES

// Jobber is an interface that is implemented to run jobs.
type Jobber interface {
	RunJob(jobRoutine int)
}

//** INIT FUNCTION

// init is called when the system is inited.
func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

//** PUBLIC FUNCTIONS

// New creates a new JobPool.
func New(numberOfRoutines int, queueCapacity int32) (jobPool *JobPool) {
	// Create the job queue.
	jobPool = &JobPool{
		priorityJobQueue:     list.New(),
		normalJobQueue:       list.New(),
		queueChannel:         make(chan *queueJob),
		dequeueChannel:       make(chan *dequeueJob),
		shutdownQueueChannel: make(chan string),
		jobChannel:           make(chan string, queueCapacity),
		shutdownJobChannel:   make(chan struct{}),
		queuedJobs:           0,
		activeRoutines:       0,
		queueCapacity:        queueCapacity,
	}

	// Launch the job routines to process work.
	for jobRoutine := 0; jobRoutine < numberOfRoutines; jobRoutine++ {
		// Add the routine to the wait group.
		jobPool.shutdownWaitGroup.Add(1)

		// Start the job routine.
		go jobPool.jobRoutine(jobRoutine)
	}

	// Start the queue routine to capture and provide jobs.
	go jobPool.queueRoutine()

	return jobPool
}

//** PUBLIC MEMBER FUNCTIONS

// Shutdown will release resources and shutdown all processing.
func (jobPool *JobPool) Shutdown(goRoutine string) (err error) {
	defer catchPanic(&err, goRoutine, "Shutdown")

	writeStdout(goRoutine, "Shutdown", "Started")
	writeStdout(goRoutine, "Shutdown", "Queue Routine")

	jobPool.shutdownQueueChannel <- "Shutdown"
	<-jobPool.shutdownQueueChannel

	close(jobPool.shutdownQueueChannel)
	close(jobPool.queueChannel)
	close(jobPool.dequeueChannel)

	writeStdout(goRoutine, "Shutdown", "Shutting Down Job Routines")

	// Close the channel to shut things down
	close(jobPool.shutdownJobChannel)
	jobPool.shutdownWaitGroup.Wait()

	close(jobPool.jobChannel)

	writeStdout(goRoutine, "Shutdown", "Completed")
	return err
}

// QueueJob queues a job to be processed.
func (jobPool *JobPool) QueueJob(goRoutine string, jober Jobber, priority bool) (err error) {
	defer catchPanic(&err, goRoutine, "QueueJob")

	// Create the job object to queue.
	job := queueJob{
		jober,            // Jobber Interface.
		priority,         // Priority.
		make(chan error), // Result Channel.
	}

	defer close(job.resultChannel)

	// Queue the job
	jobPool.queueChannel <- &job
	err = <-job.resultChannel

	return err
}

// QueuedJobs will return the number of jobs items in queue.
func (jobPool *JobPool) QueuedJobs() int32 {
	return atomic.AddInt32(&jobPool.queuedJobs, 0)
}

// ActiveRoutines will return the number of routines performing work.
func (jobPool *JobPool) ActiveRoutines() int32 {
	return atomic.AddInt32(&jobPool.activeRoutines, 0)
}

//** PRIVATE FUNCTIONS

// catchPanic is used to catch any Panic and log exceptions to Stdout. It will also write the stack trace
//  err: A reference to the err variable to be returned to the caller. Can be nil.
func catchPanic(err *error, goRoutine string, functionName string) {
	if r := recover(); r != nil {
		// Capture the stack trace.
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		writeStdoutf(goRoutine, functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

// writeStdout is used to write a system message directly to stdout.
func writeStdout(goRoutine string, functionName string, message string) {
	log.Printf("%s : %s : %s\n", goRoutine, functionName, message)
}

// writeStdoutf is used to write a formatted system message directly stdout.
func writeStdoutf(goRoutine string, functionName string, format string, a ...interface{}) {
	writeStdout(goRoutine, functionName, fmt.Sprintf(format, a...))
}

//** PRIVATE MEMBER FUNCTIONS

// queueRoutine performs the thread safe queue related processing.
func (jobPool *JobPool) queueRoutine() {
	for {
		select {
		case <-jobPool.shutdownQueueChannel:
			writeStdout("Queue", "queueRoutine", "Going Down")
			jobPool.shutdownQueueChannel <- "Down"
			return

		case queueJob := <-jobPool.queueChannel:
			// Enqueue the job
			jobPool.queueRoutineEnqueue(queueJob)
			break

		case dequeueJob := <-jobPool.dequeueChannel:
			// Dequeue a job
			jobPool.queueRoutineDequeue(dequeueJob)
			break
		}
	}
}

// queueRoutineEnqueue places a job on either the normal or priority queue.
func (jobPool *JobPool) queueRoutineEnqueue(queueJob *queueJob) {
	defer catchPanic(nil, "Queue", "queueRoutineEnqueue")

	// If the queue is at capacity don't add it.
	if atomic.AddInt32(&jobPool.queuedJobs, 0) == jobPool.queueCapacity {
		queueJob.resultChannel <- fmt.Errorf("Job Pool At Capacity")
		return
	}

	if queueJob.priority == true {
		jobPool.priorityJobQueue.PushBack(queueJob)
	} else {
		jobPool.normalJobQueue.PushBack(queueJob)
	}

	// Increment the queued work count.
	atomic.AddInt32(&jobPool.queuedJobs, 1)

	// Tell the caller the work is queued.
	queueJob.resultChannel <- nil

	// Tell the job routine to wake up.
	jobPool.jobChannel <- "Wake Up"
}

// queueRoutineDequeue remove a job from the queue.
func (jobPool *JobPool) queueRoutineDequeue(dequeueJob *dequeueJob) {
	defer catchPanic(nil, "Queue", "queueRoutineDequeue")

	var nextJob *list.Element

	if jobPool.priorityJobQueue.Len() > 0 {
		nextJob = jobPool.priorityJobQueue.Front()
		jobPool.priorityJobQueue.Remove(nextJob)
	} else {
		nextJob = jobPool.normalJobQueue.Front()
		jobPool.normalJobQueue.Remove(nextJob)
	}

	// Decrement the queued work count.
	atomic.AddInt32(&jobPool.queuedJobs, -1)

	// Cast the list element back to a Job.
	job := nextJob.Value.(*queueJob)

	// Give the caller the work to process.
	dequeueJob.ResultChannel <- job
}

// jobRoutine performs the actual processing of jobs.
func (jobPool *JobPool) jobRoutine(jobRoutine int) {
	for {
		select {
		// Shutdown the job routine.
		case <-jobPool.shutdownJobChannel:
			writeStdout(fmt.Sprintf("JobRoutine %d", jobRoutine), "jobRoutine", "Going Down")
			jobPool.shutdownWaitGroup.Done()
			return

		// Perform the work.
		case <-jobPool.jobChannel:
			jobPool.doJobSafely(jobRoutine)
			break
		}
	}
}

// dequeueJob pulls a job from the queue.
func (jobPool *JobPool) dequeueJob() (job *queueJob, err error) {
	defer catchPanic(&err, "jobRoutine", "dequeueJob")

	// Create the job object to queue.
	requestJob := dequeueJob{
		ResultChannel: make(chan *queueJob), // Result Channel.
	}

	defer close(requestJob.ResultChannel)

	// Dequeue the job
	jobPool.dequeueChannel <- &requestJob
	job = <-requestJob.ResultChannel

	return job, err
}

// doJobSafely will executes the job within a safe context.
func (jobPool *JobPool) doJobSafely(jobRoutine int) {
	defer catchPanic(nil, "jobRoutine", "doJobSafely")
	defer atomic.AddInt32(&jobPool.activeRoutines, -1)

	// Update the active routine count.
	atomic.AddInt32(&jobPool.activeRoutines, 1)

	// Dequeue a job
	queueJob, err := jobPool.dequeueJob()
	if err != nil {
		writeStdoutf("Queue", "jobpool.JobPool", "doJobSafely", "ERROR : %s", err)
		return
	}

	// Perform the job.
	queueJob.RunJob(jobRoutine)
}
