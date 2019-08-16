package redis_cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/giskook/vav-common/base"
	"github.com/gomodule/redigo/redis"
)

const ()

func (r *redis_cli) GetTwisInfo(id, status_key, access_server_uuid, stream_media string) (*base.VavmsInfo, error) {
	c := r.get_conn()
	defer c.Close()
	c.Send("HMGET", id, ACODEC, VCODEC, SAMPLING_RATE)
	c.Send("GET", status_key)
	c.Send("HGET", "twis_"+id, "ttl")
	c.Send("LRANGE", stream_media, "0", "-1")
	c.Flush()
	avs, err := redis.Strings(c.Receive())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sim %s get avs audio and video format error %s ", id, err.Error()))
	}

	status, err := redis.String(c.Receive())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sim %s get status error %s ", id, err.Error()))
	}

	ttl, err := redis.String(c.Receive())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sim %s status %s get data type and ttl error %s", id, status, err.Error()))
	}

	srvs, err := redis.Values(c.Receive())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sim %s get stream media error %s ", id, err.Error()))
	}

	srv_single := new(base.StreamMedia)
	for _, srv := range srvs {
		err = json.Unmarshal(srv.([]byte), srv_single)
		if err != nil {
			continue
		}
		if srv_single.AccessUUID == access_server_uuid {
			break
		}
	}

	return &base.VavmsInfo{
		Acodec:       avs[0],
		Vcodec:       avs[1],
		SamplingRate: avs[2],
		DataType:     status,
		TTL:          ttl,
		DomainInner:  srv_single.DomainInner,
		DomainOuter:  srv_single.DomainOuter,
	}, nil
}
