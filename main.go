package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/arykalin/ehser-notifier/users"
	"gopkg.in/yaml.v2"

	"go.uber.org/zap"

	"github.com/spf13/pflag"

	"github.com/arykalin/ehser-notifier/mailer"
	"github.com/arykalin/ehser-notifier/notifier"
	sheet "github.com/arykalin/ehser-notifier/scheet"
	"github.com/arykalin/ehser-notifier/telegram"
)

type Config struct {
	AnswersSheetID   string `yaml:"sheet_answers"`
	SkipSheet        int    `yaml:"sheet_skip"`
	MailUser         string `yaml:"mail_user"`
	MailPassword     string `yaml:"mail_password"`
	MailHost         string `yaml:"mail_host"`
	MailPort         int    `yaml:"mail_port"`
	MailDebugAddress string `yaml:"mail_debug_address"`
	MailCCAddress    string `yaml:"mail_cc_address"`
	SentFile         string `yaml:"sent_file"`
	TeleToken        string `yaml:"telegram_token"`
	TeleChatID       int64  `yaml:"telegram_chat_id"`
}

func main() {
	pathConfig := pflag.StringP("path", "c", "./config.yml", "path to config file")
	help := pflag.BoolP("help", "h", false, "show help")
	pflag.Parse()

	b, err := ioutil.ReadFile(*pathConfig)
	if err != nil {
		log.Fatalf("can't read file")
	}

	if *help {
		pflag.PrintDefaults()
		os.Exit(0)
	}

	var config Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		log.Fatalf("can't unmarshal config: %s", err)
	}

	sLoggerConfig := zap.NewDevelopmentConfig()
	sLoggerConfig.DisableStacktrace = true
	sLoggerConfig.DisableCaller = true
	sLogger, err := sLoggerConfig.Build()
	if err != nil {
		panic(err)
	}
	logger := sLogger.Sugar()

	s, err := sheet.NewSheetService(logger)
	if err != nil {
		log.Fatalf("failed to init sheet client: %s", err)
	}

	// Add users from first form
	// get answers
	formUsers := users.NewUsers()

	err = addUsersForm(s, config, logger, formUsers)

	err = formUsers.DumpUsers()
	if err != nil {
		logger.Fatalf("failed to dump users map: %s", err)
	}

	newMailer := mailer.NewMailer(
		logger,
		config.MailUser,
		config.MailPassword,
		config.MailHost,
		config.MailPort,
		config.MailDebugAddress,
		config.MailCCAddress,
	)
	newTeleLog := telegram.NewTelegramLog(
		config.TeleChatID,
		config.TeleToken,
	)
	n := notifier.NewNotifier(
		logger,
		newMailer,
		config.SentFile,
		newTeleLog,
	)

	err = n.Notify(formUsers.GetUsers())
	if err != nil {
		logger.Fatalf("error notify: %s", err)
	}
}

func addUsersForm(s sheet.Sheet, config Config, logger *zap.SugaredLogger, formUsers users.UsersInt) error {
	spreadsheet, err := s.GetSheet(config.AnswersSheetID)
	if err != nil {
		logger.Fatalf("failed to get sheet data: %s", err)
	}
	gotSheet, err := spreadsheet.SheetByIndex(0)
	if err != nil {
		logger.Fatalf("failed to get sheet data: %s", err)
	}
	sheetConfig := users.SheetConfig{
		MailIdx:             1,
		NameIdx:             2,
		AgeIdx:              3,
		PhoneIdx:            4,
		SocialLinkIdx:       5,
		CityIdx:             6,
		ExperienceAnswerIdx: 9,
		Skip:                config.SkipSheet,
	}
	err = formUsers.AddUsers(gotSheet, &sheetConfig)
	if err != nil {
		logger.Fatalf("failed to make users map: %s", err)
	}
	return err
}
