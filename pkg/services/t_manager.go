package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"go.etcd.io/bbolt"
)

const (
	typeMapBucketName        string = "known-types"
	typeMapKey                      = "type-map"
	tManagerBucketNameFormat string = "%ss"
)

// StoreManager is used for managing various types persistently
type StoreManager struct {
	baseDir     string
	logger      zerolog.Logger
	db          *bbolt.DB
	initialized bool
	typeMap     map[string]bool
	mu          sync.RWMutex
}

func NewStoreManager(baseDir string,
	logger zerolog.Logger) *StoreManager {
	return &StoreManager{
		baseDir:     baseDir,
		logger:      logger,
		initialized: false,
	}
}

func (p *StoreManager) Put(key string, target []byte, t reflect.Type) error {

	name := cleanTypeName(t.String())
	_, ok := p.typeMap[name]
	if !ok {
		err := p.db.Update(func(tx *bbolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		p.typeMap[name] = true
	}

	err := p.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), target)
	})
	if err != nil {
		p.logger.Err(err).Msg("can't update port manager database")
		return err
	}

	return nil
}

func (p *StoreManager) Get(key string, t reflect.Type) ([]byte, error) {
	name := cleanTypeName(t.String())
	_, ok := p.typeMap[name]
	if !ok {
		return nil, errors.New("there are no types stored")
	}

	var target []byte
	err := p.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(name))
		payload := b.Get([]byte(key))
		if payload == nil {
			return fmt.Errorf("%s not found", key)
		}

		target = payload

		return nil
	})

	return target, err
}

func (p *StoreManager) GetAll(t reflect.Type) (map[string][]byte, error) {
	name := cleanTypeName(t.String())
	_, ok := p.typeMap[name]
	if !ok {
		return nil, errors.New("there are no types stored")
	}

	targets := make(map[string][]byte)
	err := p.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(name))
		if b == nil {
			return fmt.Errorf("cannot find %s bucket", name)
		}

		return b.ForEach(func(k, v []byte) error {
			key := fmt.Sprintf("%s", k)
			targets[key] = v
			return nil
		})
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot get all keys in the bucket")
	}

	return targets, err
}

func (p *StoreManager) Flush() error {
	defer p.mu.Unlock()
	p.mu.Lock()

	if err := p.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(typeMapBucketName))
		if err != nil {
			return err
		}
		payload, err := json.Marshal(p.typeMap)
		if err != nil {
			return err
		}
		return b.Put([]byte(typeMapKey), payload)
	}); err != nil {
		p.logger.Err(err).Msg("can't serialize the type map")
		return err
	}
	return nil
}

func (p *StoreManager) Start(retry bool) error {
	if p.initialized {
		return nil
	}
	return p.init()
}

func (p *StoreManager) Restart(retry bool) error {
	if err := p.Stop(retry); err != nil {
		return err
	}
	if err := p.Start(retry); err != nil {
		return err
	}
	return nil
}

func (p *StoreManager) Stop(retry bool) error {
	if err := p.Flush(); err != nil {
		return err
	}

	if err := p.db.Close(); err != nil {
		p.logger.Err(err).Msg("can't close port manager database")
		return err
	}
	return nil
}

func (p *StoreManager) init() error {
	dbPath, err := dbPath(p.baseDir)
	if err != nil {
		return err
	}

	p.db, err = bbolt.Open(dbPath, 0600, &bbolt.Options{})
	if err != nil {
		p.logger.Err(err).Msg("error opening database")
		return err
	}

	return p.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(typeMapBucketName))
		if err != nil {
			return err
		}

		val := b.Get([]byte(typeMapKey))
		if val == nil {
			p.typeMap = make(map[string]bool)
			p.initialized = true
			return nil
		}

		if err := json.Unmarshal(val, &p.typeMap); err != nil {
			return fmt.Errorf("cannot load type map into the manager: %w", err)
		}

		p.initialized = true
		return nil
	})
}

func dbPath(baseDir string) (string, error) {
	dbRoot := filepath.Join(baseDir, "pleiades")
	if err := os.MkdirAll(dbRoot, os.ModePerm); err != nil {
		return "", err
	}
	return filepath.Join(dbRoot, "store.db"), nil
}

func cleanTypeName(typeName string) string {
	return strings.ReplaceAll(typeName, "*", "")
}
