package services

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/hashicorp/consul/api"
	"github.com/lni/dragonboat/v3/logger"
	"go.etcd.io/bbolt"
	"go.uber.org/fx"
	"wraith/conf"
)

const (
	directoryDbName            string = "directories.db"
	directoryManagerBucketName string = "directories"
)

func ProvideDirectoryManager() fx.Option {
	return fx.Provide(NewDirectoryManager, conf.NewConsulClient)
}

// DirectoryManager is how the host manages directory layouts.
type DirectoryManager struct {
	env    *conf.EnvironmentConfig
	logger logger.ILogger
	client *api.Client
	db     *bbolt.DB
	done   <-chan struct{}
}

func NewDirectoryManager(env *conf.EnvironmentConfig, logger logger.ILogger, client *api.Client) *DirectoryManager {
	return &DirectoryManager{env: env, logger: logger, client: client, done: make(<-chan struct{}, 1)}
}

func (d *DirectoryManager) dbPath() string {
	return filepath.Join(d.env.BaseDir, string(WraithCoreDirectory), directoryDbName)
}

func (d *DirectoryManager) Start() {
	var err error
	d.db, err = bbolt.Open(d.dbPath(), 0600, &bbolt.Options{})
	if err != nil {
		d.logger.Panicf("error opening directory manager database")
	}

	if err := d.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(directoryDbName))
		return err
	}); err != nil {
		d.logger.Panicf("error creating directory manager bucket: %s", err)
	}
}

func (d *DirectoryManager) Stop() {
	if err := d.db.Close(); err != nil {
		d.logger.Errorf("can't close directory manager database: %s", err)
	}

	// stop the background listener
	<-d.done
}

func (d *DirectoryManager) NewDirectory(dir DirectoryType, id int) (string, error) {
	if (dir == ShardDirectory || dir == WriteAheadLogDirectory) && id <= 0 {
		return "", errors.New("cannot create a shard or wal directory without an id")
	}

	dirPath := DirectoryPath{
		BaseDir: d.env.BaseDir,
	}

	if err := dirPath.SetType(dir, id); err != nil {
		d.logger.Errorf("can't set directory type: %s", err)
		return "", err
	}

	path, err := dirPath.RenderDirectoryPath()
	if err != nil {
		d.logger.Errorf("can't render directory path: %s", err)
		return "", err
	}

	if err := os.MkdirAll(path, 0644); err != nil {
		d.logger.Errorf("can't create directory tree: %s", err)
		return "", err
	}

	return path, nil
}
