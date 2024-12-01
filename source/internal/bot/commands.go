package bot

import (
	"Grade_Portal_TelegramBot/config"
	"Grade_Portal_TelegramBot/internal/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config) {
	switch update.Message.Command() {
	case "start":
		handlers.HandleStart(bot, update)
	case "getotp":
		handlers.HandleOTP(bot, update, update.Message.CommandArguments(), cfg)
	case "register":
		handlers.HandleRegister(bot, update, update.Message.CommandArguments(), cfg)
	case "resetpassword":
		handlers.HandleRegister(bot, update, update.Message.CommandArguments(), cfg)
	case "login":
		handlers.HanldeLogin(bot, update, update.Message.CommandArguments(), cfg)
	case "help":
		handlers.HandleHelp(bot, update)
	case "info":
		handlers.HandleInfo(bot, update, cfg)
	case "grade":
		handlers.HandleGrade(bot, update, update.Message.CommandArguments(), cfg)
	case "allgrade":
		handlers.HandleAllGrade(bot, update, cfg)
	case "clear":
		handlers.HandleClear(bot, update)
	case "history":
		handlers.HandleHistory(bot, update)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Lệnh không hợp lệ. Dùng /help để xem danh sách lệnh.")
		bot.Send(msg)
	}
}
