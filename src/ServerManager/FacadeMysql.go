package ServerManager

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"net/http"
	"strings"
	"sync"
	"time"
	"utils"
)

var MysqlLocker sync.Mutex

type FacadeMysql struct {
	strDbConn string
	dbconn    *sql.DB
}

type ServerItem struct {
	Servername        string
	Insid             string
	Innerip           string
	Outerip           string
	Path              string
	ConfPath          string
	Port              int32
	Status            int32
	Lastupdatetime    string
	StatusString      string
	ConfigFileContent string
	Versin            string
	Other             string
}
type NodeItem struct {
	NodeId         int
	Innerip        string
	Outerip        string
	LastUpdateTime int64
	Remarks        string
}

type ProcessItem struct {
	Innerip        string
	Outerip        string
	ServerName     string
	Insid          string
	Path           string
	Port           string
	ConfigContent  string
	Status         string
	Lastupdatetime string
	Other          string
}

type FileItem struct {
	FileName string
	Time     int64
	Content  []byte
}

type ConfigTemplate struct {
	TemplateName string
	TemplateType int
	EditTime     int64
	Content      string
	Remarks      string
}

type RpcTestItem struct {
	Module string
	Object string
	Function     string
	Data      string
}

type KeyValue struct {
	Key   string
	Value string
}

var G_MysqlConf string
func checkErr(err error) {
	if err != nil {
		fmt.Errorf(err.Error())
		GetFacadeMysql().Init(GetFacadeMysql().strDbConn)
	}
}
func (this *FacadeMysql) Init(strconn string) int {

	dbconn, err := sql.Open("mysql", strconn)
	checkErr(err)
	this.strDbConn = strconn
	this.dbconn = dbconn
	return 0
}

func (this *FacadeMysql) UpdateLastTime(strServerName string, strInsId string, nTime uint64) int {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("update serverlist set lastupdatetime=? where servername=? and insid=?")
	checkErr(err)

	res, err := stmt.Exec(nTime, strServerName, strInsId)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	utils.Debugln("affect==", affect)
	return 0
}

func (this *FacadeMysql) MkNewProcess(strInnerIp string, strServerName string, strInsId string, strConfContent string) int {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("insert into processlist set Innerip=?,Insid=?,ConfigContent=?,Servername=? ")
	checkErr(err)

	res, err := stmt.Exec(strInnerIp, strInsId, strConfContent, strServerName)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	utils.Debugln("affect==", affect)
	return 0
}

func (this *FacadeMysql) GetProcessConfigContent(strServerName string, strInsId string) (int, string) {
	retTemplateList := this.GetAllConfigTemplateList()
	g_all_ConfigTemplate := make(map[string]ConfigTemplate)
	for _, oneTemplate := range retTemplateList {
		g_all_ConfigTemplate[oneTemplate.TemplateName] = oneTemplate
	}
	utils.Debugln("GetProcessConfigContent",strServerName,strInsId)
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("select  ConfigContent from  processlist where Servername =? and Insid=? ")
	checkErr(err)
	strTemplate := ""
	rows, err := stmt.Query(strServerName, strInsId)
	if err == nil {
		for rows.Next() {
			err = rows.Scan(&strTemplate)
		}
	}
	var strResult = ""
	if strTemplate != "" {
		spliteResult := strings.Split(strTemplate, "$")
		for i, Value := range spliteResult {
			if i%2 == 1 {
				if _, OK := g_all_ConfigTemplate[Value]; !OK {
					return -1, ""
				}
				spliteResult[i] = g_all_ConfigTemplate[Value].Content
			}
		}

		for i := 0; i < len(spliteResult); i++ {
			strResult = strResult + spliteResult[i]
		}
	}

	/*
		var Content = $("#ConfigEditWind_content").val()
		var spliteResult = Content.split("$")
		for (var i = 0; i < spliteResult.length; i++) {
		if (i % 2 == 1) {
		if (!(spliteResult[i] in g_all_ConfigTemplate)) {
		alert("模版编译失败")
		return ""
		} else {
		spliteResult[i] = g_all_ConfigTemplate[spliteResult[i]].Content
		}
		}
		}
		var strResult = ""
		for (var i = 0; i < spliteResult.length; i++) {
		strResult = strResult + spliteResult[i]
		}
		alert(strResult)*/
	return 0,strResult
}

