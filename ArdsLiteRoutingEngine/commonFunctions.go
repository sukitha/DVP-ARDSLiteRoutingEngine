package main

func AppendIfMissingReq(dataList []Request, i Request) []Request {
	for _, ele := range dataList {
		if ele == i {
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
