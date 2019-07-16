package redis_cli

import (
	"encoding/json"
	"github.com/giskook/vav-common/base"
	"github.com/gomodule/redigo/redis"
)

const (
	PLAY_TYPE_V      int    = 1
	PLAY_TYPE_A      int    = 2
	LIVE_TYPE        string = "live_type"   // 0 none 1 video 2 audio 3 both
	LIVE_STATUS      string = "live_status" // http 1 rtp set 2
	PLAYBACK_TYPE    string = "play_back_type"
	PLAYBACK_STATUS  string = "play_back_status"
	PLAY_STATUS_INIT int    = 1
	PLAY_STATUS_OK   int    = 2
)

func (r *redis_cli) GetVavmsInfo(id, id_channel, access_server_uuid, stream_media string) (*base.VavmsInfo, error) {
	c := r.get_conn()
	defer c.Close()
	c.Send("HMGET", id, ACODEC, VCODEC)
	c.Send("HMGET", id_channel, LIVE_TYPE, LIVE_STATUS, PLAYBACK_TYPE, PLAYBACK_STATUS)
	c.Send("LRANGE", stream_media, "0", "-1")
	c.Flush()
	av, err := redis.Strings(c.Receive())
	if err != nil {
		return nil, err
	}

	status, err := redis.Ints(c.Receive())
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
		Acodec:         av[0],
		Vcodec:         av[1],
		LiveType:       status[0],
		LiveStatus:     status[1],
		PlayBackType:   status[2],
		PlayBackStatus: status[3],
		DomainInner:    srv_single.DomainInner,
		DomainOuter:    srv_single.DomainOuter,
	}, nil
}
