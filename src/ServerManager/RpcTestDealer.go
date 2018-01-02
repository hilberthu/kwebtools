package ServerManager

import (
	"3rdparty/src/protorpc"
	"utils"
	//"net/http"
	"fmt"
	"net"
)

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import (
	"3rdparty/src/snappy"
	math "math"
	"time"
)

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type GameHeartbeatData struct {
	Ver              *int32  `protobuf:"varint,1,opt,name=ver" json:"ver,omitempty"`
	SessionID        *string `protobuf:"bytes,2,opt,name=sessionID" json:"sessionID,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GameHeartbeatData) Reset()         { *m = GameHeartbeatData{} }
func (m *GameHeartbeatData) String() string { return proto.CompactTextString(m) }
func (*GameHeartbeatData) ProtoMessage()    {}

func (m *GameHeartbeatData) GetVer() int32 {
	if m != nil && m.Ver != nil {
		return *m.Ver
	}
	return 0
}

func (m *GameHeartbeatData) GetSessionID() string {
	if m != nil && m.SessionID != nil {
		return *m.SessionID
	}
	return ""
}

type GameActionData struct {
	Module           *int32  `protobuf:"varint,1,opt,name=module" json:"module,omitempty"`
	FuncName         *string `protobuf:"bytes,2,opt,name=funcName" json:"funcName,omitempty"`
	Uid              *int64  `protobuf:"varint,3,opt,name=uid" json:"uid,omitempty"`
	SessionId        *string `protobuf:"bytes,4,opt,name=sessionId" json:"sessionId,omitempty"`
	Data             *string `protobuf:"bytes,5,opt,name=data" json:"data,omitempty"`
	ConnId           *uint32 `protobuf:"varint,6,opt,name=connId" json:"connId,omitempty"`
	IP               *uint64 `protobuf:"varint,7,opt" json:"IP,omitempty"`
	ObjName          *string `protobuf:"bytes,8,opt,name=objName" json:"objName,omitempty"`
	SerialID         *uint64 `protobuf:"varint,9,opt,name=serialID" json:"serialID,omitempty"`
	ServerID         *uint32 `protobuf:"varint,10,opt,name=serverID" json:"serverID,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GameActionData) Reset()         { *m = GameActionData{} }
func (m *GameActionData) String() string { return proto.CompactTextString(m) }
func (*GameActionData) ProtoMessage()    {}

func (m *GameActionData) GetModule() int32 {
	if m != nil && m.Module != nil {
		return *m.Module
	}
	return 0
}

func (m *GameActionData) GetFuncName() string {
	if m != nil && m.FuncName != nil {
		return *m.FuncName
	}
	return ""
}

func (m *GameActionData) GetUid() int64 {
	if m != nil && m.Uid != nil {
		return *m.Uid
	}
	return 0
}

func (m *GameActionData) GetSessionId() string {
	if m != nil && m.SessionId != nil {
		return *m.SessionId
	}
	return ""
}

func (m *GameActionData) GetData() string {
	if m != nil && m.Data != nil {
		return *m.Data
	}
	return ""
}

func (m *GameActionData) GetConnId() uint32 {
	if m != nil && m.ConnId != nil {
		return *m.ConnId
	}
	return 0
}

func (m *GameActionData) GetIP() uint64 {
	if m != nil && m.IP != nil {
		return *m.IP
	}
	return 0
}

func (m *GameActionData) GetObjName() string {
	if m != nil && m.ObjName != nil {
		return *m.ObjName
	}
	return ""
}

func (m *GameActionData) GetSerialID() uint64 {
	if m != nil && m.SerialID != nil {
		return *m.SerialID
	}
	return 0
}

func (m *GameActionData) GetServerID() uint32 {
	if m != nil && m.ServerID != nil {
		return *m.ServerID
	}
	return 0
}

func init() {
}

func DealWithRpcTest(strIpPort string, Module int, Object, Function, Data string, nServerId uint32, nUid int64) string {
	var responseAction GameActionData
	var sendAction GameActionData
	conn, err := net.Dial("tcp", strIpPort)
	if err != nil {
		fmt.Println("dial err", strIpPort)
		return "Dial strIpPort error"
	}
	defer conn.Close()
	//先握手
	//{"IP":"","Module":3,"Object":"","Function":"Handshake","UID":4001,"Session":"20eb0763-c9b0-4920-922d-43f6f9dd855f","Data":"","CallBack":"","SerialID":296598297,"ServerID":2001}
	sendAction.Reset()
	sendAction.FuncName = proto.String("Handshake")
	sendAction.ServerID = proto.Uint32(nServerId)
	sendAction.Uid = proto.Int64(nUid)
	sendAction.Module = proto.Int(3)
	sendAction.Data = proto.String("kingsoft123")
	sendAction.SessionId = proto.String("session")
	sendAction.ConnId = proto.Uint32(1000001)
	sendAction.SerialID = proto.Uint64(uint64(time.Now().Unix()))

	//先完成握手
	err = protorpc.WriteRequest_Out(conn, 10, "GameService.GameAction", &sendAction)
	if err != nil {
		return "handshack failed1!" + err.Error()
	}
	var readBuffer []byte = make([]byte, 10000000)
	readBuffer, err = protorpc.RecvFrame_Out(conn)
	readBuffer, err = protorpc.RecvFrame_Out(conn)
	pbRequest, err := snappy.Decode(nil, readBuffer)
	fmt.Println("RecvFrame_Out Decode", string(pbRequest), err)
	if err != nil {
		return "handshack failed1!"
	}
	fmt.Println("RecvFrame_Out Decode", string(readBuffer), err)
	err = proto.Unmarshal(pbRequest, &responseAction)
	fmt.Println("responseAction==", responseAction.String())

	//再发送协议
	sendAction.ObjName = proto.String(Object)
	sendAction.FuncName = proto.String(Function)
	sendAction.Module = proto.Int(Module)
	sendAction.Data = proto.String(Data)
	sendAction.SerialID = proto.Uint64(uint64(time.Now().Unix()) + 1)
	err = protorpc.WriteRequest_Out(conn, 10, "GameService.GameAction", &sendAction)
	if err != nil {
		return "excute failed1!" + err.Error()
	}
	readBuffer, err = protorpc.RecvFrame_Out(conn)
	if err != nil {
		return "excute failed1!" + err.Error()
	}
	readBuffer, err = protorpc.RecvFrame_Out(conn)
	if err != nil {
		return "excute failed1!" + err.Error()
	}
	pbRequest, err = snappy.Decode(nil, readBuffer)
	if err != nil {
		return "excute failed!-----" + err.Error()
	}
	err = proto.Unmarshal(pbRequest, &responseAction)

	utils.Debugln("DealWithRpcTest---------------", responseAction.String())
	//time.Sleep(time.Second*3)
	return *responseAction.Data
}
