package boot

import (
	"database/sql"
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

// GetTallyReporter initializes (only once) the metrics reporter and a root scope
// and returns the reporter. It also registers general purpose metrics in the scope
// like DBStats (if sql.DB is provided) and other useful stats related with the
// infrastructure (e.g. CPU).
func GetTallyReporter(db *sql.DB) promreporter.Reporter {
	metricsOnce.Do(func() {
		reporter = promreporter.NewReporter(promreporter.Options{})

		// Default scope for the application. All the metrics are reported every
		// second to a cached reporter.
		scope, _ = tally.NewRootScope(tally.ScopeOptions{
			Prefix:         strings.Replace(GetConfig().AppName, "-", "_", -1),
			Tags:           map[string]string{},
			CachedReporter: reporter,
			Separator:      promreporter.DefaultSeparator,
		}, time.Second)

		// If sql.DB is provided register all the metrics collected in sql.DBStats
		// in its own Tally scope and report them periodically.
		if db != nil {
			registerDBStatsAndReport(db, reporter)
		}
	})

	return reporter
}

// GetTallyScope returns the scope attached to the reporter. Adapter layers can
// access the scope to register new metrics.
func GetTallyScope() tally.Scope {
	return scope
}

// registerDBStatsAndReport registers sql.DBStats as Tally Gauges and starts an endless
// goroutine to keep them updated.
func registerDBStatsAndReport(db *sql.DB, reporter promreporter.Reporter) {
	s, _ := tally.NewRootScope(tally.ScopeOptions{
		CachedReporter: reporter,
	}, time.Second)

	_, _ = reporter.RegisterGauge("sql_dbstats_max_open_connections", nil, "Maximum number of open connections to the database.")
	_, _ = reporter.RegisterGauge("sql_dbstats_open_connections", nil, "The number of established connections both in use and idle.")
	_, _ = reporter.RegisterGauge("sql_dbstats_in_use", nil, "The number of connections currently in use.")
	_, _ = reporter.RegisterGauge("sql_dbstats_idle", nil, "The number of idle connections.")
	_, _ = reporter.RegisterGauge("sql_dbstats_wait_count", nil, "The total number of connections waited for.")
	_, _ = reporter.RegisterGauge("sql_dbstats_wait_duration", nil, "The total time blocked waiting for a new connection.")
	_, _ = reporter.RegisterGauge("sql_dbstats_max_idle_closed", nil, "The total number of connections closed due to SetMaxIdleConns.")
	_, _ = reporter.RegisterGauge("sql_dbstats_max_lifetime_closed", nil, "The total number of connections closed due to SetConnMaxLifetime.")
	_, _ = reporter.RegisterGauge("sql_dbstats_max_idletime_closed", nil, "The total number of connections closed due to SetConnMaxIdleTime.")

	go func() {
		for range time.Tick(3 * time.Second) {
			s.Gauge("sql_dbstats_max_open_connections").Update(float64(db.Stats().MaxOpenConnections))
			s.Gauge("sql_dbstats_open_connections").Update(float64(db.Stats().OpenConnections))
			s.Gauge("sql_dbstats_in_use").Update(float64(db.Stats().InUse))
			s.Gauge("sql_dbstats_idle").Update(float64(db.Stats().Idle))
			s.Gauge("sql_dbstats_wait_count").Update(float64(db.Stats().WaitCount))
			s.Gauge("sql_dbstats_wait_duration").Update(float64(db.Stats().WaitDuration))
			s.Gauge("sql_dbstats_max_idle_closed").Update(float64(db.Stats().MaxIdleClosed))
			s.Gauge("sql_dbstats_max_lifetime_closed").Update(float64(db.Stats().MaxLifetimeClosed))
			s.Gauge("sql_dbstats_max_idletime_closed").Update(float64(db.Stats().MaxIdleTimeClosed))
		}
	}()
}
