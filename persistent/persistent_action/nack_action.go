package persistent_action

// Nack_Action represents a reason why message was nack-ed.
type Nack_Action int32

const (
	// Nack_Unknown means that the client does not know what action to take. Let the server decide.
	Nack_Unknown Nack_Action = 0
	// Nack_Park means to park the message and do not resend. Put it on poison queue.
	Nack_Park Nack_Action = 1
	// Nack_Retry means explicitly retry the message.
	Nack_Retry Nack_Action = 2
	// Nack_Skip means to skip this message, do not resend and do not put in poison queue.
	Nack_Skip Nack_Action = 3
	// Nack_Stop is meanT to stop the subscription.
	// Currently, in EventStoreDB 21.6 this functionality is not implemented.
	Nack_Stop Nack_Action = 4
)
