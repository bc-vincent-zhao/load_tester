package main

import (
	"flag"
	"fmt"
	"os/exec"
	"time"
)

type tsOpts struct {
	outputPrefix string
}

func tsCmd() command {

	fs := flag.NewFlagSet("ts", flag.ExitOnError)
	opts := &tsOpts{}

	fs.StringVar(&opts.outputPrefix, "outputPrefix", "ts", "output prefix to store test result")

	return command{fs, func(spec Spec, args []string) error {
		fs.Parse(args)
		return ts(spec, opts)
	}}
}

func ts(spec Spec, opts *tsOpts) error {
	now := time.Now().Unix()
	for i, ep := range spec.Endpoints {
		args := []string{
			"attack",
			fmt.Sprintf("-targets=%s", ep.TargetsFile),
			fmt.Sprintf("-duration=%s", spec.Duration),
			fmt.Sprintf("-rate=%d", ep.RequestRate),
			fmt.Sprintf("-workers=%d", ep.RequestRate),
			fmt.Sprintf("-output=%s_%s_@%d_%d.bin", opts.outputPrefix, ep.Name, ep.RequestRate, now),
			fmt.Sprintf("-header=X-Auth-Token:%s", spec.AuthToken),
		}
		if ep.CustomHeader != "" {
			args = append(args, fmt.Sprintf("-header=%s", ep.CustomHeader))
		}

		cmd := exec.Command("vegeta", args...)
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
