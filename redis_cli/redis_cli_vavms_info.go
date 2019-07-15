package redis_cli

import (
	"encoding/json"
	"github.com/giskook/vav-common/base"
	"github.com/gomodule/redigo/redis"
)

func (r *redis_cli) GetVavmsInfo(id, id_channel, access_server_uuid, stream_media string) (*base.VavmsInfo, error) {
	c := r.get_conn()
	defer c.Close()
	c.Send("HMGET", id, ACODEC, VCODEC)
	c.Send("GEt", id_channel)
	c.Send("LRANGE", stream_media, "0", "-1")
	c.Flush()
	av, err := redis.Strings(c.Receive())
	if err != nil {
		return nil, err
	}

	play_type, err := redis.String(c.Receive())
	if err != nil {
		return nil, err
	}

	srvs, err := redis.Values(c.Receive())
	if err != nil {
		return nil, err
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
		Acodec:      av[0],
		Vcodec:      av[1],
		PlayType:    play_type,
		DomainInner: srv_single.DomainInner,
		DomainOuter: srv_single.DomainOuter,
	}, nil
}
