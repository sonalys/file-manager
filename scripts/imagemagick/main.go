package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"

	flags "github.com/jessevdk/go-flags"
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
	output, err := cmd.CombinedOutput()
	log.Print(string(output))
	if err != nil {
		os.Exit(1)
	}
}
