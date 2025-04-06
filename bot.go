package line

import (
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/pkg/errors"
	"github.com/yanun0323/line/internal"
)

type Bot interface {
	//	ListenAndServe listens on the TCP network address addr and then calls [Serve] to handle requests on incoming connections. Accepted connections are configured to enable TCP keep-alives.
	//
	//	If addr is blank, ":http" is used.
	//
	//	ListenAndServe always returns a non-nil error. After [Server.Shutdown] or [Server.Close], the returned error is [ErrServerClosed].
	ListenAndServe(addr string, callbackPath string) error

	//	SetJoinEventHandler sets the handler for join events.
	SetJoinEventHandler(func(EventJoin) error)

	//	SetLeaveEventHandler sets the handler for leave events.
	SetLeaveEventHandler(func(EventLeave) error)

	//	SetMemberJoinedEventHandler sets the handler for member joined events.
	SetMemberJoinedEventHandler(func(EventMemberJoined) error)

	//	SetMemberLeftEventHandler sets the handler for member left events.
	SetMemberLeftEventHandler(func(EventMemberLeft) error)

	//	SetMessageEventHandler sets the handler for message events.
	SetMessageEventHandler(func(EventMessage) error)
}

type bot struct {
	channelSecret string

	joinEventHandler         func(EventJoin) error
	leaveEventHandler        func(EventLeave) error
	memberJoinedEventHandler func(EventMemberJoined) error
	memberLeftEventHandler   func(EventMemberLeft) error
	messageEventHandler      func(EventMessage) error
}

// NewBot creates a new bot which is used to handle events from LINE.
//
// # Example:
//
//	// create a new bot
//	bot, err := line.NewBot("CHANNEL_SECRET")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// create a new notifier
//	notifier, err := line.NewNotifier("CHANNEL_ACCESS_TOKEN")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// set the message event handler
//	bot.SetMessageEventHandler(func(event EventMessage) error {
//		_, err := notifier.ReplyMessage(event.Data.ReplyToken, "Hello, world!")
//		if err != nil {
//			return err
//		}
//
//		return nil
//	})
//
//	// listen and serve
//	if err := bot.ListenAndServe(":8080", "/callback"); err != nil {
//		log.Fatal(err)
//	}
func NewBot(channelSecret string) (Bot, error) {
	return &bot{
		channelSecret: channelSecret,
	}, nil
}

func (b *bot) ListenAndServe(addr string, callbackPath string) error {
	mux := http.NewServeMux()
	mux.HandleFunc(callbackPath, b.HandleEvent)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server.ListenAndServe()
}

func (b *bot) SetJoinEventHandler(handler func(EventJoin) error) {
	b.joinEventHandler = handler
}

func (b *bot) SetLeaveEventHandler(handler func(EventLeave) error) {
	b.leaveEventHandler = handler
}

func (b *bot) SetMemberJoinedEventHandler(handler func(EventMemberJoined) error) {
	b.memberJoinedEventHandler = handler
}

func (b *bot) SetMemberLeftEventHandler(handler func(EventMemberLeft) error) {
	b.memberLeftEventHandler = handler
}

func (b *bot) SetMessageEventHandler(handler func(EventMessage) error) {
	b.messageEventHandler = handler
}

func (b *bot) HandleEvent(w http.ResponseWriter, req *http.Request) {
	// log.Print("/callback called...")
	cb, err := webhook.ParseRequest(b.channelSecret, req)
	if err != nil {
		log.Printf("parse request, err: %+v", err)

		if errors.Is(err, webhook.ErrInvalidSignature) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	// log.Print("Handling events...")
	for _, event := range cb.Events {
		// log.Printf("Start handling event: %T", event)
		var (
			err error
		)
		switch e := event.(type) {
		case webhook.JoinEvent:
			err = invoke(b.joinEventHandler, EventJoin{
				WebhookEventID: e.WebhookEventId,
				Source:         b.getSource(e.Source),
				Timestamp:      e.Timestamp,
				Data: eventJoinData{
					ReplyToken: e.ReplyToken,
				},
			})
		case webhook.LeaveEvent:
			err = invoke(b.leaveEventHandler, EventLeave{
				WebhookEventID: e.WebhookEventId,
				Source:         b.getSource(e.Source),
				Timestamp:      e.Timestamp,
				Data:           eventLeaveData{},
			})
		case webhook.MemberJoinedEvent:
			err = invoke(b.memberJoinedEventHandler, EventMemberJoined{
				WebhookEventID: e.WebhookEventId,
				Source:         b.getSource(e.Source),
				Timestamp:      e.Timestamp,
				Data: eventMemberJoinedData{
					ReplyToken: e.ReplyToken,
					JoinedMemberIDs: mapping(e.Joined.Members, func(members []webhook.UserSource) []string {
						ids := make([]string, len(members))
						for i, m := range members {
							ids[i] = m.UserId
						}
						return ids
					}),
				},
			})
		case webhook.MemberLeftEvent:
			err = invoke(b.memberLeftEventHandler, EventMemberLeft{
				WebhookEventID: e.WebhookEventId,
				Source:         b.getSource(e.Source),
				Timestamp:      e.Timestamp,
				Data: eventMemberLeftData{
					LeftMemberIDs: mapping(e.Left.Members, func(members []webhook.UserSource) []string {
						ids := make([]string, len(members))
						for i, m := range members {
							ids[i] = m.UserId
						}
						return ids
					}),
				},
			})
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				err = invoke(b.messageEventHandler, EventMessage{
					WebhookEventID: e.WebhookEventId,
					Source:         b.getSource(e.Source),
					Timestamp:      e.Timestamp,
					Data: eventMessageData{
						ReplyToken:      e.ReplyToken,
						MessageID:       message.Id,
						Text:            message.Text,
						QuoteToken:      message.QuoteToken,
						QuotedMessageID: message.QuotedMessageId,
					},
				})
			default:
				err = errors.Errorf("unsupported message content: %T", message)
			}
		default:
			err = errors.Errorf("unsupported event: %T", event)
		}

		if err != nil {
			log.Printf("%sERROR%s handle event, err: %+v", internal.ColorRed, internal.ColorReset, err)
		}
	}
}

func (b *bot) getSource(s webhook.SourceInterface) source {
	switch ss := s.(type) {
	case webhook.UserSource:
		return source{
			Type:   SourceTypeUser,
			UserID: ss.UserId,
		}
	case webhook.GroupSource:
		return source{
			Type:    SourceTypeGroup,
			UserID:  ss.UserId,
			GroupID: ss.GroupId,
		}
	case webhook.RoomSource:
		return source{
			Type:   SourceTypeRoom,
			UserID: ss.UserId,
			RoomID: ss.RoomId,
		}
	default:
		return source{
			Type: SourceTypeNotFound,
		}
	}
}
