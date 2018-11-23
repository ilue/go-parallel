package parallel

import (
	"runtime"
	"sync"
)

type parallelForJob struct {
	wg  *sync.WaitGroup
	f   func(int, int)
	arg int
}

func (self *parallelForJob) Run(workerId int) {
	defer self.wg.Done()
	self.f(self.arg, workerId)
}

var (
	numWorkers = runtime.NumCPU()
	jobCh = make(chan *parallelForJob)
)

func init() {
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			for job := range jobCh {
				job.Run(id)
			}
		}(i)
	}
}

func NumWorkers() int {
	return numWorkers
}

func ParallelFor(first, last int, f func(int, int)) {
	var wg sync.WaitGroup
	wg.Add(last - first)

	for i := first; i < last; i++ {
		jobCh <- &parallelForJob{&wg, f, i}
	}

	wg.Wait()
}
