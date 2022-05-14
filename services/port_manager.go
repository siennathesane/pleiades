package services

import (
	"bytes"
	"encoding/binary"
	"path/filepath"
	"unsafe"

	"github.com/hashicorp/consul/api"
	"github.com/lni/dragonboat/v3/logger"
	"go.etcd.io/bbolt"
	"wraith/conf"
)

const (
	portDbName            string = "ports.db"
	portManagerBucketName string = "ports"
)

// HostPortMapping is used to save which port on the host maps to which raft cluster
type HostPortMapping struct {
	ClusterId uint64
	Port      int
}

// PortManager is used for managing the port mappings between
// raft clusters and the host
type PortManager struct {
	PutMapping <-chan HostPortMapping
	env        conf.EnvironmentConfig
	logger     logger.ILogger
	client     *api.Client
	db         *bbolt.DB
	done       <-chan struct{}
}

func NewPortManager(env conf.EnvironmentConfig, logger logger.ILogger, client *api.Client) *PortManager {
	return &PortManager{PutMapping: make(<-chan HostPortMapping, 10), env: env, logger: logger, client: client, done: make(<-chan struct{}, 1)}
}

func (p *PortManager) dbPath() string {
	return filepath.Join(p.env.BaseDir, "wraith", "ports.db")
}

func (p *PortManager) run() {
	for {
		select {
		case update := <-p.PutMapping:
			err := p.db.Update(func(tx *bbolt.Tx) error {
				b := tx.Bucket([]byte(portManagerBucketName))

				// always use the cluster id as the key
				clusterIdBuf := make([]byte, unsafe.Sizeof(update.ClusterId))
				binary.LittleEndian.PutUint64(clusterIdBuf, update.ClusterId)

				var portBuf bytes.Buffer
				if err := binary.Write(&portBuf, binary.LittleEndian, update.Port); err != nil {
					return err
				}

				return b.Put(clusterIdBuf, portBuf.Bytes())
			})
			if err != nil {
				p.logger.Errorf("can't update port manager database: %s", err)
			}
		case <-p.done:
			return
		}
	}
}

func (p *PortManager) Start() error {
	var err error
	p.db, err = bbolt.Open(p.dbPath(), 0600, &bbolt.Options{})
	if err != nil {
		p.logger.Errorf("error opening port manager database")
		return err
	}

	// start our background listener
	go p.run()

	return p.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(portManagerBucketName))
		return err
	})
}

func (p *PortManager) Stop() {
	if err := p.db.Close(); err != nil {
		p.logger.Errorf("can't close port manager database: %s", err)
	}

	// stop the background listener
	<-p.done
}
