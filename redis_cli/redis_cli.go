package redis_cli

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"sync"
	"time"
)

const (
	DATA_TYPE_AUDIO_VIDEO              string = "0"
	DATA_TYPE_VIDEO                    string = "1"
	DATA_TYPE_TWO_WAY_INTERCOM         string = "2"
	DATA_TYPE_LISTEN                   string = "3"
	DATA_TYPE_BROADCAST                string = "4"
	DATA_TYPE_TRANSPARENT_TRANSMISSION string = "5"
)

type Conf struct {
	Addr         string
	Passwd       string
	MaxIdle      int
	ConnTimeOut  time.Duration
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

type redis_cli struct {
	conf *Conf
	pool *redis.Pool
}

var instance *redis_cli
var once sync.Once

func (r *redis_cli) dial() (redis.Conn, error) {
	opt_conn_timeout := redis.DialConnectTimeout(r.conf.ConnTimeOut)
	opt_read_timeout := redis.DialReadTimeout(r.conf.ReadTimeOut)
	opt_write_timeout := redis.DialWriteTimeout(r.conf.WriteTimeOut)
	opt_dial_passwd := redis.DialPassword(r.conf.Passwd)
	c, err := redis.Dial("tcp", r.conf.Addr, opt_conn_timeout, opt_read_timeout, opt_write_timeout, opt_dial_passwd)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return c, err
}

func (r *redis_cli) test_on_borrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}

	_, err := c.Do("PING")

	return err
}

func (r *redis_cli) Init(conf *Conf) {
	r.conf = conf
	r.pool = redis.NewPool(r.dial, r.conf.MaxIdle)
	r.pool.TestOnBorrow = r.test_on_borrow
}

func (socket *redis_cli) get_conn() redis.Conn {
	return socket.pool.Get()
}

func (socket *redis_cli) Close() {
	socket.pool.Close()
}

func GetInstance() *redis_cli {
	once.Do(func() {
		instance = &redis_cli{}
	})

	return instance
}
