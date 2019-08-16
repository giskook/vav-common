package redis_cli

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func (r *redis_cli) Pub(topic, content string) error {
	c := r.get_conn()
	defer c.Close()

	sub_count, err := redis.Int(c.Do("PUBLISH", topic, content))
	if err != nil {
		return err
	}
	if sub_count == 0 {
		return errors.New("there is no subscribers")
	}

	return nil
}

func (r *redis_cli) Sub(topic string, exit chan struct{}, on_message func([]byte) error) error {
	c := r.get_conn()
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe(topic)

	for {
		select {
		case <-exit:
			return errors.New("closed")
		default:
			switch v := psc.Receive().(type) {
			case redis.Message:
				fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
				err := on_message(v.Data)
				if err != nil {
					return err
				}
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				return v
			}
		}
	}
}
