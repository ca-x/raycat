package main

import (
	"context"
	"errors"
	"github.com/ServiceWeaver/weaver"
	"raycat/internal/pkg/readkit"
)

var _ subFileSourceProvider = (*subFileSource)(nil)

var (
	subFileSourceDirNotConfiguredError = errors.New("sub file source dir not configured")
)

type subFileSourceConfig struct {
	weaver.AutoMarshal
	Dir string `toml:"dir"`
}

type subFileSourceProvider interface {
	UpdateSub(ctx context.Context) ([]byte, error)
}

type subFileSource struct {
	weaver.Implements[subFileSourceProvider]
	weaver.WithConfig[subFileSourceConfig]
}

func (s *subFileSource) UpdateSub(ctx context.Context) ([]byte, error) {
	if s.Config().Dir != "" {
		s.Logger(ctx).Warn("skip file source,for dir was not config")
		return nil, subFileSourceDirNotConfiguredError
	}
	content, err := readkit.ReadAll(s.Config().Dir)
	if err != nil {
		s.Logger(ctx).Error("failed to read dir", err)
		return nil, err
	}
	return content, nil
}
