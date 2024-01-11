package helpers

import (
	"database/sql"
	_ "errors"
	"fmt"
	_ "io/ioutil"
	_ "net/http"
	_ "net/url"
	"strconv"
	"strings"
	_ "strings"
	"time"
	"ttnmwastemanagementsystem/configs"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/models"

	"gopkg.in/gomail.v2"
)

const MAIL_TAG string = "helper/MAIL"

type Mail struct{}

type MailConfig struct {
	MailHost        string
	MailPort        int
	MailUsername    string
	MailPassword    string
	MailEncryption  string
	MailFromAddress string
	MailTo          string
	MailNoReply     string
	MailFromName    string
}

func GetMailConfig() (*MailConfig, error) {

	mailHost := configs.EnvConfigs.MailHost
	mailPort := configs.EnvConfigs.MailPort
	mailUsername := configs.EnvConfigs.MailUsername
	mailPassword := configs.EnvConfigs.MailPassword
	mailEncryption := configs.EnvConfigs.MailEncryption
	mailFromAddress := configs.EnvConfigs.MailFromAddress

	mailConfig := MailConfig{
		MailHost:        mailHost,
		MailPort:        mailPort,
		MailUsername:    mailUsername,
		MailPassword:    mailPassword,
		MailEncryption:  mailEncryption,
		MailFromAddress: mailFromAddress,
	}

	return &mailConfig, nil
}

// retuns true if we have managed to send the mail, else false
func (mail Mail) SendPasswordResetMail(email string, name string) bool {
	functions := Functions{}
	token := functions.TokenGenerator()

	var timeOut = configs.EnvConfigs.AccountRecoveryTokenExpirationTime

	//insert this to the database
	_, err := gen.REPO.DB.NamedExec(`UPDATE users set recovery_token=:token,recovery_sent_at=:sent_at where email=:email`,
		map[string]interface{}{
			"email":   email,
			"token":   token,
			"sent_at": time.Now(),
		})

	if err != nil {
		logger.Log(MAIL_TAG, fmt.Sprint("Failed to insert reset token,", err), logger.LOG_LEVEL_ERROR)
		return false
	}

	mailConfig, err := GetMailConfig()
	if err != nil {
		return false
	}
	logger.Log(MAIL_TAG, fmt.Sprint("Mail config", mailConfig.MailPassword), logger.LOG_LEVEL_INFO)

	passwordResetLink := configs.EnvConfigs.AppURL + "/auth/challenge/enter_new_password/web?token=" + token

	templateString := functions.FileToString("./templates/auth/password_reset.html")
	if templateString == "" {
		logger.Log(MAIL_TAG, "Failed to load password reset template", logger.LOG_LEVEL_ERROR)
		return false
	}

	templateString = strings.ReplaceAll(templateString, "{{.password_reset_link}}", passwordResetLink)
	templateString = strings.ReplaceAll(templateString, "{{.password_reset_timeout}}", strconv.Itoa(int(timeOut))+" Hrs")

	templateString = functions.ReplaceTemplateWithOrganizationInformation(templateString)
	msg := gomail.NewMessage()
	msg.SetHeader("From", mailConfig.MailFromAddress)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Password reset link")
	msg.SetBody("text/html", templateString)

	n := gomail.NewDialer(mailConfig.MailHost, int(mailConfig.MailPort), mailConfig.MailUsername, mailConfig.MailPassword)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		logger.Log(MAIL_TAG, fmt.Sprint("Error sending email ::", err), logger.LOG_LEVEL_ERROR)
		if configs.CanDebug() {
			fmt.Println("Error sending email :: ", err)
		}
		return false
	}
	return true
}

