package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golangci/golangci-lint/pkg/printers"
)

func main() {
	results, err := readResults(bufio.NewReader(os.Stdin))
	if err != nil {
		fmt.Printf("reading results: %s", err)
		os.Exit(1)
	}

	writeServiceMessages(os.Stdout, results)
}

func readResults(r io.Reader) (*printers.JSONResult, error) {
	result := &printers.JSONResult{}

	err := json.NewDecoder(r).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("decoding input: %w", err)
	}

	return result, nil
}

type lintTest struct {
	name string
	enabled bool
	issues []string
}

func (lt *lintTest) getName() string {
	return fmt.Sprintf("linter: %s", lt.name)
}

func (lt *lintTest) failed() bool {
	return len(lt.issues) > 0
}

const (
	timestampFormat = "2006-01-02T15:04:05.000"
	testStarted     = "##teamcity[testStarted timestamp='%s' name='%s']"
	testStdErr      = "##teamcity[testStdErr timestamp='%s' name='%s' out='%s']"
	testFailed      = "##teamcity[testFailed timestamp='%s' name='%s']"
	testIgnored     = "##teamcity[testIgnored timestamp='%s' name='%s']"
	testFinished    = "##teamcity[testFinished timestamp='%s' name='%s']"
)

func getNow() string {
	return time.Now().Format(timestampFormat)
}

func mustFprintln(w io.Writer, a ...interface{}) {
	_, err := fmt.Fprintln(w, a...)
	if err != nil {
		panic("writing to io.Writer: " + err.Error())
	}
}

func writeServiceMessages(w io.Writer, results *printers.JSONResult) {
	linterTests := map[string]*lintTest{}

	for _, linter := range results.Report.Linters {
		linterTests[linter.Name] = &lintTest{
			name: linter.Name,
			enabled: linter.Enabled,
		}
	}

	for _, issue := range results.Issues {
		linterTests[issue.FromLinter].issues = append(
			linterTests[issue.FromLinter].issues,
			fmt.Sprintf("%s:%v - %s", issue.FilePath(), issue.Line(), issue.Text),
		)
	}


	for _, test := range linterTests {
		mustFprintln(w, fmt.Sprintf(testStarted, getNow(), test.getName()))

		if !test.enabled {
			mustFprintln(w, fmt.Sprintf(testIgnored, getNow(), test.getName()))
		} else {
			if test.failed() {
				for _, issue := range test.issues {
					mustFprintln(w, fmt.Sprintf(testStdErr, getNow(), test.getName(), issue))
				}
				mustFprintln(w, fmt.Sprintf(testFailed, getNow(), test.getName()))
			} else {
				mustFprintln(w, fmt.Sprintf(testFinished, getNow(), test.getName()))
			}
		}
	}
}
