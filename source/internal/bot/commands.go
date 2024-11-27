package bot
import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"Grade_Portal_TelegramBot/internal/handlers"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		handlers.HandleStart(bot, update)
	case "getOTP":
		handlers.HandleOTP(bot, update, update.Message.CommandArguments())
	case "register":
		handlers.HandleRegister(bot, update, update.Message.CommandArguments())
	case "resetPassWord":
		handlers.HandleRegister(bot, update, update.Message.CommandArguments())
	case "login":
		handlers.HanldeLogin(bot, update, update.Message.CommandArguments())
	case "help":
		handlers.HandleHelp(bot, update)
	case "info":
		handlers.HandleInfo(bot, update)
	case "grade":
		handlers.HandleGrade(bot, update, update.Message.CommandArguments())
	case "allGrade":
		handlers.HandleAllGrade(bot, update)
	case "clear":
		handlers.HandleClear(bot, update)
	case "history":
		handlers.HandleHistory(bot, update)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Lệnh không hợp lệ. Dùng /help để xem danh sách lệnh.")
		bot.Send(msg)
	}
}