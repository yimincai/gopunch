package bot

import (
	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Session  *discordgo.Session
	Message  *discordgo.Message
	Args     []string
	Handler  *CommandHandler
	Commands []Command
}
