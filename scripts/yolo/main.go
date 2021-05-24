package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/model"
	"gopkg.in/yaml.v2"
)

func read(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("couldn't read the file %s", path))
	}
	return content, nil
}

func readYaml(path string, dest interface{}) error {
	dump, err := read(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(dump, dest)
	if err != nil {
		return err
	}

	return nil
}

func readTxt(path string) ([][]string, error) {
	dump, err := read(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(dump), "\n")
	buffer := make([][]string, len(lines))
	for i := range lines {
		buffer[i] = strings.Split(lines[i], " ")
	}
	return buffer, nil
}

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
	parameters := []string{"detect.py", "--save-conf", "--save-txt", "--source", file, "--img-size", "1280"}
	parameters = append(parameters, config.Parameters...)
	cmd := exec.Command("python3", parameters...)

	err := cmd.Run()
	if err != nil {
		logrus.Error(errors.Wrap(err, fmt.Sprintf("failed to execute script: %s", cmd.String())))
		os.Exit(1)
	}

	scriptPath := fmt.Sprintf("/usr/src/app/runs/detect/exp/labels/%s.txt", data.Filename)
	tags, err := readTxt(scriptPath)
	if err != nil {
		logrus.Error(errors.Wrap(err, fmt.Sprintf("failed to deserialize '%s'", scriptPath)))
		os.Exit(1)
	}

	var dataset struct {
		Names []string `yaml:"names"`
	}

	err = readYaml("/usr/src/app/data/coco128.yaml", &dataset)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to deserialize coco128.yaml"))
		os.Exit(1)
	}

	for i := range tags {
		classNum, _ := strconv.Atoi(tags[i][0])
		if classNum > len(dataset.Names) {
			logrus.Errorf("class mismatch: has %d classes, want class[%d]. data: %v", len(dataset.Names), classNum, tags[i])
			os.Exit(1)
		}
		tags[i][0] = dataset.Names[classNum]
	}

	output := model.ScriptOutput{
		Metadata:      tags,
		ExecutionTime: time.Since(t1),
		LastRun:       t1,
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
