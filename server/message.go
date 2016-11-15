package server

import (
	"encoding/xml"
)

const (
	// 消息类型
	MsgTypeText       = "text"       // 文本消息
	MsgTypeImage      = "image"      // 图片消息
	MsgTypeVoice      = "voice"      // 语音消息
	MsgTypeVideo      = "video"      // 视频消息
	MsgTypeShortVideo = "shortvideo" // 小视频消息
	MsgTypeLocation   = "location"   // 地理位置消息
	MsgTypeLink       = "link"       // 链接消息
	// 事件类型
	EvtTypeSubscribe   = "subscribe"   // 关注事件/用户未关注时扫描带参数二维码事件
	EvtTypeUnsubscribe = "unsubscribe" // 取消关注事件
	EvtTypeScan        = "SCAN"        // 用户已关注时扫描带参数二维码事件
	EvtTypeLocation    = "LOCATION"    // 上报地理位置事件
	EvtTypeClick       = "CLICK"       // 自定义菜单拉取消息事件
	EvtTypeView        = "VIEW"        // 自定义菜单跳转链接事件
)

type ReqMsg struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      string
	ToUserName   string // 开发者微信号
	FromUserName string // 发送方OpenID
	CreateTime   int64  // 消息创建时间
	MsgType      string // 消息类型
	// 普通消息参数
	Content      string  // 文本消息内容
	PicURL       string  `xml:"PicUrl"`  // 图片链接
	MediaID      string  `xml:"MediaId"` // 图片/语音/视频消息媒体ID
	Format       string  // 语音格式
	Recognition  string  // 语音识别结果
	ThumbMediaID string  `xml:"ThumbMediaId"` // 视频消息缩略图的媒体ID
	LocationX    float64 `xml:"Location_X"`   // 地理位置维度
	LocationY    float64 `xml:"Location_Y"`   // 地理位置经度
	Scale        int     // 地图缩放大小
	Label        string  // 地理位置信息
	Title        string  // 消息标题
	Description  string  // 消息描述
	URL          string  `xml:"Url"`   // 消息链接
	MsgID        int64   `xml:"MsgId"` // 消息ID
	// 事件推送参数
	Event     string  // 事件类型
	EventKey  string  // 事件KEY值
	Ticket    string  // 二维码的Ticket
	Latitude  float64 // 地理位置纬度
	Longitude float64 // 地理位置经度
	Precision float64 // 地理位置精度
}

type CDATA struct {
	Value string `xml:",cdata"`
}

type msgBase struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   CDATA
	MsgType      CDATA
}

type encMsg struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATA
	MsgSignature CDATA
	TimeStamp    CDATA
	Nonce        CDATA
}

type textMsg struct {
	*msgBase
	Content CDATA
}

type imageMsg struct {
	*msgBase
	MediaID CDATA `xml:"Image>MediaId"`
}

type voiceMsg struct {
	*msgBase
	MediaID CDATA `xml:"Voice>MediaId"`
}

type videoMsg struct {
	*msgBase
	Video struct {
		MediaID     CDATA `xml:"MediaId"`
		Title       CDATA
		Description CDATA
	}
}

type musicMsg struct {
	*msgBase
	Music struct {
		Title        CDATA
		Description  CDATA
		MusicURL     CDATA `xml:"MusicUrl"`
		HQMusicURL   CDATA `xml:"HQMusicUrl"`
		ThumbMediaID CDATA `xml:"ThumbMediaId"`
	}
}

type newsMsg struct {
	*msgBase
	ArticleCount CDATA
	Articles     []*Article `xml:">item"`
}

type transfer2CustomerService struct {
	*msgBase
	KfAccount CDATA `xml:"TransInfo>KfAccount"`
}

type Article struct {
	Title       CDATA
	Description CDATA
	PicURL      CDATA `xml:"PicUrl"`
	URL         CDATA `xml:"Url"`
}
