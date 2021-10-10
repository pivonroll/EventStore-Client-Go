package persistent

import (
	"github.com/pivonroll/EventStore-Client-Go/protos/persistent"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
	"github.com/pivonroll/EventStore-Client-Go/stream_revision"
)

// SubscriptionGroupForStreamAllRequest is a struct with all data necessary to create a persistent subscription group.
type SubscriptionGroupForStreamAllRequest struct {
	GroupName string // name of the persistent subscription group
	// AllPosition
	// AllPositionStart
	// AllPositionEnd
	Position stream_revision.IsReadPositionAll // position from which we want to start to receive events from a stream $all
	// CreateRequestAllNoFilter
	// CreateRequestAllFilter
	Filter   isCreateRequestAllFilter      // filter for messages from stream $all
	Settings CreateOrUpdateRequestSettings // setting for a persistent subscription group
}

func (request SubscriptionGroupForStreamAllRequest) build() *persistent.CreateReq {
	streamOption := &persistent.CreateReq_Options_All{
		All: &persistent.CreateReq_AllOptions{},
	}

	buildCreateRequestPosition(request.Position, streamOption.All)
	request.Filter.build(streamOption.All)
	protoSettings := request.Settings.buildCreateRequestSettings()

	result := &persistent.CreateReq{
		Options: &persistent.CreateReq_Options{
			StreamOption: streamOption,
			GroupName:    request.GroupName,
			Settings:     protoSettings,
		},
	}

	return result
}

func buildCreateRequestPosition(
	position stream_revision.IsReadPositionAll,
	protoOptions *persistent.CreateReq_AllOptions) {
	switch position.(type) {
	case stream_revision.ReadPositionAllStart:
		protoOptions.AllOption = &persistent.CreateReq_AllOptions_Start{
			Start: &shared.Empty{},
		}
	case stream_revision.ReadPositionAll:
		protoOptions.AllOption = &persistent.CreateReq_AllOptions_Position{
			Position: &persistent.CreateReq_Position{
				CommitPosition:  position.(stream_revision.ReadPositionAll).CommitPosition,
				PreparePosition: position.(stream_revision.ReadPositionAll).PreparePosition,
			},
		}
	case stream_revision.ReadPositionAllEnd:
		protoOptions.AllOption = &persistent.CreateReq_AllOptions_End{
			End: &shared.Empty{},
		}
	}
}

type isCreateRequestAllFilter interface {
	isCreateRequestAllFilter()
	build(*persistent.CreateReq_AllOptions)
}

type CreateRequestAllNoFilter struct{}

func (c CreateRequestAllNoFilter) isCreateRequestAllFilter() {
}

func (c CreateRequestAllNoFilter) build(protoOptions *persistent.CreateReq_AllOptions) {
	protoOptions.FilterOption = &persistent.CreateReq_AllOptions_NoFilter{
		NoFilter: &shared.Empty{},
	}
}

type CreateRequestAllFilter struct {
	FilterBy CreateRequestAllFilterByType
	// CreateRequestAllFilterByRegex
	// CreateRequestAllFilterByPrefix
	Matcher isCreateRequestAllFilterMatcher
	// CreateRequestAllFilterWindowMax
	// CreateRequestAllFilterWindowCount
	Window                       isCreateRequestAllFilterWindow
	CheckpointIntervalMultiplier uint32
}

type CreateRequestAllFilterByType string

const (
	CreateRequestAllFilterByStreamIdentifier CreateRequestAllFilterByType = "CreateRequestAllFilterByStreamIdentifier"
	CreateRequestAllFilterByEventType        CreateRequestAllFilterByType = "CreateRequestAllFilterByEventType"
)

func (c CreateRequestAllFilter) isCreateRequestAllFilter() {
}

func (c CreateRequestAllFilter) build(protoOptions *persistent.CreateReq_AllOptions) {
	filter := &persistent.CreateReq_AllOptions_Filter{
		Filter: &persistent.CreateReq_AllOptions_FilterOptions{
			CheckpointIntervalMultiplier: c.CheckpointIntervalMultiplier,
		},
	}

	c.Window.build(filter.Filter)
	protoMatcher := c.Matcher.build()

	switch c.FilterBy {
	case CreateRequestAllFilterByEventType:
		filter.Filter.Filter = &persistent.CreateReq_AllOptions_FilterOptions_EventType{
			EventType: protoMatcher,
		}
	case CreateRequestAllFilterByStreamIdentifier:
		filter.Filter.Filter = &persistent.CreateReq_AllOptions_FilterOptions_StreamIdentifier{
			StreamIdentifier: protoMatcher,
		}
	}
	protoOptions.FilterOption = filter
}

type isCreateRequestAllFilterMatcher interface {
	isCreateRequestAllFilterByStreamIdentifierMatcher()
	build() *persistent.CreateReq_AllOptions_FilterOptions_Expression
}

type CreateRequestAllFilterByRegex struct {
	Regex string
}

func (c CreateRequestAllFilterByRegex) isCreateRequestAllFilterByStreamIdentifierMatcher() {
}

func (c CreateRequestAllFilterByRegex) build() *persistent.CreateReq_AllOptions_FilterOptions_Expression {
	return &persistent.CreateReq_AllOptions_FilterOptions_Expression{
		Regex: c.Regex,
	}
}

type CreateRequestAllFilterByPrefix struct {
	Prefix []string
}

func (c CreateRequestAllFilterByPrefix) isCreateRequestAllFilterByStreamIdentifierMatcher() {
}

func (c CreateRequestAllFilterByPrefix) build() *persistent.CreateReq_AllOptions_FilterOptions_Expression {
	return &persistent.CreateReq_AllOptions_FilterOptions_Expression{
		Prefix: c.Prefix,
	}
}

type isCreateRequestAllFilterWindow interface {
	isCreateRequestAllFilterWindow()
	build(*persistent.CreateReq_AllOptions_FilterOptions)
}

type CreateRequestAllFilterWindowMax struct {
	Max uint32
}

func (c CreateRequestAllFilterWindowMax) isCreateRequestAllFilterWindow() {
}

func (c CreateRequestAllFilterWindowMax) build(options *persistent.CreateReq_AllOptions_FilterOptions) {
	options.Window = &persistent.CreateReq_AllOptions_FilterOptions_Max{
		Max: c.Max,
	}
}

type CreateRequestAllFilterWindowCount struct{}

func (c CreateRequestAllFilterWindowCount) isCreateRequestAllFilterWindow() {
}

func (c CreateRequestAllFilterWindowCount) build(options *persistent.CreateReq_AllOptions_FilterOptions) {
	options.Window = &persistent.CreateReq_AllOptions_FilterOptions_Count{
		Count: &shared.Empty{},
	}
}
