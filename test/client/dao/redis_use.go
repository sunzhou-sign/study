package dao

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type redisClient struct {
	conn redis.Conn
}

func GetClient() *redisClient {
	cli := &redisClient{
		conn: redisPool.Get(),
	}

	if err := cli.Ping(); err != nil {
		Init()
		cli = &redisClient{
			conn: redisPool.Get(),
		}
	}

	return cli
}

func (r *redisClient) Get(key string) (string, error) {
	v, err := redis.String(r.conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return v, nil
}

func (r *redisClient) sPop(key string, count int) ([]string, error) {
	data, err := redis.Strings(r.conn.Do("SPOP", key, count))
	if err != nil {
		fmt.Println("SPOP err: ", err.Error())
		return nil, err
	}
	return data, err
}

func (r *redisClient) Del(key string) error {
	_, err := r.conn.Do("DEL", key)
	if err != nil {
		fmt.Println("Del err: ", err.Error())
	}
	return err
}

func (r *redisClient) Exists(key string) (bool, error) {
	flag, err := redis.Bool(r.conn.Do("exists", key))

	return flag, err
}

func (r *redisClient) HSetWithExpTs(key, field, value string, expTs int) error {
	_, err := r.conn.Do("HSET", key, field, value)
	if err != nil {
		fmt.Println("HSET err: ", err.Error())
		return err
	}
	if expTs > 0 {
		_, err = r.conn.Do("EXPIRE", key, expTs)
	}

	return err
}

func (r *redisClient) HExists(key, field string) (bool, error) {
	flag, err := redis.Bool(r.conn.Do("HEXISTS", key, field))

	return flag, err
}

func (r *redisClient) Close() error {
	return r.conn.Close()
}

func (r *redisClient) Ping() error {
	return r.conn.Err()
}

func (r *redisClient) Set(key string, value string, expTs int) {
	_, err := r.conn.Do("SET", key, value)
	if err != nil {
		fmt.Println("SET err: ", err.Error())
		return
	}

	if expTs > 0 {
		r.conn.Do("EXPIRE", key, expTs)
	}
}
