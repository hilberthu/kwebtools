var g_curGroupName;
var g_action;
var g_nodelist = {};
var g_createTableData;
var g_selectedNode;
var g_spinner = null
var g_all_Templates = {};
var g_all_ConfigTemplate = {}
var g_all_ConfigTemplate_renderObj = {}
var g_all_ProcessInfo = {}
var g_rpclist= {}
var g_rpcid = 1
var g_ServerId = 1
var g_Uid = 20001
var g_IpPort = "10.20.104.50:20003"

function sleep(numberMillis) {
    var now = new Date();
    var exitTime = now.getTime() + numberMillis;
    while (true) {
        now = new Date();
        if (now.getTime() > exitTime)
            return;
    }
}


function onClickServerGroup(ServerGroupName) {
    var data = { "ServerGroupName": ServerGroupName };
    g_curGroupName = ServerGroupName;
    myajax("getOneGroupServerList", JSON.stringify(data));
    $("#publish_table").hide();


}

function myajax(action, data) {
    g_action = action;
    strUrl = "ajax?" + "action=" + action + "&data=" + data;
    $.get(strUrl, dealWithAjaxRetData);
}

function myajaxV2(action, data,afterActionName) {
    g_action = afterActionName;
    strUrl = "ajax?" + "action=" + action + "&data=" + data;
    $.get(strUrl, dealWithAjaxRetData);
}

function UnicodeToUtf8(unicode) {
    var uchar;
    var utf8str = "";
    var i;

    for (i = 0; i < unicode.length; i += 2) {
        uchar = (unicode[i] << 8) | unicode[i + 1]; //UNICODE为2字节编码，一次读入2个字节  
        utf8str = utf8str + String.fromCharCode(uchar); //使用String.fromCharCode强制转换  
    }
    return utf8str;
}

function onClickOpenUploadDlg() {
    $('#temp_template').html(g_all_Templates["fileuploadTemplate.html"])
    $("#upLoadDlg").modal()
}

function onClickCreateNew() {
    var data = { "ServerGroupName": "ServerGroupName" };
    myajax("getNodeList", JSON.stringify(data));
}
//strNode string, strPath string, strServerName string, strInsId string, nPort int32
function OnSubmitCreateNewServer() {
    var data = { "ServerGroupName": "ServerGroupName" };
    data.Node = $("#InputNode").val();
    data.ServerName = $("#InputServerName").val();
    data.InsId = $("#InputInsId").val();
    data.Path = $("#InputPath").val();
    data.Port = $("#InputPort").val();
    data.OutIp = $("#InputOuterIp").val();
    data.Other = $("#InputOther").val();
    //data.SererName =  $("#InputServerName").val();
    //alert(JSON.stringify(data))
    myajax("createNewServer", JSON.stringify(data));

}

function OnClickNode(strNode) {
    $("#InputNode").val(strNode);
}

/*
  strNode, _ = mapParam["Node"]
  strPath, _ = mapParam["Path"]
  strServerName, _ = mapParam["ServerName"]
  strInsId, _ = mapParam["InsId"]
*/
function onClickStop(innerIp, path, ServerName, InsId) {
    var data = {};
    data.Node = innerIp;
    data.Path = path;
    data.ServerName = ServerName
    data.InsId = InsId
    myajax("stopServer", JSON.stringify(data));
    ShowSpin()
}

function onClickStart(innerIp, path, ServerName, InsId) {
    var data = {};
    data.Node = innerIp;
    data.Path = path;
    data.ServerName = ServerName
    data.InsId = InsId
    myajax("startServer", JSON.stringify(data));
    ShowSpin()
}

function ShowSpin() {
    var opts = {
        lines: 13 // The number of lines to draw
            ,
        length: 28 // The length of each line
            ,
        width: 14 // The line thickness
            ,
        radius: 42 // The radius of the inner circle
            ,
        scale: 1 // Scales overall size of the spinner
            ,
        corners: 1 // Corner roundness (0..1)
            ,
        color: '#000' // #rgb or #rrggbb or array of colors
            ,
        opacity: 0.25 // Opacity of the lines
            ,
        rotate: 0 // The rotation offset
            ,
        direction: 1 // 1: clockwise, -1: counterclockwise
            ,
        speed: 1 // Rounds per second
            ,
        trail: 60 // Afterglow percentage
            ,
        fps: 20 // Frames per second when using setTimeout() as a fallback for CSS
            ,
        zIndex: 2e9 // The z-index (defaults to 2000000000)
            ,
        className: 'spinner' // The CSS class to assign to the spinner
            ,
        top: '50%' // Top position relative to parent
            ,
        left: '50%' // Left position relative to parent
            ,
        shadow: false // Whether to render a shadow
            ,
        hwaccel: false // Whether to use hardware acceleration
            ,
        position: 'absolute' // Element positioning
    }
    if (g_spinner == null) {
        var target = document.getElementById('right_table_serverlist')
        g_spinner = new Spinner(opts).spin(target);
    }

}

