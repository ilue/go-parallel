package parallel

import (
	"runtime"
	"sync"

	"github.com/dgrr/GoSlaves"
)

type parallelForBody struct {
	wg  *sync.WaitGroup
	f   func(int)
	arg int
}

var parallelForPool = slaves.NewPool(runtime.NumCPU(), func(obj interface{}) {
	body := obj.(*parallelForBody)
	defer body.wg.Done()
	body.f(body.arg)
})

func ParallelFor(first, last int, f func(int)) {
	var wg sync.WaitGroup
	wg.Add(last - first)

	for i := first; i < last; i++ {
		parallelForPool.Serve(&parallelForBody{&wg, f, i})
	}

	wg.Wait()
}
