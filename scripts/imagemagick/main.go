package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/model"
)

type Arguments struct {
	FromPath string `long:"from" required:"true"`
	ToPath   string `long:"to" required:"true"`
	Quality  int64  `long:"quality" short:"q" default:"85"`
}

func main() {
	var opts Arguments
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	cmd := exec.Command("convert", "-quality", strconv.FormatInt(opts.Quality, 10), opts.FromPath, opts.ToPath)
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	if err := os.Remove(opts.FromPath); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

	output = model.ScriptOutput{
		MovedTo: opts.ToPath,
	}

	serialized, err := json.Marshal(output)
	if err := os.Remove(opts.FromPath); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

	fmt.Print(serialized)
}