function HideSpin() {
    g_spinner.spin()
    g_spinner = null
}


function onClickShowServerList() {
    var data = { "ServerGroupName": g_curGroupName };
    myajax("getOneGroupServerList", JSON.stringify(data));
    $("#rightTable").show();
    $("#publish_table").hide();

}

function onClickShowPubList() {
    var data = { "ServerGroupName": g_curGroupName };
    $("#rightTable").hide();
    myajax("getPublishServerList", JSON.stringify(data));

}

function onClickGetPakages() {

}

//点击节点列表
function onClickNodeList() {
    myajax("getNodeList", "null");
    //GetAllCheckedNodes()
}

//点击文件管理
function onClickFileManager() {
    myajax("getFileList", "null");
    //GetAllCheckedNodes()
}

function onClickConfigTemplateManager() {
    myajax("getConfigTemplateList", "null");

    //GetAllCheckedNodes()
}


//保存所有备注信息
function SaveRemaks() {
    var tbModify = []
    $('input[name="nodecheck"]:checked').each(function() { //遍历每一个名字为nodes的复选框，其中选中的执行函数      
        var oneItem = {}
        var spliteResult = $(this).val().split("_")
        oneItem.key = spliteResult[1]
        var id = "nodeRemarks_" + oneItem.key
        oneItem.value = $("input[ name='" + id + "']").val()
        tbModify.push(oneItem)
    })
    myajax("updateNodeRemarks", JSON.stringify(tbModify));
    //http://blog.csdn.net/paincupid/article/details/50923271
    /*for (;iterator!=$("#nodemanager_list_tbody").lastChild;iterator = iterator.next())
    {
      console.log(iterator.firstChild.firstChild.firstChild.html());
    }*/
}

function DelNodes() {
    var bDelete = confirm("删除是不可恢复的，你确认要删除吗？");
    if (bDelete == false) {
        return
    }
    var tbModify = []
    $('input[name="nodecheck"]:checked').each(function() { //遍历每一个名字为nodes的复选框，其中选中的执行函数      
        var oneItem = {}
        var spliteResult = $(this).val().split("_")
        oneItem.key = spliteResult[1]
        var id = "nodeRemarks_" + oneItem.key
        oneItem.value = $("input[ name='" + id + "']").val()
        tbModify.push(oneItem)
    })
    myajax("DelNodes", JSON.stringify(tbModify));
    //http://blog.csdn.net/paincupid/article/details/50923271
    /*for (;iterator!=$("#nodemanager_list_tbody").lastChild;iterator = iterator.next())
    {
      console.log(iterator.firstChild.firstChild.firstChild.html());
    }*/
}

function PushFiles() {

}

formatTime = function(time) {
    var unixTimestamp = new Date(time * 1000);
    return unixTimestamp.toLocaleString();
}

function deleteUploadFile(filename) {
    var bDelete = confirm("删除是不可恢复的，你确认要删除吗？");
    if (bDelete == false) {
        return
    }
    myajax("deleteUploadFile", filename);
}


function updateConfigEditWind(configTemplateNamme) {
    var oneConfigItem = g_all_ConfigTemplate[configTemplateNamme]
    $("#ConfigEditWind_TemplateName").val(configTemplateNamme)
    $("#ConfigEditWind_Remarks").val(oneConfigItem.Remarks)
    if (oneConfigItem.TemplateType == "进程模版") {
        $("input[name='optionsRadiosinline']").eq(1).attr("checked", "checked");
        $("input[name='optionsRadiosinline']").eq(0).removeAttr("checked");
        $("input[name='optionsRadiosinline']").eq(1).click();
    } else {
        $("input[name='optionsRadiosinline']").eq(0).attr("checked", "checked");
        $("input[name='optionsRadiosinline']").eq(1).removeAttr("checked");
        $("input[name='optionsRadiosinline']").eq(0).click();
    }
    $("#ConfigEditWind_content").val(oneConfigItem.Content)

}

