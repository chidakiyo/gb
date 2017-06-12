package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"

	"go.uber.org/ratelimit"
	"time"
	"net/http"
	"io/ioutil"
)

func main() {

	app := createNewApp()

	app.Action = gbAction

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.After = func(c *cli.Context) error {
		return nil
	}

	app.Run(os.Args)
}

func createNewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "gb"
	app.Usage = "Usage: gb [options] [http[s]://]hostname[:port]/path"
	app.Version = "0.0.1"
	app.HideHelp = true

	app.Flags = []cli.Flag{
		// -n オプション : 総リクエスト数
		cli.Int64Flag{
			Name:  "requests, n",
			Usage: "Number of requests to perform",
		},
		// -c オプション : 同時接続数（並列数）
		cli.Int64Flag{
			Name:  "concurrency, c",
			Usage: "Number of multiple requests to make at a time",
		},
	}
	return app
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

	rl := ratelimit.New(2) // per second

	prev := time.Now()
	for i := 0; i < 100; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
}

type RequestResult struct {
	Body string
	IsError bool
	ErrorMsg string
}

func getRequest(url string, resultChan chan RequestResult) {

	// request
	res, err := http.Get(url)

	// ハンドリング
	go func() {
		b, _ := ioutil.ReadAll(res.Body) // TODO error
		result := RequestResult{}
		result.Body = string(b)
		if err != nil {
			result.IsError = true
			result.ErrorMsg = err.Error()
		}
		resultChan <- result
	} ()
}





