package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func SingleHandling(ReqClass, ReqType, ReqCategory, sessionId string, resourceIds []string) string {
	return SelectHandlingResource(ReqClass, ReqType, ReqCategory, sessionId, resourceIds)
}

func SelectHandlingResource(ReqClass, ReqType, ReqCategory, sessionId string, resourceIds []string) string {
	for _, key := range resourceIds {
		fmt.Println(key)
		strResObj := RedisGet(key)
		fmt.Println(strResObj)

		var resObj Resource
		json.Unmarshal([]byte(strResObj), &resObj)

		conInfo := GetConcurrencyInfo(resObj.Company, resObj.Tenant, resObj.ResourceId, ReqClass, ReqType, ReqCategory)
		metaData := GetReqMetaData(resObj.Company, resObj.Tenant, ReqClass, ReqType, ReqCategory)
		resState := GetResourceState(resObj.Company, resObj.Tenant, resObj.ResourceId)

		if resState == "Available" && conInfo.RejectCount < metaData.MaxRejectCount {
			ClearSlotOnMaxRecerved(ReqClass, ReqType, ReqCategory, sessionId, resObj, metaData)

			var tagArray = make([]string, 8)

			tagArray[0] = fmt.Sprintf("company_%d", resObj.Company)
			tagArray[1] = fmt.Sprintf("tenant_%d", resObj.Tenant)
			tagArray[2] = fmt.Sprintf("class_%s", ReqClass)
			tagArray[3] = fmt.Sprintf("type_%s", ReqType)
			tagArray[4] = fmt.Sprintf("category_%s", ReqCategory)
			tagArray[5] = fmt.Sprintf("state_%s", "Available")
			tagArray[6] = fmt.Sprintf("resourceid_%s", resObj.ResourceId)
			tagArray[7] = fmt.Sprintf("objtype_%s", "CSlotInfo")

			tags := fmt.Sprintf("tag:*%s*", strings.Join(tagArray, "*"))
			fmt.Println(tags)
			availableSlots := RedisSearchKeys(tags)

			for _, tagKey := range availableSlots {
				strslotKey := RedisGet(tagKey)
				fmt.Println(strslotKey)

				strslotObj := RedisGet(strslotKey)
				fmt.Println(strslotObj)

				var slotObj CSlotInfo
				json.Unmarshal([]byte(strslotObj), &slotObj)

				slotObj.State = "Reserved"
				slotObj.SessionId = sessionId
				slotObj.OtherInfo = "Inbound"

				if ReserveSlot(slotObj) == true {
					fmt.Println("Return resource Data:", resObj.OtherInfo)
					return resObj.OtherInfo
				}
			}
		}

	}
	return "No matching resources at the moment"
}
