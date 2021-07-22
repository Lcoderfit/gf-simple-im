package api

import (
	"fmt"
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
	"github.com/gogf/gf/util/gconv"
	"time"
)

// 聊天管理器
var Chat = &chatApi{}

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

// @summary 聊天室首页
// @description 聊天室首页，只显示模板内容。如果当前用户未登录，则引导跳转到名称设置页面
// @tags 聊天室
// @produce html
// @success 200 {string} string "执行结果" 
// @router /chat/index [GET]
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

// @summary 设置聊天名称页面
// @description 展示设置聊天名称页面，在该页面设置名称，成功后再跳转到聊天室页面。
// @tags 聊天室
// @produce html
// @success 200 {string} string "执行成功后跳转到聊天室页面"
// @router /chat/setname [GET]
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

// @summary WebSocket接口
// @description 通过WebSocket连接该接口发送任意数据
// @tags 聊天室
// @router /chat/websocket [POST]
func (a *chatApi) WebSocket(r *ghttp.Request) {
	// 从请求获取到数据后，先进行json解析将数据传给msg
	// 然后通过msg定义的v标签进行参数校验
	msg := &model.ChatMsg{}

	// 初始化WebSocket请求
	var (
		ws  *ghttp.WebSocket
		err error
	)
	// 将当前请求作为一个WebSocket请求
	ws, err = r.WebSocket()
	if err != nil {
		g.Log().Error(err)
	}

	name := r.Session.GetString("chat_name")
	if name == "" {
		// 直接发送请求过来的机器的IP地址, 假设客户端发送请求到Proxy1，Proxy1转发请求到Proxy2
		// Proxy2转发请求到服务器，则RemoteAddr是Proxy2的Ip
		name = r.Request.RemoteAddr
	}

	// 设置当前WebSocket对应用户昵称
	names.Add(name)
	users.Set(ws, name)

	// 向所有客户端发送上线消息
	a.writeUserListToClient()

	for {
		// 阻塞读取ws数据
		_, msgByte, err := ws.ReadMessage()
		if err != nil {
			// 如果读取失败，那么断开链接，这里清除用户信息
			names.Remove(name)
			users.Remove(ws)
			// -------省略失败重连机制，直接下线
			a.writeUserListToClient()
			break
		}
		// JSON参数解析
		if err := gjson.DecodeTo(msgByte, msg); err != nil {
			a.write(ws, model.ChatMsg{
				Type: "error",
				Data: "消息格式不正确: " + err.Error(),
				From: "",
			})
			continue
		}
		// 数据校验
		if err := g.Validator().Ctx(r.Context()).CheckStruct(msg); err != nil {
			a.write(ws, model.ChatMsg{
				Type: "error",
				Data: "消息格式不正确：" + err.Error(),
				From: "",
			})
			// 校验不通过就重试
			continue
		}
		// 发送发来自于name自己
		msg.From = name

		// 设置日志的策略??
		g.Log().Cat("chat").Println(msg)

		// ws操作类型
		switch msg.Type {
		case "send":
			intervalKey := fmt.Sprintf("%p", ws)
			// 三个参数分别为key, value和过期时间
			// 如果key存在，则返回false，否则返回true
			// key再一秒内存在，表示再一秒内点击发送了多次，提示发送过于频繁
			if ok, _ := cache.SetIfNotExist(intervalKey, struct{}{}, sendInterval); !ok {
				a.write(ws, model.ChatMsg{
					Type: "error",
					Data: "你的消息发送过于频繁，请休息下再重试",
					From: "",
				})
				continue
			}
			// 有消息时，群发消息
			if msg.Data != nil {
				if err = a.writeGroup(model.ChatMsg{
					Type: "send",
					// 对html特殊字符进行编码
					// gconv.String() 将接口类型转换为字符串
					Data: ghtml.SpecialChars(gconv.String(msg.Data)),
					From: ghtml.SpecialChars(msg.From),
				}); err != nil {
					g.Log().Error(err)
				}
			}
		}

	}

}

// 向单个用户发送消息
// 内部方法不会自动注册到路由中
func (a *chatApi) write(ws *ghttp.WebSocket, msg model.ChatMsg) error {
	b, err := gjson.Encode(msg)
	if err != nil {
		return err
	}
	return ws.WriteMessage(ghttp.WS_MSG_TEXT, b)
}

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
