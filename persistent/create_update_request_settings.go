package persistent

import (
	"fmt"

	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/persistent"
)

const SUBSCRIBER_COUNT_UNLIMITED = 0

type ConsumerStrategy int32

const (
	ConsumerStrategy_RoundRobin       ConsumerStrategy = 0
	ConsumerStrategy_DispatchToSingle ConsumerStrategy = 1
	ConsumerStrategy_Pinned           ConsumerStrategy = 2
)

// SubscriptionGroupSettings are settings for a persistent subscription group.
//
// You can read more at https://developers.eventstore.com/clients/grpc/persistent-subscriptions.html#persistent-subscription-settings
type SubscriptionGroupSettings struct {
	ResolveLinks          bool             // whether the subscription should resolve link events to their linked events.
	ExtraStatistics       bool             // whether to track latency statistics on this subscription.
	MaxRetryCount         int32            // the maximum number of retries (due to timeout) before a message is considered to be parked.
	MinCheckpointCount    int32            // the minimum number of messages to process before a checkpoint may be written.
	MaxCheckpointCount    int32            // the maximum number of messages not checkpointed before forcing a checkpoint.
	MaxSubscriberCount    int32            // the maximum number of subscribers allowed.
	LiveBufferSize        int32            // the size of the buffer (in-memory) listening to live messages as they happen before paging occurs.
	ReadBatchSize         int32            // the number of events read at a time when paging through history.
	HistoryBufferSize     int32            // the number of events to cache when paging through history.
	NamedConsumerStrategy ConsumerStrategy // the strategy to use for distributing events to client consumers.
	// MessageTimeoutInMs
	// MessageTimeoutInTicks
	MessageTimeout isCreateRequestMessageTimeout // the amount of time after which to consider a message as timed out and retried.
	// CheckpointAfterMs
	// CheckpointAfterTicks
	CheckpointAfter isCreateRequestCheckpointAfter // the amount of time to try to checkpoint after.
}

func (settings SubscriptionGroupSettings) buildCreateRequestSettings() *persistent.CreateReq_Settings {
	result := &persistent.CreateReq_Settings{
		ResolveLinks:          settings.ResolveLinks,
		ExtraStatistics:       settings.ExtraStatistics,
		MaxRetryCount:         settings.MaxRetryCount,
		MinCheckpointCount:    settings.MinCheckpointCount,
		MaxCheckpointCount:    settings.MaxCheckpointCount,
		MaxSubscriberCount:    settings.MaxSubscriberCount,
		LiveBufferSize:        settings.LiveBufferSize,
		ReadBatchSize:         settings.ReadBatchSize,
		HistoryBufferSize:     settings.HistoryBufferSize,
		NamedConsumerStrategy: consumerStrategyProto(settings.NamedConsumerStrategy),
	}

	settings.MessageTimeout.buildCreateRequestSettings(result)
	settings.CheckpointAfter.buildCreateRequestSettings(result)

	return result
}

func (settings SubscriptionGroupSettings) buildUpdateRequestSettings() *persistent.UpdateReq_Settings {
	result := &persistent.UpdateReq_Settings{
		ResolveLinks:          settings.ResolveLinks,
		ExtraStatistics:       settings.ExtraStatistics,
		MaxRetryCount:         settings.MaxRetryCount,
		MinCheckpointCount:    settings.MinCheckpointCount,
		MaxCheckpointCount:    settings.MaxCheckpointCount,
		MaxSubscriberCount:    settings.MaxSubscriberCount,
		LiveBufferSize:        settings.LiveBufferSize,
		ReadBatchSize:         settings.ReadBatchSize,
		HistoryBufferSize:     settings.HistoryBufferSize,
		NamedConsumerStrategy: updateRequestConsumerStrategyProto(settings.NamedConsumerStrategy),
	}

	settings.MessageTimeout.buildUpdateRequestSettings(result)
	settings.CheckpointAfter.buildUpdateRequestSettings(result)

	return result
}

// DefaultRequestSettings are the default values for any persistent subscription group.
// Settings are recommended by EventStoreDB.
// More info can be found at https://developers.eventstore.com/clients/grpc/persistent-subscriptions.html#persistent-subscription-settings
var DefaultRequestSettings = SubscriptionGroupSettings{
	ResolveLinks:          false,
	ExtraStatistics:       false,
	MaxRetryCount:         10,
	MinCheckpointCount:    10,
	MaxCheckpointCount:    10 * 1000,
	MaxSubscriberCount:    SUBSCRIBER_COUNT_UNLIMITED,
	LiveBufferSize:        500,
	ReadBatchSize:         20,
	HistoryBufferSize:     500,
	NamedConsumerStrategy: ConsumerStrategy_RoundRobin,
	MessageTimeout: MessageTimeoutInMs{
		MilliSeconds: 30_000,
	},
	CheckpointAfter: CheckpointAfterMs{
		MilliSeconds: 2_000,
	},
}

type isCreateRequestMessageTimeout interface {
	isCreateRequestMessageTimeout()
	buildCreateRequestSettings(*persistent.CreateReq_Settings)
	buildUpdateRequestSettings(*persistent.UpdateReq_Settings)
}

