package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

//MultipleHandling return a array of handling resources
func MultipleHandling(ardsLbIP, ardsLbPort, ServerType, RequestType, sessionID string, selectedResources SelectedResource, nuOfResRequested, reqCompany, reqTenant int) (handlingResult string, handlingResource []string) {
	return SelectMultipleHandlingResource(ardsLbIP, ardsLbPort, ServerType, RequestType, sessionID, selectedResources, nuOfResRequested, reqCompany, reqTenant)
}

//SelectMultipleHandlingResource find and select the resources for handle the request accourding to their availablility
func SelectMultipleHandlingResource(ardsLbIP, ardsLbPort, ServerType, RequestType, sessionID string, selectedResources SelectedResource, nuOfResRequested, reqCompany, reqTenant int) (handlingResult string, handlingResource []string) {
	selectedResList := make([]string, 0)
	selectedResKeyList := make([]string, 0)
	resourceIds := append(selectedResources.Priority, selectedResources.Threshold...)
	log.Println("///////////////////////////////////////selectedResources/////////////////////////////////////////////////")
	log.Println("Priority:: ", selectedResources.Priority)
	log.Println("Threshold:: ", selectedResources.Threshold)
	log.Println("ResourceIds:: ", resourceIds)
	for _, key := range resourceIds {
		log.Println(key)
		strResObj := RedisGet(key)
		log.Println(strResObj)

		var resObj Resource
		json.Unmarshal([]byte(strResObj), &resObj)

		conInfo, cErr := GetConcurrencyInfo(resObj.Company, resObj.Tenant, resObj.ResourceId, RequestType)

		log.Println("conInfo.RejectCount:: ", conInfo.RejectCount)
		log.Println("conInfo.IsRejectCountExceeded:: ", conInfo.IsRejectCountExceeded)

		if cErr == nil {
			metaData, mErr := GetReqMetaData(reqCompany, reqTenant, ServerType, RequestType)
			log.Println("metaData.MaxRejectCount:: ", metaData.MaxRejectCount)

			if mErr == nil {
				resState, resMode, sErr := GetResourceState(resObj.Company, resObj.Tenant, resObj.ResourceId)
				if sErr == nil {

					if resState == "Available" && resMode == "Inbound" && conInfo.RejectCount < metaData.MaxRejectCount && conInfo.IsRejectCountExceeded == false {
						ClearSlotOnMaxRecerved(ardsLbIP, ardsLbPort, ServerType, RequestType, sessionID, resObj)

						var tagArray = make([]string, 8)

						tagArray[0] = fmt.Sprintf("company_%d:", resObj.Company)
						tagArray[1] = fmt.Sprintf("tenant_%d:", resObj.Tenant)
						tagArray[4] = fmt.Sprintf("handlingType_%s:", RequestType)
						tagArray[5] = fmt.Sprintf("state_%s:", "Available")
						tagArray[6] = fmt.Sprintf("resourceid_%s:", resObj.ResourceId)
						tagArray[7] = fmt.Sprintf("objtype_%s", "CSlotInfo")

						tags := fmt.Sprintf("tag:*%s*", strings.Join(tagArray, "*"))
						log.Println(tags)
						availableSlots := RedisSearchKeys(tags)

						for _, tagKey := range availableSlots {
							strslotKey := RedisGet(tagKey)
							log.Println(strslotKey)

							strslotObj := RedisGet(strslotKey)
							log.Println(strslotObj)

							var slotObj CSlotInfo
							json.Unmarshal([]byte(strslotObj), &slotObj)

							slotObj.State = "Reserved"
							slotObj.SessionId = sessionID
							slotObj.OtherInfo = "Inbound"
							slotObj.MaxReservedTime = metaData.MaxReservedTime
							slotObj.MaxAfterWorkTime = metaData.MaxAfterWorkTime
							slotObj.MaxFreezeTime = metaData.MaxFreezeTime
							slotObj.TempMaxRejectCount = metaData.MaxRejectCount

							if ReserveSlot(ardsLbIP, ardsLbPort, slotObj) == true {
								log.Println("Return resource Data:", conInfo.RefInfo)
								selectedResList = AppendIfMissingString(selectedResList, conInfo.RefInfo)
								selectedResKeyList = AppendIfMissingString(selectedResKeyList, key)
								if len(selectedResList) == nuOfResRequested {
									selectedResListString, _ := json.Marshal(selectedResList)
									handlingResult = string(selectedResListString)
									handlingResource = selectedResKeyList
								}
							}
						}
					}
				}
			}
		}

	}
	handlingResult = "No matching resources at the moment"
	handlingResource = make([]string, 0)
	return
}
