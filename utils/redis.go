package utils

import (
	"JD/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func SetConform(n string, x string) bool {
	//参照 官方示例 写个redis 短信验证 过期检验 很简单 很烂
	c := RedisPool.Get()
	defer c.Close()
	_, err := c.Do("SET", n, x, "EX", 300)
	if err != nil {
		return false
	}
	return true
}
func GetConform(Number string, ConformCode string) error {
	//
	c := RedisPool.Get()
	//
	defer c.Close()
	ok, err := redis.Bool(c.Do("EXISTS", Number))
	if !ok {
		err = errors.New("验证码过期或未发送")
		return err
	}
	v, err := redis.String(c.Do("GET", Number))
	if err == redis.ErrNil {
		err = errors.New("验证码超时")
		return err
	}
	if v != ConformCode {
		err = errors.New("验证码错误")
		return err
	}
	return nil
}

func SetToken(clientId int, token string) bool {
	c := RedisPool.Get()
	defer c.Close()
	//与token 验证时间一致
	_, err := redis.String(c.Do("SET", clientId, token, "EX", 1800))
	if err != nil {
		return false
	}
	_, err = redis.String(c.Do("SET", token, "", "EX", 1800))

	return true

}
func ConformToken(token string) bool {
	c := RedisPool.Get()
	defer c.Close()
	ok, err := redis.Bool(c.Do("EXISTS", token))
	if err != nil {
		return false
	}
	return ok
}
func GET(client string) (string, error) {
	c := RedisPool.Get()
	defer c.Close()

	token, err := redis.String(c.Do("GET", client))
	if err != nil {
		return "", err
	}
	return token, nil
}
func DeleteToken(token string) bool {
	c := RedisPool.Get()
	defer c.Close()
	ok, err := redis.Bool(c.Do("DEL", token))
	if err != nil {
		return false
	}
	return ok
}

func HashSet(info models.Request) {
	c := RedisPool.Get()
	jsoninfo, _ := json.Marshal(info)
	defer c.Close()
	ok, err := c.Do("SET", info.ClientId, jsoninfo)
	if err != nil {

	}
	fmt.Println(ok)
}
func HashGet(ClintId string) (models.Request, error) {
	c := RedisPool.Get()
	defer c.Close()
	info, err := redis.String(c.Do("GET", ClintId))
	if err != nil {
		return models.Request{}, err
	}
	fmt.Println(info)
	var clientInfo models.Request
	err = json.Unmarshal([]byte(info), &clientInfo)
	return clientInfo, nil
}
