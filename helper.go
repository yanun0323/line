package line

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