// MessageTimeoutInMs is the amount of time (in milliseconds) after which to consider a message as timed out and retried.
type MessageTimeoutInMs struct {
	MilliSeconds int32
}

func (c MessageTimeoutInMs) isCreateRequestMessageTimeout() {
}

func (c MessageTimeoutInMs) buildCreateRequestSettings(
	protoSettings *persistent.CreateReq_Settings,
) {
	protoSettings.MessageTimeout = &persistent.CreateReq_Settings_MessageTimeoutMs{
		MessageTimeoutMs: c.MilliSeconds,
	}
}

func (c MessageTimeoutInMs) buildUpdateRequestSettings(
	protoSettings *persistent.UpdateReq_Settings,
) {
	protoSettings.MessageTimeout = &persistent.UpdateReq_Settings_MessageTimeoutMs{
		MessageTimeoutMs: c.MilliSeconds,
	}
}

// MessageTimeoutInTicks is the amount of time (in .NET ticks) after which to consider a message as timed out and retried.
// A single tick represents one hundred nanoseconds.
type MessageTimeoutInTicks struct {
	Ticks int64
}

func (c MessageTimeoutInTicks) isCreateRequestMessageTimeout() {
}

func (c MessageTimeoutInTicks) buildCreateRequestSettings(
	protoSettings *persistent.CreateReq_Settings,
) {
	protoSettings.MessageTimeout = &persistent.CreateReq_Settings_MessageTimeoutTicks{
		MessageTimeoutTicks: c.Ticks,
	}
}

func (c MessageTimeoutInTicks) buildUpdateRequestSettings(
	protoSettings *persistent.UpdateReq_Settings,
) {
	protoSettings.MessageTimeout = &persistent.UpdateReq_Settings_MessageTimeoutTicks{
		MessageTimeoutTicks: c.Ticks,
	}
}

type isCreateRequestCheckpointAfter interface {
	isCreateRequestCheckpointAfter()
	buildCreateRequestSettings(*persistent.CreateReq_Settings)
	buildUpdateRequestSettings(*persistent.UpdateReq_Settings)
}

// CheckpointAfterTicks is the amount of time (in .NET ticks) to try to checkpoint after.
// A single tick represents one hundred nanoseconds.
type CheckpointAfterTicks struct {
	Ticks int64
}

func (c CheckpointAfterTicks) isCreateRequestCheckpointAfter() {
}

func (c CheckpointAfterTicks) buildCreateRequestSettings(
	protoSettings *persistent.CreateReq_Settings,
) {
	protoSettings.CheckpointAfter = &persistent.CreateReq_Settings_CheckpointAfterTicks{
		CheckpointAfterTicks: c.Ticks,
	}
}

func (c CheckpointAfterTicks) buildUpdateRequestSettings(
	protoSettings *persistent.UpdateReq_Settings,
) {
	protoSettings.CheckpointAfter = &persistent.UpdateReq_Settings_CheckpointAfterTicks{
		CheckpointAfterTicks: c.Ticks,
	}
}

// CheckpointAfterMs is the amount of time (in milliseconds) to try to checkpoint after.
type CheckpointAfterMs struct {
	MilliSeconds int32
}

func (c CheckpointAfterMs) isCreateRequestCheckpointAfter() {
}

func (c CheckpointAfterMs) buildCreateRequestSettings(protoSettings *persistent.CreateReq_Settings) {
	protoSettings.CheckpointAfter = &persistent.CreateReq_Settings_CheckpointAfterMs{
		CheckpointAfterMs: c.MilliSeconds,
	}
}

func (c CheckpointAfterMs) buildUpdateRequestSettings(
	protoSettings *persistent.UpdateReq_Settings,
) {
	protoSettings.CheckpointAfter = &persistent.UpdateReq_Settings_CheckpointAfterMs{
		CheckpointAfterMs: c.MilliSeconds,
	}
}

func consumerStrategyProto(strategy ConsumerStrategy) persistent.CreateReq_ConsumerStrategy {
	switch strategy {
	case ConsumerStrategy_DispatchToSingle:
		return persistent.CreateReq_DispatchToSingle
	case ConsumerStrategy_Pinned:
		return persistent.CreateReq_Pinned
	case ConsumerStrategy_RoundRobin:
		return persistent.CreateReq_RoundRobin
	default:
		panic(fmt.Sprintf("Could not map strategy %v to proto", strategy))
	}
}

func updateRequestConsumerStrategyProto(
	strategy ConsumerStrategy,
) persistent.UpdateReq_ConsumerStrategy {
	switch strategy {
	case ConsumerStrategy_DispatchToSingle:
		return persistent.UpdateReq_DispatchToSingle
	case ConsumerStrategy_Pinned:
		return persistent.UpdateReq_Pinned
	case ConsumerStrategy_RoundRobin:
		return persistent.UpdateReq_RoundRobin
	default:
		panic(fmt.Sprintf("Could not map strategy %v to proto", strategy))
	}
}
