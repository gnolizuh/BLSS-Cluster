package main

import (
	"github.com/urfave/cli"
	"os"
	"fmt"
	"net/http"
	"time"
	"log"
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
			Usage: "listen address",
		},
		cli.StringFlag{
			Name:  "method,m",
			Value: "POST",
			Usage: "HTTP request method: GET/POST",
		},
	}
	myApp.Action = func(c *cli.Context) error {
		config := Config{}
		config.Listen = c.String("listen")
		config.Method = c.String("method")

		fmt.Printf("%s %d#0: %s/%s\n",
			time.Now().Format("2006/01/02 15:04:05"), os.Getpid(),
			myApp.Name, myApp.Version)

		m := NewManager(&config)

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
