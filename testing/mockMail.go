package testing

import (
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
