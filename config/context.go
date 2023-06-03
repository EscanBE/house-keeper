package config

import (
	"github.com/EscanBE/go-app-name/database"
	"github.com/EscanBE/go-app-name/database/postgres"
	"github.com/EscanBE/go-lib/logging"
	"github.com/EscanBE/go-lib/telegram/bot"
	libutils "github.com/EscanBE/go-lib/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ApplicationExecutionContext hold the working context of the application entirely.
// It contains application configuration, logger, as well as connection pool,...
type ApplicationExecutionContext struct {
	Config   ApplicationConfig
	Bot      *bot.TelegramBot
	Database database.Database
	Logger   logging.Logger
}

// NewContext inits `ApplicationExecutionContext` with needed information to be stored within and return it
func NewContext(conf *ApplicationConfig, bot *bot.TelegramBot) *ApplicationExecutionContext {
	logger := logging.NewDefaultLogger()

	db, err2 := postgres.NewPostgresDatabase(conf.Endpoints.Database, logger)
	libutils.ExitIfErr(err2, "unable to create new database client")

	ctx := &ApplicationExecutionContext{
		Config:   *conf,
		Bot:      bot,
		Database: db,
		Logger:   logger,
	}

	err := logger.ApplyConfig(conf.Logging)
	libutils.ExitIfErr(err, "failed to apply logging config")

	return ctx
}

func (aec ApplicationExecutionContext) SendTelegramLogMessage(msg string) (*tgbotapi.Message, error) {
	if aec.Bot == nil {
		return nil, nil
	}
	m, err := aec.SendTelegramMessage(tgbotapi.NewMessage(aec.Config.TelegramConfig.LogChannelID, msg))
	if err != nil {
		aec.Logger.Error("Failed to send telegram log message", "type", "log", "error", err.Error())
	}
	return m, err
}

func (aec ApplicationExecutionContext) SendTelegramErrorMessage(msg string) (*tgbotapi.Message, error) {
	if aec.Bot == nil {
		return nil, nil
	}
	m, err := aec.SendTelegramMessage(tgbotapi.NewMessage(aec.Config.TelegramConfig.ErrChannelID, msg))
	if err != nil {
		aec.Logger.Error("Failed to send telegram error message", "type", "error", "error", err.Error())
		return nil, err
	}
	return m, nil
}

func (aec ApplicationExecutionContext) SendTelegramError(err error) (*tgbotapi.Message, error) {
	return aec.SendTelegramErrorMessage(err.Error())
}

func (aec ApplicationExecutionContext) SendTelegramMessage(c tgbotapi.Chattable) (*tgbotapi.Message, error) {
	if aec.Bot == nil {
		return nil, nil
	}
	m, err := aec.Bot.Send(c)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