//检查并保存模版
function CompileAndSaveConfigTemplate() {
    var data = {}
    data.TemplateName = $("#ConfigEditWind_TemplateName").val()
    if ($("input[name='optionsRadiosinline']").eq(1).prop("checked") == true) {
        data.TemplateType = 2
    } else {
        data.TemplateType = 1
    }
    data.Content = $("#ConfigEditWind_content").val()
    data.Remarks = $("#ConfigEditWind_Remarks").val()
    myajax("SaveConfigTemplate", JSON.stringify(data))

}


function CompileConfigTemplate(ElememntId) {

    var Content = $(ElememntId).val()
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
    alert(strResult)

}

function delConfigEditWind(strTemplateKey) {
    myajax("DelConfigTemplate", strTemplateKey)
    sleep(1000)
    myajax("getConfigTemplateList", "")
}

//点击文件推送，初始化模态框
function onClickPushFile() {
    console.log("alian")
}

function onClickProcessList() {
    /*var ProcessList = {}
    ProcessList.processList = []
    var oneItem = {}
    oneItem.Innerip = "127.0.0.1"
    oneItem.Outerip = "127.0.0.1"
    oneItem.ServerName = "127.0.0.1"
    oneItem.Insid = "127.0.0.1"
    oneItem.Path = "127.0.0.1"
    oneItem.Port = "127.0.0.1"
    oneItem.ConfigContent = "127.0.0.1"
    oneItem.Status = "ok"
    oneItem.Lastupdatetime = "10000"
    oneItem.Other = "other"
    ProcessList.processList.push(oneItem)
    $('#temp_template').html(g_all_Templates["ProcessListTemplate.html"])

    var source = $('#processList-template').html();
    var myTemplate = Handlebars.compile(source);
    var out = myTemplate(ProcessList)
    //TemplateName TemplateType LastEditTime
    $('#rightlist').html(out);
    //$('#rightlist').append(g_all_Templates["ConfigTemplateEdit.html"])
    $('#temp_template').html("")*/

    //myajax("getNodeList", "")
    myajaxV2("getNodeList","","getNodeList2")
    //myajax("getProcessList", "")
    //myajaxV2("getNodeList","","getNodeList2")
    //myajax("getProcessList", "")


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


function CompileAndSaveProcessConfigContent() {
    var data = {}
    var options = $("#ProcessEdit_ServerNameList option:selected");
    data.ServerName = options.val();
    options = $("#ProcessEdit_InnerIpList option:selected");
    data.Innerip = options.val();
    data.Insid = $("#ProcessEdit_InsId").val();
    data.ConfigContent = $("#ProcessConfigFileEdit_content").val();
    myajax("updateProcess", JSON.stringify(data))
}

function CompileAndAddProcessConfigContent() {
    var data = {}
    var options = $("#ProcessEdit_ServerNameList option:selected");
    data.ServerName = options.val();
    options = $("#ProcessEdit_InnerIpList option:selected");
    data.Innerip = options.val();
    data.Insid = $("#ProcessEdit_InsId").val();
    data.ConfigContent = $("#ProcessConfigFileEdit_content").val();
    myajax("mkNewProcess", JSON.stringify(data))
}

function onEditProcess(ServerName,Insid,Innerip){

    strSelect = "option[value!='"+ServerName+"']"
    $(strSelect).attr("selected",false)
    strSelect = "option[value!='"+Innerip+"']"
    $(strSelect).attr("selected",false)

    $("#ProcessEdit_InsId").val(Insid)
    var strSelect = "option[value='"+ServerName+"']"
    $(strSelect).attr("selected",true)

    strSelect = "option[value='"+Innerip+"']"
    $(strSelect).attr("selected",true)
    var strContent = g_all_ProcessInfo[ServerName][Insid].ConfigContent
    var tempObj
    try{
        tempObj = JSON.parse(strContent)
    }
    catch (err){
        $("#ProcessConfigFileEdit_content").val(strContent)
        return
    }
    if (tempObj != null){
        strContent = JSON.stringify(tempObj,null,4)
        $("#ProcessConfigFileEdit_content").val(strContent)
    }
}

function onDeleteProcess(ServerName,Insid) {
    var bDelete = confirm("删除是不可恢复的，你确认要删除吗？");
    if (bDelete == false) {
        return
    }
    var data = {}
    data.ServerName = ServerName
    data.Insid = Insid
    myajax("deleteprocess",JSON.stringify(data))
}

function onClickRPCTest(){
    myajax("getrpclist","")

}

function OnClickNewRpc(){

    $("#rpctest_list_tbody").append("<tr id=\"rpcItem_"+parseInt(g_rpcid)+"\">\n" +
        "    <td><input type=\"text\" id=\"Module_"+parseInt(g_rpcid) +"\" class=\"form-control\" ></td>\n" +
        "    <td><input type=\"text\" id=\"Object_"+parseInt(g_rpcid) +"\" class=\"form-control\" ></td>\n" +
        "    <td><input type=\"text\" id=\"Function_"+parseInt(g_rpcid) +"\" class=\"form-control\" ></td>\n" +
        "    <td><input type=\"text\" id=\"Data_"+parseInt(g_rpcid) + "\" class=\"form-control\" ></td>\n" +
        "    <td>\n" +
        "    <button type=\"button\" class=\"btn btn-primary\" href=\"javascript:;\" onclick=\"OnClickExcuteRpc("+parseInt(g_rpcid)+")\"> 执行</button>\n" +
        "    <button type=\"button\" class=\"btn btn-primary\" href=\"javascript:;\" onclick=\"OnClickSaveRpc("+parseInt(g_rpcid)+")\"> 保存</button>\n" +
        "    <button type=\"button\" class=\"btn btn-primary\" href=\"javascript:;\" onclick=\"OnClickDelRpc("+parseInt(g_rpcid)+")\"> 删除</button>\n" +
        "    </td>\n" +
        "</tr>");
    g_rpcid =  g_rpcid +1

}

function OnClickExcuteRpc(rpcId){
    var strModule = "#Module_"+parseInt(rpcId)
    var strObject = "#Object_"+parseInt(rpcId)
    var strFunction= "#Function_"+parseInt(rpcId)
    var strData= "#Data_"+parseInt(rpcId)
    var rpcItem = {}
    rpcItem.Module = $(strModule).val()
    rpcItem.Object = $(strObject).val()
    rpcItem.Function = $(strFunction).val()
    rpcItem.Data= $(strData).val()
    rpcItem.ServerId =  $("#Rpc_ServerId").val()
    rpcItem.Uid = $("#Rpc_Uid").val()
    rpcItem.IpPort = $("#Rpc_IpPort").val()
    var bOk = checkParamEmpty(rpcItem)
    if (!bOk){
        return
    }
    myajax("slgrpctest",JSON.stringify(rpcItem))

}
function OnClickSaveRpc(rpcId) {
    var strModule = "#Module_"+parseInt(rpcId)
    var strObject = "#Object_"+parseInt(rpcId)
    var strFunction= "#Function_"+parseInt(rpcId)
    var strData= "#Data_"+parseInt(rpcId)
    var rpcItem = {}
    rpcItem.Module = $(strModule).val()
    rpcItem.Object = $(strObject).val()
    rpcItem.Function = $(strFunction).val()
    rpcItem.Data= $(strData).val()
    var bOk = checkParamEmpty(rpcItem)
    if (!bOk){
        return
    }
    myajax("saverpclist",JSON.stringify(rpcItem))
}

function checkParamEmpty(rpcItem) {
    for (key in rpcItem){
        if(rpcItem[key] == ""){
            alert(key + " 不能为空!")
            return false
        }
    }
    return true
}

var g_deletedItemId = ""
function OnClickDelRpc(rpcId){
    var strModule = "#Module_"+parseInt(rpcId)
    var strObject = "#Object_"+parseInt(rpcId)
    var strFunction= "#Function_"+parseInt(rpcId)
    var strData= "#Data_"+parseInt(rpcId)
    var rpcItem = {}
    rpcItem.Module = $(strModule).val()
    rpcItem.Object = $(strObject).val()
    rpcItem.Function = $(strFunction).val()
    rpcItem.Data= $(strData).val()
    var bOk = checkParamEmpty(rpcItem)
    if (!bOk){
        return
    }
    myajax("delrpclist",JSON.stringify(rpcItem))
    g_deletedItemId = "#rpcItem_"+parseInt(rpcId)

}


//需要拉取的模版列表
var Templatefilelist = ["nodelistTemplate.html", "ProcessListTemplate.html", "ServerManagerPage.html", "fileManage.html", "fileuploadTemplate.html", "configTemplate.html", "ConfigTemplateEdit.html", "ProcessEditTemplate.html","RpcTest.html"]
var g_ServerNameList = ["kGameLogicServer", "kGameGatewayServer", "kMessageCenterServer", "kChatServer", "kServerManager", "kRankServer", "kArenaServer","kLoginServer"]

myajax("getTemplates", JSON.stringify(Templatefilelist))
//sleep(1000)
//myajax("getNodeList","")
//g_action = "getNodeList2"
//sleep(1000)



