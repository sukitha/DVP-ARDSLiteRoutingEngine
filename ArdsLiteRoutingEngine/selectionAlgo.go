package main

import (
	"encoding/json"
	"fmt"
)

func IsAttributeAvailable(reqAttributeInfo []ReqAttributeData, resAttributeInfo []ResAttributeData) (isAttrAvailable, isThreshold bool) {
	isAttrAvailable = false
	isThreshold = false

	for _, reqAtt := range reqAttributeInfo {
		if len(reqAtt.AttributeCode) > 0 {
			attCode := reqAtt.AttributeCode[0]

			for _, resAtt := range resAttributeInfo {
				if attCode == resAtt.Attribute && resAtt.HandlingType == reqAtt.HandlingType {
					isAttrAvailable = true

					if resAtt.Percentage > 0 && resAtt.Percentage <= 25 {
						isThreshold = true
					}
					return
				}
			}
		}
	}
	return
}

func BasicSelectionAlgo(Company, Tenant int, SessionId string) []string {

	fmt.Println(Company)
	fmt.Println(Tenant)
	fmt.Println(SessionId)

	var result = BasicSelection(Company, Tenant, SessionId)
	return result

}

func BasicThresholdSelectionAlgo(Company, Tenant int, SessionId string) []string {

	fmt.Println(Company)
	fmt.Println(Tenant)
	fmt.Println(SessionId)

	var result = BasicThresholdSelection(Company, Tenant, SessionId)
	return result

}

func WeightBaseSelectionAlgo(Company, Tenant int, SessionId string) []string {

	fmt.Println(Company)
	fmt.Println(Tenant)
	fmt.Println(SessionId)

	var result = WeightBaseSelection(Company, Tenant, SessionId)
	return result

}

func GetConcurrencyInfo(_company, _tenant int, _resId, _category string) (ciObj ConcurrencyInfo, err error) {
	key := fmt.Sprintf("ConcurrencyInfo:%d:%d:%s:%s", _company, _tenant, _resId, _category)
	fmt.Println(key)
	var strCiObj string
	strCiObj, err = RedisGet_v1(key)
	fmt.Println(strCiObj)

	json.Unmarshal([]byte(strCiObj), &ciObj)

	return
}
