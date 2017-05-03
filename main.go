package main

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/apex/go-apex"
	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"fmt"
)

type response struct {
	Text string `json:"err"`
}

type KeyValueCmd struct {
	name  string
	key   string
	value string
	cmd   string
}

func ParseKeyValue(message *tgbotapi.Message) (*KeyValueCmd, error) {

	var kvc KeyValueCmd
	space := strings.Index(message.Text, " ")
	if space < 0 {
		kvc.cmd = message.Text
	} else {
		kvc.cmd = message.Text[:space]
	}

	kvc.name = fmt.Sprintf("%d.%s.%s", message.From.ID, message.From.FirstName, message.From.LastName)

	Info("name:", kvc.name)

	equal := strings.Index(message.Text, "=")
	if equal < 0 {
		kvc.key = message.Text[space+1:]
	} else {
		kvc.key = message.Text[space+1 : equal]
		kvc.value = message.Text[equal+1:]
	}

	Info(kvc)
	return &kvc, nil
}

func (kvc *KeyValueCmd) Get() (*string, error) {

	if len(kvc.name) == 0 || len(kvc.key) == 0 {
		err := errors.New("lack name and key. ex:/get name.key")

		return nil, err
	}

	params := dynamodb.GetItemInput{
		TableName: aws.String("demo"),
		Key: map[string]*dynamodb.AttributeValue{ // Required
			"name": { // Required
				S: aws.String(kvc.name),
			},
		},
	}
	var err error
	var text string
	resp, err := DynOp(params)
	if err != nil {
		Error(err)
		return nil, err
	} else {
		output, exist := resp.(dynamodb.GetItemOutput)
		if !exist {
			err = errors.New(fmt.Sprintf("%s.%s is not exist", kvc.name, kvc.key))
			Error(err, output)
			return nil, err
		}

		key, exist := output.Item["key"]
		if !exist {
			err = errors.New(fmt.Sprintf("%s is not exist", kvc.key))
			Error(err, output.Item)
			return nil, err
		}

		value, exist := output.Item["value"]
		if !exist {
			err = errors.New(fmt.Sprintf("%s is not exist", kvc.value))
			Error(err, output.Item)
			return nil, err
		}

		text = fmt.Sprintf("%s=%s", aws.StringValue(key.S), aws.StringValue(value.S))

	}

	return &text, err

}

func (kvc *KeyValueCmd) Set() error {

	if len(kvc.name) == 0 || len(kvc.key) == 0 || len(kvc.value) == 0 {
		err := errors.New("lack name and key. ex:/set name.key=value")
		return err
	}

	params := dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{ // Required
			"name": {
				S: aws.String(kvc.name),
			},
		},
		TableName:        aws.String("demo"),
		UpdateExpression: aws.String("SET #key = :keyStr,#value = :valueStr"),
		ExpressionAttributeNames: map[string]*string{ // Required
			"#key":   aws.String("key"),
			"#value": aws.String("value"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":keyStr":   {S: aws.String(kvc.key)},
			":valueStr": {S: aws.String(kvc.value)},
		},
	}
	var err error
	resp, err := DynOp(params)
	if err != nil {
		Error(err)
		return err
	} else {
		output, exist := resp.(dynamodb.UpdateItemOutput)
		if !exist {
			err = errors.New(fmt.Sprintf("%s.%s is not exist", kvc.name, kvc.key))
			return err
		}

		Info(output.String())
	}

	return err

}

func main() {

	Info("Start")

	bot, err := tgbotapi.NewBotAPI("239897563:AAE1E8Y-g_AcoylQCT2YyJDpEg4Ga9vlCsTk")
	if err != nil {
		Fatal(err)
	}

	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {

		var update tgbotapi.Update

		if err := json.Unmarshal(event, &update); err != nil {
			return nil, err
		}

		if update.Message == nil {
			return response{Text: "not from telegram"}, nil
		}
		var msg tgbotapi.MessageConfig

		kvc, err := ParseKeyValue(update.Message)
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			return nil, err
		}

		Infof("[%s] cmd=%s,name=%s,key=%s", update.Message.From.UserName, kvc.cmd, kvc.name, kvc.key)

		switch kvc.cmd {
		case "/start":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "GientBot will watch your key-value and trigger an action with a specify condition. It has states, thus GiantBot hold a state machine for you.")
		case "/get":

			text, err := kvc.Get()
			if err != nil {
				err := err.Error()
				text = &err
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, *text)

		case "/set":

			var text string
			err := kvc.Set()
			if err != nil {
				err := err.Error()
				text = err
			} else {
				text = "ok"
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
		case "/keyboard":
			keylayout := [][]tgbotapi.KeyboardButton{
				{{Text: "/set key=value"}, {Text: "/get key"}},
			}

			rkm := tgbotapi.ReplyKeyboardMarkup{
				Keyboard:        keylayout,
				ResizeKeyboard:  true,
				OneTimeKeyboard: false,
				Selective:       false}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "demo keyboard")
			msg.ReplyMarkup = &rkm

		default:

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		}

		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		var resp response
		if err == nil {
			resp.Text = "nil"
		} else {
			resp.Text = err.Error()
		}

		return resp, nil
	})

	Info("Exit")
}
