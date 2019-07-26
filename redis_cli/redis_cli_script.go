package redis_cli

import (
	"errors"
	"github.com/gomodule/redigo/redis"
)

func (r *redis_cli) DoScript(script string, args ...string) (int, error) {
	c := r.get_conn()
	defer c.Close()

	count := len(args)
	s := redis.NewScript(count, script)

	switch count {
	case 1:
		return redis.Int(s.Do(c, args[0]))
		break
	case 2:
		return redis.Int(s.Do(c, args[0], args[1]))
		break
	case 3:
		return redis.Int(s.Do(c, args[0], args[1], args[2]))
		break
	case 4:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3]))
		break
	case 5:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4]))
		break
	case 6:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5]))
		break
	case 7:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6]))
		break
	case 8:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7]))
		break
	case 9:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8]))
		break
	case 10:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9]))
		break
	case 11:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10]))
		break
	}

	return 0, errors.New("over the max args")
}
