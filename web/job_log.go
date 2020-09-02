package web

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/shunfei/cronsun"
)

func EnsureJobLogIndex() {
	cronsun.GetDb().WithC(cronsun.Coll_JobLog, func(c *mgo.Collection) error {
		c.EnsureIndex(mgo.Index{
			Key: []string{"beginTime"},
		})
		c.EnsureIndex(mgo.Index{
			Key: []string{"hostname"},
		})
		c.EnsureIndex(mgo.Index{
			Key: []string{"ip"},
		})

		return nil
	})
}

type JobLog struct{}

func (jl *JobLog) GetDetail(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.String(http.StatusBadRequest, "empty log id.")
		return
	}

	if !bson.IsObjectIdHex(id) {
		c.String(http.StatusBadRequest, "invalid ObjectId.")
		return
	}

	logDetail, err := cronsun.GetJobLogById(bson.ObjectIdHex(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == mgo.ErrNotFound {
			statusCode = http.StatusNotFound
		}
		c.String(statusCode, err.Error())
		return
	}

	c.JSON(http.StatusOK, logDetail)
}

func searchText(field string, keywords []string) (q []bson.M) {
	for _, k := range keywords {
		k = strings.TrimSpace(k)
		if len(k) == 0 {
			continue
		}
		q = append(q, bson.M{field: bson.M{"$regex": bson.RegEx{Pattern: k, Options: "i"}}})
	}

	return q
}

func (jl *JobLog) GetList(c *gin.Context) {
	hostnames := getStringArrayFromQuery("hostnames", ",", c.Request)
	ips := getStringArrayFromQuery("ips", ",", c.Request)
	names := getStringArrayFromQuery("names", ",", c.Request)
	ids := getStringArrayFromQuery("ids", ",", c.Request)
	begin := getTime(c.Query("begin"))
	end := getTime(c.Query("end"))
	page := getPage(c.Query("page"))
	failedOnly := c.Query("failedOnly") == "true"
	pageSize := getPageSize(c.Query("pageSize"))
	orderBy := "-beginTime"

	query := bson.M{}
	var textSearch = make([]bson.M, 0, 2)
	textSearch = append(textSearch, searchText("hostname", hostnames)...)
	textSearch = append(textSearch, searchText("name", names)...)

	if len(ips) > 0 {
		query["ip"] = bson.M{"$in": ips}
	}

	if len(ids) > 0 {
		query["jobId"] = bson.M{"$in": ids}
	}

	if !begin.IsZero() {
		query["beginTime"] = bson.M{"$gte": begin}
	}
	if !end.IsZero() {
		query["endTime"] = bson.M{"$lt": end.Add(time.Hour * 24)}
	}

	if failedOnly {
		query["success"] = false
	}

	if len(textSearch) > 0 {
		query["$or"] = textSearch
	}

	var pager struct {
		Total int               `json:"total"`
		List  []*cronsun.JobLog `json:"list"`
	}
	var err error
	if c.Request.FormValue("latest") == "true" {
		var latestLogList []*cronsun.JobLatestLog
		latestLogList, pager.Total, err = cronsun.GetJobLatestLogList(query, page, pageSize, orderBy)
		for i := range latestLogList {
			latestLogList[i].JobLog.Id = bson.ObjectIdHex(latestLogList[i].RefLogId)
			pager.List = append(pager.List, &latestLogList[i].JobLog)
		}
	} else {
		pager.List, pager.Total, err = cronsun.GetJobLogList(query, page, pageSize, orderBy)
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	pager.Total = int(math.Ceil(float64(pager.Total) / float64(pageSize)))
	c.JSON(http.StatusInternalServerError, pager)
}
