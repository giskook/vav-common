package redis_cli

import (
	"encoding/json"
	gkbase "github.com/giskook/go/base"
	"github.com/giskook/vav-common/base"
	"github.com/gomodule/redigo/redis"
)

const (
	// 0  fail 1 ok
	SCRIPT_STREAM_MEDIA_DEL = `local key 
	local result 
	key = redis.call("LINDEX", KEYS[1], KEYS[2])
	result = redis.call("LREM", KEYS[1], 1, key)
	return result
	`
	// 0  fail 1 ok
	SCRIPT_STREAM_MEDIA_UPDATE = `local key 
	local result 
	key = redis.call("LINDEX", KEYS[1], KEYS[2]) 
	result = redis.call("LINSERT", KEYS[1], "AFTER", key, KEYS[3])
	if tonumber(result) <= 0 then 
		return 0
	end
	result = redis.call("LREM", KEYS[1], 1, key)
	return result
	`
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

func (r *redis_cli) DelStreamMedia(stream_media, index string) bool {
	c := r.get_conn()
	defer c.Close()

	s := redis.NewScript(2, SCRIPT_STREAM_MEDIA_DEL)
	result, err := redis.Int(s.Do(c, stream_media, index))
	if err != nil {
		gkbase.ErrorCheck(err)
		return false
	}

	switch result {
	case 0:
		return false
	}

	return true
}

func (r *redis_cli) UpdateStreamMedia(stream_media, index string, new_stream_media *base.StreamMedia) bool {
	c := r.get_conn()
	defer c.Close()

	data, err := json.Marshal(new_stream_media)
	if err != nil {
		gkbase.ErrorCheck(err)
		return false
	}
	s := redis.NewScript(3, SCRIPT_STREAM_MEDIA_UPDATE)
	result, err := redis.Int(s.Do(c, stream_media, index, string(data)))
	if err != nil {
		gkbase.ErrorCheck(err)
		return false
	}

	switch result {
	case 0:
		return false
	}

	return true
}
