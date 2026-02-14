package analyzer_test

import (
	"testing"

	"github.com/Grisha1Kadetov/BeautyLogs/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzerCustomLogger(t *testing.T) {
	t.Parallel()
	config := analyzer.DefaultConfig
	config.Loggers = []analyzer.LoggerInfo{
		{Pkg: "custom-logger", Funcs: []string{"CustomPrint"}},
	}
	a := analyzer.NewCustom(config)
	analysistest.Run(t, analysistest.TestData(), a, "custom-logger")
}

func TestAnalyzerLowercase(t *testing.T) {
	t.Parallel()
	config := analyzer.DefaultConfig
	config.OnlyEng = false
	a := analyzer.NewCustom(config)
	analysistest.Run(t, analysistest.TestData(), a, "lowercase")
}

func TestAnalyzerLowercaseOff(t *testing.T) {
	t.Parallel()
	config := analyzer.DefaultConfig
	config.Lowercase = false
	a := analyzer.NewCustom(config)
	analysistest.Run(t, analysistest.TestData(), a, "lowercase-off")
}

func TestAnalyzerOnlyEng(t *testing.T) {
	t.Parallel()
	config := analyzer.DefaultConfig
	a := analyzer.NewCustom(config)
	analysistest.Run(t, analysistest.TestData(), a, "onlyeng")
}

func TestAnalyzerOnlyEngOff(t *testing.T) {
	t.Parallel()
	config := analyzer.DefaultConfig
	config.OnlyEng = false
	a := analyzer.NewCustom(config)
	analysistest.Run(t, analysistest.TestData(), a, "onlyeng-off")
}

func TestAnalyzerSensitive(t *testing.T) {
	t.Parallel()
	config := analyzer.DefaultConfig
	a := analyzer.NewCustom(config)
	analysistest.Run(t, analysistest.TestData(), a, "sensitive")
}

func TestAnalyzerSensitiveOff(t *testing.T) {
	t.Parallel()
	config := analyzer.DefaultConfig
	config.Sensitive = false
	a := analyzer.NewCustom(config)
	analysistest.Run(t, analysistest.TestData(), a, "sensitive-off")
}

func TestAnalyzerSpecial(t *testing.T) {
	t.Parallel()
	config := analyzer.DefaultConfig
	config.IgnoreSpecial = map[rune]any{'#': true}
	a := analyzer.NewCustom(config)
	analysistest.Run(t, analysistest.TestData(), a, "special")
}

func TestAnalyzerSpecialOff(t *testing.T) {
	t.Parallel()
	config := analyzer.DefaultConfig
	config.Special = false
	a := analyzer.NewCustom(config)
	analysistest.Run(t, analysistest.TestData(), a, "special-off")
}
