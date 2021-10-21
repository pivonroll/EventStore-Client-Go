package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

type isEmit interface {
	isEmit()
}

// EmitEnabled setting determines whether a projection can emit events and any projection
// that calls emit() or linkTo() requires it.
// Read more at https://developers.eventstore.com/server/v21.6/projections/projections-config.html#emit-enabled
type EmitEnabled struct {
	EmitEnabled bool
}

func (u EmitEnabled) isEmit() {
}

// NoEmit should be used to remove emit set from a projection.
type NoEmit struct{}

func (u NoEmit) isEmit() {
}

// UpdateRequest represents a set of options which we want to update for a specific projection.
// Projection is referenced by name.
type UpdateRequest struct {
	EmitOption     isEmit
	Query          string
	ProjectionName string
}

func (updateConfig UpdateRequest) build() *projections.UpdateReq {
	if strings.TrimSpace(updateConfig.ProjectionName) == "" {
		panic("Failed to build UpdateRequest. Trimmed projection name is an empty string")
	}

	if strings.TrimSpace(updateConfig.Query) == "" {
		panic("Failed to build UpdateRequest. Trimmed query is an empty string")
	}

	result := &projections.UpdateReq{
		Options: &projections.UpdateReq_Options{
			Name:  updateConfig.ProjectionName,
			Query: updateConfig.Query,
		},
	}

	switch updateConfig.EmitOption.(type) {
	case EmitEnabled:
		emitOption := updateConfig.EmitOption.(EmitEnabled)
		result.Options.EmitOption = &projections.UpdateReq_Options_EmitEnabled{
			EmitEnabled: emitOption.EmitEnabled,
		}
	case NoEmit:
		result.Options.EmitOption = &projections.UpdateReq_Options_NoEmitOptions{
			NoEmitOptions: &shared.Empty{},
		}
	}

	return result
}
