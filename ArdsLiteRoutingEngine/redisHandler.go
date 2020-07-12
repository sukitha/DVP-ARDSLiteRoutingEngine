package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/mediocregopher/radix.v2/sentinel"
	"github.com/mediocregopher/radix.v2/util"
	uuid "github.com/satori/go.uuid"
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
var useMsgQueue string
var routingEngineId string
var redisMode string
var redisClusterName string
var sentinelHosts string
var sentinelPort string
var ardsServiceHost string
var ardsServicePort string
var useAmqpAdapter string
var useDynamicPort string

var sentinelPool *sentinel.Client
var redisPool *pool.Pool

func errHndlr(err error) {
	if err != nil {
		log.Println("error:", err)
	}
}

func errHndlrNew(errorFrom, command string, err error) {
	if err != nil {
		log.Println("error:", errorFrom, ":: ", command, ":: ", err)
	}
}

func GetDirPath() string {
	envPath := os.Getenv("GO_CONFIG_DIR")
	if envPath == "" {
		envPath = "./"
	}
	log.Println(envPath)
	return envPath
}

func GetDefaultConfig() Configuration {
	confPath := filepath.Join(dirPath, "conf.json")
	log.Println("GetDefaultConfig config path: ", confPath)
	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		log.Println(operr)
		panic(operr)
	}

	defconfiguration := Configuration{}
	deferr := json.Unmarshal(content, &defconfiguration)

	if deferr != nil {
		panic(deferr)
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
	redisMode = defconfiguration.RedisMode
	redisClusterName = defconfiguration.RedisClusterName
	sentinelHosts = defconfiguration.SentinelHosts
	sentinelPort = defconfiguration.SentinelPort
	ardsServiceHost = defconfiguration.ArdsServiceHost
	ardsServicePort = defconfiguration.ArdsServicePort
	useAmqpAdapter = defconfiguration.UseAmqpAdapter
	useDynamicPort = defconfiguration.UseDynamicPort
}

func InitiateRedis() {
	dirPathtest, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println(dirPathtest)
	dirPath = GetDirPath()
	confPath := filepath.Join(dirPath, "custom-environment-variables.json")
	log.Println("InitiateRedis config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		log.Println(operr)
	}

	envconfiguration := EnvConfiguration{}
	enverr := json.Unmarshal(content, &envconfiguration)

	if enverr != nil {
		log.Println("error:", enverr)
		LoadDefaultConfig()
	} else {
		var converr1 error
		var converr2 error
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
		useMsgQueue = os.Getenv(envconfiguration.UseMsgQueue)
		accessToken = os.Getenv(envconfiguration.AccessToken)
		routingEngineId = os.Getenv(envconfiguration.RoutingEngineId)
		redisMode = os.Getenv(envconfiguration.RedisMode)
		redisClusterName = os.Getenv(envconfiguration.RedisClusterName)
		sentinelHosts = os.Getenv(envconfiguration.SentinelHosts)
		sentinelPort = os.Getenv(envconfiguration.SentinelPort)
		ardsServiceHost = os.Getenv(envconfiguration.ArdsServiceHost)
		ardsServicePort = os.Getenv(envconfiguration.ArdsServicePort)
		useAmqpAdapter = os.Getenv(envconfiguration.UseAmqpAdapter)
		useDynamicPort = os.Getenv(envconfiguration.UseDynamicPort)

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
		if useMsgQueue == "" {
			useMsgQueue = defConfig.UseMsgQueue
		}
		if accessToken == "" {
			accessToken = defConfig.AccessToken
		}
		if routingEngineId == "" {
			routingEngineId = defConfig.RoutingEngineId
		}
		if redisMode == "" {
			redisMode = defConfig.RedisMode
		}
		if redisClusterName == "" {
			redisClusterName = defConfig.RedisClusterName
		}
		if sentinelHosts == "" {
			sentinelHosts = defConfig.SentinelHosts
		}
		if sentinelPort == "" {
			sentinelPort = defConfig.SentinelPort
		}
		if ardsServiceHost == "" {
			ardsServiceHost = defConfig.ArdsServiceHost
		}
		if ardsServicePort == "" {
			ardsServicePort = defConfig.ArdsServicePort
		}
		if useAmqpAdapter == "" {
			useAmqpAdapter = defConfig.UseAmqpAdapter
		}
		if useDynamicPort == "" {
			useDynamicPort = defConfig.UseDynamicPort
		}

		redisIp = fmt.Sprintf("%s:%s", redisIp, redisPort)
	}

	log.Println("RoutingEngineId:", routingEngineId)
	log.Println("RedisMode:", redisMode)
	log.Println("RedisIp:", redisIp)
	log.Println("RedisDb:", redisDb)
	log.Println("LocationDb:", locationDb)
	log.Println("SentinelHosts:", sentinelHosts)
	log.Println("SentinelPort:", sentinelPort)
	log.Println("useMsgQueue:", useMsgQueue)
	log.Println("useAmqpAdapter:", useAmqpAdapter)

	var err error

	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		if redisPassword != "" {
			if err = client.Cmd("AUTH", redisPassword).Err; err != nil {
				client.Close()
				return nil, err
			}
		}
		if err = client.Cmd("select", redisDb).Err; err != nil {
			client.Close()
			return nil, err
		}
		return client, nil
	}

	if redisMode == "sentinel" {

		sentinelIps := strings.Split(sentinelHosts, ",")

		if len(sentinelIps) > 1 {
			sentinelIp := fmt.Sprintf("%s:%s", sentinelIps[0], sentinelPort)
			sentinelPool, err = sentinel.NewClientCustom("tcp", sentinelIp, 20, df, redisClusterName)
			if err != nil {
				log.Println("InitiateSentinel ::", err)
			}
		} else {
			log.Println("Not enough sentinel servers")
		}

	} else {
		redisPool, err = pool.NewCustom("tcp", redisIp, 10, df)

		if err != nil {
			errHndlrNew("InitiateRedis", "InitiatePool", err)
		}
	}

}

