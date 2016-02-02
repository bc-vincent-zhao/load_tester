package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type saturateOpts struct {
	outputPrefix string
}

func saturateCmd() command {

	fs := flag.NewFlagSet("saturate", flag.ExitOnError)
	opts := &saturateOpts{}

	fs.StringVar(&opts.outputPrefix, "outputPrefix", "saturate", "output prefix to store test result")

	return command{fs, func(spec Spec, args []string) error {
		fs.Parse(args)
		return saturate(spec, opts)
	}}
}

func saturate(spec Spec, opts *saturateOpts) error {
	now := time.Now().Unix()
	for i, ep := range spec.Endpoints {
		var (
			output []byte
			err    error
		)
		cmd := exec.Command("wrk", "-t4", "-c100", "-d5m", "--script="+ep.Script, ep.Url)
		fmt.Printf("running %v \n", cmd)
		if output, err = cmd.Output(); err != nil {
			return err
		}

		file, err := os.Create(fmt.Sprintf("%s_%s_%d.txt", opts.outputPrefix, ep.Name, now))
		if err != nil {
			return err
		}

		_, err = file.Write(output)
		if err != nil {
			return err
		}

		// sleep a minute between each endpoint to give
		// already queued requests time to finish
		if i != len(spec.Endpoints)-1 {
			time.Sleep(time.Minute)
		}
	}

	return nil
}
