# 微信公众平台SDK for Go

[![GoDoc](https://godoc.org/github.com/go-with/wxmp?status.svg)](https://godoc.org/github.com/go-with/wxmp)

这是一个使用Go语言编写的微信公众平台SDK

## 举个栗子

被动回复用户消息

```Go
package main

import (
	"log"
	"net/http"

	"github.com/go-with/wxmp/server"
)

const (
	appId          = "" // 应用ID
	token          = "" // 令牌
	encodingAesKey = "" // 消息加解密密钥
)

func main() {
	store := server.NewMemoryStore()
	store.SetToken(appId, token)
	store.SetEncodingAESKey(appId, encodingAesKey)

	h := server.New(store)
	// 注册文本消息处理器
	h.OnMsg(server.MsgTypeText, TextMsgHandler)
	// 注册关注事件处理器
	h.OnEvt(server.EvtTypeSubscribe, SubscribeEvtHandler)

	http.Handle("/wxmp", h)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// 文本消息处理器
func TextMsgHandler(c *server.Context) {
	// 回音墙
	content := c.ReqMsg.Content
	err := c.ReplyTextMsg(content)
	if err != nil {
		log.Print(err)
	}
}

// 关注事件处理器
func SubscribeEvtHandler(c *server.Context) {
	err := c.ReplyTextMsg("Hey guy")
	if err != nil {
		log.Print(err)
	}
}

```

## 相关链接

[微信公众平台](https://mp.weixin.qq.com/)

[微信公众平台开发者文档](http://mp.weixin.qq.com/wiki)

[微信公众平台接口调试工具](http://mp.weixin.qq.com/debug/)

[微信公众平台接口测试帐号申请](http://mp.weixin.qq.com/debug/cgi-bin/sandbox?t=sandbox/login)

## 许可协议

[The MIT License (MIT)](LICENSE)