// Redis String Methods
func RedisSet(key, value string) string {

	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisSet", r)
		}
		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}

	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	strObj, _ := client.Cmd("set", key, value).Str()
	//log.Println(strObj)
	return strObj

}

func RedisGet(key string) string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisGet", r)
		}
		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}

	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	strObj, _ := client.Cmd("get", key).Str()
	//log.Println(strObj)
	return strObj
	/*if redisMode == "instance" {
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
		log.Println(strObj)
		return strObj
	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		strObj, _ := client.Cmd("get", key).Str()
		log.Println(strObj)
		return strObj
	}*/

}

func RedisGet_v1(key string) (strObj string, err error) {
	var client *redis.Client

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisGet", r)
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("Recovered in RedisGet: %v", r)
			}
		}
		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}

	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	strObj, err = client.Cmd("get", key).Str()
	//log.Println(strObj)
	return

	/*if redisMode == "instance" {

		client, err1 := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
		errHndlr(err1)
		defer client.Close()
		//authServer
		client.Cmd("auth", redisPassword)
		//errHndlr(authE.Err)
		// select database
		r := client.Cmd("select", redisDb)
		errHndlr(r.Err)

		strObj, err = client.Cmd("get", key).Str()
		log.Println(strObj)
		return

	} else {
		client, err2 := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err2)
		defer sentinelPool.PutMaster(redisClusterName, client)

		strObj, err = client.Cmd("get", key).Str()
		log.Println(strObj)
		return
	}*/

}

