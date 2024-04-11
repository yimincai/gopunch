package commands

// type CommandDeleteUser struct {
// 	Cfg  *config.Config
// 	Repo repository.Repository
// }
//
// func (c *CommandDeleteUser) Invokes() []string {
// 	return []string{"Ping"}
// }
//
// func (c *CommandDeleteUser) Description() string {
// 	return "pong!"
// }
//
// func (c *CommandDeleteUser) Exec(ctx *bot.Context) (err error) {
// 	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "pong!")
// 	if err != nil {
// 		logger.Errorf("Error sending message: %s", err)
// 	}
//
// 	logger.Infof("Command Executed: %v, UserID: %s", c.Invokes(), ctx.Message.Author.ID)
// 	return
// }
