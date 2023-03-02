package wechat

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/swgloomy/gutil/glog"
	"os"
	"wechat-chatting-service/chatgpt"
)

func Login() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.IsSendByFriend() {
			_, err := msg.ReplyText(chatgpt.AcquireContent(msg.Content))
			if err != nil {
				glog.Error("package:wechat func:Login MessageHandler ReplyText run err! err: %+v \n", err)
				return
			}
		}
	}
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	if err := bot.Login(); err != nil {
		glog.Error("package:wechat func:Login login run err! err: %+v \n", err)
		return
	}

	err := bot.Block()
	if err != nil {
		glog.Error("package:wechat func:Login Block err! application exit! err: %+v \n", err)
		os.Exit(1)
		return
	}
}
