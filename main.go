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
	AuthToken  string
	UrlsFile   string `yaml:"rand_urls_file"`
	BodiesFile string `yaml:"rand_bodies_file"`
	Endpoints  []Endpoint
}

type Endpoint struct {
	Name        string
	Method      string
	Url         string
	Script      string
	TargetsFile string `yaml:"targets_file"`
	RequestRate int    `yaml:"request_rate"`
}

type command struct {
	fs *flag.FlagSet
	fn func(args []string) error
}

func main() {

	commands := map[string]command{
		"randomize": randomizeCmd(),
		"saturate":  saturateCmd(),
		"sustain":   sustainCmd(),
	}

	fs := flag.NewFlagSet("loadtest", flag.ExitOnError)
	reps := fs.Int("reps", 5, "# of repititions to run a single load test")

	fs.Usage = func() {
		fmt.Println("Usage: loadtest [global flags] <command> [command flags]")
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

	if args[0] == "randomize" {
		*reps = 1
	}

	// repeat load testing commands multiple times
	// to make sure out test setup is not flawed
	// and we can use the sample to identify outliers
	for i := 0; i < *reps; i++ {
		if err := cmd.fn(args[1:]); err != nil {
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
