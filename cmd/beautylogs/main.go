package main

import (
	"github.com/Grisha1Kadetov/BeautyLogs/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
