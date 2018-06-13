package main

import (
	"flag"
	"github.com/prometheus/prometheus/promql"
	"github.com/hashicorp/go-multierror"
	"log"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"os"
)

type TestCase struct {
}

func (tc TestCase) Fatal(args ...interface{}) {
	log.Fatal(args)
}

func (tc TestCase) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

// var testfilesGlob string

func init() {
	// flag.StringVar(&testfilesGlob, "t", "*.test", "test file glob pattern")
}

func collectTestFiles(globPatterns []string) ([]string, error) {
	if len(globPatterns) == 0 {
		return filepath.Glob("*.test")
	}

	var (
		err error
		filePaths []string
	)

	for _, pattern := range globPatterns {
		files, newErr := filepath.Glob(pattern)
		if err != nil {
			err = multierror.Append(err, newErr)
		}

		for _, file := range files {
			if _, err := os.Stat(file); err == nil {
				filePaths = append(filePaths, file)
			}
		}
	}

	return filePaths, err
}

func main() {
	flag.Parse()
	testFiles, err := collectTestFiles(flag.Args())

	if err != nil {
		log.Fatal(err)
	}

	testCase := TestCase{}

	var testErrs error
	for _, file := range testFiles {
		fmt.Printf("Testing %v...", file)

		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		t, err := promql.NewTest(testCase, string(content))
		if err != nil {
			log.Fatal(err)
		}
		testErr := t.Run()
		if testErr != nil {
			testErrs = multierror.Append(testErrs, testErr)
			fmt.Print("Failed\n")
		} else {
			fmt.Print("OK\n")
		}
	}

	if testErrs != nil {
		log.Fatal(testErrs)
	}
}
