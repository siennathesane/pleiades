package conf

type GossipConfig struct {
	BindAddress      string
	AdvertiseAddress string
	Seed             []string
}
