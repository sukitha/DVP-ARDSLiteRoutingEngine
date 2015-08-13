package main

import (
	"code.google.com/p/gorest"
	"encoding/json"
	"fmt"
)

type ArdsLiteRS struct {
	gorest.RestService `root:"/resourceselection/" consumes:"application/json" produces:"application/json"`
	getResource        gorest.EndPoint `method:"GET" path:"/getresource/{Company:int}/{Tenant:int}/{SessionId:string}/{ReqClass:string}/{ReqType:string}/{ReqCategory:string}/{SelectionAlgo:string}/{HandlingAlgo:string}/{OtherInfo:string}" output:"string"`
}

func (ardsLiteRs ArdsLiteRS) GetResource(Company, Tenant int, SessionId, ReqClass, ReqType, ReqCategory, SelectionAlgo, HandlingAlgo, OtherInfo string) string {

	const longForm = "Jan 2, 2006 at 3:04pm (MST)"

	fmt.Println("Company:", Company)
	fmt.Println("Tenant:", Tenant)
	fmt.Println("SessionId:", SessionId)
	fmt.Println("OtherInfo:", OtherInfo)

	byt := []byte(OtherInfo)
	var otherInfo string
	json.Unmarshal(byt, &otherInfo)

	var result = SelectResources(Company, Tenant, SessionId, ReqClass, ReqType, ReqCategory, SelectionAlgo, HandlingAlgo, otherInfo)

	return result

}

func GetRequestedResCount(ReqOtherInfo string) int {
	fmt.Println("ReqOtherInfo:", ReqOtherInfo)
	var requestedResCount = 1

	var resCount MultiResCount
	err := json.Unmarshal([]byte(ReqOtherInfo), &resCount)

	if err != nil {
		println(err)
		requestedResCount = 1
	} else {
		requestedResCount = resCount.ResourceCount
	}

	return requestedResCount
}

func SelectResources(Company, Tenant int, SessionId, ReqClass, ReqType, ReqCategory, SelectionAlgo, HandlingAlgo, OtherInfo string) string {
	var selectionResult = make([]string, 0)
	var handlingResult = ""
	switch SelectionAlgo {
	case "BASIC":
		selectionResult = BasicSelectionAlgo(Company, Tenant, SessionId)
	case "WEIGHTBASE":
		selectionResult = WeightBaseSelectionAlgo(Company, Tenant, SessionId)
	default:
		selectionResult = make([]string, 0)
	}

	switch HandlingAlgo {
	case "SINGLE":
		handlingResult = SingleResourceAlgo(ReqClass, ReqType, ReqCategory, SessionId, selectionResult)
	case "MULTIPLE":
		fmt.Println("ReqOtherInfo:", OtherInfo)
		resCount := GetRequestedResCount(OtherInfo)
		fmt.Println("GetRequestedResCount:", resCount)
		handlingResult = MultipleHandling(ReqClass, ReqType, ReqCategory, SessionId, selectionResult, resCount)
	default:
		handlingResult = ""
	}

	return handlingResult
}