func RedisSearchKeys(pattern string) []string {
	//log.Println("Start RedisSearchKeys")
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisSearchKeys", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	matchingKeys := make([]string, 0)

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("ScanAndGetKeys", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("ScanAndGetKeys", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	log.Println("Start ScanAndGetKeys:: ", pattern)
	scanResult := util.NewScanner(client, util.ScanOpts{Command: "SCAN", Pattern: pattern, Count: 1000})

	for scanResult.HasNext() {
		//log.Println("next:", scanResult.Next())
		matchingKeys = AppendIfMissing(matchingKeys, scanResult.Next())
	}

	//log.Println("Scan Result:: ", matchingKeys)
	return matchingKeys
	/*if redisMode == "instance" {

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

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		strObj, _ := client.Cmd("keys", pattern).List()
		return strObj
	}*/

}

func RedisSetNx(key, value string, timeout int) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisSetNx", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	tmpValue, _ := client.Cmd("set", key, value, "nx", "ex", timeout).Str()
	if tmpValue == "OK" {
		log.Println("GetRLock: ", true)
		return true
	} else {
		log.Println("GetRLock: ", false)
		return false
	}

	/*if redisMode == "instance" {

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
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		tmpValue, _ := client.Cmd("set", key, value, "nx", "ex", timeout).Str()
		if tmpValue == "OK" {
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}
	}*/

}

/*func RedisSetEx(key, value string, timeSec int) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisSetEx", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", redisIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)

	strObj, _ := client.Cmd("setex", key, timeSec, value).Bool()
	log.Println("setex: ", strObj)
	return strObj
}*/

func RedisRemoveRLock(key, value string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisRemoveRLock", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	luaScript := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
	tmpValue, _ := client.Cmd("eval", luaScript, 1, key, value).Int()
	if tmpValue == 1 {
		log.Println("GetRLock: ", true)
		return true
	} else {
		log.Println("GetRLock: ", false)
		return false
	}

	/*if redisMode == "instance" {

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
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		luaScript := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
		tmpValue, _ := client.Cmd("eval", luaScript, 1, key, value).Int()
		if tmpValue == 1 {
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}
	}*/

}

func RedisRemove(key string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisRemove", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	tmpValue, _ := client.Cmd("del", key).Int()

	if tmpValue == 1 {
		//log.Println("GetRLock: ", true)
		return true
	} else {
		//log.Println("GetRLock: ", false)
		return false
	}

	/*if redisMode == "instance" {

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
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		tmpValue, _ := client.Cmd("del", key).Int()

		if tmpValue == 1 {
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}
	}*/

}

func RedisCheckKeyExist(key string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in CheckKeyExist", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	tmpValue, sErr := client.Cmd("exists", key).Int()
	errHndlr(sErr)

	if tmpValue == 1 {
		//log.Println("GetRLock: ", true)
		return true
	} else {
		//log.Println("GetRLock: ", false)
		return false
	}

	/*if redisMode == "instance" {

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
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		tmpValue, sErr := client.Cmd("exists", key).Int()
		errHndlr(sErr)

		if tmpValue == 1 {
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}
	}*/

}

// Redis Hashes Methods

func RedisHashGetAll(hkey string) map[string]string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisHashGetAll", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	strHash, _ := client.Cmd("hgetall", hkey).Map()
	return strHash

	/*if redisMode == "instance" {

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

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		strHash, _ := client.Cmd("hgetall", hkey).Map()
		return strHash
	}*/

}

func RedisHashGetValue(hkey, queueId string) string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisHashGetValue", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	strHash, _ := client.Cmd("hget", hkey, queueId).Str()
	return strHash

	/*if redisMode == "instance" {

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

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		strHash := client.Cmd("hget", hkey, queueId).String()
		return strHash
	}*/

}

func RedisHashSetField(hkey, field, value string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisHashSetField", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	tmpValue, _ := client.Cmd("hset", hkey, field, value).Int()

	if tmpValue == 1 {
		//log.Println("GetRLock: ", true)
		return true
	} else {
		//log.Println("GetRLock: ", false)
		return false
	}

	/*if redisMode == "instance" {

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
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		tmpValue, _ := client.Cmd("hset", hkey, field, value).Int()

		if tmpValue == 1 {
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}
	}*/

}

func RedisRemoveHashField(hkey, field string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisRemoveHashField", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	tmpValue, _ := client.Cmd("hdel", hkey, field).Int()

	if tmpValue == 1 {
		//log.Println("GetRLock: ", true)
		return true
	} else {
		//log.Println("GetRLock: ", false)
		return false
	}

	/*if redisMode == "instance" {

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
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		tmpValue, _ := client.Cmd("hdel", hkey, field).Int()

		if tmpValue == 1 {
			log.Println("GetRLock: ", true)
			return true
		} else {
			log.Println("GetRLock: ", false)
			return false
		}
	}*/

}

// Redis List Methods

func RedisListLpop(lname string) string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisListLpop", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	lpopItem, _ := client.Cmd("lpop", lname).Str()
	log.Println(lpopItem)
	return lpopItem

	/*if redisMode == "instance" {

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
		log.Println(lpopItem)
		return lpopItem

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		lpopItem, _ := client.Cmd("lpop", lname).Str()
		log.Println(lpopItem)
		return lpopItem
	}*/

}

