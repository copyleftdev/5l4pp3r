package storage

import (
	"context"
	"fmt"

	"github.com/copyleftdev/5l4pp3r/internal/model"
)

type Storage interface {
	InitSchema(ctx context.Context) error
	StoreSystemInfo(ctx context.Context, si *model.SystemInfo) error
	StoreNetworkInterface(ctx context.Context, ni *model.NetworkInterface) error
	StoreConfigFile(ctx context.Context, cf *model.ConfigFile) error
	Commit(ctx context.Context) error
	Close() error
}

func NewStorage(dbType, uri string) (Storage, error) {
	switch dbType {
	case "sqlite":
		return NewSQLiteStorage(uri)
	case "postgres":
		return NewPostgresStorage(uri)
	default:
		return nil, ErrUnsupportedDB
	}
}

var ErrUnsupportedDB = fmt.Errorf("unsupported database type")
