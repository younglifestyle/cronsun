package cronjob

import (
	"github.com/longcron/cronjob/db"
)

var (
	mgoDB *db.Mdb
)

func GetDb() *db.Mdb {
	return mgoDB
}
