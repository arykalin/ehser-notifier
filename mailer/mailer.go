package mailer

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/arykalin/ehser-notifier/users"
	"go.uber.org/zap"

	"gopkg.in/gomail.v2"
)

type mailer struct {
	logger       *zap.SugaredLogger
	User         string
	Password     string
	SMTPHost     string
	SMTPPort     int
	DebugAddress string
	CCAddress    string
	dryRun       bool
}

const (
	track = "Среда как зона поиска, влияния и поддержки (трек для педагогов)"
)

type Mailer interface {
	SendGreeting(user users.User) error
}

const (
	tmpltEsherRisk = "mailer/mails/template_esher_protivorechie.html"
)

func (m mailer) SendGreeting(user users.User) (err error) {

	body, subj, err := m.makeBodyAndSubj(user)
	if err != nil {
		return err
	}

	if body == "" {
		return fmt.Errorf("body is empty")
	}

	if subj == "" {
		return fmt.Errorf("subject is empty")
	}

	err = m.sendMail(user, subj, body)
	if err != nil {
		return err
	}
	return nil
}

func (m mailer) sendMail(user users.User, subj string, body string) (err error) {
	if m.dryRun {
		return nil
	}
	//if debug address is set send mail to it
	var toEmail string
	if m.DebugAddress != "" {
		toEmail = m.DebugAddress
	} else {
		toEmail = user.Email
	}

	gm := gomail.NewMessage()
	gm.SetHeader("From", m.User)
	gm.SetHeader("To", toEmail)
	if m.CCAddress != "" {
		gm.SetHeader("Cc", m.CCAddress)
	}
	gm.SetHeader("Subject", subj)
	gm.SetBody("text/html", body)

	d := gomail.NewDialer(m.SMTPHost, 587, m.User, m.Password)

	// Send the email.
	if err = d.DialAndSend(gm); err != nil {
		return err
	}
	return err
}

func (m mailer) makeBodyAndSubj(user users.User) (body string, subj string, err error) {
	subj = fmt.Sprintf("Регистрация на зимнюю школу ЭШЭР Интуиция")
	body, err = m.ParseTemplate(tmpltEsherRisk, user)
	if err != nil {
		return body, subj, err
	}
	m.logger.Debugw("user template", "user", user.Email, "template", tmpltEsherRisk)
	return body, subj, err
}

func (m *mailer) ParseTemplate(templateFileName string, data interface{}) (body string, err error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return body, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return body, err
	}
	return buf.String(), err
}

func NewMailer(
	logger *zap.SugaredLogger,
	user string,
	password string,
	host string,
	port int,
	debugAddress string,
	ccAddress string,
) Mailer {
	return &mailer{
		logger:       logger,
		User:         user,
		Password:     password,
		SMTPHost:     host,
		SMTPPort:     port,
		DebugAddress: debugAddress,
		CCAddress:    ccAddress,
		dryRun:       false,
	}
}
