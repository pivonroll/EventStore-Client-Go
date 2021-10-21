package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

type UpdateOptionsEmitOptionType string

const (
	UpdateOptionsEmitOptionEnabledType UpdateOptionsEmitOptionType = "UpdateOptionsEmitOptionEnabledType"
	UpdateOptionsEmitOptionNoEmitType  UpdateOptionsEmitOptionType = "UpdateOptionsEmitOptionNoEmitType"
)

type UpdateOptionsEmitOption interface {
	GetType() UpdateOptionsEmitOptionType
}

type UpdateOptionsEmitOptionEnabled struct {
	EmitEnabled bool
}

func (u UpdateOptionsEmitOptionEnabled) GetType() UpdateOptionsEmitOptionType {
	return UpdateOptionsEmitOptionEnabledType
}

type UpdateOptionsEmitOptionNoEmit struct{}

func (u UpdateOptionsEmitOptionNoEmit) GetType() UpdateOptionsEmitOptionType {
	return UpdateOptionsEmitOptionNoEmitType
}

type UpdateOptionsRequest struct {
	EmitOption UpdateOptionsEmitOption
	Query      string
	Name       string
}

func (updateConfig *UpdateOptionsRequest) build() *projections.UpdateReq {
	if strings.TrimSpace(updateConfig.Name) == "" {
		panic("Failed to build UpdateOptionsRequest. Trimmed name is an empty string")
	}

	if strings.TrimSpace(updateConfig.Query) == "" {
		panic("Failed to build UpdateOptionsRequest. Trimmed query is an empty string")
	}

	result := &projections.UpdateReq{
		Options: &projections.UpdateReq_Options{
			Name:  updateConfig.Name,
			Query: updateConfig.Query,
		},
	}

	if updateConfig.EmitOption.GetType() == UpdateOptionsEmitOptionNoEmitType {
		result.Options.EmitOption = &projections.UpdateReq_Options_NoEmitOptions{
			NoEmitOptions: &shared.Empty{},
		}
	} else if updateConfig.EmitOption.GetType() == UpdateOptionsEmitOptionEnabledType {
		emitOption := updateConfig.EmitOption.(UpdateOptionsEmitOptionEnabled)
		result.Options.EmitOption = &projections.UpdateReq_Options_EmitEnabled{
			EmitEnabled: emitOption.EmitEnabled,
		}
	}

	return result
}
