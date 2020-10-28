# switch-to-go-m7

## Prerequisites
* Go 1.5
* Graphviz (https://graphviz.org/), for converting dot file to a svg/image/etc

## Usage

```
$ go build -o crawl ./cmd/crawl
$ ./crawl -depth=2 -out out.dot
$ dot -T svg out.svg out.dot
```