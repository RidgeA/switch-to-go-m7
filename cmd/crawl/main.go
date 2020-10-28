package main

import (
	"context"
	"github.com/RidgeA/switch-to-go-m7/internal/crawler"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"
)

var templ = template.Must(template.New("graph").Parse(`
digraph G {
	{{range .}}"{{.From}}" -> "{{.To}}"
	{{end}}
}
`))

func main() {
	o := parseOptions(os.Args[1:])

	file, err := os.Create(o.out)
	if err != nil {
		log.Fatal(err)
	}

	cr := crawler.NewCrawler("https://uk.wikipedia.org/wiki/%D0%93%D0%BE%D0%BB%D0%BE%D0%B2%D0%BD%D0%B0_%D1%81%D1%82%D0%BE%D1%80%D1%96%D0%BD%D0%BA%D0%B0",
		o.depth,
		crawler.WithLinkSkipper(skipUrls()),
		crawler.WithLinkSkipper(func(s string) bool { return !strings.HasPrefix(s, "https://uk.wikipedia.org/wiki/") }),
	)

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	connections := cr.Crawl(ctx)

	err = templ.Execute(file, connections)
	if err != nil {
		log.Fatal(err)
	}

}

func skipUrls() func(string) bool {
	basePath := "https://uk.wikipedia.org/wiki/"
	urls := []string{
		"Файл",
		"Портал",
		"Вікіпедія",
		"Шаблон",
		"Спеціальна",
		"Обговорення",
		"Категорія",
		"Special",
		"Довідка",
	}
	return func(s string) bool {
		for _, segment := range urls {
			if strings.HasPrefix(s, basePath+url.QueryEscape(segment)+":") {
				return true
			}
		}
		return false
	}
}
