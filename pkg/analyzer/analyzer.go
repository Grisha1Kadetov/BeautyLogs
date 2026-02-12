package analyzer

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

type StringList []string

func (s *StringList) String() string {
	return strings.Join(*s, ",")
}

func (s *StringList) Set(v string) error {
	if *s == nil {
		*s = []string{}
	}
	
	for _, part := range strings.Split(v, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			*s = append(*s, part)
		}
	}

	return nil
}

type RuneMap map[rune]any

func (rm *RuneMap) String() string {
	runes := make([]string, len(*rm))

	i := 0
	for k := range *rm {
		runes[i] = string(k)
		i++
	}
	return strings.Join(runes, ",")
}

func (rm *RuneMap) Set(v string) error {
	if *rm == nil {
		*rm = make(map[rune]any)
	}

	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}

	for _, r := range v {
		(*rm)[r] = true
	}
	return nil
}

type Config struct {
	onlyEng       bool
	lowercase     bool
	special       bool
	sensitive     bool
	sensitiveKeys StringList
	ignoreSpecial RuneMap
}

var Analyzer = New()

func New() *analysis.Analyzer {
	config := Config{
		onlyEng:       true,
		lowercase:     true,
		special:       true,
		sensitive:     true,
		sensitiveKeys: StringList{"password", "secret", "token", "key", "credential"},
		ignoreSpecial: make(RuneMap),
	}

	a := &analysis.Analyzer{
		Name: "beautylogs",
		Doc:  "checks for printf-like functions and their naming",
		Run:  func(p *analysis.Pass) (any, error) { return run(p, &config) },
	}

	a.Flags.BoolVar(&config.onlyEng, "only-eng", config.onlyEng, "require english/latin only")
	a.Flags.BoolVar(&config.lowercase, "lowercase", config.lowercase, "require lowercase first letter")
	a.Flags.BoolVar(&config.special, "special-char", config.special, "check special chars and emoji")
	a.Flags.BoolVar(&config.sensitive, "sensitive", config.sensitive, "check sensitive keywords")
	a.Flags.Var(&config.sensitiveKeys, "sensitive-keys", "comma-separated list of sensitive keywords to check")
	a.Flags.Var(&config.ignoreSpecial, "ignore-special-chars", "string of runes to ignore (e.g. ':;🙂')")

	return a
}

func NewCustom(config *Config) *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "beautylogs-custom",
		Doc:  "checks for printf-like functions and their naming with custom config",
		Run:  func(p *analysis.Pass) (any, error) { return run(p, config) },
	}
	return a
}