func (this *FacadeMysql) UpdateProcess(ColumMap map[string]string) int {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	strSql := "update processlist set "
	//nIndex := 0
	for key, value := range ColumMap {
		if key != "ServerName" && key != "Insid" {
			strSql = strSql + key + "=" + "'" + value + "'"
				strSql += ","
		}
	}
	strSql = strSql[0:len(strSql)-1]
	strSql += " where ServerName='" + ColumMap["ServerName"] + "' and Insid='" + ColumMap["Insid"] + "'"
	utils.Debugln("strSql=======", strSql)
	rusult, err := this.dbconn.Exec(strSql)
	utils.Debugln("rusult=======", rusult)
	if err != nil {
		utils.Debugln("err=======", err.Error())
	}
	return 0
}

/*
DROP TABLE IF EXISTS `processlist`;
CREATE TABLE `processlist` (
  `Innerip` varchar(32) NOT NULL,
  `Outerip` varchar(32) NOT NULL,
  `ServerName` varchar(64) NOT NULL,
  `Insid` varchar(8) NOT NULL,
  `Path` varchar(256) DEFAULT NULL,
  `Port` int(11) DEFAULT NULL,
  `ConfigContent` text,
  `Status` int(11) DEFAULT NULL,
  `Lastupdatetime` int(64) DEFAULT NULL,
  `Other` varchar(156) DEFAULT NULL,
  PRIMARY KEY (`Insid`,`ServerName`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
*/
func (this *FacadeMysql) GetProcessList() (retList []ProcessItem) {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("select * from processlist")
	checkErr(err)

	rows, err := stmt.Query()
	if err == nil {
		for rows.Next() {
			var stOneProcessItem ProcessItem
			err = rows.Scan(&stOneProcessItem.ServerName, &stOneProcessItem.Insid, &stOneProcessItem.Innerip, &stOneProcessItem.Outerip, &stOneProcessItem.ConfigContent, &stOneProcessItem.Path, &stOneProcessItem.Port, &stOneProcessItem.Status, &stOneProcessItem.Lastupdatetime, &stOneProcessItem.Other)
			retList = append(retList, stOneProcessItem)

		}
	}
	return retList
}

func (this *FacadeMysql) DelProcessList (strServerName,strInsid string) int{
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("delete  from processlist where ServerName=? and Insid=?")
	checkErr(err)
	_, err = stmt.Exec(strServerName,strInsid)
	if err == nil {
		return 0
	}
	return -1
}

//server自发现
func (this *FacadeMysql) ServerDiscovery(strInnerIp string, strServerName string, strInsId string, strPath string, port int32, status int) int {

	bIsExists := this.CheckIsServerExists(strServerName, strInsId)
	if bIsExists {
		utils.Debugln("Server is exist in db", strServerName, strInsId)
		return 0
	}
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("insert into serverlist values(?,?,?,?,?,?,?,?,?)")

	res, err := stmt.Exec(strInnerIp, strServerName, strInsId, strPath, port, status, time.Now().Unix(), "", "")
	if err != nil {
		utils.Debugln("Server is exist in db", strServerName, strInsId, err.Error())
		return 0
	}

	affect, err := res.RowsAffected()
	utils.Debugln("affect==", affect)
	return 0
}

func (this *FacadeMysql) GetServerListByServerName(strServerName string) []ServerItem {
	return nil
}

func (this *FacadeMysql) GetAllServerName() []string {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("select DISTINCT servername from  serverlist")
	checkErr(err)
	var retList []string
	rows, err := stmt.Query()
	if err == nil {
		for rows.Next() {
			var ServerName string
			err = rows.Scan(&ServerName)
			if err == nil {
				retList = append(retList, ServerName)
			}
		}
	}
	return retList
}

