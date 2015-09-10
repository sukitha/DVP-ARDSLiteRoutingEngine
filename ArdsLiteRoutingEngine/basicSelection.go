package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func BasicSelection(_company, _tenent int, _sessionId string) []string {
	requestKey := fmt.Sprintf("Request:%d:%d:%s", _company, _tenent, _sessionId)
	fmt.Println(requestKey)

	strResObj := RedisGet(requestKey)
	fmt.Println(strResObj)

	var reqObj RequestSelection
	json.Unmarshal([]byte(strResObj), &reqObj)

	var resourceConcInfo = make([]ConcurrencyInfo, 0)
	var matchingResources = make([]string, 0)
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
			fmt.Println("attCode", att)
			tagArray = AppendIfMissingString(tagArray, fmt.Sprintf("attribute_%s", att))
		}

		tags := fmt.Sprintf("tag:*%s*", strings.Join(tagArray, "*"))
		fmt.Println(tags)
		val := RedisSearchKeys(tags)
		lenth := len(val)
		fmt.Println(lenth)

		for _, match := range val {
			strResKey := RedisGet(match)
			splitVals := strings.Split(strResKey, ":")
			if len(splitVals) == 4 {
				concInfo := GetConcurrencyInfo(reqObj.Company, reqObj.Tenant, splitVals[3], reqObj.Category)
				resourceConcInfo = append(resourceConcInfo, concInfo)
				//matchingResources = AppendIfMissing(matchingResources, strResKey)
				//fmt.Println(strResKey)
			}
		}

		sort.Sort(timeSlice(resourceConcInfo))

		for _, res := range resourceConcInfo {
			resKey := fmt.Sprintf("Resource:%d:%d:%s", reqObj.Company, reqObj.Tenant, res.ResourceId)
			matchingResources = AppendIfMissingString(matchingResources, resKey)
			fmt.Println(resKey)
		}

	}

	return matchingResources

}
