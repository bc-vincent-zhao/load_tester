package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Spec struct {
	BaseUrl     string `yaml:"base_url"`
	AuthToken   string `yaml:"auth_token"`
	RandCount   int    `yaml:"rand_count"`
	MaxBodySize int    `yaml:"max_body_size"`
	UrlsFile    string `yaml:"rand_urls_file"`
	BodiesFile  string `yaml:"rand_bodies_file"`
	Duration    string `yaml:"duration"`
	Endpoints   []Endpoint
}

type Endpoint struct {
	Name         string
	Method       string
	Url          string
	Script       string
	TargetsFile  string `yaml:"targets_file"`
	RequestRate  int    `yaml:"request_rate"`
	CustomHeader string `yaml:"custom_header"`
}

type command struct {
	fs *flag.FlagSet
	fn func(spec Spec, args []string) error
}

func main() {

	commands := map[string]command{
		"saturate":  saturateCmd(),
		"histogram": histogramCmd(),
		"ts":        tsCmd(),
	}

	fs := flag.NewFlagSet("load_test", flag.ExitOnError)
	reps := fs.Int("reps", 1, "# of repititions to run a single load test")
	specFile := fs.String("specFile", "", "spec yaml for all endpoints to test")
	randGen := fs.Bool("generateRandom", true, "generate new randomized url and request body for each repitition")

	fs.Usage = func() {
		fmt.Println("Usage: load_test [global flags] <command> [command flags]")
		fmt.Printf("\n global flags:\n")
		fs.PrintDefaults()
		for name, cmd := range commands {
			fmt.Printf("\n%s command:\n", name)
			cmd.fs.PrintDefaults()
		}
	}

	fs.Parse(os.Args[1:])

	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	var cmd command
	var ok bool
	if cmd, ok = commands[args[0]]; !ok {
		log.Fatalf("Unknown command %s", args[0])
	}

	// read in the test spec
	spec, err := readSpec(*specFile)
	if err != nil {
		log.Fatal(err)
	}

	// repeat load testing commands multiple times
	// to make sure out test setup is not flawed
	// and we can use the sample to identify outliers
	for i := 0; i < *reps; i++ {

		if *randGen {
			// generate randomized urls and data for each test iteration
			// to make sure we hit longest execution path instead of some cache
			if err := randomize(spec); err != nil {
				log.Fatal(err)
			}
		}

		// run the actual test command
		if err := cmd.fn(spec, args[1:]); err != nil {
			log.Fatal(err)
		}
		// sleep 2 minute between each test run to give
		// already queued requests time to finish
		if i != *reps-1 {
			time.Sleep(2 * time.Minute)
		}
	}
}

func readSpec(file string) (Spec, error) {
	var spec Spec

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return spec, err
	}

	err = yaml.Unmarshal(data, &spec)
	return spec, err
}
