package commands

import (
	"fmt"

	"github.com/JiveshL-KDK/go-gator/internal/state"
)

type CommandType struct {
	name        string
	args        []string
	description string
	handler     func(s *state.State, cmd CommandInput) error
}

type CommandInput struct {
	name string
	args []string
}

func (cmd CommandInput) Run() error {

	state, err := state.InitialzeState()

	if err != nil {
		return err
	}

	command, ok := gatorCommands[cmd.name]

	if !ok {
		return fmt.Errorf("command not found!")
	}

	err = command.handler(state, cmd)

	if err != nil {
		return err
	}

	return nil

}

func NewCommand(name string, args []string) CommandInput {
	return CommandInput{name: name, args: args}
}

type CommandMap map[string]CommandType

var gatorCommands CommandMap = CommandMap{
	"login": CommandType{
		name:        "login",
		args:        []string{"<username>"},
		description: "Login as <username>",
		handler:     handleLogin,
	},
	"register": CommandType{
		name:        "register",
		args:        []string{"<username>"},
		description: "Register a <username>",
		handler:     handleRegister,
	},

	"list": CommandType{
		name:        "list",
		args:        []string{},
		description: "List all users",
		handler:     listUsers,
	},
	"reset": CommandType{
		name:        "reset",
		args:        []string{},
		description: "Remove all users",
		handler:     resetAllUsers,
	},
	"agg": CommandType{
		name:        "agg",
		args:        []string{},
		description: "Aggregate feed",
		handler:     agg,
	},

	"add": CommandType{
		name:        "add",
		args:        []string{},
		description: "add a new feed",
		handler:     middlewareLoggedIn(addFeed),
	},
	"list-feeds": CommandType{
		name:        "list-feeds",
		args:        []string{},
		description: "list feed of the current users",
		handler:     middlewareLoggedIn(listFeeds),
	},
	"follow": CommandType{
		name:        "follow",
		args:        []string{},
		description: "follow a new feed",
		handler:     middlewareLoggedIn(follow),
	},
	"following": CommandType{
		name:        "following",
		args:        []string{},
		description: "show currently followed feeds",
		handler:     middlewareLoggedIn(following),
	},
	"unfollow": CommandType{
		name:        "unfollow",
		args:        []string{},
		description: "unfollow a field",
		handler:     middlewareLoggedIn(unfollow),
	},
}