// retuns true if we have managed to send the mail, else false
func (mail Mail) SendToken(email string, name string) bool {
	functions := Functions{}
	token, token_err := functions.NumberTokenGenerator(4)
	if token_err != nil {
		return false
	}

	var timeOut = configs.EnvConfigs.AcoountVerificationTokenExpirationTime
	emailVerificationToken := models.EmailVerificationToken{}
	err := gen.REPO.DB.Get(&emailVerificationToken, gen.REPO.DB.Rebind("select * from email_verification_token  where email=?"), email)

	if err != nil && err != sql.ErrNoRows {
		return false
	}

	if err != nil { //no rows
		gen.REPO.DB.NamedExec(`insert into email_verification_token(email,token,created_at)values(:email,:token,:created_at)`,
			map[string]interface{}{
				"email":      email,
				"token":      token,
				"created_at": time.Now(),
			})
	} else {
		//update
		gen.REPO.DB.NamedExec(`UPDATE email_verification_token set token=:token,created_at=:created_at where email=:email`,
			map[string]interface{}{
				"email":      email,
				"token":      token,
				"created_at": time.Now(),
			})

	}

	mailConfig, err := GetMailConfig()
	if err != nil {
		logger.Log(MAIL_TAG, fmt.Sprint("Error getting mail config", err.Error()), logger.LOG_LEVEL_ERROR)
		return false
	}
	logger.Log(MAIL_TAG, fmt.Sprint("Mail config - ", mailConfig), logger.LOG_LEVEL_INFO)

	templateString := functions.FileToString("./templates/auth/verification_token_email.html")
	if templateString == "" {
		logger.Log(MAIL_TAG, "Failed to load email template", logger.LOG_LEVEL_ERROR)
		return false
	}

	templateString = strings.ReplaceAll(templateString, "{{.code_expiration}}", strconv.Itoa(int(timeOut))+" Hrs")
	templateString = strings.ReplaceAll(templateString, "{{.code}}", fmt.Sprint(token))
	templateString = functions.ReplaceTemplateWithOrganizationInformation(templateString)

	msg := gomail.NewMessage()
	msg.SetHeader("From", mailConfig.MailFromAddress)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Verify email")
	msg.SetBody("text/html", templateString)

	n := gomail.NewDialer(mailConfig.MailHost, int(mailConfig.MailPort), mailConfig.MailUsername, mailConfig.MailPassword)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		logger.Log(MAIL_TAG, fmt.Sprint("Error sending email ::", err), logger.LOG_LEVEL_ERROR)
		if configs.CanDebug() {
			fmt.Println("Error sending email :: ", err)
		}
		return false
	}
	return true
}

// retuns true if we have managed to send the mail, else false
func (mail Mail) SendMailVerification(email string, name string) bool {
	functions := Functions{}
	token := functions.TokenGenerator()

	var timeOut = configs.EnvConfigs.AcoountVerificationTokenExpirationTime

	//insert this to the database
	_, err := gen.REPO.DB.NamedExec(`UPDATE users set confirmation_token=:token,confirmation_sent_at=:sent_at where email=:email`,
		map[string]interface{}{
			"email":   email,
			"token":   token,
			"sent_at": time.Now(),
		})

	if err != nil {
		logger.Log(MAIL_TAG, fmt.Sprint("Failed to insert reset  confirmation token,", err), logger.LOG_LEVEL_ERROR)
		return false
	}

	mailConfig, err := GetMailConfig()
	if err != nil {
		logger.Log(MAIL_TAG, fmt.Sprint("Error getting mail config", err.Error()), logger.LOG_LEVEL_ERROR)
		return false
	}
	logger.Log(MAIL_TAG, fmt.Sprint("Mail config - ", mailConfig), logger.LOG_LEVEL_INFO)

	emailVerifyLink := configs.EnvConfigs.AppURL + "/auth/challenge/verify_email/web?token=" + token

	templateString := functions.FileToString("./templates/auth/verify_email.html")
	if templateString == "" {
		logger.Log(MAIL_TAG, "Failed to load email verify template", logger.LOG_LEVEL_ERROR)
		return false
	}

	templateString = strings.ReplaceAll(templateString, "{{.email_verify_link}}", emailVerifyLink)
	templateString = strings.ReplaceAll(templateString, "{{.email_verify_timeout}}", strconv.Itoa(int(timeOut))+" Hrs")
	templateString = functions.ReplaceTemplateWithOrganizationInformation(templateString)

	msg := gomail.NewMessage()
	msg.SetHeader("From", mailConfig.MailFromAddress)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Verify email")
	msg.SetBody("text/html", templateString)

	n := gomail.NewDialer(mailConfig.MailHost, int(mailConfig.MailPort), mailConfig.MailUsername, mailConfig.MailPassword)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		logger.Log(MAIL_TAG, fmt.Sprint("Error sending email ::", err), logger.LOG_LEVEL_ERROR)
		if configs.CanDebug() {
			fmt.Println("Error sending email :: ", err)
		}
		return false
	}
	return true
}
