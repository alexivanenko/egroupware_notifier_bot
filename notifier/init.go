// notifier package contains loop function for checking GMail account, filtered email
// and send notification to Telegram private chat
package notifier

import (
	"log"
	"strconv"

	"github.com/alexivanenko/egroupware_notifier_bot/config"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var chatId int64
var bot *tgbotapi.BotAPI

func init() {
	//Telegram Chat ID (because we send messages to only one private chat we stored chat ID in the config)
	id, err := strconv.ParseInt(config.String("bot", "chat_id"), 10, 64)
	if err != nil {
		log.Panic(err)
	} else {
		chatId = id
	}

	//Init Telegram Bot
	botApi, err := tgbotapi.NewBotAPI(config.String("bot", "token"))
	if err != nil {
		log.Panic(err)
	} else {
		bot = botApi
	}

	if config.Is("bot", "debug") {
		bot.Debug = true
	} else {
		bot.Debug = false
	}

	config.Log("Authorized on account " + bot.Self.UserName)
}
