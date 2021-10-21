package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

type EmitType string

const (
	EmitEnabledType EmitType = "EmitEnabledType"
	NoEmitType      EmitType = "NoEmitType"
)

type IsEmit interface {
	GetType() EmitType
}

type EmitEnabled struct {
	EmitEnabled bool
}

func (u EmitEnabled) GetType() EmitType {
	return EmitEnabledType
}

type NoEmit struct{}

func (u NoEmit) GetType() EmitType {
	return NoEmitType
}

// UpdateRequest are projection's name and it's options which we want to update.
type UpdateRequest struct {
	EmitOption     IsEmit
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

	if updateConfig.EmitOption.GetType() == NoEmitType {
		result.Options.EmitOption = &projections.UpdateReq_Options_NoEmitOptions{
			NoEmitOptions: &shared.Empty{},
		}
	} else if updateConfig.EmitOption.GetType() == EmitEnabledType {
		emitOption := updateConfig.EmitOption.(EmitEnabled)
		result.Options.EmitOption = &projections.UpdateReq_Options_EmitEnabled{
			EmitEnabled: emitOption.EmitEnabled,
		}
	}

	return result
}
