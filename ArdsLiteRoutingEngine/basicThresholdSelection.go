package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
)

func BasicThresholdSelection(_company, _tenent int, _requests []Request) (result []SelectionResult) {
	//requestKey := fmt.Sprintf("Request:%d:%d:%s", _company, _tenent, _sessionId)
	//log.Println(requestKey)
	//
	//strReqObj := RedisGet(requestKey)
	//log.Println(strReqObj)
	//
	//var reqObj RequestSelection
	//json.Unmarshal([]byte(strReqObj), &reqObj)

	var selectedResources = make([]SelectionResult, len(_requests))

	for i, reqObj := range _requests {

		selectedResources[i].Request = reqObj.SessionId

		var resourceConcInfo = make([]ConcurrencyInfo, 0)
		var resourceThresholdConcInfo = make([]ConcurrencyInfo, 0)
		var matchingResources = make([]string, 0)
		var matchingThresholdResources = make([]string, 0)
		if len(reqObj.AttributeInfo) > 0 {
			var tagArray = make([]string, 3)

			tagArray[0] = fmt.Sprintf("company_%d", reqObj.Company)
			tagArray[1] = fmt.Sprintf("tenant_%d", reqObj.Tenant)
			//tagArray[2] = fmt.Sprintf("class_%s", reqObj.Class)
			//tagArray[3] = fmt.Sprintf("type_%s", reqObj.Type)
			//tagArray[4] = fmt.Sprintf("category_%s", reqObj.Category)
			tagArray[2] = fmt.Sprintf("objtype_%s", "Resource")

			attInfo := make([]string, 0)

			for _, value := range reqObj.AttributeInfo {
				for _, att := range value.AttributeCode {
					attInfo = AppendIfMissingString(attInfo, att)
				}
			}

			sort.Sort(ByStringValue(attInfo))
			for _, att := range attInfo {
				//log.Println("attCode", att)
				tagArray = AppendIfMissingString(tagArray, fmt.Sprintf("attribute_%s", att))
			}

			tags := fmt.Sprintf("tag:*%s*", strings.Join(tagArray, "*"))
			log.Println(tags)
			val := RedisSearchKeys(tags)
			lenth := len(val)
			log.Println(lenth)

			for _, match := range val {
				strResKey := RedisGet(match)
				strResObj := RedisGet(strResKey)
				//log.Println(strResObj)

				var resObj Resource
				json.Unmarshal([]byte(strResObj), &resObj)

				_attAvailable, _isThreshold := IsAttributeAvailable(reqObj.AttributeInfo, resObj.ResourceAttributeInfo, reqObj.RequestType)

				if resObj.ResourceId != "" && _attAvailable {
					concInfo, err := GetConcurrencyInfo(resObj.Company, resObj.Tenant, resObj.ResourceId, reqObj.RequestType)
					if err != nil {
						log.Println("Error in GetConcurrencyInfo")
					} else {
						if _isThreshold {
							resourceThresholdConcInfo = append(resourceThresholdConcInfo, concInfo)
						} else {
							resourceConcInfo = append(resourceConcInfo, concInfo)
						}
					}
					//matchingResources = AppendIfMissing(matchingResources, strResKey)
					//log.Println(strResKey)
				}
			}

			sort.Sort(timeSlice(resourceConcInfo))
			sort.Sort(timeSlice(resourceThresholdConcInfo))

			for _, res := range resourceConcInfo {
				resKey := fmt.Sprintf("Resource:%d:%d:%s", res.Company, res.Tenant, res.ResourceId)
				matchingResources = AppendIfMissingString(matchingResources, resKey)
				//log.Println(resKey)
			}

			for _, res := range resourceThresholdConcInfo {
				resKey := fmt.Sprintf("Resource:%d:%d:%s", res.Company, res.Tenant, res.ResourceId)
				matchingThresholdResources = AppendIfMissingString(matchingThresholdResources, resKey)
				//log.Println(resKey)
			}

		}
		selectedResources[i].Resources.Priority = matchingResources
		selectedResources[i].Resources.Threshold = matchingThresholdResources
	}
	return selectedResources

}
