package main

import (
	"context"
	"github.com/ServiceWeaver/weaver"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/sourcegraph/conc/stream"
	"raycat/internal/pkg/fetcher"
	"raycat/internal/pkg/tinypool"
	"time"
)

var _ subURLSourceProvider = (*subURLSource)(nil)

var (
	urlSubLru = expirable.NewLRU[string, []byte](20, nil, time.Hour*4)
	urlBuf    = tinypool.New(tinypool.BufReset)
)

type subUrlSourceConfig struct {
	weaver.AutoMarshal
	TimeoutSeconds int      `toml:"timeout_seconds"`
	UrlSubs        []string `toml:"url_subs"`
}

type subURLSourceProvider interface {
	UpdateUrlSub(ctx context.Context) ([]byte, error)
}

type subURLSource struct {
	weaver.Implements[subURLSourceProvider]
	weaver.WithConfig[subUrlSourceConfig]
}

func (s *subURLSource) UpdateUrlSub(ctx context.Context) ([]byte, error) {
	config := s.Config()

	buf := urlBuf.Get()
	defer urlBuf.Free(buf)

	st := stream.New().WithMaxGoroutines(len(config.UrlSubs))
	for _, sub := range config.UrlSubs {
		sub := sub
		if subContent, found := urlSubLru.Get(sub); found {
			buf.Write(subContent)
			continue
		}
		st.Go(func() stream.Callback {
			client := fetcher.NewClient(config.TimeoutSeconds)
			content, err := client.Fetch(sub)
			if err != nil {
				s.Logger(ctx).Error("failed to fetch url sub from source", "url", sub, "error", err)
				return func() {
					urlSubLru.Remove(sub)
				}
			}
			return func() {
				urlSubLru.Add(sub, content)
				buf.Write(content)
			}
		})
	}
	st.Wait()
	b := buf.Bytes()
	return b, nil
}
