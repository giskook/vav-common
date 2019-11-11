package redis_cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/giskook/vav-common/base"
	"github.com/gomodule/redigo/redis"
)

const ()

func (r *redis_cli) GetVavmsInfo(id, channel, access_server_uuid, stream_media string) (*base.VavmsInfo, error) {
	c := r.get_conn()
	defer c.Close()
	c.Send("HMGET", id, ACODEC, VCODEC, SAMPLING_RATE)
	c.Send("GET", id+"_"+channel+"_status")
	c.Send("HMGET", id+"_"+channel+"_live", "data_type", "ttl")
	c.Send("HMGET", id+"_"+channel+"_back", "data_type", "ttl")
	c.Send("LRANGE", stream_media, "0", "-1")
	c.Flush()
	avs, err := redis.Strings(c.Receive())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sim %s channel %s get avs audio and video format error %s ", id, channel, err.Error()))
	}

	status, err := redis.String(c.Receive())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sim %s channel %s get status error %s ", id, channel, err.Error()))
	}
	var data_type, ttl string

	live_data_type_ttl, err_live := redis.Strings(c.Receive())
	back_data_type_ttl, err_back := redis.Strings(c.Receive())
	if status == "live" {
		if err_live != nil {
			return nil, errors.New(fmt.Sprintf("sim %s channel %s status %s get data type and ttl error %s", id, channel, status, err_live.Error()))
		}
		data_type = live_data_type_ttl[0]
		ttl = live_data_type_ttl[1]
	}
	if status == "back" {
		if err_back != nil {
			return nil, errors.New(fmt.Sprintf("sim %s channel %s status %s get data type and ttl error %s", id, channel, status, err_live.Error()))
		}
		data_type = back_data_type_ttl[0]
		ttl = back_data_type_ttl[1]
	}

	srvs, err := redis.Values(c.Receive())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sim %s channel %s get stream media error %s ", id, channel, err.Error()))
	}

	match := false
	srv_single := new(base.StreamMedia)
	for _, srv := range srvs {
		err = json.Unmarshal(srv.([]byte), srv_single)
		if err != nil {
			continue
		}
		if srv_single.AccessUUID == access_server_uuid {
			match = true
			break
		}
	}
	if !match {
		return nil, errors.New("no such stream media " + access_server_uuid)
	}

	return &base.VavmsInfo{
		Acodec:       avs[0],
		Vcodec:       avs[1],
		SamplingRate: avs[2],
		Status:       status,
		DataType:     data_type,
		TTL:          ttl,
		StreamMedia:  srv_single,
	}, nil
}
