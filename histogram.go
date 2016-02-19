package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type histogramOpts struct {
	outputPrefix string
}

func histogramCmd() command {

	fs := flag.NewFlagSet("histogram", flag.ExitOnError)
	opts := &histogramOpts{}

	fs.StringVar(&opts.outputPrefix, "outputPrefix", "hdr", "output prefix to store test result")

	return command{fs, func(spec Spec, args []string) error {
		fs.Parse(args)
		return histogram(spec, opts)
	}}
}

func histogram(spec Spec, opts *histogramOpts) error {
	now := time.Now().Unix()
	for _, ep := range spec.Endpoints {
		var (
			output []byte
			err    error
		)
		args := []string{
			fmt.Sprintf("-t%d", spec.Workers),
			fmt.Sprintf("-c%d", spec.Connections),
			fmt.Sprintf("--rate=%d", ep.RequestRate),
			fmt.Sprintf("--script=%s", ep.Script),
			fmt.Sprintf("-d%s", spec.Duration),
			"--latency",
			ep.Url,
		}

		cmd := exec.Command("wrk2", args...)
		fmt.Printf("running %v \n", cmd)
		if output, err = cmd.Output(); err != nil {
			return err
		}

		file, err := os.Create(fmt.Sprintf("%s_%s_@%d_%d.hdrm", opts.outputPrefix, ep.Name, ep.RequestRate, now))
		if err != nil {
			return err
		}

		_, err = file.Write(output)
		if err != nil {
			return err
		}

		// sleep a minute between each endpoint to give
		// already queued requests time to finish
		time.Sleep(time.Minute)
	}
	return nil
}
