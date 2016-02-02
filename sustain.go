package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type sustainOpts struct {
	outputPrefix string
}

func sustainCmd() command {

	fs := flag.NewFlagSet("sustain", flag.ExitOnError)
	opts := &sustainOpts{}

	fs.StringVar(&opts.outputPrefix, "outputPrefix", "sustain", "output prefix to store test result")

	return command{fs, func(spec Spec, args []string) error {
		fs.Parse(args)
		return sustain(spec, opts)
	}}
}

func sustain(spec Spec, opts *sustainOpts) error {
	now := time.Now().Unix()
	for _, ep := range spec.Endpoints {
		var (
			output []byte
			err    error
		)
		rate := fmt.Sprintf("--rate=%d", ep.RequestRate)
		script := fmt.Sprintf("--script=%s", ep.Script)

		cmd := exec.Command("wrk2", "-t4", "-c100", "-d5m", rate, script, "--latency", ep.Url)
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

	for i, ep := range spec.Endpoints {
		targets := fmt.Sprintf("-targets=%s", ep.TargetsFile)
		header := fmt.Sprintf("-header=X-Auth-Token: %s", spec.AuthToken)
		rate := fmt.Sprintf("-rate=%d", ep.RequestRate)
		workers := fmt.Sprintf("-workers=%d", ep.RequestRate)
		output := fmt.Sprintf("-output=%s_%s_@%d_%d.bin", opts.outputPrefix, ep.Name, ep.RequestRate, now)

		cmd := exec.Command("vegeta", "attack", targets, header, "-duration=5m", rate, workers, output)
		fmt.Printf("running %v \n", cmd)
		if err := cmd.Run(); err != nil {
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
