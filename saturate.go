package main

import (
	"flag"
	"fmt"
	"os/exec"
	"time"
)

type saturateOpts struct {
	specFile     string
	outputPrefix string
}

func saturateCmd() command {

	fs := flag.NewFlagSet("saturate", flag.ExitOnError)
	opts := &saturateOpts{}

	fs.StringVar(&opts.specFile, "specFile", "", "spec yaml for all endpoints to test")
	fs.StringVar(&opts.outputPrefix, "outputPrefix", "saturate", "output prefix to store test result")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return saturate(opts)
	}}
}

func saturate(opts *saturateOpts) error {
	spec, err := readSpec(opts.specFile)
	if err != nil {
		return err
	}

	for i, ep := range spec.Endpoints {
		saveTo := fmt.Sprintf(">%s_run_%d.txt", opts.outputPrefix, i)
		if err = exec.Command("wrk", "-t4", "-c100", "-d5M", "--script="+ep.Script, ep.Url, saveTo).Run(); err != nil {
			return err
		}

		// sleep 2 minute between each endpoint to give
		// already queued requests time to finish
		if i != len(spec.Endpoints)-1 {
			time.Sleep(2 * time.Minute)
		}
	}

	return nil
}
