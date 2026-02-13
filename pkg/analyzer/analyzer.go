package analyzer

import (
	"fmt"
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
	return strings.Join(runes, " ")
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

type LoggerInfo struct {
	Funcs StringList `mapstructure:"funcs"`
	Pkg   string     `mapstructure:"pkg"`
}

type LoggerInfos []LoggerInfo

// -logger="fmt:Printf,Println" -logger="log/slog:Info"
func (ls *LoggerInfos) Set(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}

	pkg, rest, ok := strings.Cut(v, ":")
	if !ok {
		return fmt.Errorf("bad --logger %q: expected <pkg>:<func1,func2,...>", v)
	}
	pkg = strings.TrimSpace(pkg)
	rest = strings.TrimSpace(rest)
	if pkg == "" {
		return fmt.Errorf("bad --logger %q: empty pkg", v)
	}
	if rest == "" {
		return fmt.Errorf("bad --logger %q: empty funcs list", v)
	}

	var funcs StringList
	if err := funcs.Set(rest); err != nil {
		return err
	}
	if len(funcs) == 0 {
		return fmt.Errorf("bad --logger %q: no funcs parsed", v)
	}

	// merge
	for i := range *ls {
		if (*ls)[i].Pkg == pkg {
			(*ls)[i].Funcs = append((*ls)[i].Funcs, funcs...)
			return nil
		}
	}

	*ls = append(*ls, LoggerInfo{Pkg: pkg, Funcs: funcs})
	return nil
}

func (ls *LoggerInfos) String() string {
	if ls == nil || len(*ls) == 0 {
		return ""
	}
	parts := make([]string, 0, len(*ls))
	for _, li := range *ls {
		if li.Pkg == "" || len(li.Funcs) == 0 {
			continue
		}
		parts = append(parts, li.Pkg+":"+li.Funcs.String())
	}
	return strings.Join(parts, "; ")
}

type Config struct {
	OnlyEng       bool
	Lowercase     bool
	Special       bool
	Sensitive     bool
	SensitiveKeys StringList
	IgnoreSpecial RuneMap
	Loggers       LoggerInfos `mapstructure:"loggers"`
}

var Analyzer = New()

var DefaultConfig = Config{
	OnlyEng:       true,
	Lowercase:     true,
	Special:       true,
	Sensitive:     true,
	SensitiveKeys: StringList{"password", "secret", "token", "key", "credential"},
	IgnoreSpecial: RuneMap{':': true, ',': true, '-': true, '.': true, '_': true, '(': true, ')': true},
	Loggers: []LoggerInfo{
		{Funcs: StringList{"Printf", "Println", "Print"}, Pkg: "fmt"},
		{Pkg: "go.uber.org/zap", Funcs: StringList{
			"Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal", "Log",
			"Debugf", "Infof", "Warnf", "Errorf", "DPanicf", "Panicf", "Fatalf",
			"Debugln", "Infoln", "Warnln", "Errorln", "DPanicln", "Panicln", "Fatalln",
			"Debugw", "Infow", "Warnw", "Errorw", "DPanicw", "Panicw", "Fatalw",
			"Logf", "Logln", "Logw",
		}},
		{Pkg: "log/slog", Funcs: StringList{
			"Debug", "Info", "Warn", "Error",
			"DebugContext", "InfoContext", "WarnContext", "ErrorContext",
			"Log", "LogAttrs",
		}},
	},
}

func New() *analysis.Analyzer {
	config := Config{
		OnlyEng:   true,
		Lowercase: true,
		Special:   true,
		Sensitive: true,
	}

	a := &analysis.Analyzer{
		Name: "beautylogs",
		Doc:  "checks for printf-like functions and their naming",
		Run:  func(p *analysis.Pass) (any, error) { return run(p, config) },
	}

	a.Flags.BoolVar(&config.OnlyEng, "only-eng", config.OnlyEng, "require english/latin only")
	a.Flags.BoolVar(&config.Lowercase, "lowercase", config.Lowercase, "require lowercase first letter")
	a.Flags.BoolVar(&config.Special, "special-char", config.Special, "check special chars and emoji")
	a.Flags.Var(&config.IgnoreSpecial, "ignore-special-chars", "string of runes to ignore (e.g. ':;')")
	a.Flags.BoolVar(&config.Sensitive, "sensitive", config.Sensitive, "check sensitive keywords")
	a.Flags.Var(&config.SensitiveKeys, "sensitive-keys", "comma-separated list of sensitive keywords to check")
	a.Flags.Var(&config.Loggers, "logger", `logger targets in form "<pkg>:<f1,f2,...>". repeatable`)

	return a
}

func NewCustom(config Config) *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "beautylogs-custom",
		Doc:  "checks for printf-like functions and their naming with custom config",
		Run:  func(p *analysis.Pass) (any, error) { return run(p, config) },
	}
	return a
}
