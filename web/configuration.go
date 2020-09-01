package web

import (
	"github.com/shunfei/cronsun/conf"
	"net/http"
)

type Configuration struct{}

func (cnf *Configuration) Configuratios(W http.ResponseWriter, R *http.Request) {
	r := struct {
		Alarm             bool `json:"alarm"`
		LogExpirationDays int  `json:"log_expiration_days"`
	}{
		Alarm: conf.Config.Mail.Enable,
	}

	if conf.Config.Web.LogCleaner.EveryMinute > 0 {
		r.LogExpirationDays = conf.Config.Web.LogCleaner.ExpirationDays
	}

	outJSON(W, r)
}
