package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"raycat/internal/pkg/tinypool"

	"github.com/ServiceWeaver/weaver"
)

//go:generate weaver generate ./...

var (
	bufPool = tinypool.New[bytes.Buffer](tinypool.BufReset)
)

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

// app implements the main component, the entry point to a Service Weaver app.
type app struct {
	weaver.Implements[weaver.Main]
	configure weaver.Ref[subConfigureProvider]
	fileSub   weaver.Ref[subFileSourceProvider]
	urlSub    weaver.Ref[subURLSourceProvider]
	lis       weaver.Listener `weaver:"lis"`
}

// serve serves HTTP traffic.
func serve(ctx context.Context, app *app) error {
	http.HandleFunc("/x", subShareHandlerApp(app))
	app.Logger(ctx).Info("Listening on...", "address", app.lis)
	return http.Serve(app.lis, nil)
}

func subShareHandlerApp(app *app) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := base64.NewEncoder(base64.StdEncoding, w)
		defer encoder.Close()

		buf := bufPool.Get()
		defer bufPool.Free(buf)

		privateToken := r.URL.Query().Get("token")

		subFilePaths, _ := app.configure.Get().GetSubFilePaths(context.Background(), privateToken)
		urlSubPaths, timeout, _ := app.configure.Get().GetUrlSubs(context.Background(), privateToken)

		fileSub, err := app.fileSub.Get().UpdateFileSub(context.Background(), subFilePaths)
		if err != nil {
			app.Logger(context.Background()).Error("failed to get file sub update", "error", err)
		}

		urlSub, err := app.urlSub.Get().UpdateUrlSub(context.Background(), urlSubPaths, timeout)
		if err != nil {
			app.Logger(context.Background()).Error("failed to get url sub update", "error", err)
		}

		if len(fileSub) > 0 {
			if _, err = buf.Write(fileSub); err != nil {
				app.Logger(context.Background()).Error("failed to write file sub to buffer", "error", err)
			}
		}

		if len(urlSub) > 0 {
			if _, err = buf.Write(urlSub); err != nil {
				app.Logger(context.Background()).Error("failed to write url sub to buffer", "error", err)
			}
		}

		_, err = io.Copy(encoder, buf)
		if err != nil {
			app.Logger(context.Background()).Error("failed to copy file to http response", "error", err)
		}
	}
}
