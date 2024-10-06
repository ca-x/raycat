package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
)

//go:generate weaver generate ./...

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

// app implements the main component, the entry point to a Service Weaver app.
type app struct {
	weaver.Implements[weaver.Main]
	fileSub weaver.Ref[subFileSourceProvider]
	lis     weaver.Listener `weaver:"lis"`
}

// serve serves HTTP traffic.
func serve(ctx context.Context, app *app) error {
	http.HandleFunc("/x", subShareHandlerApp(app))
	app.Logger(ctx).Info("Listening on...", "address", app.lis)
	return http.Serve(app.lis, nil)
}

func subShareHandlerApp(app *app) func(w http.ResponseWriter, _ *http.Request) {
	fileSub, err := app.fileSub.Get().UpdateSub(context.Background())
	if err != nil {
		app.Logger(context.Background()).Error("failed to get file sub update", "error", err)
		return nil
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		encoder := base64.NewEncoder(base64.StdEncoding, w)
		defer encoder.Close()
		_, err = io.Copy(encoder, bytes.NewReader(fileSub))
		if err != nil {
			app.Logger(context.Background()).Error("failed to copy file sub to http response", "error", err)
		}
	}
}
