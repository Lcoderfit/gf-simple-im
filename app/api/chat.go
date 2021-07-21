package api

import (
	"gf-simple-im/app/model"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/encoding/ghtml"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"time"
)

// 聊天管理器
const Chat = &chatApi{}

type chatApi struct{}

const (
	// 允许客户端发送聊天消息的间隔时间
	sendInterval = time.Second
)

var (
	users = gmap.New(true)       // 使用默认的并发安全的map (true表示是否为并发安全的)
	names = gset.NewStrSet(true) // 使用并发安全的set，用于用户唯一性校验 (true表示是否为并发安全的)
	cache = gcache.New()         // 使用特定的缓存对象，不适用全局缓存对象
)

// 首页
func (a *chatApi) Index(r *ghttp.Request) {
	view := r.GetView()
	//
	if r.Session.Contains("chat_name") {
		view.Assign("tplMain", "chat/include/chat.html")
	} else {
		view.Assign("tplMain", "chat/include/main.html")
	}
	r.Response.WriteTpl("chat/index.html")
}

// 设置
func (a *chatApi) SetName(r *ghttp.Request) {
	var (
		apiReq *model.ChatApiSetNameReq
	)
	// 只对表单字段或正文内容进行解析
	if err := r.ParseForm(&apiReq); err != nil {
		r.Session.Set("chat_name_error", gerror.Current(err).Error())
		r.Response.RedirectBack()
	}
	// 对所有html字符内容进行编码
	name := ghtml.Entities(apiReq.Name)
	r.Session.Set("chat_name_temp", name)
	if names.Contains(name) {
		r.Session.Set("chat_name_error", "用户昵称已被占用")
		// 响应重定向, 将客户端重定向到引用页，参数可以指定重定向的状态码
		r.Response.RedirectBack()
	} else {
		r.Session.Set("chat_name", name)
		r.Session.Remove("chat_name_temp", "chat_name_error")
		r.Response.RedirectTo("/chat")
	}
}

// WebSocket接口
func (a *chatApi) WebSocket(r *ghttp.Request) {
	msg := &model.ChatMsg{}
	// 初始化WebSocket请求
	var (
		ws  *ghttp.WebSocket
		err error
	)
	ws, err = r.WebSocket()
	if err != nil {
		g.Log().Error(err)
		return
	}

	name := r.Session.GetString("chat_name")
	if name == "" {
		name = r.Request.RemoteAddr
	}
}

//
//// @summary WebSocket接口
//// @description 通过WebSocket连接该接口发送任意数据。
//// @tags    聊天室
//// @router  /chat/websocket [POST]
//func (a *chatApi) WebSocket(r *ghttp.Request) {
//	msg := &model.ChatMsg{}
//
//	// 初始化WebSocket请求
//	var (
//		ws  *ghttp.WebSocket
//		err error
//	)
//	ws, err = r.WebSocket()
//	if err != nil {
//		g.Log().Error(err)
//		return
//	}
//
//	name := r.Session.GetString("chat_name")
//	if name == "" {
//		name = r.Request.RemoteAddr
//	}
//
//	// 初始化时设置用户昵称为当前链接信息
//	names.Add(name)
//	users.Set(ws, name)
//
//	// 初始化后向所有客户端发送上线消息
//	a.writeUserListToClient()
//
//	for {
//		// 阻塞读取WS数据
//		_, msgByte, err := ws.ReadMessage()
//		if err != nil {
//			// 如果失败，那么表示断开，这里清除用户信息
//			// 为简化演示，这里不实现失败重连机制
//			names.Remove(name)
//			users.Remove(ws)
//			// 通知所有客户端当前用户已下线
//			a.writeUserListToClient()
//			break
//		}
//		// JSON参数解析
//		if err := gjson.DecodeTo(msgByte, msg); err != nil {
//			a.write(ws, model.ChatMsg{
//				Type: "error",
//				Data: "消息格式不正确: " + err.Error(),
//				From: "",
//			})
//			continue
//		}
//		// 数据校验
//		if err := g.Validator().Ctx(r.Context()).CheckStruct(msg); err != nil {
//			a.write(ws, model.ChatMsg{
//				Type: "error",
//				Data: gerror.Current(err).Error(),
//				From: "",
//			})
//			continue
//		}
//		msg.From = name
//
//		// 日志记录
//		g.Log().Cat("chat").Println(msg)
//
//		// WS操作类型
//		switch msg.Type {
//		// 发送消息
//		case "send":
//			// 发送间隔检查
//			intervalKey := fmt.Sprintf("%p", ws)
//			if ok, _ := cache.SetIfNotExist(intervalKey, struct{}{}, sendInterval); !ok {
//				a.write(ws, model.ChatMsg{
//					Type: "error",
//					Data: "您的消息发送得过于频繁，请休息下再重试",
//					From: "",
//				})
//				continue
//			}
//			// 有消息时，群发消息
//			if msg.Data != nil {
//				if err = a.writeGroup(
//					model.ChatMsg{
//						Type: "send",
//						Data: ghtml.SpecialChars(gconv.String(msg.Data)),
//						From: ghtml.SpecialChars(msg.From),
//					}); err != nil {
//					g.Log().Error(err)
//				}
//			}
//		}
//	}
//}
//

// 向单个用户发送消息
// 内部方法不会自动注册到路由中
func (a *chatApi) write(ws *ghttp.WebSocket, msg model.ChatMsg) error {
	b, err := gjson.Encode(msg)
	if err != nil {
		return err
	}
	return ws.WriteMessage(ghttp.WS_MSG_TEXT, b)
}

//// 向客户端写入消息。
//// 内部方法不会自动注册到路由中。
//func (a *chatApi) write(ws *ghttp.WebSocket, msg model.ChatMsg) error {
//	msgBytes, err := gjson.Encode(msg)
//	if err != nil {
//		return err
//	}
//	return ws.WriteMessage(ghttp.WS_MSG_TEXT, msgBytes)
//}

// 向所有客户端发送消息
func (a *chatApi) writeGroup(msg model.ChatMsg) error {
	b, err := gjson.Encode(msg)
	if err != nil {
		return err
	}
	// 设置读锁，保证并发读的安全
	users.RLockFunc(func(m map[interface{}]interface{}) {
		for user := range m {
			// 每个登录的用户都会建立一个websocket连接，对每个建立websocket连接的用户发送消息
			// WS_MSG_TEXT表示数据会被编码成utf8个是的文本消息
			user.(*ghttp.WebSocket).WriteMessage(ghttp.WS_MSG_TEXT, b)
		}
	})
	return nil
}

// 向客户端返回用户列表
// 内部方法不会自动注册到路由中
func (a *chatApi) writeUserListToClient() error {
	// 设置一个有序列表
	array := garray.NewSortedStrArray()
	// 遍历set中的元素，如果函数f返回true则继续迭代，否则停止迭代
	names.Iterator(func(v string) bool {
		// 将所有用户名都放入有序列表中
		array.Add(v)
		return true
	})

	// 将用户列表信息发送给所有用户
	if err := a.writeGroup(model.ChatMsg{
		Type: "list",
		Data: array.Slice(),
		From: "",
	}); err != nil {
		return err
	}
	return nil
}
