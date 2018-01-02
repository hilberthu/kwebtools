dealWithAjaxRetData = function(data, status) {
    if ("getOneGroupServerList" == g_action) {
        $("#serverlist_tbody").empty();
        $("#serverlist_tbody").append(data);
        $("#rightTable").show();
    }

    if ("getNodeList" == g_action) {
        var NodeList = JSON.parse(data)
        g_nodelist = NodeList
        $('#temp_template').html(g_all_Templates["nodelistTemplate.html"])
        var source = document.getElementById("NodeList-template").innerHTML;
        var myTemplate = Handlebars.compile(source);
        if (NodeList != null && NodeList.length != null) {
            for (var i = NodeList.length - 1; i >= 0; i--) {
                NodeList[i].StatusString = formatTime(NodeList[i].LastUpdateTime)
                if (NodeList[i].LastUpdateTime > Date.parse(new Date()) / 1000 - 300) {
                    NodeList[i].alertstyle = "btn-success"
                } else {
                    NodeList[i].alertstyle = "btn-danger"
                }
                //NodeList[i].StatusString = "ok"

            }
            g_nodelist = NodeList
        }

        var obj1 = {}
        obj1.NodeList = NodeList
        //将json对象用刚刚注册的Handlebars模版封装，得到最终的html，插入到基础table中。
        var out = myTemplate(obj1)
        $('#rightlist').html(out);
        $('#temp_template').html("")
    }

    if ("createNewServer" == g_action) {
        alert(data)
    }

    if ("stopServer" == g_action) {
        alert(data)
        HideSpin();
    }

    if ("startServer" == g_action) {
        alert(data)
        HideSpin();
    }

    if ("getPublishServerList" == g_action) {
        $("#publish_serverlist_tbody").empty();
        $("#publish_serverlist_tbody").append(data);
        $("#publish_table").show();

    }

    //getTemplates
    if ("getTemplates" == g_action) {
        g_all_Templates = JSON.parse(data)
        g_all_Templates = JSON.parse(g_all_Templates.data)
        myajax("getConfigTemplateList","")
    }

    if ("updateNodeRemarks" == g_action) {
        var tbResult = JSON.parse(data);
        if (tbResult.ret == 0) {
            alert("success")
        } else {
            alert("failed! " + tbResult.data)
        }
    }
    //getFileList
    if ("getFileList" == g_action) {
        //console.log(g_all_Templates["fileManage.html"])
        $('#temp_template').html(g_all_Templates["fileManage.html"])
        var source = $('#FileList-template').html();
        var myTemplate = Handlebars.compile(source);
        var retObj = JSON.parse(data)
        var renderObj = {}
        renderObj.FileList = JSON.parse(retObj.data)
        if (renderObj.FileList != null) {
            for (var i = 0; i < renderObj.FileList.length; i++) {
                renderObj.FileList[i].Time = formatTime(renderObj.FileList[i].Time)
            }
        }

        //将json对象用刚刚注册的Handlebars模版封装，得到最终的html，插入到基础table中。
        var out = myTemplate(renderObj)
        $('#rightlist').html(out);
        $('#temp_template').html("")
    }

    if ("deleteUploadFile" == g_action) {
        myajax("getFileList", "null");
    }

    if ("getConfigTemplateList" == g_action) {
        $('#temp_template').html(g_all_Templates["configTemplate.html"])
        var source = $('#ConfigTemplateList-template').html();
        var myTemplate = Handlebars.compile(source);
        var retObj = JSON.parse(data)
        var renderObj = {}
        renderObj.ConfigTemplateList = JSON.parse(retObj.data)
        for (var i = renderObj.ConfigTemplateList.length - 1; i >= 0; i--) {
            if (renderObj.ConfigTemplateList[i].TemplateType == 1) {
                renderObj.ConfigTemplateList[i].TemplateType = "小模版"
            } else {
                renderObj.ConfigTemplateList[i].TemplateType = "进程模版"
            }
            renderObj.ConfigTemplateList[i].EditTime = formatTime(renderObj.ConfigTemplateList[i].EditTime)
            g_all_ConfigTemplate[renderObj.ConfigTemplateList[i].TemplateName] = renderObj.ConfigTemplateList[i]
        }
        var out = myTemplate(renderObj)
        g_all_ConfigTemplate_renderObj = renderObj
        console.log("out==", out)
        //TemplateName TemplateType LastEditTime
        $('#rightlist').html(out);
        $('#rightlist').append(g_all_Templates["ConfigTemplateEdit.html"])
        $('#table_allConfigTemplate').attr('id', 'table_allConfigTemplate_new')
        $('#temp_template').html("")
    }

    if ("SaveConfigTemplate" == g_action) {
        var retObj = JSON.parse(data)
        if (retObj.ret == 0) {
            $('#temp_template').html(g_all_Templates["configTemplate.html"])
            var source = $('#ConfigTemplateList-template').html();
            var myTemplate = Handlebars.compile(source);
            var retObj = JSON.parse(data)
            var renderObj = {}
            renderObj.ConfigTemplateList = JSON.parse(retObj.data)
            for (var i = renderObj.ConfigTemplateList.length - 1; i >= 0; i--) {
                if (renderObj.ConfigTemplateList[i].TemplateType == 1) {
                    renderObj.ConfigTemplateList[i].TemplateType = "小模版"
                } else {
                    renderObj.ConfigTemplateList[i].TemplateType = "进程模版"
                }
                renderObj.ConfigTemplateList[i].EditTime = formatTime(renderObj.ConfigTemplateList[i].EditTime)
                g_all_ConfigTemplate[renderObj.ConfigTemplateList[i].TemplateName] = renderObj.ConfigTemplateList[i]
            }
            var out = myTemplate(renderObj)
            $('#temp_template').html(out)
            $('#table_allConfigTemplate_new').html($('#table_allConfigTemplate').html())
            $('#temp_template').html("")
        }

    }

    if ("DelNodes" == g_action) {
        myajax("getNodeList", "")
    }

    if ("getProcessList" == g_action) {
        var dataObj = JSON.parse(data)
        var ProcessList = {}
        ProcessList.processList = JSON.parse(dataObj.data)
        if( ProcessList.processList==null){
            ProcessList.processList = {}
        }
        for (var i = ProcessList.processList.length - 1; i >= 0; i--) {
            //ProcessList.processList[i].Lastupdatetime = "adfadfadfadfadf"
            if (g_all_ProcessInfo[ProcessList.processList[i].ServerName]==null){
                g_all_ProcessInfo[ProcessList.processList[i].ServerName] = {}
            }
            if (parseInt(ProcessList.processList[i].Lastupdatetime) > Date.parse(new Date()) / 1000 - 300) {
                ProcessList.processList[i].alertstyle = "btn-success"
            } else {
                ProcessList.processList[i].alertstyle = "btn-danger"
            }
            ProcessList.processList[i].Lastupdatetime = formatTime(parseInt(ProcessList.processList[i].Lastupdatetime))
            ProcessList.processList[i].Path = decodeURIComponent( ProcessList.processList[i].Path)
            g_all_ProcessInfo[ProcessList.processList[i].ServerName][ProcessList.processList[i].Insid] = ProcessList.processList[i]
        }
        $('#temp_template').html(g_all_Templates["ProcessListTemplate.html"])

        var source = $('#processList-template').html();
        var myTemplate = Handlebars.compile(source);
        var out = myTemplate(ProcessList)
        //TemplateName TemplateType LastEditTime
        $('#rightlist').html(out);
        //$('#rightlist').append(g_all_Templates["ConfigTemplateEdit.html"])
        $('#temp_template').html("")

        $('#temp_template').html(g_all_Templates["ProcessEditTemplate.html"])
        source = $('#ProcessEdit-template').html();
        myTemplate = Handlebars.compile(source);
        ProcessList = {}
        ProcessList.ServerNameList = g_ServerNameList
        ProcessList.InnerIpList = g_nodelist
        out = myTemplate(ProcessList)

        $('#rightlist').append(out)
        console.log("alian")
    }
    if("deleteprocess"==g_action){
        onClickProcessList()
    }
    if("getNodeList2" == g_action){
        var NodeList = JSON.parse(data)
        g_nodelist = NodeList
        myajax("getProcessList", "")
    }

    if("getrpclist" == g_action){
        g_rpcid = 1
        g_rpclist = JSON.parse(JSON.parse(data)["data"])
        for (var i = g_rpclist.length - 1; i >= 0; i--){
            g_rpclist[i].rpcid =g_rpcid
            g_rpcid = g_rpcid + 1
        }
        $('#temp_template').html(g_all_Templates["RpcTest.html"])

        var source = $('#RpcTest-template').html();
        var myTemplate = Handlebars.compile(source);
        var RpcList = {}
        RpcList.RpcList = g_rpclist
        RpcList.ServerId = g_ServerId
        RpcList.IpPort = g_IpPort
        RpcList.Uid = g_Uid
        var out = myTemplate(RpcList)
        //TemplateName TemplateType LastEditTime
        $('#rightlist').html(out);
        //$('#rightlist').append(g_all_Templates["ConfigTemplateEdit.html"])
        $('#temp_template').html("")
    }

    if("saverpclist" == g_action){
        var retObj = JSON.parse(data)
        if (retObj["ret"] == 0){
            alert("sucess")
        }else{
            alert("error = "+retObj["data"])
        }
    }

    if("delrpclist" == g_action){
        var retObj = JSON.parse(data)
        if (retObj["ret"] == 0){
            alert("sucess")
            $(g_deletedItemId).remove()
        }else{
            alert("error = "+retObj["data"])
        }
    }

    if("slgrpctest" == g_action){
        var retObj = JSON.parse(data)

        if (retObj["ret"] == 0){
            var obj2 = JSON.parse(retObj.data)
            JSON.stringify(obj2, null, "\t");
            $("#RpcTest_RetContent").val(JSON.stringify(obj2, null, "\t"))
        }else{
            alert("error = "+retObj["data"])
        }

    }
    //alert("Data: " + data + "\nStatus: " + status);
}