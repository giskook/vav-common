package redis_cli

import (
	"github.com/gomodule/redigo/redis"
)

const (
	ACODEC        string = "acodec"
	VCODEC        string = "vcodec"
	SAMPLING_RATE string = "sampling_rate"
)

func (r *redis_cli) SetVehicleStreamFormat(id, acodec, vcodec, sampling_rate string) error {
	c := r.get_conn()
	defer c.Close()
	_, err := c.Do("HMSET", id, ACODEC, acodec, VCODEC, vcodec, SAMPLING_RATE, sampling_rate)
	if err != nil {
		return err
	}

	return nil
}

func (r *redis_cli) GetVehicleStreamFormat(id string) (string, string, error) {
	c := r.get_conn()
	defer c.Close()

	codecs, err := redis.Strings(c.Do("HMGET", id, ACODEC, VCODEC))
	if err != nil {
		return "", "", err
	}

	return codecs[0], codecs[1], nil
}