func (this *FacadeMysql) CheckIsServerExists(strServerName string, strInsId string) bool {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("select * from  serverlist where servername=? and insid=?")
	checkErr(err)
	var retList []ServerItem
	rows, err := stmt.Query(strServerName, strInsId)
	if err == nil {
		for rows.Next() {
			var stOneServerItem ServerItem
			err = rows.Scan(&stOneServerItem.Innerip, &stOneServerItem.Servername, &stOneServerItem.Insid, &stOneServerItem.Path, &stOneServerItem.Port, &stOneServerItem.Status, &stOneServerItem.Lastupdatetime, &stOneServerItem.Outerip, &stOneServerItem.Other)
			if err == nil {
				retList = append(retList, stOneServerItem)
			} else {
				utils.Debugln("Scan error", err.Error())
			}
		}
	} else {
		utils.Debugln("stmt.Query error", err.Error())
	}
	if len(retList) == 0 {
		return false
	}
	return true
}

//获取所有服务器节点
func (this *FacadeMysql) GetNodeList() (retList []NodeItem) {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("select * from  nodelist order by InnerIp")
	checkErr(err)
	rows, err := stmt.Query()
	if err == nil {
		for rows.Next() {
			var stOneNodeItem NodeItem
			err = rows.Scan(&stOneNodeItem.Innerip, &stOneNodeItem.Outerip, &stOneNodeItem.LastUpdateTime, &stOneNodeItem.Remarks)
			retList = append(retList, stOneNodeItem)

		}
	} else {
		utils.Debugln("stmt.Query error", err.Error())
	}
	utils.Debugln("GetNodeList retList==", retList)
	return retList

}

//获取上传的文件列表
func (this *FacadeMysql) GetFileList() (retList []FileItem) {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("select filename,time from  filelist order by filename")
	checkErr(err)
	rows, err := stmt.Query()
	if err == nil {
		for rows.Next() {
			var stOneFileItem FileItem
			err = rows.Scan(&stOneFileItem.FileName, &stOneFileItem.Time)
			checkErr(err)
			retList = append(retList, stOneFileItem)

		}
	} else {
		utils.Debugln("stmt.Query error", err.Error())
	}
	utils.Debugln("GetFileList retList==", retList)
	return retList

}

func (this *FacadeMysql) DeleteUploadFile(filename string) error {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("delete  from  filelist where filename=?")
	checkErr(err)
	_, err = stmt.Exec(filename)
	if err == nil {
		return nil
	} else {
		return err
	}
}

func (this *FacadeMysql) InsertFile2List(filename string) error {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("insert into filelist values(?,?,?)")
	checkErr(err)
	t := time.Now()
	_, err = stmt.Exec(filename, t.Unix(), "")
	if err == nil {
		return nil
	} else {
		checkErr(err)
		return err
	}
}

func (this *FacadeMysql) UpdateNodeListRemarks(updateList []KeyValue) error {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	for i := 0; i < len(updateList); i++ {
		stmt, err := this.dbconn.Prepare("update nodelist set remarks=? where InnerIp=?")
		if err != nil {
			return err
		}

		res, err := stmt.Exec(updateList[i].Value, updateList[i].Key)
		affect, err := res.RowsAffected()
		if err != nil {
			return err
		}
		utils.Debugln("affect==", affect)
	}
	return nil
}

func (this *FacadeMysql) DelNodes(updateList []KeyValue) error {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	for i := 0; i < len(updateList); i++ {
		stmt, err := this.dbconn.Prepare("delete from nodelist  where InnerIp=?")
		if err != nil {
			checkErr(err)
			return err
		}

		res, err := stmt.Exec(updateList[i].Key)
		affect, err := res.RowsAffected()
		if err != nil {
			checkErr(err)
			return err
		}
		utils.Debugln("affect==", affect)
	}
	return nil
}

