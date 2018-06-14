package main

import (
	"testing"
	"os"
	"log"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/prometheus/prometheus/promql"
	"github.com/hashicorp/go-multierror"
	"strings"
)

type StubTestCase struct{}

func (stc StubTestCase) Fatal(args ...interface{}) {
	log.Fatal(args)
}
func (stc StubTestCase) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

// Fixture is a Prometheus metric
type Metric string

// Instant is an instant, e.g., 0m, 1m, 1d, etc
type Instant string

type InstantMetricFixtures map[Instant][]Metric

type ExpressionSource struct {
	FromFile    string `yaml:"fromFile"`
	FromLiteral string `yaml:"fromLiteral"`
}

func (es ExpressionSource) Get() (string, error) {
	if es.FromFile != "" {
		result, err := ioutil.ReadFile(es.FromFile)  // TODO(kevinjqiu): cache the result

		if err != nil {
			return "", err
		}
		return strings.Replace(string(result), "\n", " ", -1), nil
	}
	if es.FromLiteral != "" {
		return es.FromLiteral, nil
	}
	return "", nil
}

type Evaluation struct {
	At   Instant          `yaml:"at"`
	Expr ExpressionSource `yaml:"expr"`
}

type Assertion struct {
	Eval     Evaluation `yaml:"eval"`
	Expected []Metric   `yaml:"expected"`
}

type PromQLExprTestCase struct {
	Description string                `yaml:"description"`
	Fixtures    InstantMetricFixtures `yaml:"fixtures"`
	Assertions  []Assertion           `yaml:"assertions"`
}

func (pqltc PromQLExprTestCase) generateCommands() (string, error) {
	lines := []string{}

	lines = append(lines, "clear")
	for instant, fixtures := range pqltc.Fixtures {
		lines = append(lines, fmt.Sprintf("load %s", instant))
		for _, fixture := range fixtures {
			lines = append(lines, fmt.Sprintf("    %s", fixture))
		}
	}

	lines = append(lines, "")

	for _, assertion := range pqltc.Assertions {
		expr, err := assertion.Eval.Expr.Get()
		if err != nil {
			return "", err
		}
		lines = append(lines, fmt.Sprintf("eval instant at %s %s", assertion.Eval.At, expr))
		for _, expectedMetric := range assertion.Expected {
			lines = append(lines, fmt.Sprintf("    %s", expectedMetric))
		}
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n"), nil
}

func (pqltc PromQLExprTestCase) Run() error {
	commands, err := pqltc.generateCommands()
	if err != nil {
		return err
	}
	fmt.Printf("   %s.....", pqltc.Description)
	pt, err := promql.NewTest(StubTestCase{}, commands)
	if err != nil {
		fmt.Printf("ERROR\n")
		log.Fatal(err)
	}
	testErr := pt.Run()
	if testErr != nil {
		fmt.Printf("FAIL\n")
		return testErr
	}
	fmt.Printf("OK\n")
	return nil
}

type PromQLExprTest struct {
	Name      string               `yaml:"name"`
	TestCases []PromQLExprTestCase `yaml:"testCases"`
}

func (pqltest PromQLExprTest) Run() error {
	var err error

	for _, tc := range pqltest.TestCases {
		newErr := tc.Run()
		if newErr != nil {
			err = multierror.Append(err, newErr)
		}
	}

	return err
}

func parsePromQLTestCase(filePath string) (PromQLExprTest, error) {
	var testCase PromQLExprTest

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return testCase, err
	}

	err = yaml.UnmarshalStrict(content, &testCase)
	if err != nil {
		return testCase, err
	}
	return testCase, nil
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

	if len(filePaths) == 0 {
		log.Println("WARNING: no test files specified")
		return
	}

	var testErrs error
	for _, file := range filePaths {
		fmt.Printf("%v\n", file)

		testCase, err := parsePromQLTestCase(file)
		if err != nil {
			log.Fatal(err)
		}

		testErr := testCase.Run()

		if testErr != nil {
			testErrs = multierror.Append(testErrs, testErr)
		}
	}

	if testErrs != nil {
		log.Fatal(testErrs)
	}
}
