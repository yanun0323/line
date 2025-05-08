package line

// SourceType is the type of the source of the event.
type SourceType int

const (
	// SourceTypeNotFound represents the source type is not found.
	SourceTypeNotFound SourceType = iota
	// SourceTypeUser represents the event is from a user.
	SourceTypeUser
	// SourceTypeGroup represents the event is from a group.
	SourceTypeGroup
	// SourceTypeRoom represents the event is from a room.
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

type eventStickerData struct {
	ReplyToken string
	PackageID  string
	StickerID  string
}

// EventJoin is the event of a user joining a group or room.
type EventJoin event[eventJoinData]

// EventLeave is the event of a user leaving a group or room.
type EventLeave event[eventLeaveData]

// EventMemberJoined is the event of a user joining a group or room.
type EventMemberJoined event[eventMemberJoinedData]

// EventMemberLeft is the event of a user leaving a group or room.
type EventMemberLeft event[eventMemberLeftData]

// EventMessage is the event of a message.
type EventMessage event[eventMessageData]

// EventSticker is the event of a sticker.
type EventSticker event[eventStickerData]
