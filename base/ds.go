package base

type StreamMedia struct {
	AccessUUID  string `json:"access_uuid"`
	DomainInner string `json:"domain_inner"`
	DomainOuter string `json:"domain_outer"`
}

type VavmsInfo struct {
	Acodec         string
	Vcodec         string
	LiveType       int
	LiveStatus     int
	PlayBackType   int
	PlayBackStatus int
	DomainInner    string
	DomainOuter    string
}
