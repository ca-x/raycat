package main

import (
	"context"
	"errors"
	"github.com/ServiceWeaver/weaver"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"raycat/internal/pkg/readkit"
	"time"
)

var _ subFileSourceProvider = (*subFileSource)(nil)

var (
	subFileSourceDirNotConfiguredError = errors.New("sub file source dir not configured")
)

var (
	fileSubLru = expirable.NewLRU[string, []byte](10, nil, time.Hour*24)
)

type subFileSourceConfig struct {
	weaver.AutoMarshal
	FilePaths []string `toml:"file_paths"`
}

type subFileSourceProvider interface {
	UpdateSub(ctx context.Context) ([]byte, error)
}

type subFileSource struct {
	weaver.Implements[subFileSourceProvider]
	weaver.WithConfig[subFileSourceConfig]
}

func (s *subFileSource) UpdateSub(ctx context.Context) ([]byte, error) {
	if len(s.Config().FilePaths) == 0 {
		s.Logger(ctx).Warn("skip file source,for dir was not config")
		return nil, subFileSourceDirNotConfiguredError
	}
	allFileSubs := make([]byte, 0, 1024)
	for _, path := range s.Config().FilePaths {
		if content, found := fileSubLru.Get(path); found {
			allFileSubs = append(allFileSubs, content...)
			continue
		}
		content, err := readkit.ReadAll(path)
		if err != nil {
			s.Logger(ctx).Error("failed to read dir", err)
			return nil, err
		}
		fileSubLru.Add(path, content)
		allFileSubs = append(allFileSubs, content...)
	}
	return allFileSubs, nil
}
