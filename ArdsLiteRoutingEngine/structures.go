package main

//Configurations

type Configuration struct {
	RedisIp   string
	RedisPort string
	RedisDb   int
	Port      string
}

type EnvConfiguration struct {
	RedisIp   string
	RedisPort string
	RedisDb   string
	Port      string
}

//Request

type ReqAttributeData struct {
	AttributeCode      []string
	AttributeGroupName string
	HandlingType       string
	WeightPrecentage   string
}

type Request struct {
	Company          int
	Tenant           int
	ServerType       string
	RequestType      string
	SessionId        string
	ArriveTime       string
	Priority         string
	QueueId          string
	ReqHandlingAlgo  string
	ReqSelectionAlgo string
	ServingAlgo      string
	HandlingAlgo     string
	SelectionAlgo    string
	RequestServerUrl string
	HandlingResource string
	ResourceCount    int
	OtherInfo        string
	LbIp             string
	LbPort           string
}

type RequestSelection struct {
	Company       int
	Tenant        int
	ServerType    string
	RequestType   string
	SessionId     string
	AttributeInfo []ReqAttributeData
}

type ReqMetaData struct {
	MaxReservedTime int
	MaxRejectCount  int
}

//Resource

type ResAttributeData struct {
	Attribute    string
	HandlingType string
	Percentage   float64
}

type Resource struct {
	Company               int
	Tenant                int
	Class                 string
	Type                  string
	Category              string
	ResourceId            string
	ResourceAttributeInfo []ResAttributeData
	OtherInfo             string
}

type CSlotInfo struct {
	Company          int
	Tenant           int
	HandlingType     string
	State            string
	HandlingRequest  string
	ResourceId       string
	SlotId           int
	ObjKey           string
	SessionId        string
	LastReservedTime string
	MaxReservedTime  int
	OtherInfo        string
}

type ConcurrencyInfo struct {
	RejectCount       int
	ResourceId        string
	LastConnectedTime string
	RefInfo           string
}

type WeightBaseResourceInfo struct {
	ResourceId string
	Weight     float64
}

type MultiResCount struct {
	ResourceCount int
}

type updateCsReult struct {
	IsSuccess bool
}
