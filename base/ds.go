package base

type StreamMedia struct {
	AccessUUID  string `json:"access_uuid"`
	DomainInner string `json:"domain_inner"`
	DomainOuter string `json:"domain_outer"`
}

type VavmsInfo struct {
	Acodec         string
	Vcodec         string
	LiveType       int // 0 none 1 video 2 audio 3 both
	LiveStatus     int // http 1 rtp set 2
	PlayBackType   int
	PlayBackStatus int
	DomainInner    string
	DomainOuter    string
}
