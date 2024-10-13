package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"raycat/internal/pkg/tinypool"
	"strconv"
	"strings"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metadata"
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
	config := app.configure.Get()
	subPublishPath, err := config.GetSubPublishPath(ctx)
	if err != nil {
		app.Logger(ctx).Warn("failed to get sub publish path,will use /subscribe as default", "err", err)
		subPublishPath = "/subscribe"
	}
	if !strings.HasPrefix(subPublishPath, "/") {
		subPublishPath = "/" + subPublishPath
	}
	http.Handle(subPublishPath, weaver.InstrumentHandlerFunc("subscribeCallCount", subShareHandlerApp(app)))
	app.Logger(ctx).Info("Listening on...", "address", app.lis)
	return http.Serve(app.lis, nil)
}

func subShareHandlerApp(app *app) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,profile-web-page-url,profile-update-interval")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		encoder := base64.NewEncoder(base64.StdEncoding, w)
		defer encoder.Close()

		buf := bufPool.Get()
		defer bufPool.Free(buf)

		authParamName, err := app.configure.Get().GetSubAuthParamName(context.Background())
		if err != nil {
			app.Logger(context.Background()).Warn("failed to get sub auth param name,use default name token", "err", err)
			authParamName = "token"
		}

		privateToken := r.URL.Query().Get(authParamName)

		ctx := context.Background()
		ctx = metadata.NewContext(ctx, map[string]string{"privateToken": privateToken})
		subFilePaths, _ := app.configure.Get().GetSubFilePaths(ctx)
		urlSubPaths, timeout, _ := app.configure.Get().GetUrlSubs(ctx)

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

		// handle for client experience
		responseOpt, err := app.configure.Get().GetResponseOption(context.Background())
		if err != nil {
			app.Logger(context.Background()).Warn("failed to get response option", "error", err)
		}
		if responseOpt != nil {
			if responseOpt.ProfileWebPage != "" {
				w.Header().Set("profile-web-page-url", responseOpt.ProfileWebPage)
			}
			if responseOpt.UpdateIntervalHours > 0 {
				w.Header().Set("profile-update-interval", strconv.Itoa(responseOpt.UpdateIntervalHours))
			}
		}

		_, err = io.Copy(encoder, buf)
		if err != nil {
			app.Logger(context.Background()).Error("failed to copy file to http response", "error", err)
		}
	}
}
