package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coreos/etcd/clientv3"

	"github.com/longcron/cronjob"
	"github.com/longcron/cronjob/conf"
	"github.com/longcron/cronjob/log"
)

type Job struct{}

func (j *Job) GetJob(c *gin.Context) {
	group := c.Param("group")
	id := c.Param("id")

	job, err := cronjob.GetJob(group, id)
	var statusCode int
	if err != nil {
		if err == cronjob.ErrNotFound {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}
		c.String(statusCode, err.Error())
		return
	}

	c.JSON(http.StatusOK, job)
}

func (j *Job) DeleteJob(c *gin.Context) {
	group := c.Param("group")
	jobId := c.Param("id")

	_, err := cronjob.DeleteJob(group, jobId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusNoContent, "")
}

func (j *Job) ChangeJobStatus(c *gin.Context) {
	_, isPause := c.GetQuery("pause")
	group := c.Param("group")
	id := c.Param("id")

	job, err := j.updateJobStatus(group, id, isPause)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, job)
}

func (j *Job) updateJobStatus(group, id string, isPause bool) (*cronjob.Job, error) {
	originJob, rev, err := cronjob.GetJobAndRev(group, id)
	if err != nil {
		return nil, err
	}

	if originJob.Pause == isPause {
		log.Debugf("no modify")
		return nil, err
	}

	originJob.Pause = isPause
	b, err := json.Marshal(originJob)
	if err != nil {
		return nil, err
	}

	_, err = cronjob.DefalutClient.PutWithModRev(originJob.Key(), string(b), rev)
	if err != nil {
		return nil, err
	}

	return originJob, nil
}

func (j *Job) BatchChangeJobStatus(c *gin.Context) {
	var jobIds []string

	err := c.BindJSON(&jobIds)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	op := c.Param("op")
	var isPause bool
	switch op {
	case "pause":
		isPause = true
	case "start":
	default:
		c.String(http.StatusBadRequest, "Unknow batch operation.")
		return
	}

	var updated int
	for i := range jobIds {
		id := strings.Split(jobIds[i], "/") // [Group, ID]
		if len(id) != 2 || id[0] == "" || id[1] == "" {
			continue
		}

		_, err = j.updateJobStatus(id[0], id[1], isPause)
		if err != nil {
			continue
		}
		updated++
	}

	c.String(http.StatusOK, fmt.Sprintf("%d of %d updated.", updated, len(jobIds)))
}

