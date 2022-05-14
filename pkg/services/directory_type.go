package services

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
)

type DirectoryType string

const (
	WraithCoreDirectory        DirectoryType = "wraith"
	PortManagerDirectory       DirectoryType = "port-manager"
	DirectoryManagerDirectory  DirectoryType = "directory-manager"
	WriteAheadLogDirectory     DirectoryType = "wal"
	ShardDirectory             DirectoryType = "shard"
	WriteAheadLogBaseDirectory string        = "wals"
	ShardBaseDirectory         string        = "shards"
)

// DirectoryConfigModel is how a database is represented on disk
type DirectoryConfigModel struct {
	BaseDir  string        `json:"base_dir"`
	Type     DirectoryType `json:"type"`
	Id       int           `json:"id,omitempty"`
	FullPath string        `json:"full_path"`
}

// DirectoryPath is how a database is represented
type DirectoryPath struct {
	BaseDir string
	dType   DirectoryType
	id      int
}

func (d *DirectoryPath) Type() DirectoryType {
	return d.dType
}

func (d *DirectoryPath) SetType(dir DirectoryType, id int) error {
	if (d.dType == ShardDirectory || d.dType == WriteAheadLogDirectory) && id <= 0 {
		return errors.New("cannot render a directory which requires an id")
	}
	d.dType = dir
	d.id = id
	return nil
}

func (d *DirectoryPath) RenderDirectoryPath() (string, error) {
	b, err := filepath.Abs(d.BaseDir)
	if err != nil {
		return "", err
	}

	if d.dType == "" {
		return "", errors.New("cannot render without a directory type")
	}

	switch d.dType {
	case WraithCoreDirectory:
		b = filepath.Join(b, string(WraithCoreDirectory))
	case PortManagerDirectory:
		b = filepath.Join(b, string(PortManagerDirectory))
	case DirectoryManagerDirectory:
		b = filepath.Join(b, string(DirectoryManagerDirectory))
	case ShardDirectory:
		b = filepath.Join(b, ShardBaseDirectory, strconv.Itoa(d.id))
	case WriteAheadLogDirectory:
		b = filepath.Join(b, WriteAheadLogBaseDirectory, strconv.Itoa(d.id))
	}

	return b, nil
}

func (d *DirectoryPath) GetFullDatabasePath() (string, error) {
	switch d.dType {
	case PortManagerDirectory:
		dbp, err := d.RenderDirectoryPath()
		if err != nil {
			return "", err
		}
		return filepath.Join(dbp, string(PortManagerDatabase)), nil
	case DirectoryManagerDirectory:
		dbp, err := d.RenderDirectoryPath()
		if err != nil {
			return "", err
		}
		return filepath.Join(dbp, string(DirectoryManagerDatabase)), nil
	case ShardDirectory:
		dbp, err := d.RenderDirectoryPath()
		if err != nil {
			return "", err
		}
		return filepath.Join(dbp, fmt.Sprintf(string(ShardDatabaseBaseFormat), d.id)), nil
	case WriteAheadLogDirectory:
		dbp, err := d.RenderDirectoryPath()
		if err != nil {
			return "", err
		}
		return filepath.Join(dbp, fmt.Sprintf(string(WriteAheadLogDatabaseBaseFormat), d.id)), nil
	}
	return "", nil
}
