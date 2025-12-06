package steps

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

// resolveFeaturePaths attempts to find the features directory regardless
// of where the test is executed (repo root, tests/, or tests/steps/)
func resolveFeaturePaths() []string {
	// Allows overriding via environment variable if desired
	if p := os.Getenv("GODOG_FEATURES"); p != "" {
		return []string{p}
	}
	pathCandidates := []string{
		filepath.Join("..", "features"),
	}
	var paths []string
	for _, c := range pathCandidates {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			paths = append(paths, c)
			break
		}
	}
	// fallback: mantém "features" mesmo que não exista (Godog avisará)
	if len(paths) == 0 {
		paths = []string{"features"}
	}
	return paths
}

var opt = godog.Options{
	Output: colors.Colored(os.Stdout), // colorido no terminal
	Format: "pretty",                  // mostra os steps
	Paths:  resolveFeaturePaths(),     // pasta(s) de features resolvida(s)
	Tags:   "",                        // opcional: ex. "@wip && ~@ignore"62.5
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()

	status := godog.TestSuite{
		Name:                 "orchestrator",
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: InitializeTestSuite,
		Options:              &opt,
	}.Run()

	// if st := m.Run(); st > status {
	// 	status = st
	// }

	os.Exit(status)
}
