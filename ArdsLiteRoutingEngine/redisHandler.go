package main

import (
	"encoding/json"
	"fmt"
	"github.com/fzzy/radix/redis"
	"os"
	"strconv"
	"time"
)

var redisIp string
var redisDb int
var ardsUrl string
var resCsUrl string

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}

func GetDefaultConfig() Configuration {
	fileDefault, _ := os.Open("conf.json")
	defdecoder := json.NewDecoder(fileDefault)
	defconfiguration := Configuration{}
	deferr := defdecoder.Decode(&defconfiguration)

	if deferr != nil {
		fmt.Println("error:", deferr)
		defconfiguration.RedisIp = "127.0.0.1:6379"
		defconfiguration.RedisDb = 6
		defconfiguration.ArdsContinueUrl = "http://localhost:2221/continueArds/continue"
		defconfiguration.ResourceCSlotUrl = "http://localhost:2225/DVP/API/1.0.0.0/ARDS/resource"
	}

	return defconfiguration
}

func LoadDefaultConfig() {
	fileDefault, _ := os.Open("conf.json")
	defdecoder := json.NewDecoder(fileDefault)
	defconfiguration := Configuration{}
	deferr := defdecoder.Decode(&defconfiguration)

	if deferr != nil {
		fmt.Println("error:", deferr)
		redisIp = "127.0.0.1:6379"
		redisDb = 6
		ardsUrl = "http://localhost:2221/continueArds/continue"
		resCsUrl = "http://localhost:2225/DVP/API/1.0.0.0/ARDS/resource"
	} else {
		redisIp = defconfiguration.RedisIp
		redisDb = defconfiguration.RedisDb
		ardsUrl = defconfiguration.ArdsContinueUrl
		resCsUrl = defconfiguration.ResourceCSlotUrl
	}
}

func InitiateRedis() {
	fileEnv, _ := os.Open("custom-environment-variables.json")
	envdecoder := json.NewDecoder(fileEnv)
	envconfiguration := EnvConfiguration{}
	enverr := envdecoder.Decode(&envconfiguration)
	if enverr != nil {
		fmt.Println("error:", enverr)
		LoadDefaultConfig()
	} else {
		var converr error
		defConfig := GetDefaultConfig()
		redisIp = os.Getenv(envconfiguration.RedisIp)
		redisDb, converr = strconv.Atoi(os.Getenv(envconfiguration.RedisDb))
		ardsUrl = os.Getenv(envconfiguration.ArdsContinueUrl)
		resCsUrl = os.Getenv(envconfiguration.ResourceCSlotUrl)

		if redisIp == "" {
			redisIp = defConfig.RedisIp
		}
		if redisDb == 0 || converr != nil {
			redisDb = defConfig.RedisDb
		}
		if ardsUrl == "" {
			ardsUrl = defConfig.ArdsContinueUrl
		}
		if resCsUrl == "" {
			resCsUrl = defConfig.ResourceCSlotUrl
		}
	}

	fmt.Println("RedisIp:", redisIp)
	fmt.Println("RedisDb:", redisDb)
	fmt.Println("ArdsUrl:", ardsUrl)
	fmt.Println("ResCsUrl:", resCsUrl)

}

// Redis String Methods
func RedisGet(key string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisGet", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	strObj, _ := client.Cmd("get", key).Str()
	fmt.Println(strObj)
	return strObj
}

func RedisSearchKeys(pattern string) []string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSearchKeys", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	strObj, _ := client.Cmd("keys", pattern).List()
	return strObj
}

// Redis Hashes Methods

func RedisHashGetAll(hkey string) map[string]string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashGetAll", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	strHash, _ := client.Cmd("hgetall", hkey).Hash()
	return strHash
}

func RedisHashSetField(hkey, field, value string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashSetField", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	result, _ := client.Cmd("hset", hkey, field, value).Bool()
	return result
}

func RedisRemoveHashField(hkey, field string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisRemoveHashField", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	result, _ := client.Cmd("hdel", hkey, field).Bool()
	return result
}

// Redis List Methods

func RedisListLpop(lname string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpop", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	lpopItem, _ := client.Cmd("lpop", lname).Str()
	fmt.Println(lpopItem)
	return lpopItem
}

func RedisListLpush(lname, value string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpush", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	result, _ := client.Cmd("lpush", lname, value).Bool()
	return result
}
