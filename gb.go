package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gb"
	app.Usage = "Usage: gb [options] [http[s]://]hostname[:port]/path"
	app.Version = "0.0.1"
	app.HideHelp = true

	app.Flags = []cli.Flag{
		cli.Int64Flag{
			Name:  "requests, n",
			Usage: "Number of requests to perform",
		},
		cli.Int64Flag{
			Name:  "concurrency, c",
			Usage: "Number of multiple requests to make at a time",
		},
	}

	app.Action = gbAction

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.After = func(c *cli.Context) error {
		return nil
	}

	app.Run(os.Args)
}

func gbAction(c *cli.Context) {

	// グローバルオプション
	var requests = c.GlobalInt64("requests")
	var concurrency = c.GlobalInt64("concurrency")

	// パラメータ
	var paramFirst = ""
	if len(c.Args()) > 0 {
		paramFirst = c.Args().First() // c.Args()[0] と同じ意味
	}

	fmt.Printf("%d, %d\n", requests, concurrency)

	fmt.Printf("Hello world! %s\n", paramFirst)
}
