package main

import (
	"flag"
	"log"
)

type (
	options struct {
		out   string
		depth int
	}
)

func parseOptions(args []string) options {
	o := options{}
	fs := flag.NewFlagSet("", flag.ExitOnError)

	fs.StringVar(&o.out, "out", "", "output file name, required")
	fs.IntVar(&o.depth, "depth", 1, "lookup depth, default - 1")

	_ = fs.Parse(args)

	if o.out == "" {
		log.Fatal("output file name required")
	}

	return o
}
