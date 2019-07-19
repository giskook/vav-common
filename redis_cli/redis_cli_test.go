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
	GetInstance().SetStreamMedia("vavms_stream_media", []*base.StreamMedia{
		&base.StreamMedia{
			AccessUUID:  "vavms1",
			DomainInner: "rtmp://127.0.0.1:8080/myapp",
			DomainOuter: "rtmp://192.168.2.122:8080/myapp",
		},
		&base.StreamMedia{
			AccessUUID:  "vavms2",
			DomainInner: "rtmp://127.0.0.1:8080/myapp",
			DomainOuter: "rtmp://192.168.2.121:8080/myapp",
		},
		&base.StreamMedia{
			AccessUUID:  "vavms3",
			DomainInner: "rtmp://127.0.0.1:8080/myapp",
			DomainOuter: "rtmp://192.168.2.123:8080/myapp",
		},
		&base.StreamMedia{
			AccessUUID:  "vavms4",
			DomainInner: "rtmp://127.0.0.1:8080/myapp",
			DomainOuter: "rtmp://192.168.2.124:8080/myapp",
		},
	})

	sm, err := GetInstance().GetStreamMedia("vavms_stream_media", "0", "-1")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range sm {
		t.Log(*v)
	}
	result := GetInstance().DelStreamMedia("vavms_stream_media", "-1")
	t.Log(result)
	result = GetInstance().DelStreamMedia("vavms_stream_media", "5")
	t.Log(result)
	sm, err = GetInstance().GetStreamMedia("vavms_stream_media", "0", "-1")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range sm {
		t.Log(*v)
	}
	stream_media := &base.StreamMedia{
		AccessUUID:  "vavms4",
		DomainInner: "rtmp://127.0.0.1:8080/myapp",
		DomainOuter: "rtmp://192.168.2.124:8080/myapp",
	}
	result = GetInstance().UpdateStreamMedia("vavms_stream_media", "1", stream_media)
	t.Log(result)
	result = GetInstance().UpdateStreamMedia("vavms_stream_media", "5", stream_media)
	t.Log(result)
	sm, err = GetInstance().GetStreamMedia("vavms_stream_media", "0", "-1")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range sm {
		t.Log(*v)
	}
}

func TestVechicleChan(t *testing.T) {
	init_redis()
	err := GetInstance().SetVehicleChan("15226563111_3", PLAYBACK_TYPE, "2")
	if err != nil {
		t.Fatal(err)
	}
	v, err := GetInstance().GetVehicleChan("15226563111_3", PLAYBACK_TYPE)
	if err != nil {
		t.Fatal(err)
	}
	err = GetInstance().SetVehicleChan("15226563111_3", PLAYBACK_STATUS, "0")
	if err != nil {
		t.Fatal(err)
	}
	v, err = GetInstance().GetVehicleChan("15226563111_3", PLAYBACK_STATUS)
	if err != nil {
		t.Fatal(err)
	}
	err = GetInstance().SetVehicleChan("15226563111_3", LIVE_TYPE, "1")
	if err != nil {
		t.Fatal(err)
	}
	v, err = GetInstance().GetVehicleChan("15226563111_3", LIVE_TYPE)
	if err != nil {
		t.Fatal(err)
	}
	err = GetInstance().SetVehicleChan("15226563111_3", LIVE_STATUS, "1")
	if err != nil {
		t.Fatal(err)
	}
	v, err = GetInstance().GetVehicleChan("15226563111_3", LIVE_STATUS)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func TestVehicleStreamFormat(t *testing.T) {
	init_redis()
	err := GetInstance().SetVehicleStreamFormat("15226563111", "g726", "h264")
	if err != nil {
		t.Fatal(err)
	}

	audio_format, vedio_format, err := GetInstance().GetVehicleStreamFormat("15226563111")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(audio_format, vedio_format)
}

func TestGetVavmsInfo(t *testing.T) {
	init_redis()

	vavms_info, err := GetInstance().GetVavmsInfo("15226563111", "15226563111_3", "vavms2", "vavms_stream_media")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vavms_info)
	vavms_info, err = GetInstance().GetVavmsInfo("15226563111", "15226563111_3", "vavms1", "vavms_stream_media")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vavms_info)
}
