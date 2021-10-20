package event_reader

import (
	"github.com/google/uuid"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

func toProtoUUID(id uuid.UUID) *shared.UUID {
	return &shared.UUID{
		Value: &shared.UUID_String_{
			String_: id.String(),
		},
	}
}

