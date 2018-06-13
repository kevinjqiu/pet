package main

import (
	"github.com/hashicorp/go-multierror"
	"log"
	"path/filepath"
	"os"
	"flag"
	"fmt"
	"os/exec"
	"encoding/json"
	"encoding/base64"
)

const EnvVarTestFilePathsB64 = "TEST_FILE_PATHS_B64"

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

func encodeTestFilePaths(testFiles []string) (string, error) {
	jsonBody, err := json.Marshal(testFiles)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonBody), err
}

func decodeTestFilePaths(encodedTestFilePaths string) ([]string, error) {
	jsonBody, err := base64.StdEncoding.DecodeString(encodedTestFilePaths)
	if err != nil {
		return []string{}, err
	}

	var filePaths []string
	err = json.Unmarshal(jsonBody, &filePaths)

	if err != nil {
		return []string{}, err
	}
	return filePaths, nil
}

func main() {
	flag.Parse()
	testFiles, err := collectTestFiles(flag.Args())

	encodedTestFilePaths, err := encodeTestFilePaths(testFiles)

	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("go", "test")

	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", EnvVarTestFilePathsB64, encodedTestFilePaths))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
/*
func TestAlertingRule(t *testing.T) {
	suite, err := promql.NewTest(t, `
		load 5m
			http_requests{job="app-server", instance="0", group="canary", severity="overwrite-me"}	75 85  95 105 105  95  85
			http_requests{job="app-server", instance="1", group="canary", severity="overwrite-me"}	80 90 100 110 120 130 140
	`)
	testutil.Ok(t, err)
	defer suite.Close()

	err = suite.Run()
	testutil.Ok(t, err)

	expr, err := promql.ParseExpr(`http_requests{group="canary", job="app-server"} < 100`)
	testutil.Ok(t, err)

	rule := NewAlertingRule(
		"HTTPRequestRateLow",
		expr,
		time.Minute,
		labels.FromStrings("severity", "{{\"c\"}}ritical"),
		nil, nil,
	)
	result := promql.Vector{
		{
			Metric: labels.FromStrings(
				"__name__", "ALERTS",
				"alertname", "HTTPRequestRateLow",
				"alertstate", "pending",
				"group", "canary",
				"instance", "0",
				"job", "app-server",
				"severity", "critical",
			),
			Point: promql.Point{V: 1},
		},
		{
			Metric: labels.FromStrings(
				"__name__", "ALERTS",
				"alertname", "HTTPRequestRateLow",
				"alertstate", "pending",
				"group", "canary",
				"instance", "1",
				"job", "app-server",
				"severity", "critical",
			),
			Point: promql.Point{V: 1},
		},
		{
			Metric: labels.FromStrings(
				"__name__", "ALERTS",
				"alertname", "HTTPRequestRateLow",
				"alertstate", "firing",
				"group", "canary",
				"instance", "0",
				"job", "app-server",
				"severity", "critical",
			),
			Point: promql.Point{V: 1},
		},
		{
			Metric: labels.FromStrings(
				"__name__", "ALERTS",
				"alertname", "HTTPRequestRateLow",
				"alertstate", "firing",
				"group", "canary",
				"instance", "1",
				"job", "app-server",
				"severity", "critical",
			),
			Point: promql.Point{V: 1},
		},
	}

	baseTime := time.Unix(0, 0)

	var tests = []struct {
		time   time.Duration
		result promql.Vector
	}{
		{
			time:   0,
			result: result[:2],
		}, {
			time:   5 * time.Minute,
			result: result[2:],
		}, {
			time:   10 * time.Minute,
			result: result[2:3],
		},
		{
			time:   15 * time.Minute,
			result: nil,
		},
		{
			time:   20 * time.Minute,
			result: nil,
		},
		{
			time:   25 * time.Minute,
			result: result[:1],
		},
		{
			time:   30 * time.Minute,
			result: result[2:3],
		},
	}

	for i, test := range tests {
		t.Logf("case %d", i)

		evalTime := baseTime.Add(test.time)

		res, err := rule.Eval(suite.Context(), evalTime, EngineQueryFunc(suite.QueryEngine(), suite.Storage()), nil)
		testutil.Ok(t, err)

		for i := range test.result {
			test.result[i].T = timestamp.FromTime(evalTime)
		}
		testutil.Assert(t, len(test.result) == len(res), "%d. Number of samples in expected and actual output don't match (%d vs. %d)", i, len(test.result), len(res))

		sort.Slice(res, func(i, j int) bool {
			return labels.Compare(res[i].Metric, res[j].Metric) < 0
		})
		testutil.Equals(t, test.result, res)

		for _, aa := range rule.ActiveAlerts() {
			testutil.Assert(t, aa.Labels.Get(model.MetricNameLabel) == "", "%s label set on active alert: %s", model.MetricNameLabel, aa.Labels)
		}
	}
}
*/
