package main

import (
	"fmt"
	"net"
)

func AppendIfMissingReq(dataList []Request, i Request) []Request {
	for _, ele := range dataList {
		if ele.SessionId == i.SessionId {
			return dataList
		}
	}
	return append(dataList, i)
}

func AppendIfMissingString(dataList []string, i string) []string {
	for _, ele := range dataList {
		if ele == i {
			return dataList
		}
	}
	return append(dataList, i)
}

func CheckExistingString(dataList []string, i string) bool {
	for _, ele := range dataList {
		if ele == i {
			return true
		}
	}
	return false
}

func CreateHost(_ip, _port string) string {
	testIp := net.ParseIP(_ip)
	if testIp.To4() == nil {
		return _ip
	} else {
		return fmt.Sprintf("%s:%s", _ip, _port)
	}
}

func GetSelectedResourceForRequest(records []SelectionResult, sessionId string) (resourceForRequest SelectedResource, isExisting bool) {
	for _, record := range records {
		if record.Request == sessionId {
			resourceForRequest = record.Resources
			isExisting = true
			return
		}
	}
	isExisting = false
	return
}
