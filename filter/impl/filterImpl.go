package impl

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

//var (
//	server string = "127.0.0.1:6379"
//	n uint = 100000
//	fp float64 = 0.01
//)

// redis连接池
//var pool *redis.Pool
func PoolInit(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

/*
// 测试使用，测试时按分钟切换
func GetVersion() int {
	return time.Now().Minute()%2
}
*/
// 获取当前BloomFilter版本,按月份切换
func GetVersion() int {
	var version int
	month_num := map[string]int{
		"January":   1,
		"February":  2,
		"March":     3,
		"April":     4,
		"May":       5,
		"June":      6,
		"July":      7,
		"August":    8,
		"September": 9,
		"October":   10,
		"November":  11,
		"December":  12,
	}
	month := time.Now().Month().String()
	version = month_num[month] % 2
	return version
}

// 按key设置存储中的values
func SetFunc(pool *redis.Pool, key string, value []byte) error {
	db := pool.Get()
	defer db.Close()
	_, err := redis.String(db.Do("SET", key, value))
	if err != nil {
		return err
	}
	return err
}

// 按key从存储中获取values
func GetFunc(pool *redis.Pool, key string) ([]byte, error) {
	db := pool.Get()
	defer db.Close()
	res, err := redis.Bytes(db.Do("GET", key))
	if err != nil {
		return res, err
	}
	return res, err
}
