package main

import (
	"ServerManager"
	"fmt"
	"html/template"
	"io"
	//"io/ioutil"
	"os"
	//"db"
	"encoding/json"
	//"github.com/golang/protobuf/proto"
	"net/http"
	//"protocol"
	"io/ioutil"
	"sync"
	//"strings"
	"errors"
	"reflect"
	"time"
	"utils"
	"strconv"
)

var ServerGroups []string

type oneServerItem struct {
	StrIp       string
	StrPort     string
	StrInstance string
	StrPath     string
	Status      string
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var g_agentFuncMap sync.Map

type ServerConf struct {
	MysqlConf string
	HttpPort  string
}
var G_StConf ServerConf
func loadConf() {
	fi, err := os.Open("conf.json")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	utils.Debugln("loadConf:", string(fd))
	err = json.Unmarshal(fd, &G_StConf)
	if err != nil {
		utils.Debugln("error:", err)
		os.Exit(0)
	}
	return
}

func StartWeb() error {
	loadConf()
	ServerManager.G_MysqlConf = G_StConf.MysqlConf
	http.Handle("/", http.FileServer(http.Dir("resources/")))
	http.HandleFunc("/index", leftbarHandler)
	http.HandleFunc("/ajax", ajaxHandler)
	http.HandleFunc("/FileUpload", fileUpload)
	http.HandleFunc("/agent", AgentDealer)
	err := http.ListenAndServe(G_StConf.HttpPort, nil)
	return err

}

func leftbarHandler(w http.ResponseWriter, r *http.Request) {

	//ServerGroups := GetServerGroups()
	t, err := template.ParseFiles("template/leftbar.html")
	if err != nil {
		utils.Debugln("ParseFiles error", err.Error())
	}

	err = t.Execute(w, nil)
	if err != nil {
		utils.Debugln("ParseFiles error", err.Error())
	}

}

func fileUpload(w http.ResponseWriter, r *http.Request) {
	utils.Debugln("fileUpload---", *r)
	utils.Debugln()
	file, fileHead, err := r.FormFile("file_data")

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	strFileName := fileHead.Filename
	defer file.Close()
	os.MkdirAll("uploadFold", os.ModeDir)
	f, err := os.Create("uploadFold/" + strFileName)
	defer f.Close()
	io.Copy(f, file)
	ServerManager.GetFacadeMysql().InsertFile2List(strFileName)
	//fmt.Fprintf(w, "????????: %d", file.(Sizer).Size())
	retMap := GetRetMap()

	fmt.Fprint(w, RetMap2String(retMap))
}
func FunctionMapCall(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func ajaxHandler(w http.ResponseWriter, r *http.Request) {
	funcs := make(map[string]interface{})
	funcs["getNodeList"] = getNodeList
	funcs["getTemplates"] = getTemplates
	funcs["updateNodeRemarks"] = updateNodeRemarks
	funcs["DelNodes"] = DelNodes
	funcs["getFileList"] = getFileList
	funcs["deleteUploadFile"] = DeleteUploadFile
	funcs["getConfigTemplateList"] = GetAllConfigTemplates
	funcs["SaveConfigTemplate"] = SaveConfigTemplate
	funcs["DelConfigTemplate"] = DelConfigTemplate
	funcs["updateProcess"] = updateProcess
	funcs["mkNewProcess"] = mkNewProcess
	funcs["getProcessList"] = getProcessList
	funcs["deleteprocess"] = deleteProcess
	funcs["getrpclist"] = getrpclist
	funcs["saverpclist"] = saverpclist
	funcs["delrpclist"] = delrpclist
	funcs["slgrpctest"] = slgrpctest

	strAction := r.FormValue("action")
	strData := r.FormValue("data")
	utils.Debugln("ajaxHandler---------------", strAction, strData)
	FunctionMapCall(funcs, strAction, w, r, strData)
}


func delrpclist(w http.ResponseWriter, r *http.Request, strData string){
	retMap := GetRetMap()
	mapParam := make(map[string]string)
	json.Unmarshal([]byte(strData), &mapParam)
	var stRpcTestItem ServerManager.RpcTestItem
	stRpcTestItem.Module = mapParam["Module"]
	stRpcTestItem.Object = mapParam["Object"]
	stRpcTestItem.Function = mapParam["Function"]
	stRpcTestItem.Data = mapParam["Data"]
	err:=ServerManager.GetFacadeMysql().DelRpcList(&stRpcTestItem)
	if err != nil{
		retMap["ret"] = "-1"
		retMap["data"] = string(err.Error())
	}
	fmt.Fprint(w, RetMap2String(retMap))
}

func saverpclist(w http.ResponseWriter, r *http.Request, strData string){
	retMap := GetRetMap()
	mapParam := make(map[string]string)
	json.Unmarshal([]byte(strData), &mapParam)
	var stRpcTestItem ServerManager.RpcTestItem
	stRpcTestItem.Module = mapParam["Module"]
	stRpcTestItem.Object = mapParam["Object"]
	stRpcTestItem.Function = mapParam["Function"]
	stRpcTestItem.Data = mapParam["Data"]
	err:=ServerManager.GetFacadeMysql().Save2RpcList(&stRpcTestItem)
	if err != nil{
		retMap["ret"] = "-1"
		retMap["data"] = string(err.Error())
	}
	fmt.Fprint(w, RetMap2String(retMap))
}



func getrpclist(w http.ResponseWriter, r *http.Request, strData string){
	retMap := GetRetMap()
	rpcList:=ServerManager.GetFacadeMysql().GetRpcList()
	buffer,_:=json.Marshal(rpcList)
	retMap["data"] = string(buffer)
	fmt.Fprint(w, RetMap2String(retMap))
}

func slgrpctest(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	mapParam := make(map[string]string)
	json.Unmarshal([]byte(strData), &mapParam)
    strIpPort:=mapParam["IpPort"]
	strRequestData:=mapParam["Data"]
	strServerId:=mapParam["ServerId"]
	strUid:=mapParam["Uid"]
	strObject := mapParam["Object"]
	strFunction := mapParam["Function"]
    nModule,_:=strconv.Atoi(mapParam["Module"])
	nServerId,_:=strconv.ParseInt(strServerId,10,32)
	nUid,_:=strconv.ParseInt(strUid,10,64)

	strOut:=ServerManager.DealWithRpcTest(strIpPort,nModule,strObject,strFunction,strRequestData,uint32(nServerId),nUid)
	retMap["data"] = strOut
	fmt.Fprint(w, RetMap2String(retMap))

	//DealWithRpcTest(strIp string,,strObj string,uint32 , nServerId uint32,nUid int64) string
}

func deleteProcess(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	mapParam := make(map[string]string)
	json.Unmarshal([]byte(strData), &mapParam)
	retList := ServerManager.GetFacadeMysql().DelProcessList(mapParam["ServerName"], mapParam["Insid"])
	retBuff, _ := json.Marshal(retList)
	retMap["data"] = string(retBuff)
	fmt.Fprint(w, RetMap2String(retMap))
}

func getProcessList(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	mapParam := make(map[string]string)
	json.Unmarshal([]byte(strData), &mapParam)
	retList := ServerManager.GetFacadeMysql().GetProcessList()
	retBuff, _ := json.Marshal(retList)
	retMap["data"] = string(retBuff)
	fmt.Fprint(w, RetMap2String(retMap))
}

func mkNewProcess(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	mapParam := make(map[string]string)
	json.Unmarshal([]byte(strData), &mapParam)
	ServerManager.GetFacadeMysql().MkNewProcess(mapParam["Innerip"], mapParam["ServerName"], mapParam["Insid"], mapParam["ConfigContent"])
	fmt.Fprint(w, RetMap2String(retMap))
}

func updateProcess(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	mapParam := make(map[string]string)
	json.Unmarshal([]byte(strData), &mapParam)
	ServerManager.GetFacadeMysql().UpdateProcess(mapParam)
	fmt.Fprint(w, RetMap2String(retMap))
}

func DelConfigTemplate(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	ServerManager.GetFacadeMysql().DelConfigTemplate(strData)
	fmt.Fprint(w, RetMap2String(retMap))
}

//保存配置模版
func SaveConfigTemplate(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	var stConfigTemplate ServerManager.ConfigTemplate
	json.Unmarshal([]byte(strData), &stConfigTemplate)
	stConfigTemplate.EditTime = time.Now().Unix()
	err := ServerManager.GetFacadeMysql().ReplaceConfigTemplateList(stConfigTemplate)
	if err != nil {
		retMap["ret"] = "-1"
		retMap["data"] = string(err.Error())
	} else {
		ConfigTemplateList := ServerManager.GetFacadeMysql().GetAllConfigTemplateList()
		strRetData, _ := json.Marshal(ConfigTemplateList)
		retMap["data"] = string(strRetData)
	}
	fmt.Fprint(w, RetMap2String(retMap))
}

//获取所有配置模版
func GetAllConfigTemplates(w http.ResponseWriter, r *http.Request, strData string) {
	//GetAllConfigTemplateList
	retMap := GetRetMap()
	ConfigTemplateList := ServerManager.GetFacadeMysql().GetAllConfigTemplateList()
	strRetData, _ := json.Marshal(ConfigTemplateList)
	retMap["data"] = string(strRetData)
	fmt.Fprint(w, RetMap2String(retMap))
}

//删除上传的文件
func DeleteUploadFile(w http.ResponseWriter, r *http.Request, strData string) {

	strFileName := "uploadFold/" + strData
	os.Remove(strFileName)
	ServerManager.GetFacadeMysql().DeleteUploadFile(strData)
	retMap := GetRetMap()
	fmt.Fprint(w, RetMap2String(retMap))
}

//获取文件列表
func getFileList(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	FileList := ServerManager.GetFacadeMysql().GetFileList()
	strRetData, _ := json.Marshal(FileList)
	retMap["data"] = string(strRetData)
	fmt.Fprint(w, RetMap2String(retMap))

	return
}

//更新节点的备注信息
func updateNodeRemarks(w http.ResponseWriter, r *http.Request, strData string) {
	var updateList []ServerManager.KeyValue
	retMap := GetRetMap()
	err := json.Unmarshal([]byte(strData), &updateList)
	if err != nil {
		utils.Debugln("JsonDeserialize fail", string(strData))
		retMap["ret"] = "-1"
		fmt.Fprint(w, RetMap2String(retMap))
	}
	err = ServerManager.GetFacadeMysql().UpdateNodeListRemarks(updateList)
	if err != nil {
		utils.Debugln("UpdateNodeList fail", err.Error())
		retMap["ret"] = "-1"
		retMap["data"] = err.Error()
		fmt.Fprint(w, RetMap2String(retMap))
	}
	fmt.Fprint(w, RetMap2String(retMap))

	return
}

func DelNodes(w http.ResponseWriter, r *http.Request, strData string) {
	var updateList []ServerManager.KeyValue
	retMap := GetRetMap()
	err := json.Unmarshal([]byte(strData), &updateList)
	if err != nil {
		utils.Debugln("JsonDeserialize fail", string(strData))
		retMap["ret"] = "-1"
		fmt.Fprint(w, RetMap2String(retMap))
	}
	for _, oneKeyValue := range updateList {
		g_nodeList.Delete(oneKeyValue.Key)
	}
	err = ServerManager.GetFacadeMysql().DelNodes(updateList)
	if err != nil {
		utils.Debugln("UpdateNodeList fail", err.Error())
		retMap["ret"] = "-1"
		retMap["data"] = err.Error()
		fmt.Fprint(w, RetMap2String(retMap))
	}
	fmt.Fprint(w, RetMap2String(retMap))

	return
}

func getOneGroupPublishServerList(w http.ResponseWriter, r *http.Request, strData string) {
	var dataMap map[string]string
	retMap := GetRetMap()
	err := json.Unmarshal([]byte(strData), &dataMap)
	if err != nil {
		utils.Debugln("JsonDeserialize fail", string(strData))
		retMap["ret"] = "-1"
		fmt.Fprint(w, RetMap2String(retMap))
	}
	t, err := template.ParseFiles("template/publish.html")
	var ServerList []ServerManager.ServerItem = ServerManager.GetFacadeMysql().GetServerListByServerName(dataMap["ServerGroupName"])

	err = t.Execute(w, ServerList)
	if err != nil {
		utils.Debugln("Execute fail", err.Error())
		retMap["ret"] = "-1"
		retMap["data"] = "Execute fail"
		fmt.Fprint(w, RetMap2String(retMap))
	}
	return
}

func getOneGroupServerList(w http.ResponseWriter, r *http.Request, strData string) {
	var dataMap map[string]string
	retMap := GetRetMap()
	err := json.Unmarshal([]byte(strData), &dataMap)
	if err != nil {
		utils.Debugln("JsonDeserialize fail", string(strData))
		retMap["ret"] = "-1"
		fmt.Fprint(w, RetMap2String(retMap))
	}
	t, err := template.ParseFiles("template/righttable.html")
	var ServerList []ServerManager.ServerItem = ServerManager.GetFacadeMysql().GetServerListByServerName(dataMap["ServerGroupName"])

	err = t.Execute(w, ServerList)
	if err != nil {
		utils.Debugln("Execute fail", err.Error())
		retMap["ret"] = "-1"
		retMap["data"] = "Execute fail"
		fmt.Fprint(w, RetMap2String(retMap))
	}
	return
}

func GetRetMap() map[string]string {
	var retDataMap map[string]string = make(map[string]string)
	retDataMap["ret"] = "0"
	retDataMap["data"] = ""
	return retDataMap
}

func RetMap2String(retMap map[string]string) string {
	v, _ := json.Marshal(retMap)
	return string(v)
}
func GetServerGroups() []string {
	var ServerGroups []string
	ServerGroups = ServerManager.GetFacadeMysql().GetAllServerName()
	return ServerGroups
}

func getNodeList(w http.ResponseWriter, r *http.Request, strData string) {
	NodeList := ServerManager.GetFacadeMysql().GetNodeList()
	strRet, _ := json.Marshal(NodeList)
	fmt.Fprint(w, string(strRet))
	/*var nodeServerpkg protocol.NodeServerPkg
	nodeServerpkg.Cmd = proto.Int32(17)
	nodeServerpkg.Data = []byte("aliantest")
	data, err := proto.Marshal(&nodeServerpkg)
	if err != nil {
		utils.Debugln("marshaling error: ", err)
	}
	utils.Debugln("nodeServerpkg==", string(data))

	FacadeMysql := db.GetFacadeMysql()
	// MkNewServer(strInnerIp string, strServerName string, strInsId string, strPath string, port int, strOuterIp string, strOther string)
	FacadeMysql.MkNewServer("10.20.104.175", "Servername", "1", "path", 1001, "outerip", "other")*/
	return
}

func AjaxCreateNewServer(w http.ResponseWriter, r *http.Request, strData string) {

}

func AjaxStopServer(w http.ResponseWriter, r *http.Request, strData string) {

}

func AjaxStartServer(w http.ResponseWriter, r *http.Request, strData string) {

}

func getTemplates(w http.ResponseWriter, r *http.Request, strData string) {
	var FileList []string
	json.Unmarshal([]byte(strData), &FileList)
	//FileList := []string{"nodelistTemplate.html", "ProcessListTemplate.html", "ServerManagerPage.html", "fileManage.html", "fileuploadTemplate.html"}
	var templatesMap map[string]string = make(map[string]string)
	retMap := GetRetMap()
	for i := 0; i < len(FileList); i++ {
		strFullName := "template/" + FileList[i]
		fi, err := os.Open(strFullName)
		if err != nil {
			panic(err)
		}
		defer fi.Close()
		fd, err := ioutil.ReadAll(fi)
		templatesMap[FileList[i]] = string(fd)
		utils.Debugln("getTemplates:", string(fd))

	}
	templatesFiles2String := RetMap2String(templatesMap)
	retMap["data"] = templatesFiles2String
	fmt.Fprint(w, RetMap2String(retMap))
}

func main() {
	StartMysqlWriteTimer()
	StartWeb()
	utils.SetLevel(4)

}

/*
func Upload(url, file string) (err error) {
    // Prepare a form that you will submit to that URL.
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    // Add your image file
    f, err := os.Open(file)
    if err != nil {
        return
    }
    defer f.Close()
    fw, err := w.CreateFormFile("image", file)
    if err != nil {
        return
    }
    if _, err = io.Copy(fw, f); err != nil {
        return
    }
    // Add the other fields
    if fw, err = w.CreateFormField("key"); err != nil {
        return
    }
    if _, err = fw.Write([]byte("KEY")); err != nil {
        return
    }
    // Don't forget to close the multipart writer.
    // If you don't close it, your request will be missing the terminating boundary.
    w.Close()

    // Now that you have a form, you can submit it to your handler.
    req, err := http.NewRequest("POST", url, &b)
    if err != nil {
        return
    }
    // Don't forget to set the content type, this will contain the boundary.
    req.Header.Set("Content-Type", w.FormDataContentType())

    // Submit the request
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        return
    }

    // Check the response
    if res.StatusCode != http.StatusOK {
        err = fmt.Errorf("bad status: %s", res.Status)
    }
    return
}
*/
