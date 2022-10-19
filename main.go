package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const VERSION  = `0.1`

var (
	ErrorLog = log.New(os.Stderr, `error#`, log.Lshortfile)
	DebugLog = log.New(os.Stdout, `debug#`, log.Lshortfile)
)

func helpText() {
	fmt.Println(`# https://github.com/vvampirius/http-slave`)
	flag.PrintDefaults()
}

func main() {
	help := flag.Bool("h", false, "print this help")
	ver := flag.Bool("v", false, "show version")
	contactUrl := flag.String("u", os.Getenv(`URL`), "contact url")
	contactInterval := flag.Int("i", 60, "contact interval (second)")
	flag.Parse()

	if *help {
		helpText()
		os.Exit(0)
	}

	if *ver {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	fmt.Printf("Starting version %s...\n", VERSION)
	for {
		if task, err := GetTask(*contactUrl); err == nil {
			go func() {
				data, exitCode, spent, err := ExecTask(task)
				if err != nil {
					data = []byte(err.Error())
					exitCode = 255
				}
				RespondTask(task.RespondUrl, data, exitCode, spent)
			}()
			if task.ImmediatelyNext { continue }
		}
		time.Sleep(time.Second * time.Duration(*contactInterval))
	}

}
