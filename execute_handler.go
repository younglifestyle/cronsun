package cronsun

import (
	"errors"
	"fmt"
	"github.com/shunfei/cronsun/log"
	"os"
	"sync"
)

var scriptMap = map[string]string{
	"SHELL":  ".sh",
	"PYTHON": ".py",
}

var scriptCmd = map[string]string{
	"SHELL":  "bash",
	"PYTHON": "python",
}

type ExecuteHandler interface {
	ParseJob(job *Job) (err error)
	Execute(jobId int32, glueType string, runParam *Job) error
}

type ScriptHandler struct {
	sync.RWMutex
}

func (j *Job) parseJob() error {
	suffix, ok := scriptMap[j.CmdType]
	if !ok {
		log.Infof("j.CmdType : %s", j.CmdType)
		return ErrSecurityInvalidCmd
	}

	path := FileSourcePath + j.ID + "_" + fmt.Sprintf("%d", j.UpdateTime) + suffix
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		j.Lock()
		defer j.Unlock()
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0750)
		if err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(FileSourcePath, os.ModePerm)
			if err == nil {
				file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0750)
				if err != nil {
					return err
				}
			}
		}

		if file != nil {
			defer file.Close()
			res, err := file.Write([]byte(j.Command))
			if err != nil {
				return err
			}
			if res <= 0 {
				return errors.New("write script file failed")
			}
		}
	}

	j.fileExecPath = path

	return nil
}
