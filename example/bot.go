package example

import (
	"fmt"
	"log"
	"os"

	"github.com/yanun0323/line"
)

func ExampleBot() {
	bot, err := line.NewBot(os.Getenv("LINE_CHANNEL_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	notifier, err := line.NewNotifier(os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.SetMessageEventHandler(func(event line.EventMessage) error {
		msg := fmt.Sprintf("Hello, This is your reply: %s", event.Data.Text)
		_, err := notifier.ReplyMessage(event.Data.ReplyToken, msg)
		if err != nil {
			return err
		}

		return nil
	})

	bot.ListenAndServe(":8080", "/callback")
}
