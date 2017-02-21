// notifier package contains loop function for checking GMail account, filtered email
// and send notification to Telegram private chat
package notifier

import (
	"github.com/alexivanenko/egroupware_notifier_bot/config"
	"github.com/alexivanenko/egroupware_notifier_bot/mail"
	"github.com/alexivanenko/egroupware_notifier_bot/model"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/context"
)

// LoadMailsAndNotify uses mail package and gets messages about EGroupware new/modified
// tasks messages from GMail Account.
// Then makes and send notification to Telegram private chat by ChatID from config.ini
func LoadMailsAndNotify() {
	ctx := context.Background()
	var notifiedTasks []*model.Task

	messages, err := mail.GetMessages(ctx)

	if err != nil {
		config.Log(err.Error())
	} else {
	messageLoop:
		for _, m := range messages {

			//Load and Check already notified tasks
			if tasks, err := model.LoadTasks(); err == nil {
				notifiedTasks = tasks
			}

			for _, task := range notifiedTasks {
				if task.Number == m.TaskNumber {
					continue messageLoop
				}
			}

			if err := send(m); err == nil {
				setNotified(m.TaskNumber)
			}
		}
	}

}

//send uses compacted message for sending notification to Telegram private chat
func send(msg mail.Message) error {
	//Create and Send Message to Telegram chat
	text := "<b>From</b>: " + msg.From + "\n<b>Subject</b>: <i>" + msg.Subject + "</i>" +
		"\n<b>Priority</b>: " + msg.Body.Priority +
		"\n<b>Summary</b>: " + msg.Body.Summary

	TMsg := tgbotapi.NewMessage(chatId, text)
	TMsg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.Send(TMsg)

	return err
}

//setNotified stores task number in DB
func setNotified(taskNumber int) {
	task := new(model.Task)
	task.Number = taskNumber
	task.Save()
}
