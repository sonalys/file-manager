package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/model"
)

func main() {
	var data model.UploadData
	var config model.ScriptConfiguration

	args := os.Args

	if err := json.Unmarshal([]byte(args[1]), &data); err != nil {
		logrus.Error(errors.Wrap(err, fmt.Sprintf("failed to read data: %s", args[1])))
		os.Exit(1)
	}
	if err := json.Unmarshal([]byte(args[2]), &config); err != nil {
		logrus.Error(errors.Wrap(err, fmt.Sprintf("failed to read config: %s", args[2])))
		os.Exit(1)
	}

	t1 := time.Now()
	file := fmt.Sprintf("/buffer/%s", data.GetFullName())
	cmd := exec.Command("python3", append([]string{"/app/main.py", file}, config.Parameters...)...)

	scriptOutput, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error(errors.Wrap(err, fmt.Sprintf("failed to execute script: %s", cmd.String())))
		os.Exit(1)
	}

	var tags []interface{}
	reader, err := os.Open("/dump.json")
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to open dump.json"))
		os.Exit(1)
	}
	dump, err := ioutil.ReadAll(reader)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to read dump.json"))
		os.Exit(1)
	}

	err = json.Unmarshal(dump, &tags)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to deserialize dump.json"))
		os.Exit(1)
	}

	output := model.ScriptOutput{
		Metadata:      tags,
		ExecutionTime: time.Since(t1),
		LastRun:       t1,
		Log:           string(scriptOutput),
	}

	if data.Metadata == nil {
		data.Metadata = make(map[string]model.ScriptOutput)
	}

	data.Metadata[config.Name] = output

	serialized, err := json.Marshal(data)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to serialize response"))
		os.Exit(1)
	}

	fmt.Print(string(serialized))
}
