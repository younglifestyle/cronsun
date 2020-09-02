package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	v3 "github.com/coreos/etcd/clientv3"
	"gopkg.in/mgo.v2/bson"

	"github.com/longcron/cronsun"
	"github.com/longcron/cronsun/conf"
	"github.com/longcron/cronsun/log"
)

type Node struct{}

var ngKeyDeepLen = len(conf.Config.Group)

func (n *Node) UpdateGroup(c *gin.Context) {
	g := cronsun.Group{}
	err := c.BindJSON(&g)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var successCode = http.StatusOK
	g.ID = strings.TrimSpace(g.ID)
	if len(g.ID) == 0 {
		successCode = http.StatusCreated
		g.ID = cronsun.NextID()
	}

	if err = g.Check(); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// @TODO modRev
	var modRev int64 = 0
	if _, err = g.Put(modRev); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(successCode, "")
}

func (n *Node) GetGroups(c *gin.Context) {
	list, err := cronsun.GetNodeGroups()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

func (n *Node) GetGroupByGroupId(c *gin.Context) {
	gid := c.Query("id")
	g, err := cronsun.GetGroupById(gid)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if g == nil {
		c.String(http.StatusNotFound, "not found")
		return
	}
	c.JSON(http.StatusOK, g)
}

func (n *Node) DeleteGroup(c *gin.Context) {
	groupId := c.Query("id")
	if len(groupId) == 0 {
		c.String(http.StatusBadRequest, "empty node ground id.")
		return
	}

	_, err := cronsun.DeleteGroupById(groupId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	gresp, err := cronsun.DefalutClient.Get(conf.Config.Cmd, v3.WithPrefix())
	if err != nil {
		errstr := fmt.Sprintf("failed to fetch jobs from etcd after deleted node group[%s]: %s", groupId, err.Error())
		log.Errorf(errstr)
		c.String(http.StatusInternalServerError, errstr)
		return
	}

	// update rule's node group
	for i := range gresp.Kvs {
		job := cronsun.Job{}
		err = json.Unmarshal(gresp.Kvs[i].Value, &job)
		key := string(gresp.Kvs[i].Key)
		if err != nil {
			log.Errorf("failed to unmarshal job[%s]: %s", key, err.Error())
			continue
		}

		update := false
		for j := range job.Rules {
			var ngs []string
			for _, gid := range job.Rules[j].GroupIDs {
				if gid != groupId {
					ngs = append(ngs, gid)
				}
			}
			if len(ngs) != len(job.Rules[j].GroupIDs) {
				job.Rules[j].GroupIDs = ngs
				update = true
			}
		}

		if update {
			v, err := json.Marshal(&job)
			if err != nil {
				log.Errorf("failed to marshal job[%s]: %s", key, err.Error())
				continue
			}
			_, err = cronsun.DefalutClient.PutWithModRev(key, string(v), gresp.Kvs[i].ModRevision)
			if err != nil {
				log.Errorf("failed to update job[%s]: %s", key, err.Error())
				continue
			}
		}
	}

	c.String(http.StatusOK, "")
}

func (n *Node) GetNodes(c *gin.Context) {
	nodes, err := cronsun.GetNodes()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	gresp, err := cronsun.DefalutClient.Get(conf.Config.Node, v3.WithPrefix(), v3.WithKeysOnly())
	if err == nil {
		connecedMap := make(map[string]bool, gresp.Count)
		for i := range gresp.Kvs {
			k := string(gresp.Kvs[i].Key)
			index := strings.LastIndexByte(k, '/')
			connecedMap[k[index+1:]] = true
		}

		for i := range nodes {
			nodes[i].Connected = connecedMap[nodes[i].ID]
		}
	} else {
		log.Errorf("failed to fetch key[%s] from etcd: %s", conf.Config.Node, err.Error())
	}

	c.JSON(http.StatusOK, nodes)
}

// DeleteNode force remove node (by ip) which state in offline or damaged.
func (n *Node) DeleteNode(c *gin.Context) {
	nodeId := c.Query("id")
	if len(nodeId) == 0 {
		c.String(http.StatusBadRequest, "node nodeId is required.")
		return
	}

	resp, err := cronsun.DefalutClient.Get(conf.Config.Node + nodeId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(resp.Kvs) > 0 {
		c.String(http.StatusBadRequest, "can not remove a running node.")
		return
	}

	err = cronsun.RemoveNode(bson.M{"_id": nodeId})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// remove node from group
	var errmsg = "failed to remove node %s from groups, please remove it manually: %s"
	resp, err = cronsun.DefalutClient.Get(conf.Config.Group, v3.WithPrefix())
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf(errmsg, nodeId, err.Error()))
		return
	}

	for i := range resp.Kvs {
		g := cronsun.Group{}
		err = json.Unmarshal(resp.Kvs[i].Value, &g)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf(errmsg, nodeId, err.Error()))
			return
		}

		var nids = make([]string, 0, len(g.NodeIDs))
		for _, nid := range g.NodeIDs {
			if nid != nodeId {
				nids = append(nids, nid)
			}
		}
		g.NodeIDs = nids

		if _, err = g.Put(resp.Kvs[i].ModRevision); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf(errmsg, nodeId, err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, "is ok")
}
