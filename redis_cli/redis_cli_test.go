package redis_cli

import (
	"github.com/giskook/vav-common/base"
	"testing"
	"time"
)

func init_redis() {
	conf := &Conf{
		Addr:         "127.0.0.1:6379",
		Passwd:       "redis",
		MaxIdle:      100,
		ConnTimeOut:  5 * time.Second,
		ReadTimeOut:  5 * time.Second,
		WriteTimeOut: 5 * time.Second,
	}
	GetInstance().Init(conf)
}

func TestSetAccessServer(t *testing.T) {
	init_redis()
	GetInstance().SetAccessServer("vavms", "192.168.2.122", "8876")
	ip, port, err := GetInstance().GetAccessServer("vavms")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip)
	t.Log(port)
	GetInstance().SetAccessServer("vavms", "192.168.2.123", "8877")
	ip, port, err = GetInstance().GetAccessServer("vavms")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip)
	t.Log(port)
}

func TestSetStreamMedia(t *testing.T) {
	init_redis()
	GetInstance().SetStreamMedia("stream_media", []*base.StreamMedia{
		&base.StreamMedia{
			AccessUUID:  "vavms1",
			DomainInner: "rtmp://127.0.0.1:8888/vavms",
			DomainOuter: "rtmp://192.168.2.121:8888/vavms",
		},
		&base.StreamMedia{
			AccessUUID:  "vavms2",
			DomainInner: "rtmp://127.0.0.1:8888/vavms",
			DomainOuter: "rtmp://192.168.2.122:8888/vavms",
		},
		&base.StreamMedia{
			AccessUUID:  "vavms3",
			DomainInner: "rtmp://127.0.0.1:8888/vavms",
			DomainOuter: "rtmp://192.168.2.123:8888/vavms",
		},
		&base.StreamMedia{
			AccessUUID:  "vavms4",
			DomainInner: "rtmp://127.0.0.1:8888/vavms",
			DomainOuter: "rtmp://192.168.2.124:8888/vavms",
		},
	})

	sm, err := GetInstance().GetStreamMedia("stream_media", "0", "-1")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range sm {
		t.Log(*v)
	}
	err = GetInstance().DelStreamMedia("stream_media", "-1")
	if err != nil {
		t.Fatal(err)
	}
	sm, err = GetInstance().GetStreamMedia("stream_media", "0", "-1")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range sm {
		t.Log(*v)
	}
}

func TestVechicleChan(t *testing.T) {
	init_redis()
	err := GetInstance().SetVehicleChan("13731143001_2", PLAYBACK_TYPE, "2")
	if err != nil {
		t.Fatal(err)
	}
	v, err := GetInstance().GetVehicleChan("13731143001_2", PLAYBACK_TYPE)
	if err != nil {
		t.Fatal(err)
	}
	err = GetInstance().SetVehicleChan("13731143001_2", PLAYBACK_STATUS, "1")
	if err != nil {
		t.Fatal(err)
	}
	v, err = GetInstance().GetVehicleChan("13731143001_2", PLAYBACK_STATUS)
	if err != nil {
		t.Fatal(err)
	}
	err = GetInstance().SetVehicleChan("13731143001_2", LIVE_TYPE, "3")
	if err != nil {
		t.Fatal(err)
	}
	v, err = GetInstance().GetVehicleChan("13731143001_2", LIVE_TYPE)
	if err != nil {
		t.Fatal(err)
	}
	err = GetInstance().SetVehicleChan("13731143001_2", LIVE_STATUS, "1")
	if err != nil {
		t.Fatal(err)
	}
	v, err = GetInstance().GetVehicleChan("13731143001_2", LIVE_STATUS)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func TestVehicleStreamFormat(t *testing.T) {
	init_redis()
	err := GetInstance().SetVehicleStreamFormat("13731143001", "g726", "h264")
	if err != nil {
		t.Fatal(err)
	}

	audio_format, vedio_format, err := GetInstance().GetVehicleStreamFormat("13731143001")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(audio_format, vedio_format)
}

func TestGetVavmsInfo(t *testing.T) {
	init_redis()

	vavms_info, err := GetInstance().GetVavmsInfo("13731143001", "13731143001_2", "vavms2", "stream_media")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vavms_info)
	vavms_info, err = GetInstance().GetVavmsInfo("13731143001", "13731143001_2", "vavms1", "stream_media")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vavms_info)
}
