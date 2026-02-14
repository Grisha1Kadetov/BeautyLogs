package plugin

import (
	"fmt"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/Grisha1Kadetov/BeautyLogs/pkg/analyzer"
)

func init() {
	register.Plugin("beautylogs", New)
}

type Settings struct {
	OnlyEng            *bool    `json:"only-eng"`
	Lowercase          *bool    `json:"lowercase"`
	Special            *bool    `json:"special-char"`
	Sensitive          *bool    `json:"sensitive"`
	SensitiveKeys      []string `json:"sensitive-keys"`
	IgnoreSpecialChars *string  `json:"ignore-special-chars"`
	// logger:
	//   - "fmt:Printf,Println,Print"
	//   - "log/slog:Info,Error"
	Logger []string `json:"logger"`
}

type Plugin struct {
	cfg analyzer.Config
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}

	config := analyzer.DefaultConfig

	if s.OnlyEng != nil {
		config.OnlyEng = *s.OnlyEng
	}
	if s.Lowercase != nil {
		config.Lowercase = *s.Lowercase
	}
	if s.Special != nil {
		config.Special = *s.Special
	}
	if s.Sensitive != nil {
		config.Sensitive = *s.Sensitive
	}

	if len(s.SensitiveKeys) > 0 {
		config.SensitiveKeys = analyzer.StringList(s.SensitiveKeys)
	}

	if s.IgnoreSpecialChars != nil {
		config.IgnoreSpecial = analyzer.RuneMap{}
		for _, r := range *s.IgnoreSpecialChars {
			config.IgnoreSpecial[r] = true
		}
	}

	if len(s.Logger) > 0 {
		config.Loggers = nil
		for _, item := range s.Logger {
			if err := (&config.Loggers).Set(item); err != nil {
				return nil, fmt.Errorf("bad logger %q: %w", item, err)
			}
		}
	}

	return &Plugin{cfg: config}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.NewCustom(p.cfg),
	}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
