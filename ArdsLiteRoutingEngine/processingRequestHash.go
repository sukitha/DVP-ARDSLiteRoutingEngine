package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

func GetAllProcessingHashes() []string {
	processingHashSearchKey := fmt.Sprintf("ProcessingHash:%s:%s", "*", "*")
	processingHashes := RedisSearchKeys(processingHashSearchKey)
	return processingHashes
}

func GetAllProcessingItems(_processingHashKey string) []Request {
	log.Println(_processingHashKey)
	keyItems := strings.Split(_processingHashKey, ":")

	company := keyItems[1]
	tenant := keyItems[2]
	strHash := RedisHashGetAll(_processingHashKey)

	processingReqObjs := make([]Request, 0)

	for k, v := range strHash {
		log.Println("k:", k, "v:", v)
		requestKey := fmt.Sprintf("Request:%s:%s:%s", company, tenant, v)
		strReqObj := RedisGet(requestKey)
		log.Println(strReqObj)

		if strReqObj == "" {
			log.Println("Start SetNextProcessingItem")
			tenantInt, _ := strconv.Atoi(tenant)
			companyInt, _ := strconv.Atoi(company)
			//SetNextProcessingItem(tenantInt, companyInt, _processingHashKey, k, v, "")
			SetNextProcessingItem_Ards(tenantInt, companyInt, _processingHashKey, k, v)
		} else {
			var reqObj Request
			json.Unmarshal([]byte(strReqObj), &reqObj)

			if reqObj.SessionId == "" {

				log.Println("Critical issue request object found empty ---> set next item " + k + "value " + v)

				tenantInt, _ := strconv.Atoi(tenant)
				companyInt, _ := strconv.Atoi(company)
				//SetNextProcessingItem(tenantInt, companyInt, _processingHashKey, k, v, "")
				SetNextProcessingItem_Ards(tenantInt, companyInt, _processingHashKey, k, v)

			} else {

				processingReqObjs = AppendIfMissingReq(processingReqObjs, reqObj)
			}
		}
	}
	return processingReqObjs
}

func GetRejectedQueueId(_queueId string) string {
	//splitQueueId := strings.Split(_queueId, ":")
	//splitQueueId[len(splitQueueId)-1] = "REJECTED"
	//return strings.Join(splitQueueId, ":")

	rejectQueueId := fmt.Sprintf("%s:REJECTED", _queueId)
	return rejectQueueId
}

func SetNextProcessingItem(tenant, company int, _processingHash, _queueId, currentSession, requestState string) {
	//u1 := uuid.NewV4().String()
	//setNextLock := fmt.Sprintf("lock.setNextLock.%s", _queueId)
	//if RedisSetNx(setNextLock, u1, 1) == true {
	eSession := RedisHashGetValue(_processingHash, _queueId)

	log.Println("Item in " + _processingHash + "set next processing item in queue " + _queueId + " with session " + currentSession + " has now in hash " + eSession)
	if eSession != "" && eSession == currentSession {
		rejectedQueueId := GetRejectedQueueId(_queueId)
		nextRejectedQueueItem := RedisListLpop(rejectedQueueId)

		if nextRejectedQueueItem == "" {
			nextQueueItem := RedisListLpop(_queueId)
			if nextQueueItem == "" {
				removeHResult := RedisRemoveHashField(_processingHash, _queueId)
				if removeHResult {
					log.Println("Remove HashField Success.." + _processingHash + "::" + _queueId)
				} else {
					log.Println("Remove HashField Failed.." + _processingHash + "::" + _queueId)
				}
			} else {
				setHResult := RedisHashSetField(_processingHash, _queueId, nextQueueItem)
				if setHResult {
					log.Println("Set HashField Success.." + _processingHash + "::" + _queueId + "::" + nextQueueItem)
				} else {
					log.Println("Set HashField Failed.." + _processingHash + "::" + _queueId + "::" + nextQueueItem)
				}
			}
		} else {
			setHResult := RedisHashSetField(_processingHash, _queueId, nextRejectedQueueItem)
			if setHResult {
				log.Println("Set HashField Success.." + _processingHash + "::" + _queueId + "::" + nextRejectedQueueItem)
			} else {
				log.Println("Set HashField Failed.." + _processingHash + "::" + _queueId + "::" + nextRejectedQueueItem)
			}
		}
	} else {

		log.Println("session Mismatched, " + requestState + " ignore setNextItem")
		//SetRequestState(company, tenant, currentSession, "QUEUED")
		/*there is a new session added to the hash,
		now the item should route on next processing
		process next item will run through status and remove if the status is not queued
		there is a possibility to lost the item if status changes has failed.
		recheck all queue status set methods for concurrency and async operations.
		*/

	}
	//} else {
	//log.Println("Set Next Processing Item Fail To Aquire Lock")
	//}

	defer func() {
		//ReleasetLock(setNextLock, u1)l
	}()
}

