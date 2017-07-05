package main

import (
	"encoding/json"
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var dirPath string
var redisIp string
var redisPort string
var redisDb int
var locationDb int
var redisPassword string
var port string
var accessToken string
var rabbitMQIp string
var rabbitMQPort string
var rabbitMQUser string
var rabbitMQPassword string
var useMsgQueue bool
var routingEngineId string

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}

func GetDirPath() string {
	envPath := os.Getenv("GO_CONFIG_DIR")
	if envPath == "" {
		envPath = "./"
	}
	fmt.Println(envPath)
	return envPath
}

func GetDefaultConfig() Configuration {
	confPath := filepath.Join(dirPath, "conf.json")
	fmt.Println("GetDefaultConfig config path: ", confPath)
	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	defconfiguration := Configuration{}
	deferr := json.Unmarshal(content, &defconfiguration)

	if deferr != nil {
		fmt.Println("error:", deferr)
		defconfiguration.RedisIp = "192.168.3.200"
		defconfiguration.RedisPort = "6379"
		defconfiguration.RedisDb = 6
		defconfiguration.LocationDb = 0
		defconfiguration.RedisPassword = "DuoS123"
		defconfiguration.Port = "2226"
		defconfiguration.RabbitMQIp = "45.55.142.207"
		defconfiguration.RabbitMQPort = "5672"
		defconfiguration.RabbitMQUser = "guest"
		defconfiguration.RabbitMQPassword = "guest"
		defconfiguration.UseMsgQueue = false
		defconfiguration.AccessToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJoZXNoYW5pbmRpa2EiLCJqdGkiOiIwZmIyNDJmZS02OGQwLTQ1MjEtOTM5NS0xYzE0M2M3MzNmNmEiLCJzdWIiOiI1NmE5ZTc1OWZiMDcxOTA3YTAwMDAwMDEyNWQ5ZTgwYjVjN2M0Zjk4NDY2ZjkyMTE3OTZlYmY0MyIsImV4cCI6MTQ1Njg5NDE5NSwidGVuYW50IjoxLCJjb21wYW55Ijo1LCJzY29wZSI6W3sicmVzb3VyY2UiOiJhbGwifSx7InJlc291cmNlIjoicmVxdWVzdHNlcnZlciIsImFjdGlvbnMiOlsicmVhZCIsIndyaXRlIiwiZGVsZXRlIl19LHsicmVzb3VyY2UiOiJyZXF1ZXN0bWV0YSIsImFjdGlvbnMiOlsicmVhZCIsIndyaXRlIiwiZGVsZXRlIl19LHsicmVzb3VyY2UiOiJhcmRzcmVzb3VyY2UiLCJhY3Rpb25zIjpbInJlYWQiLCJ3cml0ZSIsImRlbGV0ZSJdfSx7InJlc291cmNlIjoiYXJkc3JlcXVlc3QiLCJhY3Rpb25zIjpbInJlYWQiLCJ3cml0ZSIsImRlbGV0ZSJdfV0sImlhdCI6MTQ1NjI4OTM5NX0.AWZuYNtj4lHfxpTQCutswUfUsJXwTMVPUmqTjFdVXSk"
		defconfiguration.RoutingEngineId = "1"
	}

	return defconfiguration
}

func LoadDefaultConfig() {
	defconfiguration := GetDefaultConfig()

	redisIp = fmt.Sprintf("%s:%s", defconfiguration.RedisIp, defconfiguration.RedisPort)
	redisPort = defconfiguration.RedisPort
	redisDb = defconfiguration.RedisDb
	locationDb = defconfiguration.LocationDb
	redisPassword = defconfiguration.RedisPassword
	port = defconfiguration.Port
	rabbitMQIp = defconfiguration.RabbitMQIp
	rabbitMQPort = defconfiguration.RabbitMQPort
	rabbitMQUser = defconfiguration.RabbitMQUser
	rabbitMQPassword = defconfiguration.RabbitMQPassword
	useMsgQueue = defconfiguration.UseMsgQueue
	accessToken = defconfiguration.AccessToken
	routingEngineId = defconfiguration.RoutingEngineId
}

