package cronsun

var scriptMap = map[string]string{
	"GLUE_SHELL":      ".sh",
	"GLUE_PYTHON":     ".py",
	"GLUE_PHP":        ".php",
	"GLUE_NODEJS":     ".js",
	"GLUE_POWERSHELL": ".ps1",
}

var scriptCmd = map[string]string{
	"GLUE_SHELL":      "bash",
	"GLUE_PYTHON":     "python",
	"GLUE_PHP":        "php",
	"GLUE_NODEJS":     "node",
	"GLUE_POWERSHELL": "powershell",
}

type ExecuteHandler interface {
	ParseJob(job *Job) (err error)
	Execute(jobId int32, glueType string, runParam *Job) error
}
