package models

import "errors"

var ErrInvalidEmailCode = errors.New("Invalid or expired email code")

type SendEmailCommand struct {
	To       []string
	Template string
	Data     map[string]interface{}
	Massive  bool
	Info     string
}

// const SendEmailCommandTable = "tpt_dh_send_email_command"

// func (a *SendEmailCommand) TableName() string {
// 	return SendEmailCommandTable
// }

type SendResetPasswordEmailCommand struct {
	User *User
}

type ValidateResetPasswordCodeQuery struct {
	Code   string
	Result *User
}