/*func RedisListLpush(lname, value string) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisListLpush", r)
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
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RedisGeoRadius", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	locationInfoKey := fmt.Sprintf("location:%d:%d", tenant, company)
	log.Println("locationInfoKey: ", locationInfoKey)
	locationResult := client.Cmd("georadius", "positions", locationObj.Longitude, locationObj.Latitude, locationObj.Radius, locationObj.Metric, "WITHDIST", "ASC")
	log.Println(locationResult)
	return locationResult

	/*if redisMode == "instance" {

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
		log.Println("locationInfoKey: ", locationInfoKey)
		locationResult := client.Cmd("georadius", "positions", locationObj.Longitude, locationObj.Latitude, locationObj.Radius, locationObj.Metric, "WITHDIST", "ASC")
		log.Println(locationResult)
		return locationResult

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		r := client.Cmd("select", locationDb)
		errHndlr(r.Err)

		locationInfoKey := fmt.Sprintf("location:%d:%d", tenant, company)
		log.Println("locationInfoKey: ", locationInfoKey)
		locationResult := client.Cmd("georadius", "positions", locationObj.Longitude, locationObj.Latitude, locationObj.Radius, locationObj.Metric, "WITHDIST", "ASC")
		log.Println(locationResult)
		return locationResult
	}*/

}

func RoutingEngineDistribution(pubChannelName string) string {

	var client *redis.Client
	var err error
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RoutingEngineDistribution", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			log.Println("Cannot Put invalid connection")
		}

	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnReset", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnReset", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	log.Println("Sentinal Master.Addr: ", client.Addr)
	//log.Println("Sentinal Master.Network: ", client.Network)
	//log.Println("Sentinal Master.LastCritical.Error(): ", client.LastCritical.Error())

	activeRoutingKey, _ := client.Cmd("get", "ActiveRoutingEngine").Str()

	if activeRoutingKey == "" {
		u1 := uuid.NewV4()
		if RedisSetNx("ActiveRoutingEngineLock", u1.String(), 30) == true {
			if RedisSetNx("ActiveRoutingEngine", pubChannelName, 60) == true {
				RedisRemoveRLock("ActiveRoutingEngineLock", u1.String())
				return pubChannelName
			} else {
				RedisRemoveRLock("ActiveRoutingEngineLock", u1.String())
				return ""
			}
		} else {
			log.Println("Aquire ActiveRoutingEngineLock failed")
			return activeRoutingKey
		}

	} else {

		if activeRoutingKey == pubChannelName {
			expire, _ := client.Cmd("expire", "ActiveRoutingEngine", 60).Int()
			if expire == 1 {
				//log.Println("Extend Active Routing Engine Expire Time Success")
			} else {
				log.Println("Extend Active Routing Engine Expire Time Failed")
			}
		}

		return activeRoutingKey

	}

	/*if redisMode == "instance" {

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
				if RedisSetNx("ActiveRoutingEngine", pubChannelName, 60) == true {
					RedisRemoveRLock("ActiveRoutingEngineLock", u1.String())
					return pubChannelName
				} else {
					RedisRemoveRLock("ActiveRoutingEngineLock", u1.String())
					return ""
				}
			} else {
				log.Println("Aquire ActiveRoutingEngineLock failed")
				return activeRoutingKey
			}

		} else {

			if activeRoutingKey == pubChannelName {
				expire, _ := client.Cmd("expire", "ActiveRoutingEngine", 60).Int()
				if expire == 1 {
					log.Println("Extend Active Routing Engine Expire Time Success")
				} else {
					log.Println("Extend Active Routing Engine Expire Time Failed")
				}
			}

			return activeRoutingKey

		}

	} else {
		client, err := sentinelPool.GetMaster(redisClusterName)
		errHndlr(err)
		defer sentinelPool.PutMaster(redisClusterName, client)

		activeRoutingKey, _ := client.Cmd("get", "ActiveRoutingEngine").Str()

		if activeRoutingKey == "" {
			u1 := uuid.NewV4()
			if RedisSetNx("ActiveRoutingEngineLock", u1.String(), 30) == true {
				if RedisSetNx("ActiveRoutingEngine", pubChannelName, 60) == true {
					RedisRemoveRLock("ActiveRoutingEngineLock", u1.String())
					return pubChannelName
				} else {
					RedisRemoveRLock("ActiveRoutingEngineLock", u1.String())
					return ""
				}
			} else {
				log.Println("Aquire ActiveRoutingEngineLock failed")
				return activeRoutingKey
			}

		} else {

			if activeRoutingKey == pubChannelName {
				expire, _ := client.Cmd("expire", "ActiveRoutingEngine", 60).Int()
				if expire == 1 {
					log.Println("Extend Active Routing Engine Expire Time Success")
				} else {
					log.Println("Extend Active Routing Engine Expire Time Failed")
				}
			}

			return activeRoutingKey
		}

	}*/
}

func AppendIfMissing(windowList []string, i string) []string {
	for _, ele := range windowList {
		if ele == i {
			return windowList
		}
	}
	return append(windowList, i)
}
