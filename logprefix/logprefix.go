package logprefix

import "log"

// Get get prefixied logger
func Get(prefix string) *log.Logger {
	return log.New(log.Writer(), prefix, log.Flags())
}
