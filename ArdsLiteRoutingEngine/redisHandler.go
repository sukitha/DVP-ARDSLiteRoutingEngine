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

	// "github.com/mediocregopher/radix.v2/pool"
	// "github.com/mediocregopher/radix.v2/redis"
	// "github.com/mediocregopher/radix.v2/sentinel"
	// "github.com/mediocregopher/radix.v2/util"

	"github.com/mediocregopher/radix/v3"
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

var redisSentinel *radix.Sentinel
var redisPool *radix.Pool;
var redisScanClient radix.Client;


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
	customConnFunc := func(network, addr string) (radix.Client, error) {
		return radix.Dial(network, addr,
		radix.DialAuthPass(redisPassword),
	 	radix.DialSelectDB(redisDb),)
		
	}

	customPoolConnFunc := func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
		radix.DialAuthPass(redisPassword),
	 	radix.DialSelectDB(redisDb),)
		
	}
	// ,
	// 		radix.DialAuthPass(redisPassword),
	// 		radix.DialSelectDB(redisDb),

	if redisMode == "sentinel" {

		sentinelIps := strings.Split(sentinelHosts, ",")
		var ips []string;

		if len(sentinelIps) > 1 {

			for _, ip := range sentinelIps{
				//redis://user:secret@localhost:6379/0
				ipPortArray := strings.Split(ip, ":")
				sentinelIp := ip;
				if(len(ipPortArray) > 1){
					sentinelIp = fmt.Sprintf("%s:%s", ipPortArray[0], ipPortArray[1])
				}else{
					sentinelIp = fmt.Sprintf("%s:%s", ip, sentinelPort)
				}
				ips = append(ips, sentinelIp)
				
			}

			redisSentinel, err = radix.NewSentinel(redisClusterName, ips , radix.SentinelPoolFunc(customConnFunc))

			if err != nil {
				log.Println("InitiateSentinel ::", err)
			}
		} else {
			log.Println("Not enough sentinel servers")
		}


	} else {

		redisPool, err = radix.NewPool("tcp", redisIp, 10, radix.PoolConnFunc(customPoolConnFunc))

		var errc error
		redisScanClient, errc = customConnFunc("tcp",redisIp);

		if err != nil {
			errHndlrNew("InitiateRedis", "InitiatePool", err)
			os.Exit(0)
		}


		if errc != nil {
			errHndlrNew("InitiateRedis", "InitiateScanner", err)
			os.Exit(0)
		}
	}
}

// Redis String Methods

func Cmd(cmd radix.CmdAction) error{

	var err error;
	if redisMode == "sentinel" {
				
		if err := redisSentinel.Do(cmd); err != nil {
			fmt.Println(err)
		}

	} else {
			
		if err := redisPool.Do(cmd); err != nil {
			errHndlrNew("OnReset", "getConnFromPool", err)
		}
	}

	return err;

}

func RedisSet(key, value string) string {

	
	var setVar string;
	cmd := radix.Cmd(&setVar, "SET", key, value)

	Cmd(cmd);
    
	return setVar

}

func RedisGet(key string) string {
	var setVar string;
	cmd := radix.Cmd(&setVar, "GET", key)

	Cmd(cmd);
    
	return setVar

}

func RedisGet_v1(key string) (strObj string, err error) {
	
	var setVar string;
	cmd := radix.Cmd(&setVar, "GET", key)

	err = Cmd(cmd);
	strObj = setVar;
    
	return 

}

func RedisSearchKeys(pattern string) []string {

	
	scanOpts := radix.ScanOpts{
		Command: "SCAN",
		Count: 1000,
		Pattern: pattern,
	}

	matchingKeys := make([]string, 0)

	log.Println("Start ScanAndGetKeys:: ", pattern)

	var client  radix.Client;

	if redisMode == "sentinel" {
		addr, _ := redisSentinel.Addrs()
		client , _ = redisSentinel.Client(addr) 
	}else{
		client = redisScanClient;
	}

	if client != nil{
		scanner := radix.NewScanner(client , scanOpts)
		var key string;
		counter := 0;


		for scanner.Next(&key) {
			counter++
			matchingKeys = AppendIfMissing(matchingKeys, key)
		}

		if err := scanner.Close(); err != nil{
			fmt.Println(err)
			os.Exit(0)
		}
    }

	return matchingKeys
}

