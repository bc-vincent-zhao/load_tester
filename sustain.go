package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type sustainOpts struct {
	specFile     string
	outputPrefix string
}

func sustainCmd() command {

	fs := flag.NewFlagSet("sustain", flag.ExitOnError)
	opts := &sustainOpts{}

	fs.StringVar(&opts.specFile, "specFile", "", "spec yaml for all endpoints to test")
	fs.StringVar(&opts.outputPrefix, "outputPrefix", "sustain", "output prefix to store test result")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return sustain(opts)
	}}
}

func sustain(opts *sustainOpts) error {
	spec, err := readSpec(opts.specFile)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	//for _, ep := range spec.Endpoints {
	//var output []byte
	//rate := fmt.Sprintf("--rate=%d", ep.RequestRate)
	//script := fmt.Sprintf("--script=%s", ep.Script)

	//if output, err = exec.Command("wrk2", "-t4", "-c100", "-d3m", rate, script, "--latency", ep.Url).Output(); err != nil {
	//return err
	//}

	//file, err := os.Create(fmt.Sprintf("%s_%s_@%d_%d.hdrm", opts.outputPrefix, ep.Name, ep.RequestRate, now))
	//if err != nil {
	//return err
	//}

	//_, err = file.Write(output)
	//if err != nil {
	//return err
	//}

	//// sleep a minute between each endpoint to give
	//// already queued requests time to finish
	//time.Sleep(time.Minute)
	//}

	for i, ep := range spec.Endpoints {

		// create targets file from random urls and request bodies
		if err = writeTargetsFile(ep.Method, ep.UrlsFile, ep.BodiesFile, ep.TargetsFile); err != nil {
			return err
		}

		targets := fmt.Sprintf("-targets=%s", ep.TargetsFile)
		header := fmt.Sprintf("-header=X-Auth-Token: %s", spec.AuthToken)
		rate := fmt.Sprintf("-rate=%d", ep.RequestRate)
		workers := fmt.Sprintf("-workers=%d", ep.RequestRate)
		output := fmt.Sprintf("-output=%s_%s_@%d_%d.bin", opts.outputPrefix, ep.Name, ep.RequestRate, now)

		cmd := exec.Command("vegeta", "attack", targets, header, "-duration=3m", rate, workers, output)
		if err = cmd.Run(); err != nil {
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

func writeTargetsFile(method, urls, bodies, targets string) error {
	urlsFile, err := os.Open(urls)
	if err != nil {
		return err
	}
	defer urlsFile.Close()
	urlReader := bufio.NewScanner(urlsFile)

	targetsFile, err := os.Create(targets)
	if err != nil {
		return err
	}
	defer targetsFile.Close()

	if method == "PUT" || method == "POST" {
		bodiesFile, err := os.Open(bodies)
		if err != nil {
			return err
		}
		defer bodiesFile.Close()
		bodyReader := bufio.NewScanner(bodiesFile)

		for urlReader.Scan() && bodyReader.Scan() {
			// create a temp file for store body because vegeta
			// requires body to be in file instead of raw string
			// we don't delete tmp files for now because they can
			// be handy for debugging or re-run commands manually
			tmp, err := ioutil.TempFile("", "load_test")
			if err != nil {
				return err
			}
			defer tmp.Close()

			if _, err = tmp.Write(bodyReader.Bytes()); err != nil {
				return err
			}

			if _, err = targetsFile.WriteString(fmt.Sprintf("%s %s\n", method, urlReader.Text())); err != nil {
				return err
			}
			if _, err = targetsFile.WriteString(fmt.Sprintf("@%s\n", tmp.Name())); err != nil {
				return err
			}
		}
	} else if method == "GET" || method == "HEAD" || method == "DELETE" {
		for urlReader.Scan() {
			if _, err = targetsFile.WriteString(urlReader.Text()); err != nil {
				return err
			}
		}
	}

	return nil
}