func (j *Job) UpdateJob(c *gin.Context) {
	var job = &struct {
		*cronjob.Job
		OldGroup string `json:"oldGroup"`
	}{}

	err := c.BindJSON(job)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err = job.Check(); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var deleteOldKey string
	var successStr = "update success"
	if len(job.ID) == 0 {
		successStr = "created success"
		job.ID = cronjob.NextID()
	} else {
		job.OldGroup = strings.TrimSpace(job.OldGroup)
		if job.OldGroup != job.Group {
			deleteOldKey = cronjob.JobKey(job.OldGroup, job.ID)
		}
	}

	// 当前Job更新时间
	job.UpdateTime = time.Now().Unix()
	//	job.Command = `#!/bin/bash
	//echo "xxl-job: hello shell"
	//
	//echo "脚本位置：$0"
	//echo "任务参数：$1"
	//echo "分片序号 = $2"
	//echo "分片总数 = $3"
	//
	//echo "Good bye!"
	//exit 0`
	b, err := json.Marshal(job)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// remove old key
	// it should be before the put method
	if len(deleteOldKey) > 0 {
		if _, err = cronjob.DefalutClient.Delete(deleteOldKey); err != nil {
			log.Errorf("failed to remove old job key[%s], err: %s.", deleteOldKey, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	_, err = cronjob.DefalutClient.Put(job.Key(), string(b))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, successStr)
}

func (j *Job) GetGroups(c *gin.Context) {
	resp, err := cronjob.DefalutClient.Get(conf.Config.Cmd, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var cmdKeyLen = len(conf.Config.Cmd)
	var groupMap = make(map[string]bool, 8)

	for i := range resp.Kvs {
		ss := strings.Split(string(resp.Kvs[i].Key)[cmdKeyLen:], "/")
		groupMap[ss[0]] = true
	}

	var groupList = make([]string, 0, len(groupMap))
	for k := range groupMap {
		groupList = append(groupList, k)
	}

	sort.Strings(groupList)
	c.JSON(http.StatusOK, groupList)
}

func (j *Job) GetList(c *gin.Context) {
	group, groupIsExist := c.GetQuery("group")
	node, nodeIsExist := c.GetQuery("node")
	var prefix = conf.Config.Cmd
	if groupIsExist {
		prefix += group
	}

	type jobStatus struct {
		*cronjob.Job
		LatestStatus *cronjob.JobLatestLog `json:"latestStatus"`
		NextRunTime  string                `json:"nextRunTime"`
	}

	resp, err := cronjob.DefalutClient.Get(prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var nodeGroupMap map[string]*cronjob.Group
	if nodeIsExist {
		nodeGrouplist, err := cronjob.GetNodeGroups()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		nodeGroupMap = map[string]*cronjob.Group{}
		for i := range nodeGrouplist {
			nodeGroupMap[nodeGrouplist[i].ID] = nodeGrouplist[i]
		}
	}

	var jobIds []string
	var jobList = make([]*jobStatus, 0, resp.Count)
	for i := range resp.Kvs {
		job := cronjob.Job{}
		err = json.Unmarshal(resp.Kvs[i].Value, &job)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if len(node) > 0 && !job.IsRunOn(node, nodeGroupMap) {
			continue
		}
		jobList = append(jobList, &jobStatus{Job: &job})
		jobIds = append(jobIds, job.ID)
	}

	m, err := cronjob.GetJobLatestLogListByJobIds(jobIds)
	if err != nil {
		log.Errorf("GetJobLatestLogListByJobIds error: %s", err.Error())
	} else {
		for i := range jobList {
			jobList[i].LatestStatus = m[jobList[i].ID]
			nt := jobList[i].GetNextRunTime()
			if nt.IsZero() {
				jobList[i].NextRunTime = "NO!!"
			} else {
				jobList[i].NextRunTime = nt.Format("2006-01-02 15:04:05")
			}
		}
	}

	c.JSON(http.StatusOK, jobList)
}

func (j *Job) GetJobNodes(c *gin.Context) {
	group := c.Param("group")
	jobId := c.Param("id")

	job, err := cronjob.GetJob(group, jobId)
	var statusCode int
	if err != nil {
		if err == cronjob.ErrNotFound {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}
		c.String(statusCode, err.Error())
		return
	}

	var nodes []string
	var exNodes []string
	groups, err := cronjob.GetGroups("")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	for i := range job.Rules {
		inNodes := append(nodes, job.Rules[i].NodeIDs...)
		for _, gid := range job.Rules[i].GroupIDs {
			if g, ok := groups[gid]; ok {
				inNodes = append(inNodes, g.NodeIDs...)
			}
		}
		exNodes = append(exNodes, job.Rules[i].ExcludeNodeIDs...)
		inNodes = SubtractStringArray(inNodes, exNodes)
		nodes = append(nodes, inNodes...)
	}

	c.JSON(http.StatusInternalServerError, UniqueStringArray(nodes))
}

func (j *Job) JobExecute(c *gin.Context) {
	group := c.Param("group")
	jobId := c.Param("id")
	group = strings.TrimSpace(group)
	id := strings.TrimSpace(jobId)
	if len(group) == 0 || len(id) == 0 {
		c.String(http.StatusBadRequest, "Invalid job id or group.")
		return
	}

	node := c.Query("node")
	err := cronjob.PutOnce(group, id, node)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusNoContent, "")
}

func (j *Job) GetExecutingJob(c *gin.Context) {
	opt := &ProcFetchOptions{
		Groups:  getStringArrayFromQuery("groups", ",", c.Request),
		NodeIds: getStringArrayFromQuery("nodes", ",", c.Request),
		JobIds:  getStringArrayFromQuery("jobs", ",", c.Request),
	}

	gresp, err := cronjob.DefalutClient.Get(conf.Config.Proc, clientv3.WithPrefix())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var list = make([]*processInfo, 0, 8)
	for i := range gresp.Kvs {
		proc, err := cronjob.GetProcFromKey(string(gresp.Kvs[i].Key))
		if err != nil {
			log.Errorf("Failed to unmarshal Proc from key: %s", err.Error())
			continue
		}

		if !opt.Match(proc) {
			continue
		}

		val := string(gresp.Kvs[i].Value)
		var pv = &cronjob.ProcessVal{}
		err = json.Unmarshal([]byte(val), pv)
		if err != nil {
			log.Errorf("Failed to unmarshal ProcessVal from val: %s", err.Error())
			continue
		}
		proc.ProcessVal = *pv
		procInfo := &processInfo{
			Process: proc,
		}
		job, err := cronjob.GetJob(proc.Group, proc.JobID)
		if err == nil && job != nil {
			procInfo.JobName = job.Name
		} else {
			procInfo.JobName = proc.JobID
		}
		list = append(list, procInfo)
	}

	sort.Sort(ByProcTime(list))
	c.JSON(http.StatusOK, list)
}

func (j *Job) KillExecutingJob(c *gin.Context) {
	proc := &cronjob.Process{
		ID:     c.Query("pid"),
		JobID:  c.Query("job"),
		Group:  c.Query("group"),
		NodeID: c.Query("node"),
	}

	if proc.ID == "" || proc.JobID == "" || proc.Group == "" || proc.NodeID == "" {
		c.String(http.StatusBadRequest, "Invalid process info.")
		return
	}

	procKey := proc.Key()
	resp, err := cronjob.DefalutClient.Get(procKey)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(resp.Kvs) < 1 {
		c.String(http.StatusNotFound, "Porcess not found")
		return
	}

	var procVal = &cronjob.ProcessVal{}
	err = json.Unmarshal(resp.Kvs[0].Value, &procVal)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if procVal.Killed {
		c.String(http.StatusOK, "Killing process")
		return
	}

	procVal.Killed = true
	proc.ProcessVal = *procVal
	str, err := proc.Val()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = cronjob.DefalutClient.Put(procKey, str)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "Killing process")
}

type ProcFetchOptions struct {
	Groups  []string
	NodeIds []string
	JobIds  []string
}

func (opt *ProcFetchOptions) Match(proc *cronjob.Process) bool {
	if len(opt.Groups) > 0 && !InStringArray(proc.Group, opt.Groups) {
		return false
	}

	if len(opt.JobIds) > 0 && !InStringArray(proc.JobID, opt.JobIds) {
		return false

	}

	if len(opt.NodeIds) > 0 && !InStringArray(proc.NodeID, opt.NodeIds) {
		return false
	}

	return true
}

type processInfo struct {
	*cronjob.Process
	JobName string `json:"jobName"`
}

type ByProcTime []*processInfo

func (a ByProcTime) Len() int           { return len(a) }
func (a ByProcTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByProcTime) Less(i, j int) bool { return a[i].Time.After(a[j].Time) }
