package main

import (
	"testing"
	"os"
	"fmt"
	"io/ioutil"
	"github.com/prometheus/prometheus/promql"
	"github.com/hashicorp/go-multierror"
	"log"
)

type TestCase struct {
}

func (tc TestCase) Fatal(args ...interface{}) {
	log.Fatal(args)
}

func (tc TestCase) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

func testPromQL(testFiles []string) {
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
			fmt.Print("FAILED\n")
		} else {
			fmt.Print("OK\n")
		}
	}

	if testErrs != nil {
		log.Fatal(testErrs)
	}
}

func TestPromQLExpression(t *testing.T) {
	val, ok := os.LookupEnv(EnvVarTestFilePathsB64)
	if !ok {
		t.Fatalf("Environment variable %s not found!\n", EnvVarTestFilePathsB64)
	}
	filePaths, err := decodeTestFilePaths(val)
	if err != nil {
		t.Fatal(err)
	}

	testPromQL(filePaths)
}
