package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	//handler(nil, nil)
	//handler2(nil, nil)
	//handler3(nil, nil)
	handler4(nil, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 新たにgoroutineを生成してバックエンドにリクエストを投げる
	// 結果をerror channelに入れる
	errCh := make(chan error, 1)
	go func() {
		errCh <- request()
	}()

	// error channelにリクエストの結果が返ってくるのを待つ
	select {
	case err := <-errCh:
		if err != nil {
			log.Println("failed:", err)
			return
		}
	}

	log.Println("success")
}

func request() error {
	time.Sleep(3 * time.Second)
	return nil
}

func handler2(w http.ResponseWriter, r *http.Request) {
	// 新たにgoroutineを生成してバックエンドにリクエストを投げる
	// 結果をerror channelに入れる
	errCh := make(chan error, 1)
	go func() {
		errCh <- request()
	}()

	// error channelにリクエストの結果が返ってくるのを待つ
	select {
	case err := <-errCh:
		if err != nil {
			log.Println("failed:", err)
			return
		}

		// Timeout（2秒）を設定する．
		// 例えばしばらく経ってから再度リクエストをするように
		// レスポンスを返す．
	case <-time.After(2 * time.Second):
		log.Println("failed: timeout")
		return
	}

	log.Println("success")
}

func handler3(w http.ResponseWriter, r *http.Request) {
	// handlerからrequestをキャンセルするためのchannelを準備する
	doneCh := make(chan struct{}, 1)

	errCh := make(chan error, 1)
	go func() {
		errCh <- request2(doneCh)
	}()

	// 別途goroutineを準備してTimeoutを設定する
	go func() {
		<-time.After(2 * time.Second)
		// Timeout後にdoneChをクローズする
		// 参考: https://blog.golang.org/pipelines
		close(doneCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Println("failed:", err)
			return
		}
	}

	log.Println("success")
}

func request2(doneCh chan struct{}) error {
	//tr := &http.Transport{}

	// req, err := http.NewRequest("POST", backendService, nil)

	// 新たにgoroutineを生成して実際のリクエストを行う
	// 結果はerror channelに投げる
	errCh := make(chan error, 1)
	go func() {
		errCh <- request()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}

		// doneChはhandlerからのキャンセル シグナル（close(doneCh)）
		// を待ち受ける
	case <-doneCh:
		// キャンセルが実行されたら適切にリクエストを停止して
		// エラーを返す．
		//tr.CancelRequest(req)
		<-errCh
		return fmt.Errorf("canceled")
	}

	return nil
}

func handler4(w http.ResponseWriter, r *http.Request) {
	// 2秒でTimeoutするcontextを生成する
	// cancelを実行することでTimeout前にキャンセルを実行することができる
	//
	// また後述するようにGo1.7ではnet/httpパッケージでcontext
	// を扱えるようになる．例えば*http.Requestからそのリクエストの
	// contextを取得できる．
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		log.Fatalln("cancel defer!")
		cancel()
	}()

	errCh := make(chan error, 1)
	go func() {
		errCh <- request3(ctx)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Println("failed:", err)
			return
		}
	}

	log.Println("success")
}

func request3(ctx context.Context) error {
	//tr := &http.Transport{}
	//client := &http.Client{Transport: tr}

	//req, err := http.NewRequest("POST", backendService, nil)
	//if err != nil {
	//	return err
	//}

	// 新たにgoroutineを生成して実際のリクエストを行う
	// 結果はerror channelに投げる
	errCh := make(chan error, 1)
	go func() {
		//_, err := client.Do(req)
		//errCh <- err
		errCh <- request()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}

		// Timeoutが発生する，もしくはCancelが実行されると
		// Channelが返る
	case <-ctx.Done():
		//tr.CancelRequest(req)
		<-errCh
		return ctx.Err()
	}

	return nil
}
