# Line Bot SDK

**Note: This is a private repository for personal use only.**

## Overview

This repository contains a Go library for interacting with the LINE Messaging API. It provides a convenient way to handle LINE webhook events and build chatbots or messaging applications.

## Features

- Event handling for various LINE webhook events
- Type-safe event processing with Go generics
- Support for different source types (User, Group, Room)
- Handling for message, join, leave, and member events

## Usage

```go
package example

import (
	"fmt"
	"log"
	"os"

	"github.com/yanun0323/line"
)

func ExampleBot() {
	// Initialize a new bot with your channel secret
	bot, err := line.NewBot(os.Getenv("LINE_CHANNEL_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	// Create a notifier with your channel access token
	notifier, err := line.NewNotifier(os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	// Set a handler for message events
	bot.SetMessageEventHandler(func(event line.EventMessage) error {
		msg := fmt.Sprintf("Hello, This is your reply: %s", event.Data.Text)
		_, err := notifier.ReplyMessage(event.Data.ReplyToken, msg)
		if err != nil {
			return err
		}

		return nil
	})

	// Start the webhook server
	bot.ListenAndServe(":8080", "/callback")
}
```

## Project Structure

- `event.go`: Contains event type definitions and structures
- `example/bot.go`: Example implementation of a LINE bot

## License

This is a private project and not intended for distribution.
