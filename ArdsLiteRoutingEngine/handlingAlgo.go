package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func SingleResourceAlgo(ReqClass, ReqType, ReqCategory, SessionId string, ResourceIds []string) string {
	var result = SingleHandling(ReqClass, ReqType, ReqCategory, SessionId, ResourceIds)
	return result

}

func ReserveSlot(slotInfo CSlotInfo) bool {
	url := fmt.Sprintf("%s/%s/concurrencyslot", resCsUrl, slotInfo.ResourceId)
	fmt.Println("URL:>", url)

	slotInfoJson, _ := json.Marshal(slotInfo)
	var jsonStr = []byte(slotInfoJson)
	authToken := fmt.Sprintf("%d#%d", slotInfo.Tenant, slotInfo.Company)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//panic(err)
		return false
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	result := string(body)
	fmt.Println("response Body:", result)

	var resConv updateCsReult
	json.Unmarshal(body, &resConv)

	if resConv.IsSuccess == true {
		fmt.Println("Return true")
		return true
	}

	fmt.Println("Return false")
	return false
}

func ClearSlotOnMaxRecerved(reqClass, reqType, reqCategory, sessionId string, resObj Resource, metaData ReqMetaData) {
	var tagArray = make([]string, 8)

	tagArray[0] = fmt.Sprintf("company_%d", resObj.Company)
	tagArray[1] = fmt.Sprintf("tenant_%d", resObj.Tenant)
	tagArray[2] = fmt.Sprintf("class_%s", reqClass)
	tagArray[3] = fmt.Sprintf("type_%s", reqType)
	tagArray[4] = fmt.Sprintf("category_%s", reqCategory)
	tagArray[5] = fmt.Sprintf("state_%s", "Reserved")
	tagArray[6] = fmt.Sprintf("resourceid_%s", resObj.ResourceId)
	tagArray[7] = fmt.Sprintf("objtype_%s", "CSlotInfo")

	tags := fmt.Sprintf("tag:*%s*", strings.Join(tagArray, "*"))
	fmt.Println(tags)
	reservedSlots := RedisSearchKeys(tags)

	for _, tagKey := range reservedSlots {
		strslotKey := RedisGet(tagKey)
		fmt.Println(strslotKey)

		strslotObj := RedisGet(strslotKey)
		fmt.Println(strslotObj)

		var slotObj CSlotInfo
		json.Unmarshal([]byte(strslotObj), &slotObj)

		fmt.Println("Datetime Info" + slotObj.LastReservedTime)
		t, _ := time.Parse(layout, slotObj.LastReservedTime)
		t1 := int(time.Now().Sub(t).Seconds())
		t2 := metaData.MaxReservedTime
		fmt.Println(fmt.Sprintf("Time Info T1: %d", t1))
		fmt.Println(fmt.Sprintf("Time Info T2: %d", t2))
		if t1 > t2 {
			slotObj.State = "Available"
			slotObj.OtherInfo = "ClearReserved"

			ReserveSlot(slotObj)
		}
	}
}

func GetReqMetaData(_company, _tenent int, _class, _type, _category string) ReqMetaData {
	key := fmt.Sprintf("ReqMETA:%d:%d:%s:%s:%s", _company, _tenent, _class, _type, _category)
	fmt.Println(key)
	strMetaObj := RedisGet(key)
	fmt.Println(strMetaObj)

	var metaObj ReqMetaData
	json.Unmarshal([]byte(strMetaObj), &metaObj)

	return metaObj
}

func GetResourceState(_company, _tenant int, _resId string) string {
	key := fmt.Sprintf("ResourceState:%d:%d:%s", _company, _tenant, _resId)
	fmt.Println(key)
	strResStateObj := RedisGet(key)
	fmt.Println(strResStateObj)

	return strResStateObj
}
