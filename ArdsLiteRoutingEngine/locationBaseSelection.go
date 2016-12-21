package main

import (
	"encoding/json"
	"fmt"
)

func LocationBaseSelection(_company, _tenent int, _sessionId string) (result SelectionResult) {
	var matchingResources = make([]string, 0)

	requestKey := fmt.Sprintf("Request:%d:%d:%s", _company, _tenent, _sessionId)
	fmt.Println(requestKey)

	strReqObj := RedisGet(requestKey)
	fmt.Println(strReqObj)

	var reqObj RequestSelection
	json.Unmarshal([]byte(strReqObj), &reqObj)

	if reqObj.OtherInfo != "" {

		var locationObj ReqLocationData
		json.Unmarshal([]byte(reqObj.OtherInfo), &locationObj)

		if locationObj != (ReqLocationData{}) {

			locationResult := RedisGeoRadius(locationObj)

			for _, lor := range locationResult.Elems {

				resourceLocInfo, _ := lor.List()

				if len(resourceLocInfo) > 1 {
					resourceKey := fmt.Sprintf("Resource:%d:%d:%s", reqObj.Company, reqObj.Tenant, resourceLocInfo[0])
					strResObj := RedisGet(resourceKey)
					fmt.Println(strResObj)

					var resObj Resource
					json.Unmarshal([]byte(strResObj), &resObj)

					if resObj.ResourceId != "" {
						resKey := fmt.Sprintf("Resource:%d:%d:%s", resObj.Company, resObj.Tenant, resObj.ResourceId)
						if len(reqObj.AttributeInfo) > 0 {
							_attAvailable, _ := IsAttributeAvailable(reqObj.AttributeInfo, resObj.ResourceAttributeInfo)
							if _attAvailable {
								matchingResources = AppendIfMissingString(matchingResources, resKey)
							}
						} else {
							matchingResources = AppendIfMissingString(matchingResources, resKey)
						}
					}
				}
			}

			result.Priority = matchingResources
			return
		} else {
			return
		}

	} else {
		return
	}

}
