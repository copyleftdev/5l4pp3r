package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite" // SQLite driver

	"github.com/copyleftdev/5l4pp3r/internal/model"
)

type sqliteStorage struct {
	db *sql.DB
	tx *sql.Tx
}

func NewSQLiteStorage(uri string) (Storage, error) {
	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	// WAL mode for concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	return &sqliteStorage{db: db, tx: tx}, nil
}

func (s *sqliteStorage) InitSchema(ctx context.Context) error {
	schema := `
CREATE TABLE IF NOT EXISTS system_info (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  hostname TEXT NOT NULL,
  created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS network_interfaces (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  system_id INTEGER NOT NULL,
  interface_name TEXT NOT NULL,
  ip_address TEXT,
  mac_address TEXT,
  FOREIGN KEY(system_id) REFERENCES system_info(id)
);

CREATE TABLE IF NOT EXISTS config_files (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  system_id INTEGER NOT NULL,
  file_path TEXT NOT NULL,
  size INTEGER NOT NULL,
  permissions TEXT NOT NULL,
  modified_time DATETIME NOT NULL,
  data BLOB NOT NULL,
  FOREIGN KEY(system_id) REFERENCES system_info(id)
);`

	_, err := s.db.ExecContext(ctx, schema)
	return err
}

func (s *sqliteStorage) StoreSystemInfo(ctx context.Context, si *model.SystemInfo) error {
	si.CreatedAt = time.Now().UTC()
	res, err := s.tx.ExecContext(ctx, `INSERT INTO system_info(hostname, created_at) VALUES (?,?)`,
		si.Hostname, si.CreatedAt)
	if err != nil {
		return err
	}
	si.ID, _ = res.LastInsertId()
	return nil
}

func (s *sqliteStorage) StoreNetworkInterface(ctx context.Context, ni *model.NetworkInterface) error {
	if ni.SystemID == 0 {
		return fmt.Errorf("system_id not set on network interface")
	}
	_, err := s.tx.ExecContext(ctx, `INSERT INTO network_interfaces(system_id, interface_name, ip_address, mac_address) VALUES (?,?,?,?)`,
		ni.SystemID, ni.InterfaceName, ni.IPAddress, ni.MACAddress)
	return err
}

func (s *sqliteStorage) StoreConfigFile(ctx context.Context, cf *model.ConfigFile) error {
	if cf.SystemID == 0 {
		return fmt.Errorf("system_id not set on config file")
	}
	_, err := s.tx.ExecContext(ctx, `INSERT INTO config_files(system_id,file_path,size,permissions,modified_time,data) VALUES (?,?,?,?,?,?)`,
		cf.SystemID, cf.FilePath, cf.Size, cf.Permissions, cf.ModifiedTime, cf.Data)
	return err
}

func (s *sqliteStorage) Commit(ctx context.Context) error {
	return s.tx.Commit()
}

func (s *sqliteStorage) Close() error {
	return s.db.Close()
}
