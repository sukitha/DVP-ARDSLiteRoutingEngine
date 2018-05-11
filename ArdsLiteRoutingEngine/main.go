// ArdsRoutingEngine project main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/DuoSoftware/gorest"
	"github.com/satori/go.uuid"
)

const layout = "2006-01-02T15:04:05Z07:00"

func main() {
	log.Println("Starting Ards Route Engine")
	InitiateRedis()
	go InitiateService()

	if useMsgQueue {
		//-------------------Amqp Based Routing---------------------------------------
		rmqIps := strings.Split(rabbitMQIp, ",")
		currentRmqNodeIndex := 0
		rmqNodeTryCount := 0

		for {
			if len(rmqIps) > 1 {
				if rmqNodeTryCount > 30 {
					fmt.Println("Start to change RMQ node")
					if currentRmqNodeIndex == (len(rmqIps) - 1) {
						currentRmqNodeIndex = 0
					} else {
						currentRmqNodeIndex++
					}
					rmqNodeTryCount = 0
				}
			}
			rmqNodeTryCount++
			fmt.Println("Start Connecting to RMQ: ", rmqIps[currentRmqNodeIndex], " :: TryCount: ", rmqNodeTryCount)
			Worker(rmqIps[currentRmqNodeIndex])
			fmt.Println("End Worker()")
			time.Sleep(2 * time.Second)
		}
		// for {
		// 	Worker()
		// 	log.Println("End Worker()")
		// 	time.Sleep(2 * time.Second)
		// }
	} else {
		//-------------------RedisDb Based Routing---------------------------------------
		for {

			pubChannelName := fmt.Sprintf("RoutingChannel:%s", routingEngineId)
			if RoutingEngineDistribution(pubChannelName) == pubChannelName {
				availablePHashes := GetAllProcessingHashes()
				if len(availablePHashes) > 0 {
					for _, h := range availablePHashes {
						u1, _ := uuid.NewV4()
						if AcquireProcessingHashLock(h, u1.String()) == true {
							//go ExecuteRequestHash(h, u1.String())
							ExecuteRequestHash(h, u1.String())
						} else {
							time.Sleep(1 * time.Second)
						}
					}
				} else {
					log.Println("No Processing Hash Found...")
					time.Sleep(1 * time.Second)
				}
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}

}

//InitiateService start listening to the self host service port
func InitiateService() {
	listeningPort := fmt.Sprintf(":%s", port)
	gorest.RegisterService(new(ArdsLiteRS))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(listeningPort, nil)
}