func SetNextProcessingItem_Ards(tenant, company int, _processingHash, _queueId, currentSession string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in SetNextProcessingItem_Ards", r)
		}
	}()

	if _processingHash != "" && _queueId != "" && currentSession != "" {

		var setNextData SetNextData
		setNextData.ProcessingHashId = _processingHash
		setNextData.QueueId = _queueId
		setNextData.CurrentSession = currentSession

		req, _ := json.Marshal(setNextData)

		authToken := fmt.Sprintf("Bearer %s", accessToken)
		internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)

		ardsUrl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/ARDS/queue/setNextProcessingItem", CreateHost(ardsServiceHost, ardsServicePort))
		if Put(ardsUrl, string(req[:]), authToken, internalAuthToken) {
			log.Println("SetNextProcessingItem Ards Process Success")
		} else {
			log.Println("SetNextProcessingItem Ards Process Failed")
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
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ContinueArdsProcess", r)
		}
	}()
	if _request.ReqHandlingAlgo == "QUEUE" && _request.HandlingResource != "No matching resources at the moment" {
		req, _ := json.Marshal(_request)
		authToken := fmt.Sprintf("Bearer %s", accessToken)
		internalAuthToken := fmt.Sprintf("%d:%d", _request.Tenant, _request.Company)
		ardsUrl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/ARDS/continueprocess", CreateHost(_request.LbIp, _request.LbPort))
		if Post(ardsUrl, string(req[:]), authToken, internalAuthToken) {
			log.Println("Continue Ards Process Success")
			return true
		} else {
			log.Println("Continue Ards Process Failed")
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

func SetRequestState(_company, _tenant int, _sessionId, _newState string) string {
	reqStateKey := fmt.Sprintf("RequestState:%d:%d:%s", _company, _tenant, _sessionId)
	reqState := RedisSet(reqStateKey, _newState)
	return reqState
}

func ContinueProcessing(_request Request, _selectedResources SelectedResource) (continueProcessingResult bool, handlingResource []string) {
	log.Println("ReqOtherInfo:", _request.OtherInfo)
	var result string
	result, handlingResource = HandlingResources(_request.Company, _request.Tenant, _request.ResourceCount, _request.LbIp, _request.LbPort, _request.SessionId, _request.ServerType, _request.RequestType, _request.HandlingAlgo, _request.OtherInfo, _selectedResources)
	_request.HandlingResource = result
	continueProcessingResult = ContinueArdsProcess(_request)
	return
}

func AcquireProcessingHashLock(hashId, uuid string) bool {
	lockKey := fmt.Sprintf("ProcessingHashLock:%s", hashId)
	if RedisSetNx(lockKey, uuid, 60) == true {
		log.Println("lockKey: ", lockKey)
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

func ReleasetLock(hashId, uuid string) {
	lockKey := fmt.Sprintf("ProcessingHashLock:%s", hashId)

	if RedisRemoveRLock(lockKey, uuid) == true {
		log.Println("Release lock ", lockKey, "success.")
	} else {
		log.Println("Release lock ", lockKey, "failed.")
	}
	return
}

func ExecuteRequestHash(_processingHashKey, uuid string) {
	defer func() {
		//if r := recover(); r != nil {
		ReleasetLock(_processingHashKey, uuid)
		//}
	}()
	//for {
	if RedisCheckKeyExist(_processingHashKey) {
		processingItems := GetAllProcessingItems(_processingHashKey)
		if len(processingItems) > 0 {

			defaultRequest := processingItems[0]

			//switch (processingItems[0].ReqSelectionAlgo) {
			//case "LONGESTWAITING":
			//	sort.Sort(timeSliceReq(processingItems))
			//	break
			//case "PRIORITY":
			//	sort.Sort(ByReqPriority(processingItems))
			//	break
			//default:
			//	sort.Sort(timeSliceReq(processingItems))
			//	break
			//}

			//sort.Sort(timeSliceReq(processingItems))
			sort.Sort(ByReqPriority(processingItems))

			selectedResourcesForHash := SelectResources(defaultRequest.Company, defaultRequest.Tenant, processingItems, defaultRequest.SelectionAlgo)
			pickedResources := make([]string, 0)

			for _, longestWItem := range processingItems {

				//if longestWItem != (Request{}) {
				if longestWItem.SessionId != "" {

					log.Println("Execute processing hash item::", longestWItem.SessionId)
					requestState := GetRequestState(longestWItem.Company, longestWItem.Tenant, longestWItem.SessionId)
					if requestState == "QUEUED" {

						log.Println("pickedResources: ", pickedResources)

						resourceForRequest, isExist := GetSelectedResourceForRequest(selectedResourcesForHash, longestWItem.SessionId, pickedResources)

						log.Println("resourceForRequest: ", resourceForRequest)

						if isExist {
							continueProcessingResult, handlingResource := ContinueProcessing(longestWItem, resourceForRequest)
							if continueProcessingResult {
								log.Println("handlingResource: ", handlingResource)
								pickedResources = append(pickedResources, handlingResource...)
								log.Println("Continue ARDS Process Success")
							}
						} else {
							log.Println("Request not found in Selected Resource Data")
						}
					} else {
						log.Println("State of the queue item" + longestWItem.SessionId + "is not queued ->" + requestState)
						//SetNextProcessingItem(longestWItem.Tenant, longestWItem.Company, _processingHashKey, longestWItem.QueueId, longestWItem.SessionId, requestState)
						SetNextProcessingItem_Ards(longestWItem.Tenant, longestWItem.Company, _processingHashKey, longestWItem.QueueId, longestWItem.SessionId)
					}
				} else {
					log.Println("No Session Found")
				}
			}
			//ReleasetLock(_processingHashKey, uuid)
			//	return
		} else {
			log.Println("No Processing Items Found")
			//ReleasetLock(_processingHashKey, uuid)
			//	return
		}
	} else {
		log.Println("No Processing Hash Found")
		//ReleasetLock(_processingHashKey, uuid)
		//	return
	}
	//time.Sleep(200 * time.Millisecond)
	//}
}

func ExecuteRequestHashWithMsgQueue(_processingHashKey, uuid string) {
	defer func() {

		ReleasetLock(_processingHashKey, uuid)

	}()
	for RedisCheckKeyExist(_processingHashKey) {

		processingItems := GetAllProcessingItems(_processingHashKey)

		if len(processingItems) > 0 {

			defaultRequest := processingItems[0]

			sort.Sort(ByReqPriority(processingItems))

			selectedResourcesForHash := SelectResources(defaultRequest.Company, defaultRequest.Tenant, processingItems, defaultRequest.SelectionAlgo)
			pickedResources := make([]string, 0)

			for _, longestWItem := range processingItems {

				log.Println("Execute processing hash item::", longestWItem.Priority)

				if longestWItem.SessionId != "" {
					requestState := GetRequestState(longestWItem.Company, longestWItem.Tenant, longestWItem.SessionId)
					if requestState == "QUEUED" {

						resourceForRequest, isExist := GetSelectedResourceForRequest(selectedResourcesForHash, longestWItem.SessionId, pickedResources)
						if isExist {
							continueProcessingResult, handlingResource := ContinueProcessing(longestWItem, resourceForRequest)
							if continueProcessingResult {
								pickedResources = append(pickedResources, handlingResource...)
								log.Println("Continue ARDS Process Success")
							}
						} else {
							log.Println("Request not found in Selected Resource Data")
						}
					} else {

						//SetNextProcessingItem(longestWItem.Tenant, longestWItem.Company, _processingHashKey, longestWItem.QueueId, longestWItem.SessionId, requestState)
						SetNextProcessingItem_Ards(longestWItem.Tenant, longestWItem.Company, _processingHashKey, longestWItem.QueueId, longestWItem.SessionId)
					}
				} else {

					log.Println("No Session Found")
				}
			}
		} else {

			log.Println("No Processing Items Found")
		}
	}
}