func InitiateRedis() {
	dirPathtest, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(dirPathtest)
	dirPath = GetDirPath()
	confPath := filepath.Join(dirPath, "custom-environment-variables.json")
	fmt.Println("InitiateRedis config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	envconfiguration := EnvConfiguration{}
	enverr := json.Unmarshal(content, &envconfiguration)

	if enverr != nil {
		fmt.Println("error:", enverr)
		LoadDefaultConfig()
	} else {
		var converr1 error
		var converr2 error
		var converr3 error
		defConfig := GetDefaultConfig()
		redisIp = os.Getenv(envconfiguration.RedisIp)
		redisPort = os.Getenv(envconfiguration.RedisPort)
		redisDb, converr1 = strconv.Atoi(os.Getenv(envconfiguration.RedisDb))
		locationDb, converr2 = strconv.Atoi(os.Getenv(envconfiguration.LocationDb))
		redisPassword = os.Getenv(envconfiguration.RedisPassword)
		rabbitMQIp = os.Getenv(envconfiguration.RabbitMQIp)
		rabbitMQPort = os.Getenv(envconfiguration.RabbitMQPort)
		rabbitMQUser = os.Getenv(envconfiguration.RabbitMQUser)
		rabbitMQPassword = os.Getenv(envconfiguration.RabbitMQPassword)
		port = os.Getenv(envconfiguration.Port)
		useMsgQueue, converr3 = strconv.ParseBool(os.Getenv(envconfiguration.UseMsgQueue))
		accessToken = os.Getenv(envconfiguration.AccessToken)
		routingEngineId = os.Getenv(envconfiguration.RoutingEngineId)

		if redisIp == "" {
			redisIp = defConfig.RedisIp
		}
		if redisPort == "" {
			redisPort = defConfig.RedisPort
		}
		if converr1 != nil {
			redisDb = defConfig.RedisDb
		}
		if converr2 != nil {
			locationDb = defConfig.LocationDb
		}
		if redisPassword == "" {
			redisPassword = defConfig.RedisPassword
		}
		if port == "" {
			port = defConfig.Port
		}
		if rabbitMQIp == "" {
			rabbitMQIp = defConfig.RabbitMQIp
		}
		if rabbitMQPort == "" {
			rabbitMQPort = defConfig.RabbitMQPort
		}
		if rabbitMQUser == "" {
			rabbitMQUser = defConfig.RabbitMQUser
		}
		if rabbitMQPassword == "" {
			rabbitMQPassword = defConfig.RabbitMQPassword
		}
		if converr3 != nil {
			useMsgQueue = defConfig.UseMsgQueue
		}
		if accessToken == "" {
			accessToken = defConfig.AccessToken
		}
		if routingEngineId == "" {
			routingEngineId = defConfig.RoutingEngineId
		}

		redisIp = fmt.Sprintf("%s:%s", redisIp, redisPort)
	}

	fmt.Println("RedisIp:", redisIp)
	fmt.Println("RedisDb:", redisDb)
	fmt.Println("LocationDb:", locationDb)

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
	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	strObj, _ := client.Cmd("get", key).Str()
	fmt.Println(strObj)
	return strObj
}

func RedisGet_v1(key string) (strObj string, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisGet", r)
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("Recovered in RedisGet: %v", r)
			}
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()
	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	strObj, err = client.Cmd("get", key).Str()
	fmt.Println(strObj)
	return
}

func RedisSearchKeys(pattern string) []string {
	fmt.Println("Start RedisSearchKeys")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSearchKeys", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	strObj, _ := client.Cmd("keys", pattern).List()
	return strObj
}

func RedisSetNx(key, value string, timeout int) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSetNx", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	tmpValue, _ := client.Cmd("set", key, value, "nx", "ex", timeout).Str()
	if tmpValue == "OK" {
		fmt.Println("GetRLock: ", true)
		return true
	} else {
		fmt.Println("GetRLock: ", false)
		return false
	}

}

