package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func GetAllProcessingHashes() []string {
	processingHashSearchKey := fmt.Sprintf("ProcessingHash:%s:%s", "*", "*")
	processingHashes := RedisSearchKeys(processingHashSearchKey)
	return processingHashes
}

func GetAllProcessingItems(_processingHashKey string) []Request {
	fmt.Println(_processingHashKey)
	keyItems := strings.Split(_processingHashKey, ":")

	company := keyItems[1]
	tenant := keyItems[2]
	strHash := RedisHashGetAll(_processingHashKey)

	processingReqObjs := make([]Request, 0)

	for k, v := range strHash {
		fmt.Println("k:", k, "v:", v)
		requestKey := fmt.Sprintf("Request:%s:%s:%s", company, tenant, v)
		strReqObj := RedisGet(requestKey)
		fmt.Println(strReqObj)

		if strReqObj == "" {
			fmt.Println("Start SetNextProcessingItem")
			SetNextProcessingItem(_processingHashKey, k)
		} else {
			var reqObj Request
			json.Unmarshal([]byte(strReqObj), &reqObj)

			processingReqObjs = AppendIfMissingReq(processingReqObjs, reqObj)
		}
	}
	return processingReqObjs
}

func SetNextProcessingItem(_processingHash, _queueId string) {
	nextQueueItem := RedisListLpop(_queueId)
	if nextQueueItem == "" {
		removeHResult := RedisRemoveHashField(_processingHash, _queueId)
		if removeHResult {
			fmt.Println("Remove HashField Success.." + _processingHash + "::" + _queueId)
		} else {
			fmt.Println("Remove HashField Failed.." + _processingHash + "::" + _queueId)
		}
	} else {
		setHResult := RedisHashSetField(_processingHash, _queueId, nextQueueItem)
		if setHResult {
			fmt.Println("Set HashField Success.." + _processingHash + "::" + _queueId + "::" + nextQueueItem)
		} else {
			fmt.Println("Set HashField Failed.." + _processingHash + "::" + _queueId + "::" + nextQueueItem)
		}
	}
}

/*func GetLongestWaitingItem(_request []Request) Request {
	longetWaitingItem := Request{}
	reqCount := len(_request)
	longetWaitingItemArriveTime := time.Now()

	if reqCount > 0 {
		for _, req := range _request {
			arrTime, _ := time.Parse(layout, req.ArriveTime)
			if arrTime.Before(longetWaitingItemArriveTime) {
				longetWaitingItemArriveTime = arrTime
				longetWaitingItem = req
			}
		}
	}

	return longetWaitingItem
}*/

func ContinueArdsProcess(_request Request) bool {
	if _request.ReqHandlingAlgo == "QUEUE" && _request.HandlingResource != "No matching resources at the moment" {
		req, _ := json.Marshal(_request)
		authToken := fmt.Sprintf("Bearer %s", accessToken)
		internalAuthToken := fmt.Sprintf("%d:%d", _request.Tenant, _request.Company)
		ardsUrl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/ARDS/continueprocess", CreateHost(_request.LbIp, _request.LbPort))
		if Post(ardsUrl, string(req[:]), authToken, internalAuthToken) {
			fmt.Println("Continue Ards Process Success")
			return true
		} else {
			fmt.Println("Continue Ards Process Failed")
			return false
		}
	} else {
		return false
	}
}

func GetRequestState(_company, _tenant int, _sessionId string) string {
	reqStateKey := fmt.Sprintf("RequestState:%d:%d:%s", _company, _tenant, _sessionId)
	reqState := RedisGet(reqStateKey)
	return reqState
}

func ContinueProcessing(_request Request) bool {
	fmt.Println("ReqOtherInfo:", _request.OtherInfo)
	var result = SelectResources(_request.Company, _request.Tenant, _request.ResourceCount, _request.LbIp, _request.LbPort, _request.SessionId, _request.ServerType, _request.RequestType, _request.SelectionAlgo, _request.HandlingAlgo, _request.OtherInfo)
	_request.HandlingResource = result
	return ContinueArdsProcess(_request)
}

func AcquireProcessingHashLock(hashId string) bool {
	lockKey := fmt.Sprintf("ProcessingHashLock:%s", hashId)
	if RedisSetNx(lockKey, "LOCKED") == true {
		fmt.Println("lockKey: ", lockKey)
		//if RedisSetEx(lockKey, "LOCKED", 60) {
		return true
		//} else {
		//	RedisRemove(lockKey)
		//	return false
		//}
	} else {
		return false
	}
}

func ReleasetLock(hashId string) bool {
	lockKey := fmt.Sprintf("ProcessingHashLock:%s", hashId)
	return RedisRemove(lockKey)
}

func ExecuteRequestHash(_processingHashKey string) {
	for {
		if RedisCheckKeyExist(_processingHashKey) {
			processingItems := GetAllProcessingItems(_processingHashKey)
			sort.Sort(timeSliceReq(processingItems))
			for _, longestWItem := range processingItems {
				//if longestWItem != (Request{}) {
				if longestWItem.SessionId != "" {
					if GetRequestState(longestWItem.Company, longestWItem.Tenant, longestWItem.SessionId) == "QUEUED" {
						if ContinueProcessing(longestWItem) {
							//SetNextProcessingItem(_processingHashKey, longestWItem.QueueId)
							fmt.Println("Continue ARDS Process Success")
						}
					}
				}
			}
		} else {
			if ReleasetLock(_processingHashKey) == true {
				fmt.Println("Release lock ", _processingHashKey, "success.")
			} else {
				fmt.Println("Release lock ", _processingHashKey, "failed.")
			}
			return
		}
	}
}
