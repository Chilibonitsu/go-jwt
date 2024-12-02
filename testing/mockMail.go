package testing

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) Send(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func TestSendEmail(t *testing.T) {
	mockMailer := new(MockMailer)
	mockMailer.On("Send", "example@example.com", "Subject", "Body").Return(nil)
	mockMailer.AssertExpectations(t)
}

func SendEmail() {
	// Sender data.
	from := "from@gmail.com"
	password := "<Email Password>"

	// Receiver email address.
	to := []string{
		"sender@example.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("This is a test email message.")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
