package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/copyleftdev/5l4pp3r/internal/model"

	_ "github.com/lib/pq"
)

type postgresStorage struct {
	db *sql.DB
	tx *sql.Tx
}

func NewPostgresStorage(uri string) (Storage, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres db: %w", err)
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return &postgresStorage{db: db, tx: tx}, nil
}

func (p *postgresStorage) InitSchema(ctx context.Context) error {
	schema := `
CREATE TABLE IF NOT EXISTS system_info (
  id SERIAL PRIMARY KEY,
  hostname TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS network_interfaces (
  id SERIAL PRIMARY KEY,
  system_id INTEGER NOT NULL REFERENCES system_info(id),
  interface_name TEXT NOT NULL,
  ip_address TEXT,
  mac_address TEXT
);

CREATE TABLE IF NOT EXISTS config_files (
  id SERIAL PRIMARY KEY,
  system_id INTEGER NOT NULL REFERENCES system_info(id),
  file_path TEXT NOT NULL,
  size BIGINT NOT NULL,
  permissions TEXT NOT NULL,
  modified_time TIMESTAMP WITH TIME ZONE NOT NULL,
  data BYTEA NOT NULL
);`

	_, err := p.db.ExecContext(ctx, schema)
	return err
}

func (p *postgresStorage) StoreSystemInfo(ctx context.Context, si *model.SystemInfo) error {
	si.CreatedAt = time.Now().UTC()
	err := p.tx.QueryRowContext(ctx,
		`INSERT INTO system_info(hostname, created_at) VALUES ($1,$2) RETURNING id`,
		si.Hostname, si.CreatedAt).Scan(&si.ID)
	return err
}

func (p *postgresStorage) StoreNetworkInterface(ctx context.Context, ni *model.NetworkInterface) error {
	if ni.SystemID == 0 {
		return fmt.Errorf("system_id not set")
	}
	_, err := p.tx.ExecContext(ctx,
		`INSERT INTO network_interfaces(system_id, interface_name, ip_address, mac_address) VALUES ($1,$2,$3,$4)`,
		ni.SystemID, ni.InterfaceName, ni.IPAddress, ni.MACAddress)
	return err
}

func (p *postgresStorage) StoreConfigFile(ctx context.Context, cf *model.ConfigFile) error {
	if cf.SystemID == 0 {
		return fmt.Errorf("system_id not set")
	}
	_, err := p.tx.ExecContext(ctx,
		`INSERT INTO config_files(system_id,file_path,size,permissions,modified_time,data) VALUES ($1,$2,$3,$4,$5,$6)`,
		cf.SystemID, cf.FilePath, cf.Size, cf.Permissions, cf.ModifiedTime, cf.Data)
	return err
}

func (p *postgresStorage) Commit(ctx context.Context) error {
	return p.tx.Commit()
}

func (p *postgresStorage) Close() error {
	return p.db.Close()
}