/*func RedisSetEx(key, value string, timeSec int) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSetEx", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	strObj, _ := client.Cmd("setex", key, timeSec, value).Bool()
	fmt.Println("setex: ", strObj)
	return strObj
}*/

func RedisRemoveRLock(key, value string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisRemoveRLock", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)
	luaScript := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
	tmpValue, _ := client.Cmd("eval", luaScript, 1, key, value).Int()
	if tmpValue == 1 {
		fmt.Println("GetRLock: ", true)
		return true
	} else {
		fmt.Println("GetRLock: ", false)
		return false
	}
}

func RedisRemove(key string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisRemove", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	tmpValue, _ := client.Cmd("del", key).Int()

	if tmpValue == 1 {
		fmt.Println("GetRLock: ", true)
		return true
	} else {
		fmt.Println("GetRLock: ", false)
		return false
	}
}

func RedisCheckKeyExist(key string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in CheckKeyExist", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	tmpValue, sErr := client.Cmd("exists", key).Int()
	errHndlr(sErr)

	if tmpValue == 1 {
		fmt.Println("GetRLock: ", true)
		return true
	} else {
		fmt.Println("GetRLock: ", false)
		return false
	}
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

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	strHash, _ := client.Cmd("hgetall", hkey).Map()
	return strHash
}

func RedisHashGetValue(hkey, queueId string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashGetValue", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	strHash := client.Cmd("hget", hkey, queueId).String()
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

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	tmpValue, _ := client.Cmd("hset", hkey, field, value).Int()

	if tmpValue == 1 {
		fmt.Println("GetRLock: ", true)
		return true
	} else {
		fmt.Println("GetRLock: ", false)
		return false
	}
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

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	tmpValue, _ := client.Cmd("hdel", hkey, field).Int()

	if tmpValue == 1 {
		fmt.Println("GetRLock: ", true)
		return true
	} else {
		fmt.Println("GetRLock: ", false)
		return false
	}
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

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	lpopItem, _ := client.Cmd("lpop", lname).Str()
	fmt.Println(lpopItem)
	return lpopItem
}

/*func RedisListLpush(lname, value string) bool {
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
}*/

/*-----------------------------Geo methods--------------------------------------*/

func RedisGeoRadius(tenant, company int, locationObj ReqLocationData) *redis.Resp {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisGeoRadius", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)
	// select database
	r := client.Cmd("select", locationDb)
	errHndlr(r.Err)

	locationInfoKey := fmt.Sprintf("location:%d:%d", tenant, company)
	fmt.Println("locationInfoKey: ", locationInfoKey)
	locationResult := client.Cmd("georadius", "positions", locationObj.Longitude, locationObj.Latitude, locationObj.Radius, locationObj.Metric, "WITHDIST", "ASC")
	fmt.Println(locationResult)
	return locationResult
}

func RoutingEngineDistribution(pubChannelName string) string {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RoutingEngineDistribution", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	//authServer
	client.Cmd("auth", redisPassword)
	//errHndlr(authE.Err)

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	activeRoutingKey, _ := client.Cmd("get", "ActiveRoutingEngine").Str()

	if activeRoutingKey == "" {
		u1 := uuid.NewV4()
		if RedisSetNx("ActiveRoutingEngineLock", u1.String(), 30) == true {
			RedisSetNx("ActiveRoutingEngine", pubChannelName, 60)
			RedisRemoveRLock("ActiveRoutingEngineLock", u1.String())

			return pubChannelName
		} else {
			fmt.Println("Aquire ActiveRoutingEngineLock failed")
			return activeRoutingKey
		}

	} else {

		if activeRoutingKey == pubChannelName {
			expire, _ := client.Cmd("expire", "ActiveRoutingEngine", 60).Int()
			if expire == 1 {
				fmt.Println("Extend Active Routing Engine Expire Time Success")
			} else {
				fmt.Println("Extend Active Routing Engine Expire Time Failed")
			}
		}

		return activeRoutingKey

	}

}
