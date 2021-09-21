package notifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/arykalin/ehser-notifier/users"
	"go.uber.org/zap"

	"github.com/arykalin/ehser-notifier/mailer"
	"github.com/arykalin/ehser-notifier/telegram"
)

type SentFile struct {
	Sent map[string]bool `json:"mail_sent,omitempty"`
}

type notifier struct {
	logger   *zap.SugaredLogger
	mailer   mailer.Mailer
	sentFile string
	teleLog  telegram.TeleLog
}

type Notifier interface {
	Notify(users users.Users) error
}

func (n notifier) Notify(users users.Users) error {
	sentData := SentFile{}
	sentData.Sent = make(map[string]bool)
	byteValue, err := ioutil.ReadFile(n.sentFile)
	if err != nil {
		return err
	}
	//make backup
	backupFile := fmt.Sprintf("%s-%d", n.sentFile, time.Now().Unix())
	err = ioutil.WriteFile(backupFile, byteValue, 0644)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteValue, &sentData)
	if err != nil {
		return err
	}

	for mail, user := range users {
		sentRecord := fmt.Sprintf("%s", mail)
		if sent, ok := sentData.Sent[sentRecord]; ok && sent {
			n.logger.Debugw("message already sent. skip", "mail", mail)
			continue
		}
		n.logger.Debugf("sending links to user %s", mail)
		err = n.mailer.SendGreeting(user)
		if err != nil {
			msg := fmt.Sprintf("sending mail error mail: %s err: %s\n",
				mail, err)
			n.logger.Error(msg)
			terr := n.teleLog.SendMessage(msg)
			if terr != nil {
				n.logger.Errorw("sending telegram message error", "err", err)
			}
			continue
		}
		msg := fmt.Sprintf("Новая заявка\n"+
			"Имя: %s\n "+
			"Возраст: %s\n "+
			"Город: %s\n"+
			"Телефон: %s\n "+
			"Почта: %s\n "+
			"Соцсети: %s\n "+
			"Чего хочу: %s\n",
			user.FullName,
			user.Age,
			user.City,
			user.Phone,
			user.Email,
			user.SocialLink,
			user.ExperienceAnswer)
		terr := n.teleLog.SendMessage(msg)
		if terr != nil {
			n.logger.Errorw("sending telegram message error", "err", err)
		}
		sentData.Sent[sentRecord] = true
	}

	file, err := json.MarshalIndent(sentData, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(n.sentFile, file, 0644) //nolint:gosec
	return err
}

func NewNotifier(
	logger *zap.SugaredLogger,
	mailer mailer.Mailer,
	sentFile string,
	teleLog telegram.TeleLog,
) Notifier {
	return &notifier{
		logger:   logger,
		teleLog:  teleLog,
		mailer:   mailer,
		sentFile: sentFile,
	}
}
