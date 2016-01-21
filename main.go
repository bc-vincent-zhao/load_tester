package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	commands := map[string]command{
		"randomize": randomizeCmd(),
	}

	fs := flag.NewFlagSet("loadtest", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage: loadtest <command> [command flags]")
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

	if cmd, ok := commands[args[0]]; !ok {
		log.Fatalf("Unknown command %s", args[0])
	} else if err := cmd.fn(args[1:]); err != nil {
		log.Fatal(err)
	}
}

type command struct {
	fs *flag.FlagSet
	fn func(args []string) error
}
