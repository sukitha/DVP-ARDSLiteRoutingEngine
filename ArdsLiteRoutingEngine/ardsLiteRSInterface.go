package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/DuoSoftware/gorest"
)

//ArdsLiteRS structure of self hosted routing service endpoints
type ArdsLiteRS struct {
	gorest.RestService `root:"/resourceselection/" consumes:"application/json" produces:"application/json"`
	getResource        gorest.EndPoint `method:"GET" path:"/getresource/{Company:int}/{Tenant:int}/{ResourceCount:int}/{SessionID:string}/{ServerType:string}/{RequestType:string}/{SelectionAlgo:string}/{HandlingAlgo:string}/{OtherInfo:string}" output:"string"`
}

//GetResource return selected resource for the service request
func (ardsLiteRs ArdsLiteRS) GetResource(Company, Tenant, ResourceCount int, SessionID, ServerType, RequestType, SelectionAlgo, HandlingAlgo, OtherInfo string) string {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"

	log.Println("Company:", Company)
	log.Println("Tenant:", Tenant)
	log.Println("SessionId:", SessionID)
	log.Println("OtherInfo:", OtherInfo)

	byt := []byte(OtherInfo)
	var otherInfo string
	json.Unmarshal(byt, &otherInfo)

	requestKey := fmt.Sprintf("Request:%d:%d:%s", Company, Tenant, SessionID)

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
		result, _ := HandlingResources(Company, Tenant, ResourceCount, reqObj.LbIp, reqObj.LbPort, SessionID, ServerType, RequestType, HandlingAlgo, otherInfo, resourceForRequest)
		return result
	}
	return "Session Invalied"

}
