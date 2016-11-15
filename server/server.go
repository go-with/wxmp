package server

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// 消息处理器
type Handler func(*Context)

// 消息服务端
type Server struct {
	Store          Store
	DefaultHandler Handler
	MsgHandlers    map[string]Handler
	EvtHandlers    map[string]Handler
}

// 实例化消息服务端
func New(store Store) *Server {
	return &Server{
		Store:       store,
		MsgHandlers: make(map[string]Handler),
		EvtHandlers: make(map[string]Handler),
	}
}

// 注册普通消息处理器
func (srv *Server) OnEvt(typ string, h Handler) {
	srv.EvtHandlers[typ] = h
}

// 注册事件推送处理器
func (srv *Server) OnMsg(typ string, h Handler) {
	srv.MsgHandlers[typ] = h
}

func (srv *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := srv.createContext(res, req)

	srv.parseReqMsg(ctx)

	valid, err := srv.checkSignature(ctx)
	if err != nil {
		ctx.respond(500, mimePlain, []byte(err.Error()))
		return
	}
	if !valid {
		ctx.respond(403, mimePlain, []byte("signature check failed"))
		return
	}
	if ctx.Req.Method != http.MethodPost {
		echostr := ctx.Req.FormValue("echostr")
		ctx.respond(200, mimePlain, []byte(echostr))
		return
	}
	if ctx.SafeMode() {
		err = srv.decryptReqMsg(ctx)
		if err != nil {
			ctx.respond(500, mimePlain, []byte(err.Error()))
			return
		}
	}
	h := srv.matchHandler(ctx)
	if h != nil {
		h(ctx)
	}
	ctx.noReply()
}

// 匹配消息处理器
func (srv *Server) matchHandler(ctx *Context) (h Handler) {
	var ok bool
	if ctx.ReqMsg.MsgType != "event" {
		h, ok = srv.MsgHandlers[ctx.ReqMsg.MsgType]
	} else {
		h, ok = srv.EvtHandlers[ctx.ReqMsg.Event]
	}
	if !ok && srv.DefaultHandler != nil {
		h = srv.DefaultHandler
	}
	return
}

// 解密请求消息
func (srv *Server) decryptReqMsg(ctx *Context) (err error) {
	encodingAESKey, err := srv.Store.GetEncodingAESKey(ctx.AppID())
	if err != nil {
		return
	}
	ctx.aesKey = decodeAESKey(encodingAESKey)
	data, err := aesDecrypt(base64Decode(ctx.ReqMsg.Encrypt), ctx.aesKey)
	if err != nil {
		return
	}
	var msgLen int32
	err = binary.Read(bytes.NewReader(data[16:20]), binary.BigEndian, &msgLen)
	if err != nil {
		return
	}
	err = xml.Unmarshal(data[20:20+msgLen], ctx.ReqMsg)
	return
}

// 解析请求消息
func (srv *Server) parseReqMsg(ctx *Context) (err error) {
	defer ctx.Req.Body.Close()
	data, err := ioutil.ReadAll(ctx.Req.Body)
	if err != nil {
		return
	}
	err = xml.Unmarshal(data, ctx.ReqMsg)
	return
}

// 验证消息签名
func (srv *Server) checkSignature(ctx *Context) (valid bool, err error) {
	ctx.token, err = srv.Store.GetToken(ctx.AppID())
	if err != nil {
		return
	}
	nonce := ctx.Req.FormValue("nonce")
	timestamp := ctx.Req.FormValue("timestamp")
	if ctx.SafeMode() {
		msgSignature := signMsg(ctx.token, timestamp, nonce, ctx.ReqMsg.Encrypt)
		valid = ctx.Req.FormValue("msg_signature") == msgSignature
		return
	}
	signature := signMsg(ctx.token, timestamp, nonce)
	valid = ctx.Req.FormValue("signature") == signature
	return
}

// 创建上下文
func (srv *Server) createContext(res http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		writer: res,
		Req:    req,
		ReqMsg: new(ReqMsg),
	}
}
