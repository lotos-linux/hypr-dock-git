package cli

type Action struct {
	Handler     func(data ...string) error
	NeedsData   bool
	Usage       string
	Description string
}

type Command struct {
	Description string
	Actions     map[string]Action
	Default     *Action
}

// var commands map[string]Command
