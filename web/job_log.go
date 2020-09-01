package web

import (
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
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

func (jl *JobLog) GetDetail(W http.ResponseWriter, R *http.Request) {
	vars := mux.Vars(R)
	id := strings.TrimSpace(vars["id"])
	if len(id) == 0 {
		outJSONWithCode(W, http.StatusBadRequest, "empty log id.")
		return
	}

	if !bson.IsObjectIdHex(id) {
		outJSONWithCode(W, http.StatusBadRequest, "invalid ObjectId.")
		return
	}

	logDetail, err := cronsun.GetJobLogById(bson.ObjectIdHex(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == mgo.ErrNotFound {
			statusCode = http.StatusNotFound
			err = nil
		}
		outJSONWithCode(W, statusCode, err)
		return
	}

	outJSON(W, logDetail)
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

func (jl *JobLog) GetList(W http.ResponseWriter, R *http.Request) {
	hostnames := getStringArrayFromQuery("hostnames", ",", R)
	ips := getStringArrayFromQuery("ips", ",", R)
	names := getStringArrayFromQuery("names", ",", R)
	ids := getStringArrayFromQuery("ids", ",", R)
	begin := getTime(R.FormValue("begin"))
	end := getTime(R.FormValue("end"))
	page := getPage(R.FormValue("page"))
	failedOnly := R.FormValue("failedOnly") == "true"
	pageSize := getPageSize(R.FormValue("pageSize"))
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
	if R.FormValue("latest") == "true" {
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
		outJSONWithCode(W, http.StatusInternalServerError, err.Error())
		return
	}

	pager.Total = int(math.Ceil(float64(pager.Total) / float64(pageSize)))
	outJSON(W, pager)
}
