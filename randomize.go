package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	mrand "math/rand"
	"os"
	"time"
)

func randomize(spec Spec) error {
	err := generateRandUrlFile(spec.UrlsFile, spec.BaseUrl, spec.RandCount)
	if err != nil {
		return err
	}

	err = generateRandBodyFile(spec.BodiesFile, spec.RandCount, spec.MaxBodySize)
	if err != nil {
		return err
	}

	// create targets file from random urls and request bodies
	// for vegeta usage
	for _, ep := range spec.Endpoints {
		if err = writeTargetsFile(ep.Method, spec.UrlsFile, spec.BodiesFile, ep.TargetsFile); err != nil {
			return err
		}
	}

	return nil
}

// randString returns a random string with given length
func randString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
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
			if _, err = targetsFile.WriteString(fmt.Sprintf("%s %s\n", method, urlReader.Text())); err != nil {
				return err
			}
		}
	}

	return nil
}

// generate a file with randomised urls, one on each line
func generateRandUrlFile(output string, baseUrl string, count int) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

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
	defer file.Close()

	for i := 0; i < count; i++ {
		file.WriteString(randString(mrand.Intn(maxSize)) + "\n")
	}

	return nil
}
