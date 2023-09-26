// Package boot provides utility functions to initialize application dependencies
// at boot time (like loggers, database connections or metrics), normally from the
// main function. It also serves the adapter layer as a central package for accessing
// some useful initialized dependencies like application configuration (obtained from
// the environment) or the metrics scope. Don't import this package from the core layer.
package boot
