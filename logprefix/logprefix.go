package logprefix

import "log"

func Get(prefix string) *log.Logger {
	return log.New(log.Writer(), prefix, log.Flags())
}
