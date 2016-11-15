package server

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-with/util"
)

const (
	mimeXML   = "text/xml"
	mimePlain = "text/plain"
)

type Context struct {
	writer    http.ResponseWriter
	writeOnce sync.Once

	Req    *http.Request
	ReqMsg *ReqMsg

	token  string
	aesKey []byte
}

// 被动回复文本消息
func (ctx *Context) ReplyTextMsg(content string) error {
	msg := new(textMsg)
	msg.Content.Value = content
	msg.msgBase = ctx.createMsgBase("text")
	return ctx.replyMsg(msg)
}

// 被动回复图片消息
func (ctx *Context) ReplyImageMsg(mediaID string) error {
	msg := new(imageMsg)
	msg.MediaID.Value = mediaID
	msg.msgBase = ctx.createMsgBase("image")
	return ctx.replyMsg(msg)
}

// 被动回复语音消息
func (ctx *Context) ReplyVoiceMsg(mediaID string) error {
	msg := new(voiceMsg)
	msg.MediaID.Value = mediaID
	msg.msgBase = ctx.createMsgBase("voice")
	return ctx.replyMsg(msg)
}

// 被动回复视频消息
func (ctx *Context) ReplyVideoMsg(title, descr, MediaID string) error {
	msg := new(videoMsg)
	msg.Video.Title.Value = title
	msg.Video.Description.Value = descr
	msg.Video.MediaID.Value = MediaID
	msg.msgBase = ctx.createMsgBase("video")
	return ctx.replyMsg(msg)
}

// 被动回复音乐消息
func (ctx *Context) ReplyMusicMsg(title, descr, musicURL, hqMusicURL, thumbMediaID string) error {
	msg := new(musicMsg)
	msg.Music.Title.Value = title
	msg.Music.Description.Value = descr
	msg.Music.MusicURL.Value = musicURL
	msg.Music.HQMusicURL.Value = hqMusicURL
	msg.Music.ThumbMediaID.Value = thumbMediaID
	msg.msgBase = ctx.createMsgBase("music")
	return ctx.replyMsg(msg)
}

// 被动回复图文消息
func (ctx *Context) ReplyNewsMsg(articles []*Article) error {
	msg := new(newsMsg)
	msg.ArticleCount.Value = fmt.Sprintf("%d", len(articles))
	msg.Articles = articles
	msg.msgBase = ctx.createMsgBase("news")
	return ctx.replyMsg(msg)
}

// 将消息转发到客服
func (ctx *Context) Transfer2CustomerService(kfAccount ...string) error {
	msg := new(transfer2CustomerService)
	if len(kfAccount) > 0 {
		msg.KfAccount.Value = kfAccount[0]
	}
	msg.msgBase = ctx.createMsgBase("transfer_customer_service")
	return ctx.replyMsg(msg)
}

// 创建消息基础
func (ctx *Context) createMsgBase(typ string) *msgBase {
	base := new(msgBase)
	base.ToUserName.Value = ctx.ReqMsg.FromUserName
	base.FromUserName.Value = ctx.ReqMsg.ToUserName
	base.CreateTime.Value = fmt.Sprintf("%d", time.Now().Unix())
	base.MsgType.Value = typ
	return base
}

// 被动回复消息
func (ctx *Context) replyMsg(msg interface{}) (err error) {
	data, err := xml.Marshal(msg)
	if err != nil {
		return
	}
	if ctx.SafeMode() {
		var enc *encMsg
		enc, err = ctx.encryptMsg(data)
		if err != nil {
			return
		}
		data, err = xml.Marshal(enc)
		if err != nil {
			return
		}
	}
	ctx.respond(200, mimeXML, data)
	return
}

// 加密消息
func (ctx *Context) encryptMsg(data []byte) (enc *encMsg, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int32(len(data)))
	if err != nil {
		return
	}
	length := buf.Bytes()
	random := util.RandAlnum(16)
	appID := []byte(ctx.AppID())

	plain := bytes.Join([][]byte{random, length, data, appID}, nil)
	cipher, err := aesEncrypt(plain, ctx.aesKey)
	if err != nil {
		return
	}
	enc = new(encMsg)
	enc.Nonce.Value = util.RandNumStr(10)
	enc.TimeStamp.Value = fmt.Sprintf("%d", time.Now().Unix())
	enc.Encrypt.Value = base64Encode(cipher)
	enc.MsgSignature.Value = signMsg(ctx.token, enc.TimeStamp.Value, enc.Nonce.Value, enc.Encrypt.Value)
	return
}

func (ctx *Context) noReply() {
	ctx.respond(200, mimePlain, []byte("success"))
}

func (ctx *Context) respond(code int, mime string, data []byte) {
	ctx.writeOnce.Do(func() {
		ctx.writer.Header().Set("Content-Type", mime+"; charset=utf-8")
		ctx.writer.WriteHeader(code)
		ctx.writer.Write(data)
	})
	return
}

// 应用ID
func (ctx *Context) AppID() string {
	return ctx.Req.FormValue("appid")
}

// 是否为安全模式
func (ctx *Context) SafeMode() bool {
	return ctx.Req.FormValue("encrypt_type") == "aes"
}
