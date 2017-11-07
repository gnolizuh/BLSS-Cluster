package main

import (
	"github.com/urfave/cli"
	"github.com/go-redis/redis"
	"os"
	"net/http"
	"time"
	"log"
)

var (
	VERSION = "1.0"
)

func checkError(err error) {
	if err != nil {
		log.Printf("%+v\n", err)
		os.Exit(-1)
	}
}

func newRedis(config *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis,
		Password: config.Password,
		DB:       0,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}

func action(c *cli.Context) error {
	config := Config{}
	config.Listen = c.String("listen")
	config.Method = c.String("method")
	config.Log = c.String("log")
	config.Redis = c.String("redis")
	config.Expire = c.Int64("expire")
	config.Password = c.String("password")

	if config.Log != "" {
		f, err := os.OpenFile(config.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		checkError(err)
		defer f.Close()
		log.SetOutput(f)
	}

	log.Printf("%d#0: %s/%s", os.Getpid(), c.App.Name, c.App.Version)

	m := NewManager(&config)

	if client, err := newRedis(&config); err != nil {
		log.Printf("%d#0: redis(%s) connect failed, err: %s ", os.Getpid(), config.Redis, err)
		return nil
	} else {
		m.redisClient = client
	}

	s := &http.Server{
		Addr:           config.Listen,
		Handler:        m,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("%d#0: listening on: %s", os.Getpid(), s.Addr)
	log.Printf("%d#0: method: %s", os.Getpid(), config.Method)
	log.Printf("%d#0: redis addr: %s", os.Getpid(), config.Redis)

	log.Fatal(s.ListenAndServe())

	return nil
}

func main()  {
	myApp := cli.NewApp()
	myApp.Name = "blss-clu"
	myApp.Usage = "Make BLSS to be a dynamically expanding cluster"
	myApp.Version = VERSION
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen",
			Value: ":19000",
			Usage: "listen address",
		},
		cli.StringFlag{
			Name:  "method",
			Value: "POST",
			Usage: "HTTP request method: GET/POST",
		},
		cli.StringFlag{
			Name:  "log",
			Value: "",
			Usage: "specify a log file to output, default goes to stderr",
		},
		cli.StringFlag{
			Name:  "level",
			Value: "INFO",
			Usage: "specify a log level, default goes to info",
		},
		cli.StringFlag{
			Name:  "redis",
			Value: "localhost:6379",
			Usage: "redis addr, default set to localhost:6379",
		},
		cli.IntFlag{
			Name:  "expire",
			Value: 60,
			Usage: "set key expire timeout, default 0 second",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "123456",
			Usage: "redis password, default set to 123456",
		},
	}
	myApp.Action = action
	myApp.Run(os.Args)
}