func (this *FacadeMysql) UpdateNodeList(key string, valuemap map[string]string) error {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("insert into nodelist VALUES(?,?,?,'') ")
	if err != nil {
		checkErr(err)
		return err
	}
	_, err = stmt.Exec(valuemap["Innerip"], valuemap["Outerip"], time.Now().Unix())
	if err != nil {
		checkErr(err)
		stmt, err = this.dbconn.Prepare("update  nodelist set lastUpdateTime=? where Innerip=?")
		checkErr(err)
		_, err = stmt.Exec(valuemap["LastUpdateTime"], valuemap["Innerip"])
		checkErr(err)
	}
	return nil
}

func (this *FacadeMysql) GetAllConfigTemplateList() (retList []ConfigTemplate) {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("select * from configtemplates")
	checkErr(err)
	rows, err := stmt.Query()
	if err == nil {
		for rows.Next() {
			var stOneTemplateItem ConfigTemplate
			err = rows.Scan(&stOneTemplateItem.TemplateName, &stOneTemplateItem.TemplateType, &stOneTemplateItem.EditTime, &stOneTemplateItem.Content, &stOneTemplateItem.Remarks)
			checkErr(err)
			retList = append(retList, stOneTemplateItem)

		}
	} else {
		utils.Debugln("stmt.Query error", err.Error())
	}
	utils.Debugln("GetFileList retList==", retList)
	return retList

}

func (this *FacadeMysql) ReplaceConfigTemplateList(stConfigTemplate ConfigTemplate) error {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("replace into configtemplates value(?,?,?,?,?)")
	checkErr(err)
	_, err = stmt.Exec(stConfigTemplate.TemplateName, stConfigTemplate.TemplateType, stConfigTemplate.EditTime, stConfigTemplate.Content, stConfigTemplate.Remarks)
	return err

}

func (this *FacadeMysql) DelConfigTemplate(strTemplateName string) error {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("delete from configtemplates where TemplateName=?")
	checkErr(err)
	_, err = stmt.Exec(strTemplateName)
	return err

}

func (this *FacadeMysql) GetRpcList() (retList []RpcTestItem) {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("select * from rpclist")
	checkErr(err)
	rows, err := stmt.Query()
	if err == nil {
		for rows.Next() {
			var stOneRpcItem RpcTestItem
			err = rows.Scan(&stOneRpcItem.Module, &stOneRpcItem.Object, &stOneRpcItem.Function, &stOneRpcItem.Data)
			checkErr(err)
			retList = append(retList, stOneRpcItem)

		}
	} else {
		utils.Debugln("stmt.Query error", err.Error())
	}
	utils.Debugln("GetRpcList retList==", retList)
	return retList
}

func (this *FacadeMysql) Save2RpcList(stRpcTestItem *RpcTestItem) error{
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("replace into  rpclist VALUES (?,?,?,?)")
	checkErr(err)
	_, err = stmt.Exec(stRpcTestItem.Module,stRpcTestItem.Object,stRpcTestItem.Function,stRpcTestItem.Data)
	if err == nil {
		return err
	} else {
		utils.Debugln("stmt.Query error", err.Error())
	}
	return nil
}

func (this *FacadeMysql) DelRpcList(stRpcTestItem *RpcTestItem) error{
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	stmt, err := this.dbconn.Prepare("delete from  rpclist where Module=?  and Object=? and Function=? and Data=? ")
	checkErr(err)
	_, err = stmt.Exec(stRpcTestItem.Module,stRpcTestItem.Object,stRpcTestItem.Function,stRpcTestItem.Data)
	if err == nil {
		return err
	} else {
		utils.Debugln("stmt.Query error", err.Error())
	}
	return nil
}

var G_FacadeMysql *FacadeMysql = nil

func GetFacadeMysql() *FacadeMysql {
	MysqlLocker.Lock()
	defer MysqlLocker.Unlock()
	if G_FacadeMysql == nil {
		G_FacadeMysql = new(FacadeMysql)
		//G_FacadeMysql.Init("kdb:kingsoft123@tcp(10.20.104.175:3306)/kwebtool")
		G_FacadeMysql.Init(G_MysqlConf)

	}
	return G_FacadeMysql
}
