package crawler

import (
	"context"
	"github.com/RidgeA/switch-to-go-m7/internal/fetcher"
	"runtime"
)

type (
	OptionsFunc func(c *Crawler)

	// Filter function called for each url.
	// If the function returns true - the url won't be processed
	URLSkipFilter func(string) bool

	Crawler struct {
		start   URLToFetch
		filters []URLSkipFilter
		workers int
	}

	URLToFetch struct {
		Url   string
		Depth int
	}

	FetchedPage struct {
		URLToFetch
		Title string
		Links []string
	}

	Connection struct {
		From, To string
	}
)

func NewCrawler(start string, depth int, ops ...OptionsFunc) *Crawler {
	c := &Crawler{
		start: URLToFetch{
			Url:   start,
			Depth: depth,
		},
		workers: runtime.GOMAXPROCS(0),
	}

	for _, opt := range ops {
		opt(c)
	}

	return c
}

func (c *Crawler) Crawl(ctx context.Context) []Connection {

	urlsCh := make(chan URLToFetch)

	go func() {
		urlsCh <- c.start
	}()

	filteredCh := c.skip(ctx, dedup(ctx, urlsCh))

	chs := make([]<-chan FetchedPage, 0, c.workers)
	for i := 0; i < c.workers; i++ {
		chs = append(chs, c.loadPage(ctx, filteredCh))
	}
	pagesCh := merge(chs...)

	connectedLinks := make([]Connection, 0)
	titles := make(map[string]string)

	for p := range pagesCh {

		uniqueLinks := unique(p.Links)

		titles[p.Url] = p.Title
		for _, l := range uniqueLinks {
			connectedLinks = append(connectedLinks, Connection{From: p.Url, To: l})
		}

		if p.Depth >= 0 {
			go func(links []string, depth int) {
				for _, link := range links {
					urlsCh <- URLToFetch{
						Url:   link,
						Depth: depth,
					}
				}
			}(uniqueLinks, p.Depth-1)
		}
	}

	result := make([]Connection, 0)
	for _, c := range connectedLinks {
		fromTitle, exists := titles[c.From]
		if !exists {
			continue
		}

		toTitle, exists := titles[c.To]
		if !exists {
			continue
		}

		result = append(result, Connection{From: fromTitle, To: toTitle})
	}

	return result
}

func (c *Crawler) loadPage(ctx context.Context, in <-chan URLToFetch) <-chan FetchedPage {

	out := make(chan FetchedPage)

	go func() {
		defer close(out)
		for {
			select {
			case url := <-in:
				if url.Depth < 0 {
					return
				}
				title, links := fetcher.Fetch(ctx, url.Url)

				out <- FetchedPage{
					URLToFetch: url,
					Title:      title,
					Links:      links,
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func (c *Crawler) shouldSkip(url string) bool {
	for _, f := range c.filters {
		if f(url) {
			return true
		}
	}
	return false
}

func (c *Crawler) skip(ctx context.Context, input <-chan URLToFetch) <-chan URLToFetch {
	out := make(chan URLToFetch)

	go func() {
		defer close(out)
		for {
			select {
			case url := <-input:
				if !c.shouldSkip(url.Url) {
					out <- url
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func WithLinkSkipper(f ...URLSkipFilter) OptionsFunc {
	return func(c *Crawler) {
		c.filters = append(c.filters, f...)
	}
}

