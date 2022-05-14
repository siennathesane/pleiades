package services

type DatabaseName string

const (
	DirectoryManagerDatabase        DatabaseName = "directories.db"
	PortManagerDatabase             DatabaseName = "ports.db"
	WriteAheadLogDatabaseBaseFormat DatabaseName = "wal-%d.db"
	ShardDatabaseBaseFormat         DatabaseName = "shard-%d.db"
)
