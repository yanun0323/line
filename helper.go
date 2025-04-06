package line

import "fmt"

func invoke[T any](fn func(T) error, val T) error {
	if fn != nil {
		return fn(val)
	}
	return nil
}

func mapping[Input, Output any](input Input, fn func(Input) Output) Output {
	if fn != nil {
		return fn(input)
	}
	return *new(Output)
}

// NewMention create new mention string for line message
func NewMention(mentionKeyOrUserID string) string {
	return fmt.Sprintf("{%s}", mentionKeyOrUserID)
}
