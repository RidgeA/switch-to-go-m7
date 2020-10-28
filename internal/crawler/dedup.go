package crawler

import "context"

func dedup(ctx context.Context, input <-chan URLToFetch) <-chan URLToFetch {
	out := make(chan URLToFetch)
	visited := newSet()

	go func() {
		defer close(out)
		for {
			select {
			case url := <-input:
				if !visited.Has(url.Url) {
					visited.Add(url.Url)
					out <- url
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
