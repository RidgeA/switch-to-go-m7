package crawler

import "sync"

func merge(cs ...<-chan FetchedPage) <-chan FetchedPage {
	var wg sync.WaitGroup
	out := make(chan FetchedPage)

	output := func(c <-chan FetchedPage) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
