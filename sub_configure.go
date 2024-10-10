package main

import (
	"context"
	"errors"
	"github.com/ServiceWeaver/weaver"
)

var (
	responseOptionNotConfiguredError   = errors.New("no response option configured")
	subPublishPathNotConfiguredError   = errors.New("sub publish path not configured")
	subAuthParamNameNotConfiguredError = errors.New("sub auth param name not configured")
)

var _ subConfigureProvider = (*subConfigure)(nil)

type subConfigureProvider interface {
	GetSubFilePaths(ctx context.Context, privateSubToken string) ([]string, error)
	GetUrlSubs(ctx context.Context, privateSubToken string) ([]string, int, error)
	GetSubPublishPath(ctx context.Context) (string, error)
	GetSubAuthParamName(ctx context.Context) (string, error)
	GetResponseOption(ctx context.Context) (*responseOption, error)
}

type subConfig struct {
	weaver.AutoMarshal
	PublicSubFilePaths        []string `toml:"public_sub_file_paths"`
	PrivateSubFilePaths       []string `toml:"private_sub_file_paths"`
	UrlSubFetchTimeoutSeconds int      `toml:"url_sub_fetch_timeout_seconds"`
	PublicUrlSubs             []string `toml:"public_url_subs"`
	PrivateUrlSubs            []string `toml:"private_url_subs"`

	PrivateSubToken  string          `toml:"private_sub_token"`
	SubPublishPath   string          `toml:"sub_publish_path,omitempty"`
	SubAuthParamName string          `toml:"sub_auth_param_name,omitempty"`
	ResponseOption   *responseOption `toml:"response_option,omitempty"`
}

type responseOption struct {
	weaver.AutoMarshal
	//return to the client see https://www.clashverge.dev/guide/url_schemes.html#_
	UpdateIntervalHours int    `toml:"update_interval_hours,omitempty"`
	ProfileWebPage      string `toml:"profile_web_page,omitempty"`
	// subscription-userinfo: upload=1234; download=2234; total=1024000; expire=2218532293
}

type subConfigure struct {
	weaver.Implements[subConfigureProvider]
	weaver.WithConfig[subConfig]
}

func (s *subConfigure) GetSubFilePaths(ctx context.Context, privateSubToken string) ([]string, error) {
	config := s.Config()
	if privateSubToken != config.PrivateSubToken {
		return config.PublicSubFilePaths, nil
	}
	return append(config.PrivateSubFilePaths, config.PublicSubFilePaths...), nil
}

func (s *subConfigure) GetUrlSubs(ctx context.Context, privateSubToken string) ([]string, int, error) {
	config := s.Config()
	if privateSubToken != config.PrivateSubToken {
		s.Logger(ctx).Info("token check pass")
		return config.PublicUrlSubs, config.UrlSubFetchTimeoutSeconds, nil
	}
	return append(config.PrivateUrlSubs, config.PublicUrlSubs...), config.UrlSubFetchTimeoutSeconds, nil
}

func (s *subConfigure) GetSubPublishPath(ctx context.Context) (string, error) {
	if s.Config().SubPublishPath == "" {
		return "", subPublishPathNotConfiguredError
	}
	return s.Config().SubPublishPath, nil
}

func (s *subConfigure) GetSubAuthParamName(ctx context.Context) (string, error) {
	if s.Config().SubAuthParamName == "" {
		return "", subAuthParamNameNotConfiguredError
	}
	return s.Config().SubAuthParamName, nil
}
func (s *subConfigure) GetResponseOption(ctx context.Context) (*responseOption, error) {
	config := s.Config()
	if config.ResponseOption == nil {
		return nil, responseOptionNotConfiguredError
	}
	return config.ResponseOption, nil
}
