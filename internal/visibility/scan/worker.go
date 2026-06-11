package scan

import "sync"

// Worker defines a concurrent worker for grid scanning
type Worker struct {
	ID int
}

func NewWorker(id int) *Worker {
	return &Worker{ID: id}
}

func (w *Worker) ProcessGrid(ch <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for range ch {
		// Process grid point
	}
}
