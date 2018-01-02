package main

import (
	//"fmt"

	//"github.com/golang/protobuf/proto"
	"encoding/json"
	"net/http"
	"sync"
	//"protocol"
	//"strings"
	"time"
	"ServerManager"
	"fmt"
	"utils"
	//"net/url"
)

var g_nodeList sync.Map

func AgentDealer(w http.ResponseWriter, r *http.Request) {
	strAction := r.FormValue("action")
	strData := r.FormValue("data")
	utils.Debugln("Action==",strAction,"data=",strData)
	if strAction == "agent_heartbeat" {
		DealWithAgentHeartbeat(w, r, strData)
	}
	if strAction == "process_getserverconfig" {
		DealWithGetProcessConfig(w, r, strData)
	}

	if strAction == "process_heartbeat" {
		DealWithProcessHeartbeat(w, r, strData)
	}

}

func DealWithGetProcessConfig(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	var stDataMap map[string]string
	err:=json.Unmarshal([]byte(strData), &stDataMap)
	checkErr(err)
	strServerName:=stDataMap["ServerName"]
	strInsid := stDataMap["Insid"]
    nRet,strResult:=ServerManager.GetFacadeMysql().GetProcessConfigContent(strServerName,strInsid)
    if nRet ==0{
		retMap["data"] = strResult
		fmt.Fprint(w, RetMap2String(retMap))
	}else{
		retMap["ret"] = "-1"
		fmt.Fprint(w, RetMap2String(retMap))
	}

}

func DealWithAgentHeartbeat(w http.ResponseWriter, r *http.Request, strData string) {
	var stDataMap map[string]string
	err:=json.Unmarshal([]byte(strData), &stDataMap)
	checkErr(err)

	g_nodeList.Store(stDataMap["Innerip"], stDataMap)
	utils.Debugln("DealWithAgentHeartbeat",stDataMap)

}

func DealWithProcessHeartbeat(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	var stDataMap map[string]string
	err:=json.Unmarshal([]byte(strData), &stDataMap)
	checkErr(err)
	//stDataMap["Path"],_ = url.QueryUnescape(stDataMap["Path"] )
	ServerManager.GetFacadeMysql().UpdateProcess(stDataMap)
	fmt.Fprint(w, RetMap2String(retMap))

}

func  StartMysqlWriteTimer() error {
	utils.Debugln("StartMysqlWriteTimer begin")
	duration := time.Second*60
	utils.Debugln("duration=", duration)
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		utils.Debugln("StartTimer.time done", duration)
		UpdateData2Mysql()
		StartMysqlWriteTimer()
	}()
	return nil
}
func WriteNodeList2Mysql(strKey, mapValue interface{}) bool {
	ServerManager.GetFacadeMysql().UpdateNodeList(string(strKey.(string)),(mapValue.(map[string]string)))
return true
}
func UpdateData2Mysql()  {
	g_nodeList.Range(WriteNodeList2Mysql)
}