package redis_cli

import (
	"github.com/gomodule/redigo/redis"
)

func (r *redis_cli) GetValue(id string) (string, error) {
	c := r.get_conn()
	defer c.Close()

	return redis.String(c.Do("GET", id))
}

func (r *redis_cli) DelKey(id string) error {
	c := r.get_conn()
	defer c.Close()
	_, err := c.Do("DEL", id)
	return err
}

func (r *redis_cli) ExistKey(id string) (int, error) {
	c := r.get_conn()
	defer c.Close()
	exist, err := c.Do("EXISTS", id)

	return int(exist.(int64)), err
}

func (r *redis_cli) AddKey(key, value string) error {
	c := r.get_conn()
	defer c.Close()
	_, err := c.Do("SET", key, value)

	return err
}

func (r *redis_cli) PipeHGet(key string, subkeys []string) ([]string, error) {
	c := r.get_conn()
	defer c.Close()
	var index int
	var k string
	for index, k = range subkeys {
		c.Send("HGET", key, k)
	}
	c.Flush()
	result := make([]string, 0)
	var v string
	var err error
	for i := 0; i < index+1; i++ {
		v, err = redis.String(c.Receive())
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, nil
}
