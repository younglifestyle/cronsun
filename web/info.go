package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	v3 "github.com/coreos/etcd/clientv3"

	"github.com/longcron/cronjob"
	"github.com/longcron/cronjob/conf"
)

type Info struct{}

func (inf *Info) Overview(c *gin.Context) {
	var overview = struct {
		TotalJobs        int64                   `json:"totalJobs"`
		JobExecuted      *cronjob.StatExecuted   `json:"jobExecuted"`
		JobExecutedDaily []*cronjob.StatExecuted `json:"jobExecutedDaily"`
	}{}

	const day = 24 * time.Hour
	days := 7

	overview.JobExecuted, _ = cronjob.JobLogStat()
	end := time.Now()
	begin := end.Add(time.Duration(1-days) * day)
	statList, _ := cronjob.JobLogDailyStat(begin, end)
	list := make([]*cronjob.StatExecuted, days)
	cur := begin

	for i := 0; i < days; i++ {
		date := cur.Format("2006-01-02")
		var se *cronjob.StatExecuted

		for j := range statList {
			if statList[j].Date == date {
				se = statList[j]
				statList = statList[1:]
				break
			}
		}

		if se != nil {
			list[i] = se
		} else {
			list[i] = &cronjob.StatExecuted{Date: date}
		}

		cur = cur.Add(day)
	}

	overview.JobExecutedDaily = list
	gresp, err := cronjob.DefalutClient.Get(conf.Config.Cmd, v3.WithPrefix(), v3.WithCountOnly())
	if err == nil {
		overview.TotalJobs = gresp.Count
	}

	c.JSON(http.StatusOK, overview)
}
