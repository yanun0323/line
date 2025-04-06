package line

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/pkg/errors"
	"github.com/yanun0323/line/internal"
)

// LineMessageID is the ID of the message.
type LineMessageID string

// NotifyMessageOption is the option for the message.
type NotifyMessageOption struct {
	// QuoteToken is the token of the message to be quoted
	QuoteToken string
	// MentionUserID is the user ID of the message to be mentioned
	MentionUserID map[string]string
}

// Notifier is the interface for the notifier.
type Notifier interface {
	// ReplyMessage [FREE] reply message to user
	ReplyMessage(replyToken, text string, opt ...NotifyMessageOption) (LineMessageID, error)

	// SendMessage [PAID] send message to user
	SendMessage(targetID, text string, opt ...NotifyMessageOption) (LineMessageID, error)
}

type lineNotifier struct {
	bot       *messaging_api.MessagingApiAPI
	botUserID string
}

// NewNotifier creates a new notifier which is used to send message to a user/group/room.
//
// # Example:
//
//	notifier, err := line.NewNotifier("LINE_CHANNEL_ACCESS_TOKEN")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Reply message from user/group/room. It doesn't cost any money.
//	notifier.ReplyMessage("replyToken", "Hello, world!")
//
//	// Send message to user/group/room. It costs money.
//	notifier.SendMessage("targetID", "Hello, world!")
func NewNotifier(channelAccessToken string) (Notifier, error) {
	bot, err := messaging_api.NewMessagingApiAPI(
		channelAccessToken,
	)
	if err != nil {
		return nil, fmt.Errorf("connect to bot, err: %+v", err)
	}

	info, err := bot.GetBotInfo()
	if err != nil {
		return nil, fmt.Errorf("get bot info, err: %+v", err)
	}

	return &lineNotifier{
		botUserID: info.UserId,
		bot:       bot,
	}, nil
}

func (r *lineNotifier) ReplyMessage(replyToken, text string, opt ...NotifyMessageOption) (LineMessageID, error) {

	mentionMap := make(map[string]messaging_api.SubstitutionObjectInterface)
	option := NotifyMessageOption{}
	if len(opt) != 0 {
		option = opt[0]
	}

	for key, userID := range option.MentionUserID {
		mentionMap[key] = messaging_api.MentionSubstitutionObject{
			SubstitutionObject: messaging_api.SubstitutionObject{"mention"},
			Mentionee: messaging_api.UserMentionTarget{
				MentionTarget: messaging_api.MentionTarget{"user"},
				UserId:        userID,
			},
		}
	}

	res, err := r.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				internal.TextMessageV2Fix{
					Message:      messaging_api.Message{Type: "textV2"},
					Text:         text,
					Substitution: mentionMap,
					QuoteToken:   option.QuoteToken,
				},
			},
		},
	)
	if err != nil {
		return "", errors.Errorf("reply message, err: %+v", err)
	}

	if len(res.SentMessages) == 0 {
		return "", errors.New("no sent message")
	}

	return LineMessageID(res.SentMessages[0].Id), nil
}

func (r *lineNotifier) SendMessage(targetID, text string, opt ...NotifyMessageOption) (LineMessageID, error) {
	option := NotifyMessageOption{}
	if len(opt) != 0 {
		option = opt[0]
	}

	mentionMap := make(map[string]messaging_api.SubstitutionObjectInterface)
	for key, userID := range option.MentionUserID {
		mentionMap[key] = messaging_api.MentionSubstitutionObject{
			SubstitutionObject: messaging_api.SubstitutionObject{"mention"},
			Mentionee: messaging_api.UserMentionTarget{
				MentionTarget: messaging_api.MentionTarget{"user"},
				UserId:        userID,
			},
		}
	}

	req := &messaging_api.PushMessageRequest{
		To: targetID,
		Messages: []messaging_api.MessageInterface{
			internal.TextMessageV2Fix{
				Message:      messaging_api.Message{Type: "textV2"},
				Text:         text,
				Substitution: mentionMap,
				QuoteToken:   option.QuoteToken,
			},
		},
	}

	res, err := r.bot.PushMessage(req, "")
	if err != nil {
		return "", errors.Errorf("send message, err: %+v", err)
	}

	if len(res.SentMessages) == 0 {
		return "", errors.New("no sent message")
	}

	return LineMessageID(res.SentMessages[0].Id), nil
}
