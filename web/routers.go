package web

import (
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"

	"github.com/shunfei/cronsun"
)

func GetVersion(W http.ResponseWriter, R *http.Request) {
	outJSON(W, cronsun.Version)
}

func initRouters() (s *http.Server, err error) {
	jobHandler := &Job{}
	nodeHandler := &Node{}
	jobLogHandler := &JobLog{}
	infoHandler := &Info{}
	configHandler := &Configuration{}

	r := mux.NewRouter()
	subrouter := r.PathPrefix("/v1").Subrouter()
	subrouter.HandleFunc("/version", GetVersion).Methods("GET")

	// get job list
	subrouter.HandleFunc("/jobs", jobHandler.GetList).Methods("GET")
	// get a job group list
	subrouter.HandleFunc("/job/groups", jobHandler.GetGroups).Methods("GET")

	// create/update a job
	subrouter.HandleFunc("/job", jobHandler.UpdateJob).Methods("PUT")
	// pause/start
	subrouter.HandleFunc("/job/{group}-{id}", jobHandler.ChangeJobStatus).Methods("POST")
	// batch pause/start
	subrouter.HandleFunc("/jobs/{op}", jobHandler.BatchChangeJobStatus).Methods("POST")
	// get a job
	subrouter.HandleFunc("/job/{group}-{id}", jobHandler.GetJob).Methods("GET")
	// remove a job
	subrouter.HandleFunc("/job/{group}-{id}", jobHandler.DeleteJob).Methods("DELETE")
	// 获取执行该任务的Job
	subrouter.HandleFunc("/job/{group}-{id}/nodes", jobHandler.GetJobNodes).Methods("GET")

	// put once task to Node execute
	subrouter.HandleFunc("/job/{group}-{id}/execute", jobHandler.JobExecute).Methods("PUT")

	// query executing job
	subrouter.HandleFunc("/job/executing", jobHandler.GetExecutingJob).Methods("GET")

	// kill an executing job
	subrouter.HandleFunc("/job/executing", jobHandler.KillExecutingJob).Methods("DELETE")

	// get job log list
	subrouter.HandleFunc("/logs", jobLogHandler.GetList).Methods("GET")
	// get job log
	subrouter.HandleFunc("/log/{id}", jobLogHandler.GetDetail).Methods("GET")
	// 获取所有Node
	subrouter.HandleFunc("/nodes", nodeHandler.GetNodes).Methods("GET")
	// 删除节点
	subrouter.HandleFunc("/node/{ip}", nodeHandler.DeleteNode).Methods("DELETE")
	// get node group list
	subrouter.HandleFunc("/node/groups", nodeHandler.GetGroups).Methods("GET")
	// get a node group by group id
	subrouter.HandleFunc("/node/group/{id}", nodeHandler.GetGroupByGroupId).Methods("GET")
	// create/update a node group
	subrouter.HandleFunc("/node/group", nodeHandler.UpdateGroup).Methods("PUT")
	// delete a node group
	subrouter.HandleFunc("/node/group/{id}", nodeHandler.DeleteGroup).Methods("DELETE")

	subrouter.HandleFunc("/info/overview", infoHandler.Overview).Methods("GET")

	subrouter.HandleFunc("/configurations", configHandler.Configuratios).Methods("GET")

	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", newEmbeddedFileServer("", "index.html")))
	r.NotFoundHandler = NewBaseHandler(notFoundHandler)

	s = &http.Server{
		Handler: r,
	}
	return s, nil
}

type embeddedFileServer struct {
	Prefix    string
	IndexFile string
}

func newEmbeddedFileServer(prefix, index string) *embeddedFileServer {
	index = strings.TrimLeft(index, "/")
	return &embeddedFileServer{Prefix: prefix, IndexFile: index}
}

func (s *embeddedFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fp := path.Clean(s.Prefix + r.URL.Path)
	if fp == "." {
		fp = ""
	} else {
		fp = strings.TrimLeft(fp, "/")
	}

	b, err := Asset(fp)
	if err == nil {
		w.Write(b)
		return
	}

	if len(fp) > 0 {
		fp += "/"
	}
	fp += s.IndexFile

	// w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	// w.Header().Set("Expires", "0")

	b, err = Asset(fp)
	if err == nil {
		w.Write(b)
		return
	}

	_notFoundHandler(w, r)
}

func notFoundHandler(c *Context) {
	_notFoundHandler(c.W, c.R)
}

func _notFoundHandler(w http.ResponseWriter, r *http.Request) {
	const html = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>404 page not found</title>
</head>
<body>
    The page you are looking for is not found. Redirect to <a href="/ui/">Dashboard</a> after <span id="s">5</span> seconds.
</body>
<script type="text/javascript">
var s = 5;
setInterval(function(){
    s--;
    document.getElementById('s').innerText = s;
    if (s === 0) location.href = '/ui/';
}, 1000);
</script>
</html>`

	if strings.HasPrefix(strings.TrimLeft(r.URL.Path, "/"), "v1") {
		outJSONWithCode(w, http.StatusNotFound, "Api not found.")
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(html))
	}
}
