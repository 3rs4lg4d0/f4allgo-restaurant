package boot

import (
	"strings"
	"sync"
	"time"

	tally "github.com/uber-go/tally/v4"
	promreporter "github.com/uber-go/tally/v4/prometheus"
)

var scopeOnce sync.Once

var scope tally.Scope

func GetTallyReporter() promreporter.Reporter {
	var r promreporter.Reporter
	scopeOnce.Do(func() {
		r = promreporter.NewReporter(promreporter.Options{})
		scope, _ = tally.NewRootScope(tally.ScopeOptions{
			Prefix:         strings.Replace(GetConfig().AppName, "-", "_", -1),
			Tags:           map[string]string{},
			CachedReporter: r,
			Separator:      promreporter.DefaultSeparator,
		}, time.Second)
	})

	return r
}

func GetTallyScope() tally.Scope {
	return scope
}
