package main

import (
	"context"
	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
	"raycat/internal/ent"
	"time"
)

var _ subSourceManageProvider = (*subSourceManager)(nil)

type SubscribeEntry struct {
	weaver.AutoMarshal
	ID uuid.UUID
	// Kind holds the value of the "kind" field.
	Kind int
	// Location holds the value of the "location" field.
	Location string
	// UpdateTimeoutSeconds holds the value of the "update_timeout_seconds" field.
	UpdateTimeoutSeconds int
	// Latency holds the value of the "latency" field.
	Latency int64
	// ExpireAt holds the value of the "expire_at" field.
	ExpireAt time.Time
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time
}
type subSourceManageProvider interface {
	// CreateOneSubSource create a subscribe source
	CreateOneSubSource(ctx context.Context, subscribeEntry *SubscribeEntry) error
	// FindOneSubSource  find one sub source by its uuid
	FindOneSubSource(ctx context.Context, id uuid.UUID) (*SubscribeEntry, error)
	// DeleteOneSubSource  delete one sub source by its uuid
	DeleteOneSubSource(ctx context.Context, id uuid.UUID) error
	// UpdateOneSubSource update one sub source by its uuid
	UpdateOneSubSource(ctx context.Context, id uuid.UUID, updateEntry *SubscribeEntry) error
	// GetAllSubSourcesByKind  get all sub source by its kind
	GetAllSubSourcesByKind(ctx context.Context, kind int) ([]*SubscribeEntry, error)
	// GetAllSubSources get all sub sources
	GetAllSubSources(ctx context.Context) ([]*SubscribeEntry, error)
}

type subSourceManagerConfig struct {
	Driver  string `toml:"driver"`
	ConnStr string `toml:"conn_str"`
}

type subSourceManager struct {
	weaver.Implements[subSourceManageProvider]
	weaver.WithConfig[subSourceManagerConfig]
	c *ent.Client
}

func (s *subSourceManager) Init(ctx context.Context) error {
	config := s.Config()
	c, err := ent.Open(config.Driver, config.ConnStr)
	if err != nil {
		return err
	}
	s.c = c
	return nil
}
func (s *subSourceManager) FindOneSubSource(ctx context.Context, id uuid.UUID) (*SubscribeEntry, error) {
	subscribe, err := s.c.Subscribe.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &SubscribeEntry{
		ID:                   subscribe.ID,
		Kind:                 subscribe.Kind,
		Location:             subscribe.Location,
		UpdateTimeoutSeconds: subscribe.UpdateTimeoutSeconds,
		Latency:              subscribe.Latency,
		ExpireAt:             subscribe.ExpireAt,
		CreatedAt:            subscribe.CreatedAt,
	}, nil
}

func (s *subSourceManager) DeleteOneSubSource(ctx context.Context, id uuid.UUID) error {
	deleteBuilder := s.c.Subscribe.DeleteOneID(id)
	return deleteBuilder.Exec(ctx)
}

func (s *subSourceManager) UpdateOneSubSource(ctx context.Context, id uuid.UUID, updateEntry *SubscribeEntry) error {
	updateBuilder := s.c.Subscribe.UpdateOneID(id)
	return updateBuilder.
		SetKind(updateEntry.Kind).
		SetLocation(updateEntry.Location).
		SetLatency(updateEntry.Latency).
		Exec(ctx)
}

func (s *subSourceManager) GetAllSubSourcesByKind(ctx context.Context, kind int) ([]*SubscribeEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (s *subSourceManager) GetAllSubSources(ctx context.Context) ([]*SubscribeEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (s *subSourceManager) CreateOneSubSource(ctx context.Context, subscribeEntry *SubscribeEntry) error {
	//TODO implement me
	panic("implement me")
}
