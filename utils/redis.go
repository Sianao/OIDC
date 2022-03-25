package utils

import (
	"JD/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

// SmsConform 短信 验证使用
func SmsConform(n string, x string) bool {
	//参照 官方示例 写个redis 短信验证 过期检验 很简单 很烂
	c := RedisPool.Get()
	defer c.Close()
	_, err := c.Do("SET", n, x, "EX", 300)
	if err != nil {
		return false
	}
	return true
}
func GetSmsConform(Number string, ConformCode string) error {
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
func DeleteToken(token string) bool {
	c := RedisPool.Get()
	defer c.Close()
	ok, err := redis.Bool(c.Do("DEL", token))
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

func JsonSet(info models.Request) {
	c := RedisPool.Get()
	JsonInfo, _ := json.Marshal(info)
	defer c.Close()
	ok, err := c.Do("SET", info.ClientId, JsonInfo)
	if err != nil {

	}
	fmt.Println(ok)
}
func JsonGet(ClintId string) (models.Request, error) {
	c := RedisPool.Get()
	defer c.Close()
	info, err := redis.String(c.Do("GET", ClintId))
	if err != nil {
		return models.Request{}, err
	}
	var clientInfo models.Request
	err = json.Unmarshal([]byte(info), &clientInfo)
	return clientInfo, nil
}

// 类似两张表 一张表用于各个分类 用户订阅信息 存储
//另外一张 用于 用户 信息存储
// 发现 两张表 无法解决用户订阅信息的获取
//于是 再加一个 对订阅项目进行管理

func NewInfo(Info models.ChangeMessage) {
	c := RedisPool.Get()
	defer c.Close()
	// 获取所有成员
	member, _ := redis.Values(c.Do("SMEMBERS", Info.Category))
	fmt.Println(member, "fewfrwg")
	for _, v := range member {
		value, ok := v.([]byte)
		for _, v := range value {
			AddInfo(int(v), Info)
		}
		fmt.Println(ok)

	}

}
func AddNewInfo(category string, types string, info string, time time.Time) {
	s := models.ChangeMessage{}
	s.Info = info
	s.Type = types
	s.Category = category
	s.Time = time
	NewInfo(s)
}
func AddInfo(uid int, Info models.ChangeMessage) bool {
	c := RedisPool.Get()
	defer c.Close()
	Info.Msid = strconv.FormatInt(time.Now().Unix(), 10)
	Json, _ := json.Marshal(Info)
	ok, err := redis.Bool(c.Do("HMSET", byte(uid), Info.Msid, Json))
	if err != nil || !ok {
		return false
	}
	return true
}
func Subscribe(uid int, Category string) bool {
	c := RedisPool.Get()
	defer c.Close()

	ok, err := redis.Bool(c.Do("SADD", Category, uid))
	if err != nil || !ok {
		return false
	}
	return true
}

func GetInfo(Uid int) []models.ChangeMessage {
	c := RedisPool.Get()
	defer c.Close()
	var Allinfo []models.ChangeMessage
	all, err := redis.Values(c.Do("HKEYS", Uid))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var message models.ChangeMessage
	for _, v := range all {
		s, ok := v.([]byte)
		fmt.Println(s, ok)
		info, err := redis.String(c.Do("HGET", Uid, s))
		err = json.Unmarshal([]byte(info), &message)
		fmt.Println(err)
		if err != nil {
		} else {
			Allinfo = append(Allinfo, message)
		}
	}

	return Allinfo

}
func AddCategory(category string) error {
	c := RedisPool.Get()
	defer c.Close()
	ok, err := redis.Bool(c.Do("SISMEMBER", "all", category))
	redis.Bool(c.Do("SADD", category, ""))
	if err != nil || ok {
		err = errors.New("分类添加失败 或许已经存在")
		return err
	}
	ok, err = redis.Bool(c.Do("SADD", "all", category))
	if err != nil || !ok {
		err = errors.New("分类添加失败")
		return err
	}
	return nil
}

//显示所有订阅的分类

func AllCategory() []string {
	c := RedisPool.Get()
	defer c.Close()
	values, err := redis.Values(c.Do("SMEMBERS", "all"))
	var all []string
	fmt.Println(err)
	for _, v := range values {
		value, ok := v.([]byte)
		fmt.Println(ok)
		if ok {
			//这个坑了我 一直以为可以直接断言成 string
			value := string(value)
			all = append(all, value)
		}
	}
	return all
}
func MySubscribe(uid int) []string {

	all := AllCategory()
	c := RedisPool.Get()
	defer c.Close()
	var myall []string
	for _, v := range all {
		ok, _ := redis.Bool(c.Do("SISMEMBER", v, uid))
		if ok {
			myall = append(myall, v)
		}
	}
	return myall
}
func Unsubscribe(uid int, category string) error {
	c := RedisPool.Get()
	defer c.Close()
	ok, err := redis.Bool(c.Do("SREM", category, uid))
	if err != nil || !ok {
		err = errors.New("推定失败")
		return err
	}
	return nil
}
func MarkReaded(uid int, msid string) bool {
	c := RedisPool.Get()
	defer c.Close()
	ok, err := redis.Bool(c.Do("HDEL", []byte(string(uid)), msid))
	if err != nil || !ok {
		return false
	}
	return ok
}
