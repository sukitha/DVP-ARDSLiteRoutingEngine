package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

//Post Invokes third pary service with http post
func Post(serviceurl string, postData, authToken, internalAuthToken string) bool {
	log.Println("Start======================================:: ", time.Now().UTC())
	log.Println("URL:>", serviceurl)

	log.Println("PostData:>", postData)

	var jsonData = []byte(postData)
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", authToken)
	req.Header.Set("companyinfo", internalAuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//panic(err)
		return false
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	//log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	result := string(body)
	log.Println("response Body:", result)
	//log.Println("response CODE::", string(resp.StatusCode))
	log.Println("End======================================:: ", time.Now().UTC())
	if resp.StatusCode == 200 {
		log.Println("Return true")
		return true
	}

	log.Println("Return false")
	return false
}

//Put Invokes third pary service with http put
func Put(serviceurl string, postData, authToken, internalAuthToken string) bool {
	log.Println("Start======================================:: ", time.Now().UTC())
	log.Println("URL:>", serviceurl)

	log.Println("PostData:>", postData)

	var jsonData = []byte(postData)
	req, err := http.NewRequest("PUT", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", authToken)
	req.Header.Set("companyinfo", internalAuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//panic(err)
		return false
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	//log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	result := string(body)
	log.Println("response Body:", result)
	//log.Println("response CODE::", string(resp.StatusCode))
	log.Println("End======================================:: ", time.Now().UTC())
	if resp.StatusCode == 200 {
		log.Println("Return true")
		return true
	}

	log.Println("Return false")
	return false
}

//Get Invokes third pary service with http put
func Get(serviceurl, path, param string) string {
	request := fmt.Sprintf("http://%s", serviceurl)

	u, _ := url.Parse(request)
	u.Path += path
	u.Path += param

	log.Println(u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp != nil {

		response, _ := ioutil.ReadAll(resp.Body)
		tmx := string(response[:])
		log.Println(tmx)
		return tmx
	}

	return ""
}
