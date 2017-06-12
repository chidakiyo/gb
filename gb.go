package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"net/http"
	"os"
	"time"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
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

	const (
		N = 400 // 全体数
		M = 200 // 1秒あたりの処理制限
	)

	cc := make(chan int, 1000)
	go func() {
		ctx := context.Background()
		n := rate.Every(time.Second/M)
		l := rate.NewLimiter(n, M)
		for i := 0; i < N; i++ {
			if err := l.Wait(ctx); err != nil {
				fmt.Errorf("fatal; %s", err)
			}
			go get("https://maitto-app.appspot.com/")
			cc <- i
		}
		close(cc)
	}()
	for n := range cc {
		fmt.Println(n)
	}
}

func get(url string) Result {
	req, _ := http.NewRequest("GET", url, nil)

	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	//byteArray, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(byteArray))
	status := resp.StatusCode
	fmt.Printf("status : %d\n", status)

	return Result{
		Status: status,
	}

}

type Result struct {
	Status int
}
