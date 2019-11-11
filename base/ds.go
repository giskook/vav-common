package base

type StreamMedia struct {
	AccessUUID      string `json:"access_uuid"`
	RtmpApplication string `json:"rtmp_application"`
	RtmpIpInner     string `json:"rtmp_ip_inner"`
	RtmpIpOutter    string `json:"rtmp_ip_outter"`
	RtmpPortInner   string `json:"rtmp_port_inner"`
	RtmpPortOutter  string `json:"rtmp_port_outter"`
	HttpLocation    string `json:"http_location"`
	HttpIpOutter    string `json:"http_ip_outter"`
	HttpPortOutter  string `json:"http_port_outter"`
}

type VavmsInfo struct {
	Acodec       string
	Vcodec       string
	SamplingRate string
	Status       string
	DataType     string
	TTL          string
	StreamMedia  *StreamMedia
}
