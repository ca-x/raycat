package main

import (
	"context"
	"errors"
	"github.com/ServiceWeaver/weaver"
	"github.com/hashicorp/golang-lru/v2/expirable"

	"raycat/internal/pkg/readkit"
	"raycat/internal/pkg/tinypool"
	"time"
)

var _ subFileSourceProvider = (*subFileSource)(nil)

var (
	subFileSourceDirNotConfiguredError = errors.New("sub file source dir not configured")
)

var (
	fileSubLru = expirable.NewLRU[string, []byte](10, nil, time.Hour*24)
	fileBuf    = tinypool.New(tinypool.BufReset)
)

type subFileSourceProvider interface {
	UpdateFileSub(ctx context.Context, subFilePaths []string) ([]byte, error)
}

type subFileSource struct {
	weaver.Implements[subFileSourceProvider]
}

func (s *subFileSource) UpdateFileSub(ctx context.Context, subFilePaths []string) ([]byte, error) {
	if len(subFilePaths) == 0 {
		s.Logger(ctx).Warn("skip file source,for dir was not config")
		return nil, subFileSourceDirNotConfiguredError
	}
	buf := fileBuf.Get()
	defer fileBuf.Free(buf)
	for _, path := range subFilePaths {
		if content, found := fileSubLru.Get(path); found {
			if _, err := buf.Write(content); err != nil {
				s.Logger(context.Background()).Error("write buffered content to buf", "path", path, "err", err)
			}
			continue
		}
		content, err := readkit.ReadAll(path)
		if err != nil {
			s.Logger(ctx).Error("failed to read dir", err)
			return nil, err
		}
		fileSubLru.Add(path, content)
		if _, err = buf.Write(content); err != nil {
			s.Logger(context.Background()).Error("write read content to buf", "path", path, "err", err)
		}
	}
	b := buf.Bytes()
	return b, nil
}
