package boot

import (
	"strings"
	"sync"
	"time"

	tally "github.com/uber-go/tally/v4"
	promreporter "github.com/uber-go/tally/v4/prometheus"
)

var metricsOnce sync.Once

var (
	reporter promreporter.Reporter
	scope    tally.Scope
)

func GetTallyReporter() promreporter.Reporter {
	metricsOnce.Do(func() {
		reporter = promreporter.NewReporter(promreporter.Options{})
		scope, _ = tally.NewRootScope(tally.ScopeOptions{
			Prefix:         strings.Replace(GetConfig().AppName, "-", "_", -1),
			Tags:           map[string]string{},
			CachedReporter: reporter,
			Separator:      promreporter.DefaultSeparator,
		}, time.Second)
	})

	return reporter
}

func GetTallyScope() tally.Scope {
	return scope
}
