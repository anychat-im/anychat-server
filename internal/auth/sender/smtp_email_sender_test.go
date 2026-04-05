package sender

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSMTPEmailSenderDefaultsPort(t *testing.T) {
	sender, err := NewSMTPEmailSender(SMTPConfig{
		Host:        "smtp.mail.test",
		FromAddress: "noreply@example.com",
	})
	require.NoError(t, err)
	require.Equal(t, 587, sender.config.Port)
}

func TestNewSMTPEmailSenderRejectsInvalidConfig(t *testing.T) {
	_, err := NewSMTPEmailSender(SMTPConfig{
		Host:        "smtp.mail.test",
		FromAddress: "invalid-address",
	})
	require.Error(t, err)
}

func TestSMTPEmailSenderBuildMessage(t *testing.T) {
	sender, err := NewSMTPEmailSender(SMTPConfig{
		Host:        "smtp.mail.test",
		Port:        587,
		FromName:    "AnyChat 测试",
		FromAddress: "noreply@example.com",
	})
	require.NoError(t, err)

	message, err := sender.buildMessage("user@example.com", "验证码通知", "您的验证码为：123456")
	require.NoError(t, err)

	raw := string(message)
	require.Contains(t, raw, "To: user@example.com")
	require.Contains(t, raw, "Subject: =?UTF-8?")
	require.Contains(t, raw, "Content-Type: text/plain; charset=UTF-8")
	require.Contains(t, raw, "Content-Transfer-Encoding: quoted-printable")
	require.True(t, strings.Contains(raw, "123456") || strings.Contains(raw, "=31=32=33=34=35=36"))
}

func TestSMTPEmailSenderRejectsInvalidRecipientBeforeNetwork(t *testing.T) {
	sender, err := NewSMTPEmailSender(SMTPConfig{
		Host:        "smtp.mail.test",
		Port:        587,
		FromAddress: "noreply@example.com",
	})
	require.NoError(t, err)

	err = sender.Send("bad-recipient", "subject", "content")
	require.Error(t, err)
	require.Contains(t, err.Error(), "recipient")
}
