package bot

type Command interface {
	// Invokes returns a list of strings that can be used to invoke the command, all of which are case-insensitive.
	Invokes() []string
	// Description returns a short description of the command.
	Description() string
	// Exec executes the command.
	Exec(ctx *Context) error
}
