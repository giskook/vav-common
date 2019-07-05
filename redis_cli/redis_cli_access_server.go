package redis_cli

import (
	"github.com/gomodule/redigo/redis"
)

const (
	ACCESS_SERVER_IP   string = "ACCESS_SERVER_IP"
	ACCESS_SERVER_PORT string = "ACCESS_SERVER_PORT"
)

func (r *redis_cli) SetAccessServer(access_server_key, ip, port string) error {
	c := r.get_conn()
	defer c.Close()
	_, err := c.Do("HMSET", access_server_key, ACCESS_SERVER_IP, ip, ACCESS_SERVER_PORT, port)
	if err != nil {
		return err
	}

	return nil
}

func (r *redis_cli) GetAccessServer(access_server_key string) (string, string, error) {
	c := r.get_conn()
	defer c.Close()

	addr, err := redis.Strings(c.Do("HMGET", access_server_key, ACCESS_SERVER_IP, ACCESS_SERVER_PORT))
	if err != nil {
		return "", "", err
	}

	return addr[0], addr[1], nil
}
