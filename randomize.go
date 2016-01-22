package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	mrand "math/rand"
	"os"
	"time"
)

type randomizeOpts struct {
	count        int
	baseUrl      string
	outputPrefix string
	maxBodySize  int
}

func randomizeCmd() command {
	fs := flag.NewFlagSet("randomize", flag.ExitOnError)
	opts := &randomizeOpts{}

	fs.StringVar(&opts.baseUrl, "baseUrl", "", "base url for target")
	fs.StringVar(&opts.outputPrefix, "outputPrefix", "", "output prefix to store randomized urls")
	fs.IntVar(&opts.count, "count", 100000, "number of randomised urls to generate")
	fs.IntVar(&opts.maxBodySize, "maxBodySize", 100000, "max PUT/POST request body size in bytes")

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return randomize(opts)
	}}
}

func randomize(opts *randomizeOpts) error {
	urlsFile := fmt.Sprintf("%s.urls", opts.outputPrefix)
	bodiesFile := fmt.Sprintf("%s.bodies", opts.outputPrefix)

	err := generateRandUrlFile(urlsFile, opts.baseUrl, opts.count)
	if err != nil {
		return err
	}

	err = generateRandBodyFile(bodiesFile, opts.count, opts.maxBodySize)
	if err != nil {
		return err
	}

	return nil
}

// randString returns a random string with given length
func randString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// generate a file with randomised urls, one on each line
func generateRandUrlFile(output string, baseUrl string, count int) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	for i := 0; i < count; i++ {
		file.WriteString(fmt.Sprintf("%s/%d_%s\n", baseUrl, now, randString(16)))
	}

	return nil
}

// generate a file with randomised content, one on each line
func generateRandBodyFile(output string, count int, maxSize int) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		file.WriteString(randString(mrand.Intn(maxSize)) + "\n")
	}

	return nil
}
