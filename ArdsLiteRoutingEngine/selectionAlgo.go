package main

import (
	"encoding/json"
	"fmt"
)

func BasicSelectionAlgo(Company, Tenant int, SessionId string) []string {

	fmt.Println(Company)
	fmt.Println(Tenant)
	fmt.Println(SessionId)

	var result = BasicSelection(Company, Tenant, SessionId)
	return result

}

func WeightBaseSelectionAlgo(Company, Tenant int, SessionId string) []string {

	fmt.Println(Company)
	fmt.Println(Tenant)
	fmt.Println(SessionId)

	var result = WeightBaseSelection(Company, Tenant, SessionId)
	return result

}

func GetConcurrencyInfo(_company, _tenant int, _resId, _category string) ConcurrencyInfo {
	key := fmt.Sprintf("ConcurrencyInfo:%d:%d:%s:%s", _company, _tenant, _resId, _category)
	fmt.Println(key)
	strCiObj := RedisGet(key)
	fmt.Println(strCiObj)

	var ciObj ConcurrencyInfo
	json.Unmarshal([]byte(strCiObj), &ciObj)

	return ciObj
}
