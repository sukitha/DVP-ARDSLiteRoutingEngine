package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
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

// var redisSentinel *radix.Sentinel
// var redisPool *radix.Pool;
// var redisScanClient radix.Client;
// var redisScanSentinel *radix.Sentinel
// var mu sync.Mutex

var connectionOptions redis.UniversalOptions
var redisCtx = context.Background()
var rdb redis.UniversalClient


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


	if redisMode == "sentinel" {

		sentinelIps := strings.Split(sentinelHosts, ",")
		var ips []string;

		if len(sentinelIps) > 1 {

			for _, ip := range sentinelIps{
				ipPortArray := strings.Split(ip, ":")
				sentinelIp := ip;
				if(len(ipPortArray) > 1){
					sentinelIp = fmt.Sprintf("%s:%s", ipPortArray[0], ipPortArray[1])
				}else{
					sentinelIp = fmt.Sprintf("%s:%s", ip, sentinelPort)
				}
				ips = append(ips, sentinelIp)
				
			}

			connectionOptions.Addrs = ips
			connectionOptions.MasterName = redisClusterName

		} else {
			fmt.Println("Not enough sentinel servers")
			os.Exit(0)
		}


	} else {

		redisIps := strings.Split(redisIp, ",")
		var ips []string;
		if len(redisIps) > 0 {

			for _, ip := range redisIps{
				ipPortArray := strings.Split(ip, ":")
				redisAddr := ip;
				if(len(ipPortArray) > 1){
					redisAddr = fmt.Sprintf("%s:%s", ipPortArray[0], ipPortArray[1])
				}else{
					redisAddr = fmt.Sprintf("%s:%s", ip, redisPort)
				}
				ips = append(ips, redisAddr)
				
			}

			connectionOptions.Addrs = ips

		} else {
			fmt.Println("Not enough redis servers")
			os.Exit(0)
		}
	}

	connectionOptions.DB = redisDb
	connectionOptions.Password = redisPassword

	rdb = redis.NewUniversalClient(&connectionOptions)


	// res, _ := rdb.Pipelined(context.TODO(), func(pipe redis.Pipeliner) error {
	// 	pipe.Set(context.TODO(),"aaa", "aaa", 0)
	// 	pipe.Get(context.TODO(), "aaa")
	// 	pipe.Get(context.TODO(), "bbb")
	// 	pipe.HSet(context.TODO(), "TESTx", "test", "Test1")
	// 	pipe.HGetAll(context.TODO(), "TESTx")
	// 	return nil
	// })

	// fmt.Println(res)

	
}


func RedisSet(key, value string) string {

	result, _:= rdb.Set(context.TODO(), key, value, 0).Result()
	return result

}

func RedisGet(key string) string {

	//cmd := radix.Cmd(&setVar, "GET", key)

	result, _:= rdb.Get(context.TODO(),key).Result()
	return result

}

func RedisGet_v1(key string) (strObj string, err error) {
	
	//cmd := radix.Cmd(&setVar, "GET", key)
	strObj, err = rdb.Get(context.TODO(), key).Result()
	return 

}

func RedisSearchKeys(pattern string) []string {


	defer func() {

		if r := recover(); r != nil {
			fmt.Println("Recovered in ScanAndGetKeys", r)
		}
	}()



	matchingKeys := make([]string, 0)

	log.Println("Start ScanAndGetKeys:: ", pattern)

	var ctx = context.TODO()
	iter := rdb.Scan(ctx, 0, pattern, 1000).Iterator()
	for iter.Next(ctx) {
		
		AppendIfMissing(matchingKeys, iter.Val())
	}
	if err := iter.Err(); err != nil {

		fmt.Println("ScanAndGetKeys","SCAN",err)
	}


	return matchingKeys
}

func RedisSetNx(key, value string, timeout int) bool {

	//cmd := radix.Cmd(&setVar, "SET", key, value, "nx", "ex", strconv.Itoa(timeout))
	 result, _ := rdb.Do(context.TODO(),"SET", key, value, "nx", "ex", timeout).Text()
	
    if result == "OK" {
		log.Println("GetRLock: ", true)
		return true
	} else {
		log.Println("GetRLock: ", false)
		return false
	}
}


/*

IncrByXX := redis.NewScript(`
		if redis.call("GET", KEYS[1]) ~= false then
			return redis.call("INCRBY", KEYS[1], ARGV[1])
		end
		return false
	`)

n, err := IncrByXX.Run(ctx, rdb, []string{"xx_counter"}, 2).Result()
fmt.Println(n, err)

err = rdb.Set(ctx, "xx_counter", "40", 0).Err()
if err != nil {
	panic(err)
}

n, err = IncrByXX.Run(ctx, rdb, []string{"xx_counter"}, 2).Result()

*/


