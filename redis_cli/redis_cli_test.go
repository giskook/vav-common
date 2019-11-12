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

func TestSetAccessAddr(t *testing.T) {
	init_redis()
	GetInstance().SetAccessAddr("vavms", "192.168.2.122", "8876")
	ip, port, err := GetInstance().GetAccessAddr("vavms")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip)
	t.Log(port)
	GetInstance().SetAccessAddr("vavms", "192.168.2.123", "8877")
	ip, port, err = GetInstance().GetAccessAddr("vavms")
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
			AccessUUID: "vavms1",
		},
		&base.StreamMedia{
			AccessUUID: "vavms2",
		},
		&base.StreamMedia{
			AccessUUID: "vavms3",
		},
		&base.StreamMedia{
			AccessUUID: "vavms4",
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
		AccessUUID: "vavms4",
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
}

func TestVehicleStreamFormat(t *testing.T) {
	init_redis()
	err := GetInstance().SetVehicleStreamFormat("15226563111", "g726", "h264", "8000")
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

	vavms_info, err := GetInstance().GetVavmsInfo("15226563111", "3", "vavms2", "vavms_stream_media")
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

const (
	SCRIPT = `
	local result
	result = redis.call("HSET", KEYS[1], KEYS[2], KEYS[3])
	return result
	`
)

func TestScriptDo(t *testing.T) {
	init_redis()
	t.Log(GetInstance().DoScript(SCRIPT, "test_script", "test", "123"))
}

func TestPubSub(t *testing.T) {
	init_redis()
	f := func(content []byte) error {
		t.Log(content)

		return nil
	}
	go func() {
		exit := make(chan struct{})
		GetInstance().Sub("test_pub_sub", exit, f)
	}()
	time.Sleep(1000)
	t.Log(GetInstance().Pub("test_pub_sub", "hello world"))
}

func TestExists(t *testing.T) {
	init_redis()

	t.Log(GetInstance().ExistKey("abc"))

}

func TestPipeHGet(t *testing.T) {
	init_redis()
	_, err := GetInstance().PipeHGet("abc", []string{"aaa", "bbb"})
	t.Log(err)
}
