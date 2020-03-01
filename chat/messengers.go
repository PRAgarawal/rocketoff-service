package chat

// Messenger provides the ability to send a message or reply with a given chat application
type Messenger interface {
	// SendImageReply is invoked when sending a message specifically as part of a response to a command (such as a slash command
	SendImageReply(reply *CommandReply) error
}

// CommandReply is all data that a chat application potentially needs in order to reply to a command.
type CommandReply struct {
	// RequestingUserID is the chat app's unique ID for the user we will reply to
	RequestingUserID string

	// WebhookURL is the ephemeral URL created for this command interaction on the given chat application
	WebhookURL       string

	// ImageURL is the URL of the image to reply with, if this is an image reply.
	ImageURL         string
}
