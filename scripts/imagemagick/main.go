package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/model"
)

func main() {
	var data model.UploadData
	var config model.ScriptConfiguration

	args := os.Args
	json.Unmarshal([]byte(args[1]), &data)
	json.Unmarshal([]byte(args[2]), &config)

	t1 := time.Now()
	from := fmt.Sprintf("/buffer/%s", data.GetFullName())
	data.Extension = "heic"
	to := fmt.Sprintf("/buffer/%s", data.GetFullName())
	cmd := exec.Command("convert", append(config.Parameters, from, to)...)

	convertOutput, err := cmd.CombinedOutput()
	if err == nil && os.Remove(from) != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

	output := model.ScriptOutput{
		ExecutionTime: time.Since(t1),
		LastRun:       t1,
		Log:           string(convertOutput),
	}

	data.Metadata[config.Name] = output

	serialized, err := json.Marshal(data)
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

	fmt.Print(string(serialized))
}
