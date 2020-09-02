package cronjob

import (
	"fmt"
	"github.com/longcron/cronjob/conf"
	"github.com/longcron/cronjob/db"
)

const (
	BasePath       = "/data/cron-job/"
	FileSourcePath = BasePath + "filesource/"
)

var (
	initialized bool
)

func Init(baseConfFile string, watchConfiFile bool) (err error) {
	if initialized {
		return
	}

	// init id creator
	if err = initID(); err != nil {
		return fmt.Errorf("Init UUID Generator failed: %s", err)
	}

	// init config
	if err = conf.Init(baseConfFile, watchConfiFile); err != nil {
		return fmt.Errorf("Init Config failed: %s", err)
	}

	// init etcd client
	if DefalutClient, err = NewClient(conf.Config); err != nil {
		return fmt.Errorf("Connect to ETCD %s failed: %s",
			conf.Config.Etcd.Endpoints, err)
	}

	// init mongoDB
	if mgoDB, err = db.NewMdb(conf.Config.Mgo); err != nil {
		return fmt.Errorf("Connect to MongoDB %s failed: %s",
			conf.Config.Mgo.Hosts, err)
	}

	initialized = true
	return
}
