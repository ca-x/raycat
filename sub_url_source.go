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

type subURLSourceProvider interface {
	UpdateUrlSub(ctx context.Context, urlSubs []string, fetchTimeoutSeconds int) ([]byte, error)
}

type subURLSource struct {
	weaver.Implements[subURLSourceProvider]
}

func (s *subURLSource) UpdateUrlSub(ctx context.Context, urlSubs []string, fetchTimeoutSeconds int) ([]byte, error) {

	buf := urlBuf.Get()
	defer urlBuf.Free(buf)

	st := stream.New().WithMaxGoroutines(len(urlSubs))
	for _, sub := range urlSubs {
		sub := sub
		if subContent, found := urlSubLru.Get(sub); found {
			buf.Write(subContent)
			continue
		}
		st.Go(func() stream.Callback {
			client := fetcher.NewClient(fetchTimeoutSeconds)
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
