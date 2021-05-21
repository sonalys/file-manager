package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/model"
)

type Arguments struct {
	data model.UploadData `long:"data" required:"true"`
}

func main() {
	var opts Arguments
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	t1 := time.Now()
	path := opts.FromPath
	cmd := exec.Command("convert", "-quality", strconv.FormatInt(opts.Quality, 10), opts.FromPath, opts.ToPath)
	convertOutput, err := cmd.CombinedOutput()
	if err == nil && os.Remove(opts.FromPath) != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	} else {
		path = opts.ToPath
	}

	output := model.ScriptOutput{
		MovedTo:       path,
		ExecutionTime: time.Since(t1),
		LastRun:       t1,
		Log:           string(convertOutput),
	}

	serialized, err := json.Marshal(output)
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

	fmt.Print(string(serialized))
}
