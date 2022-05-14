package fsm

import (
	"strconv"
	"strings"
)

const (
	RootKeyFormat     string = "wrn:partition:service:region:account-id:resource-type/resource-id"
	PartitionIndex    int    = 1
	ServiceIndex      int    = 2
	RegionIndex       int    = 3
	AccountIdIndex    int    = 4
	ResourceTypeIndex int    = 5
	ResourceIdIndex   int    = 6
)

type PartitionType string

const (
	GlobalPartition PartitionType = "global"
)

type ServiceType string

const (
	Wraith ServiceType = "wraith"
)

type RegionType string

const (
	GlobalRegion RegionType = "*"
)

type ResourceType string

const (
	Bucket ResourceType = "bucket"
)

type WraithResourceName struct {
	Partition    PartitionType
	Service      ServiceType
	Region       RegionType
	AccountId    int
	ResourceType ResourceType
	ResourceId   string
}

func ParseWrn(wrn string) (*WraithResourceName, error) {
	s := strings.Split(wrn, ":")
	val, err := strconv.Atoi(s[AccountIdIndex])
	if err != nil {
		return nil, err
	}

	return &WraithResourceName{
		Partition:    PartitionType(s[PartitionIndex]),
		Service:      ServiceType(s[ServiceIndex]),
		Region:       RegionType(s[RegionIndex]),
		AccountId:    val,
		ResourceType: ResourceType(s[ResourceTypeIndex]),
		ResourceId:   s[ResourceIdIndex],
	}, nil
}
