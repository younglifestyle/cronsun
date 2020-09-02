package web

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/longcron/cronjob"
)

func GetVersion(c *gin.Context) {
	c.String(http.StatusOK, cronjob.Version)
}

func initRouters() (engine *gin.Engine, err error) {
	engine = gin.Default()

	jobHandler := &Job{}
	nodeHandler := &Node{}
	jobLogHandler := &JobLog{}
	infoHandler := &Info{}

	subroutine := engine.Group("/v1")
	{
		subroutine.GET("/version", GetVersion)

		// get job list
		subroutine.GET("/jobs", jobHandler.GetList)
		// get all job group list
		subroutine.GET("/all/job/groups", jobHandler.GetGroups)

		// create/update a job
		subroutine.PUT("/job", jobHandler.UpdateJob)
		// pause/start
		subroutine.POST("/job/:group/:id", jobHandler.ChangeJobStatus)
		// batch pause/start
		subroutine.POST("/jobs/:op", jobHandler.BatchChangeJobStatus)
		// get a job
		subroutine.GET("/job/:group/:id", jobHandler.GetJob)
		// remove a job
		subroutine.DELETE("/job/:group/:id", jobHandler.DeleteJob)
		// 获取执行该任务的节点
		subroutine.GET("/job/:group/:id/nodes", jobHandler.GetJobNodes)

		// put once task to Node execute
		subroutine.PUT("/job/:group/:id/execute", jobHandler.JobExecute)

		// query executing job
		subroutine.GET("/job-executing", jobHandler.GetExecutingJob)

		// kill an executing job
		subroutine.DELETE("/job-executing", jobHandler.KillExecutingJob)

		// get job log list
		subroutine.GET("/logs", jobLogHandler.GetList)
		// get job log
		subroutine.GET("/log/:id", jobLogHandler.GetDetail)
		// 获取所有Node
		subroutine.GET("/nodes", nodeHandler.GetNodes)
		// 删除节点
		subroutine.DELETE("/node", nodeHandler.DeleteNode)
		// get node group list
		subroutine.GET("/node/groups", nodeHandler.GetGroups)
		// get a node group by group id
		subroutine.GET("/node/group", nodeHandler.GetGroupByGroupId)
		// create/update a node group
		subroutine.PUT("/node/group", nodeHandler.UpdateGroup)
		// delete a node group
		subroutine.DELETE("/node/group", nodeHandler.DeleteGroup)

		subroutine.GET("/info/overview", infoHandler.Overview)
	}

	return engine, nil
}
