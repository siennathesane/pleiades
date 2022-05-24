package fsm

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	RootPrnFormat     string = "prn:partition:service:region:account-id:resource-type/resource-id"
	PartitionIndex    int    = 1
	ServiceIndex      int    = 2
	RegionIndex       int    = 3
	AccountIdIndex    int    = 4
	ResourceTypeIndex int    = 5
	ResourceIdIndex   int    = 6

	// root fsm key prefix
	fsmRootKeyCount int = 5
)

type PartitionType string

const (
	GlobalPartition PartitionType = "global"
)

type ServiceType string

const (
	Pleiades ServiceType = "pleiades"
)

type RegionType string

const (
	GlobalRegion RegionType = "*"
)

type ResourceType string

const (
	Bucket ResourceType = "bucket"
)

type PleiadesResourceName struct {
	Partition    PartitionType
	Service      ServiceType
	Region       RegionType
	AccountId    int
	ResourceType ResourceType
	ResourceId   string
}

func ParsePrn(prn string) (*PleiadesResourceName, error) {
	s := strings.Split(prn, ":")
	val, err := strconv.Atoi(s[AccountIdIndex])
	if err != nil {
		return nil, err
	}

	resourcePair := s[len(s)-1]
	rpSplit := strings.Split(resourcePair, "/")

	return &PleiadesResourceName{
		Partition:    PartitionType(s[PartitionIndex]),
		Service:      ServiceType(s[ServiceIndex]),
		Region:       RegionType(s[RegionIndex]),
		AccountId:    val,
		ResourceType: ResourceType(rpSplit[0]),
		ResourceId:   rpSplit[1],
	}, nil
}

func (prn *PleiadesResourceName) ToFsmRootPath(bucketName string) string {
	return fmt.Sprintf("/%s/%s/%s/%d/%s", prn.Partition, prn.Service, prn.Region, prn.AccountId, bucketName)
}
