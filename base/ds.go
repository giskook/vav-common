package base

type StreamMedia struct {
	AccessUUID  string `json:"access_uuid"`
	DomainInner string `json:"domain_inner"`
	DomainOuter string `json:"domain_outer"`
}

type VavmsInfo struct {
	Acodec       string
	Vcodec       string
	SamplingRate string
	Status       string
	DataType     string
	TTL          string
	DomainInner  string
	DomainOuter  string
}
