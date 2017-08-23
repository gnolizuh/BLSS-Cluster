package main

import (
	"github.com/urfave/cli"
	"os"
	"fmt"
	"net/http"
	"time"
	"log"
	"./utils"
	"./session"
)

var (
	VERSION = "1.0"
)

func main()  {
	myApp := cli.NewApp()
	myApp.Name = "BLSS-cluster"
	myApp.Usage = "Make BLSS to be a dynamically expanding cluster"
	myApp.Version = VERSION
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "listen,l",
			Value: ":19000",
			Usage: " listen address",
		},
		cli.StringFlag{
			Name:  "method",
			Value: "POST",
			Usage: "methods: GET/POST",
		},
	}
	myApp.Action = func(c *cli.Context) error {
		config := utils.Config{}
		config.Listen = c.String("listen")
		config.Method = c.String("method")

		fmt.Printf("listen: %s method: %s\n", config.Listen, config.Method)

		m := session.NewManager(&config)

		s := &http.Server{
			Addr:           config.Listen,
			Handler:        m,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		log.Fatal(s.ListenAndServe())

		return nil
	}

	myApp.Run(os.Args)
}