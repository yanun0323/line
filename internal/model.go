package internal

import (
	"encoding/json"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

type TextMessageV2Fix struct {
	messaging_api.Message

	/**
	 * Get QuickReply
	 */
	QuickReply *messaging_api.QuickReply `json:"quickReply,omitempty"`

	/**
	 * Get Sender
	 */
	Sender *messaging_api.Sender `json:"sender,omitempty"`

	/**
	 * Get Text
	 */
	Text string `json:"text"`

	/**
	 * A mapping that specifies substitutions for parts enclosed in {} within the &#39;text&#39; field.
	 */
	Substitution map[string]messaging_api.SubstitutionObjectInterface `json:"substitution,omitempty"`

	/**
	 * Quote token of the message you want to quote.
	 */
	QuoteToken string `json:"quoteToken,omitempty"`
}

func (r *TextMessageV2Fix) MarshalJSON() ([]byte, error) {

	type Alias TextMessageV2Fix
	return json.Marshal(&struct {
		*Alias

		Type string `json:"type"`
	}{
		Alias: (*Alias)(r),

		Type: "textV2",
	})
}
