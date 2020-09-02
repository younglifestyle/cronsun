package cronjob

import (
	"errors"
	"fmt"
	"github.com/longcron/cronjob/log"
	"os"
)

var scriptMap = map[string]string{
	"SHELL":  ".sh",
	"PYTHON": ".py",
}

var scriptCmd = map[string]string{
	"SHELL":  "bash",
	"PYTHON": "python",
}

func (j *Job) parseJob(serverOrNode string) error {
	suffix, ok := scriptMap[j.CmdType]
	if !ok {
		log.Infof("不支持的命令 : %s", j.CmdType)
		return ErrSecurityInvalidCmd
	}

	if serverOrNode == "node" { // 仅在Node上进行真正写文件操作
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
	}

	return nil
}
