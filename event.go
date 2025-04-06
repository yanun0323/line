package line

type SourceType int

const (
	SourceTypeNotFound SourceType = iota
	SourceTypeUser
	SourceTypeGroup
	SourceTypeRoom
)

type source struct {
	Type    SourceType
	UserID  string
	GroupID string
	RoomID  string
}

type event[Data any] struct {
	WebhookEventID string
	Source         source
	Timestamp      int64 /* milliseconds */
	Data           Data
}

type eventJoinData struct {
	ReplyToken string
}

type eventLeaveData struct {
}

type eventMemberJoinedData struct {
	ReplyToken      string
	JoinedMemberIDs []string
}

type eventMemberLeftData struct {
	LeftMemberIDs []string
}

type eventMessageData struct {
	ReplyToken      string
	MessageID       string
	Text            string
	QuoteToken      string
	QuotedMessageID string
}

type EventJoin event[eventJoinData]

type EventLeave event[eventLeaveData]

type EventMemberJoined event[eventMemberJoinedData]

type EventMemberLeft event[eventMemberLeftData]

type EventMessage event[eventMessageData]