func RedisSetNx(key, value string, timeout int) bool {

	var setVar string;
	cmd := radix.Cmd(&setVar, "SET", key, value, "nx", "ex", strconv.Itoa(timeout))

	 Cmd(cmd);
	
    if setVar == "OK" {
		log.Println("GetRLock: ", true)
		return true
	} else {
		log.Println("GetRLock: ", false)
		return false
	}
}



func RedisRemoveRLock(key, value string) bool {

    var setVar int;
	luaScript := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
	cmd := radix.Cmd(&setVar, "EVAL", luaScript, strconv.Itoa(1), key, value)

	 Cmd(cmd);
	
    if setVar == 1 {
		log.Println("GetRLock: ", true)
		return true
	} else {
		log.Println("GetRLock: ", false)
		return false
	}

}



func RedisRemove(key string) bool {
	
	
	var setVar string;
	cmd := radix.Cmd(&setVar, "DEL", key)

	 Cmd(cmd);
	
    if result, _ := strconv.Atoi(setVar); result == 1 {
		log.Println("Remove Key: ", true)
		return true
	} else {
		log.Println("Remove Key: ", false)
		return false
	}

}

func RedisCheckKeyExist(key string) bool {
	var setVar int;
	cmd := radix.Cmd(&setVar, "EXISTS", key)

	 Cmd(cmd);
	
    if setVar == 1 {
		log.Println("Remove Key: ", true)
		return true
	} else {
		log.Println("Remove Key: ", false)
		return false
	}

}

// Redis Hashes Methods

func RedisHashGetAll(hkey string) map[string]string {
	var setVar map[string]string;
	cmd := radix.Cmd(&setVar, "HGETALL", hkey)
	Cmd(cmd);
	return setVar
}

func RedisHashGetValue(hkey, queueId string) string {
	var setVar string;
	cmd := radix.Cmd(&setVar, "HGET", hkey, queueId)
	 Cmd(cmd);
	return setVar
}

func RedisHashSetField(hkey, field, value string) bool {
	
	var setVar int;
	cmd := radix.Cmd(&setVar, "HSET", hkey, field, value)
	 Cmd(cmd);

	if setVar == 1 {
		return true
	} else {
		return false
	}
}

func RedisRemoveHashField(hkey, field string) bool {
	
	var setVar int;
	cmd := radix.Cmd(&setVar, "HDEL", hkey, field)
	Cmd(cmd);
    if setVar == 1 {
		log.Println("Remove Key: ", true)
		return true
	} else {
		log.Println("Remove Key: ", false)
		return false
	}

}

// Redis List Methods

func RedisListLpop(lname string) string {
	var setVar string;
	cmd := radix.Cmd(&setVar, "LPOP", lname)
	Cmd(cmd);
	return setVar
}


/*-----------------------------Geo methods--------------------------------------*/

func RedisGeoRadius(tenant, company int, locationObj ReqLocationData) [][]string {
	var setVar [][]string;
	locationInfoKey := fmt.Sprintf("location:%d:%d", tenant, company)
	log.Println("locationInfoKey: ", locationInfoKey)
	cmd := radix.Cmd(&setVar, "georadius", "positions", fmt.Sprintf("%.6f",locationObj.Longitude) , fmt.Sprintf("%.6f",locationObj.Latitude), strconv.Itoa(locationObj.Radius)  , locationObj.Metric, "WITHDIST", "ASC")
	Cmd(cmd);
	return setVar

}

func RoutingEngineDistribution(pubChannelName string) string {


	var activeRoutingKey string;
	cmd := radix.Cmd(&activeRoutingKey, "GET", "ActiveRoutingEngine")
	Cmd(cmd);

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
			var expire int;
	        cmd := radix.Cmd(&expire, "EXPIRE", "ActiveRoutingEngine", "60")
	        Cmd(cmd);

			if expire == 1 {
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
