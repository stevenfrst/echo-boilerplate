package smtpEmail

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.elastic.co/apm"
	"log"
	"net/smtp"
	"strings"
)

func SendMail(c echo.Context,mailTo []string, subject string, message string) error {
	span, _ := apm.StartSpan(c.Request().Context(), "Sending Email", "request")
	defer span.End()

	bcc := []string{viper.GetString(`smtp.email`)}
	mime := "\r\n" + "MIME-Version: 1.0\r\n" + "Content-Type: text/html; charset=\"utf-8\"\r\n\r\n"
	body := "From: " + viper.GetString(`smtp.sender_name`) + "\n" +
		"To: " + strings.Join(mailTo, ",") + "\n" +
		"Cc: " + strings.Join(bcc, ",") + "\n" +
		"Subject: " + subject + mime + message

	auth := smtp.PlainAuth("", viper.GetString(`smtp.email`), viper.GetString(`smtp.password`), viper.GetString(`smtp.host`))
	smtpAddr := fmt.Sprintf("%s:%d", viper.GetString(`smtp.host`), viper.GetInt(`smtp.port`))

	err := smtp.SendMail(smtpAddr, auth, viper.GetString(`smtp.email`), append(mailTo, bcc...), []byte(body))
	if err != nil {
		return err
	}
	log.Println("Mail sent!")
	return nil
}
