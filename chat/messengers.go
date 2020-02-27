package chat

// Messenger provides the ability to send a message or reply with a given chat application
type Messenger interface {
	// SendImageReply is invoked when sending a message specifically as part of a response to a command (such as a slash command
	SendImageReply(reply *CommandReply) error
}

type CommandReply struct {
	RequestingUserName string
	WebhookURL         string
	ImageURL           string
}
