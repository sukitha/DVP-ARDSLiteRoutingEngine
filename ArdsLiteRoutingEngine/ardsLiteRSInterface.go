package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/DuoSoftware/gorest"
	uuid "github.com/satori/go.uuid"
)

type ArdsLiteRS struct {
	gorest.RestService `root:"/resourceselection/" consumes:"application/json" produces:"application/json"`
	getResource        gorest.EndPoint `method:"GET" path:"/getresource/{Company:int}/{Tenant:int}/{ResourceCount:int}/{SessionId:string}/{ServerType:string}/{RequestType:string}/{SelectionAlgo:string}/{HandlingAlgo:string}/{OtherInfo:string}" output:"string"`
	executeHash        gorest.EndPoint `method:"POST" path:"/executeHash/" postdata:"HashData"`
}

func (ardsLiteRs ArdsLiteRS) GetResource(Company, Tenant, ResourceCount int, SessionId, ServerType, RequestType, SelectionAlgo, HandlingAlgo, OtherInfo string) string {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"

	log.Println("Company:", Company)
	log.Println("Tenant:", Tenant)
	log.Println("SessionId:", SessionId)
	log.Println("OtherInfo:", OtherInfo)

	byt := []byte(OtherInfo)
	var otherInfo string
	json.Unmarshal(byt, &otherInfo)

	requestKey := fmt.Sprintf("Request:%d:%d:%s", Company, Tenant, SessionId)

	if RedisCheckKeyExist(requestKey) {
		strReqObj := RedisGet(requestKey)
		log.Println(strReqObj)

		var reqObj Request
		json.Unmarshal([]byte(strReqObj), &reqObj)

		tempRequestArray := make([]Request, 1)
		tempRequestArray[0] = reqObj

		pickedResources := make([]string, 0)

		selectedResources := SelectResources(Company, Tenant, tempRequestArray, SelectionAlgo)
		resourceForRequest, _ := GetSelectedResourceForRequest(selectedResources, reqObj.SessionId, pickedResources)
		result, _ := HandlingResources(Company, Tenant, reqObj.BusinessUnit, ResourceCount, reqObj.LbIp, reqObj.LbPort, SessionId, ServerType, RequestType, HandlingAlgo, otherInfo, resourceForRequest)
		return result
	}
	return "Session Invalied"

}

//ExecuteHash by amql adapter
func (ardsLiteRs ArdsLiteRS) ExecuteHash(hashData HashData) {
	log.Println("Start New amql adapter")
	log.Printf("Received a message: %s", hashData)

	hashKey := hashData.HashKey
	u1 := uuid.NewV4()
	if AcquireProcessingHashLock(hashKey, u1.String()) == true {
		go ExecuteRequestHashWithMsgQueue(hashKey, u1.String())
	}
	//dot_count := bytes.Count(d.Body, []byte("."))
	//t := time.Duration(dot_count)
	//time.Sleep(t * time.Second)
	log.Printf("Done")
	log.Println("End msgs")
	return
}

/*func SelectResources(Company, Tenant, ResourceCount int, ArdsLbIp, ArdsLbPort, SessionId, ServerType, RequestType, SelectionAlgo, HandlingAlgo, OtherInfo string) string {
	var selectionResult SelectionResult
	var handlingResult = ""
	switch SelectionAlgo {
	case "BASIC":
		selectionResult = BasicSelection(Company, Tenant, SessionId)
		break
	case "BASICTHRESHOLD":
		selectionResult = BasicThresholdSelection(Company, Tenant, SessionId)
		break
	case "WEIGHTBASE":
		selectionResult = WeightBaseSelection(Company, Tenant, SessionId)
		break
	case "LOCATIONBASE":
		selectionResult = LocationBaseSelection(Company, Tenant, SessionId)
		break
	default:
		break
	}

	switch HandlingAlgo {
	case "SINGLE":
		handlingResult = SingleResourceAlgo(ArdsLbIp, ArdsLbPort, ServerType, RequestType, SessionId, selectionResult, Company, Tenant)
	case "MULTIPLE":
		log.Println("ReqOtherInfo:", OtherInfo)
		resCount := ResourceCount
		log.Println("GetRequestedResCount:", resCount)
		handlingResult = MultipleHandling(ArdsLbIp, ArdsLbPort, ServerType, RequestType, SessionId, selectionResult, resCount, Company, Tenant)
	default:
		handlingResult = ""
	}

	return handlingResult
}*/