func RedisRemoveRLock(key, value string) bool {


	luaScript := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
	//cmd := radix.Cmd(&setVar, "EVAL", luaScript, strconv.Itoa(1), key, value)
	result, _ := rdb.Do(context.TODO(),"EVAL", luaScript, strconv.Itoa(1), key, value).Int()

	
    if result == 1 {
		log.Println("GetRLock: ", true)
		return true
	} else {
		log.Println("GetRLock: ", false)
		return false
	}

}



func RedisRemove(key string) bool {
	
	//cmd := radix.Cmd(&setVar, "DEL", key)

	result, _ := rdb.Del(context.TODO(), key).Result()

	
    if result == 1 {
		log.Println("Remove Key: ", true)
		return true
	} else {
		log.Println("Remove Key: ", false)
		return false
	}

}

func RedisCheckKeyExist(key string) bool {

	//cmd := radix.Cmd(&setVar, "EXISTS", key)

	result, _ := rdb.Exists(context.TODO(), key).Result()

	
    if result == 1 {
		log.Println("Remove Key: ", true)
		return true
	} else {
		log.Println("Remove Key: ", false)
		return false
	}

}

// Redis Hashes Methods

func RedisHashGetAll(hkey string) map[string]string {
	//var setVar map[string]string;
	//cmd := radix.Cmd(&setVar, "HGETALL", hkey)
	result, _ := rdb.HGetAll(context.TODO(), hkey).Result()
	return result
}

func RedisHashGetValue(hkey, queueId string) string {
	//var setVar string;
	//cmd := radix.Cmd(&setVar, "HGET", hkey, queueId)
	result, _ := rdb.HGet(context.TODO(), hkey, queueId).Result()
	return result
}

func RedisHashSetField(hkey, field, value string) bool {
	
	//var setVar int;
	//cmd := radix.Cmd(&setVar, "HSET", hkey, field, value)
	result, _ := rdb.HSet(context.TODO(),hkey, field, value).Result()

	if result == 1 {
		return true
	} else {
		return false
	}
}

func RedisRemoveHashField(hkey, field string) bool {
	
	
	//cmd := radix.Cmd(&setVar, "HDEL", hkey, field)
	result, _ := rdb.HDel(context.TODO(),hkey, field).Result()
	
    if result == 1 {
		log.Println("Remove Key: ", true)
		return true
	} else {
		log.Println("Remove Key: ", false)
		return false
	}

}

// Redis List Methods

func RedisListLpop(lname string) string {
	
	
	//cmd := radix.Cmd(&setVar, "LPOP", lname)
	result, _ := rdb.LPop(context.TODO(),lname).Result()
	return result
}


/*-----------------------------Geo methods--------------------------------------*/

func RedisGeoRadius(tenant, company int, locationObj ReqLocationData) [][]string {
	var setVar [][]string;
	locationInfoKey := fmt.Sprintf("location:%d:%d", tenant, company)
	log.Println("locationInfoKey: ", locationInfoKey)
	//cmd := radix.Cmd(&setVar, "georadius", "positions", fmt.Sprintf("%.6f",locationObj.Longitude) , fmt.Sprintf("%.6f",locationObj.Latitude), strconv.Itoa(locationObj.Radius)  , locationObj.Metric, "WITHDIST", "ASC")
	
	return setVar

}

func RoutingEngineDistribution(pubChannelName string) string {


	//var activeRoutingKey string;
	//cmd := radix.Cmd(&activeRoutingKey, "GET", "ActiveRoutingEngine")
	//Cmd(cmd);

	activeRoutingKey, _ := rdb.Get(context.TODO(), "ActiveRoutingEngine").Result()

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
			
	        //cmd := radix.Cmd(&expire, "EXPIRE", "ActiveRoutingEngine", "60")
	        expire, _ := rdb.Expire(context.TODO(), "ActiveRoutingEngine", 60).Result()

			if expire {
				//log.Println("Extend Active Routing Engine Expire Time Success")
			} else {
				log.Println("Extend Active Routing Engine Expire Time Failed")
			}
		}
		return activeRoutingKey
	}
}

func AppendIfMissing(windowList []string, i string) []string {
	for _, ele := range windowList {
		if ele == i {
			return windowList
		}
	}
	return append(windowList, i)
}
