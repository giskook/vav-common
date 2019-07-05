package redis_cli

import (
	"github.com/gomodule/redigo/redis"
)

func (r *redis_cli) SetVehicleChan(id, k, v string) error {
	c := r.get_conn()
	defer c.Close()
	_, err := c.Do("HSET", id, k, v)
	if err != nil {
		return err
	}

	return nil
}

func (r *redis_cli) GetVehicleChan(id, k string) (string, error) {
	c := r.get_conn()
	defer c.Close()

	v, err := redis.String(c.Do("HGET", id, k))
	if err != nil {
		return "", err
	}

	return v, nil
}
