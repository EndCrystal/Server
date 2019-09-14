package logprefix

import "log"

func LogPrefix(prefix string) (ret string) {
	ret = log.Prefix()
	log.SetPrefix(prefix)
	return ret
}
