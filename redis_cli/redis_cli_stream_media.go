package redis_cli

import (
	"encoding/json"
	"github.com/giskook/vav-common/base"
	"github.com/gomodule/redigo/redis"
)

func (r *redis_cli) SetStreamMedia(stream_media string, media []*base.StreamMedia) error {
	c := r.get_conn()
	defer c.Close()
	_, err := c.Do("DEL", stream_media)
	if err != nil {
		return err
	}
	for _, sm := range media {
		data, err := json.Marshal(sm)
		if err != nil {
			return err
		}
		err = c.Send("RPUSH", stream_media, string(data))
		if err != nil {
			return err
		}
	}
	_, err = c.Do("")
	if err != nil {
		return err
	}

	return nil
}

func (r *redis_cli) GetStreamMedia(stream_media string, start, stop string) ([]*base.StreamMedia, error) {
	c := r.get_conn()
	defer c.Close()

	srvs, err := redis.Values(c.Do("LRANGE", stream_media, start, stop))
	if err != nil {
		return nil, err
	}

	sm := make([]*base.StreamMedia, 0)

	for _, srv := range srvs {
		srv_single := new(base.StreamMedia)
		err = json.Unmarshal(srv.([]byte), srv_single)
		if err != nil {
			continue
		}
		sm = append(sm, srv_single)
	}

	return sm, nil
}

func (r *redis_cli) DelStreamMedia(stream_media, index string) error {
	c := r.get_conn()
	defer c.Close()

	ss, err := redis.String(c.Do("LINDEX", stream_media, index))
	if err != nil {
		return err
	}
	_, err = c.Do("LREM", stream_media, 1, ss)
	if err != nil {
		return err
	}

	return nil
}
