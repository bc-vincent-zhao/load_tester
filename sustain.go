package main

import "flag"

type sustainOpts struct {
	specFile     string
	outputPrefix string
}

func sustainCmd() command {

	fs := flag.NewFlagSet("sustain", flag.ExitOnError)
	opts := &sustainOpts{}

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return sustain(opts)
	}}
}

func sustain(opts *sustainOpts) error {
	return nil
}
