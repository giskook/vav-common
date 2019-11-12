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
	case 12:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11]))
		break
	case 13:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12]))
		break
	case 14:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13]))
		break
	case 15:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14]))
		break
	case 16:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15]))
		break
	case 17:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15], args[16]))
		break
	case 18:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15], args[16], args[17]))
		break
	case 19:
		return redis.Int(s.Do(c, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14], args[15], args[16], args[17], args[18]))
		break
	}

	return 0, errors.New("over the max args")
}
