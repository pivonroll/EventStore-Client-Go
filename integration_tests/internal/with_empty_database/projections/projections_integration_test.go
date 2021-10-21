package projections_integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pivonroll/EventStore-Client-Go/errors"
	"github.com/pivonroll/EventStore-Client-Go/projections/statistics"

	"github.com/stretchr/testify/require"

	"github.com/pivonroll/EventStore-Client-Go/projections"
)

func Test_CreateContinuousProjection_TrackEmittedStreamsFalse(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode: projections.ContinuousProjection{
			ProjectionName:      "MyContinuous_false",
			TrackEmittedStreams: false,
		},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)
}

func Test_CreateContinuousProjection_TrackEmittedStreamsTrue(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode: projections.ContinuousProjection{
			ProjectionName:      "MyContinuous_true",
			TrackEmittedStreams: true,
		},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)
}

func Test_CreateTransientProjection(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode: projections.TransientProjection{
			ProjectionName: "Transient",
		},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)
}

func Test_CreateOneTimeProjection(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode:  projections.OneTimeProjection{},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)
}

func Test_UpdateContinuousProjection_NoEmit(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode: projections.ContinuousProjection{
			ProjectionName:      "MyContinuous_no_emit",
			TrackEmittedStreams: false,
		},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)

	updateOptions := projections.UpdateRequest{
		EmitOption:     projections.NoEmit{},
		Query:          "fromAll().when({$init: function (s, e) {return {};}});",
		ProjectionName: "MyContinuous_no_emit",
	}

	err = client.UpdateProjection(context.Background(), updateOptions)
	require.NoError(t, err)
}

func Test_UpdateContinuousProjection_EmitFalse(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode: projections.ContinuousProjection{
			ProjectionName:      "MyContinuous_emit_false",
			TrackEmittedStreams: false,
		},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)

	updateOptions := projections.UpdateRequest{
		EmitOption:     projections.EmitEnabled{EmitEnabled: false},
		Query:          "fromAll().when({$init: function (s, e) {return {};}});",
		ProjectionName: "MyContinuous_emit_false",
	}

	err = client.UpdateProjection(context.Background(), updateOptions)
	require.NoError(t, err)
}

func Test_UpdateContinuousProjection_EmitTrue(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode: projections.ContinuousProjection{
			ProjectionName:      "MyContinuous_emit_true",
			TrackEmittedStreams: false,
		},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)

	updateOptions := projections.UpdateRequest{
		EmitOption:     projections.EmitEnabled{EmitEnabled: true},
		Query:          "fromAll().when({$init: function (s, e) {return {};}});",
		ProjectionName: "MyContinuous_emit_true",
	}

	err = client.UpdateProjection(context.Background(), updateOptions)
	require.NoError(t, err)
}

func Test_UpdateTransientProjection_NoEmit(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode: projections.TransientProjection{
			ProjectionName: "MyTransient_no_emit",
		},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)

	updateOptions := projections.UpdateRequest{
		EmitOption:     projections.NoEmit{},
		Query:          "fromAll().when({$init: function (s, e) {return {};}});",
		ProjectionName: "MyTransient_no_emit",
	}

	err = client.UpdateProjection(context.Background(), updateOptions)
	require.NoError(t, err)
}

func Test_UpdateTransientProjection_EmitFalse(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode: projections.TransientProjection{
			ProjectionName: "MyTransient_emit_false",
		},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)

	updateOptions := projections.UpdateRequest{
		EmitOption:     projections.EmitEnabled{EmitEnabled: false},
		Query:          "fromAll().when({$init: function (s, e) {return {};}});",
		ProjectionName: "MyTransient_emit_false",
	}

	err = client.UpdateProjection(context.Background(), updateOptions)
	require.NoError(t, err)
}

func Test_UpdateTransientProjection_EmitTrue(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode: projections.TransientProjection{
			ProjectionName: "MyTransient_emit_true",
		},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)

	updateOptions := projections.UpdateRequest{
		EmitOption:     projections.EmitEnabled{EmitEnabled: true},
		Query:          "fromAll().when({$init: function (s, e) {return {};}});",
		ProjectionName: "MyTransient_emit_true",
	}

	err = client.UpdateProjection(context.Background(), updateOptions)
	require.NoError(t, err)
}

func Test_AbortProjection(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	err := client.AbortProjection(context.Background(), StandardProjectionStreams)
	require.NoError(t, err)

	statisticsClient, err := client.GetProjectionStatistics(context.Background(), projections.StatisticsForProjectionByName{
		ProjectionName: StandardProjectionStreams,
	})
	require.NoError(t, err)
	require.NotNil(t, statisticsClient)

	result, stdErr := statisticsClient.Read()
	require.NoError(t, stdErr)
	require.EqualValues(t, statistics.StatusAborted, result.Status)
}

func Test_DisableProjection(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	err := client.DisableProjection(context.Background(), StandardProjectionStreams)
	require.NoError(t, err)

	statisticsClient, err := client.GetProjectionStatistics(context.Background(), projections.StatisticsForProjectionByName{
		ProjectionName: StandardProjectionStreams,
	})
	require.NoError(t, err)
	require.NotNil(t, statisticsClient)

	result, stdErr := statisticsClient.Read()
	require.NoError(t, stdErr)
	require.EqualValues(t, statistics.StatusStopped, result.Status)
}

func Test_EnableProjection(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	err := client.EnableProjection(context.Background(), StandardProjectionStreams)
	require.NoError(t, err)

	statisticsClient, err := client.GetProjectionStatistics(context.Background(), projections.StatisticsForProjectionByName{
		ProjectionName: StandardProjectionStreams,
	})
	require.NoError(t, err)
	require.NotNil(t, statisticsClient)

	result, stdErr := statisticsClient.Read()
	require.NoError(t, stdErr)
	require.EqualValues(t, statistics.StatusRunning, result.Status)
}

func Test_GetResultOfProjection(t *testing.T) {
	client, eventStreamsClient, closeFunc := initializeClientAndEventStreamsClient(t)
	defer closeFunc()

	streamName := "result_test_stream"

	const resultProjectionQuery = `
		fromStream('%s').when(
			{
				$init: function (state, ev) {
					return {
						count: 0
					};
				},
				$any: function(s, e) {
					s.count += 1;
				}
			}
		);
		`

	const eventData = `{
		message: "test",
		index: %d,
	}`

	createOptions := projections.CreateRequest{
		Mode: projections.ContinuousProjection{
			ProjectionName:      "MyContinuousProjection",
			TrackEmittedStreams: false,
		},
		Query: fmt.Sprintf(resultProjectionQuery, streamName),
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)

	testEvent := testCreateEvent()
	testEvent.EventType = "count_this"
	testEvent.Data = []byte(fmt.Sprintf(eventData, 0))

	pushEventsToStream(t, eventStreamsClient, streamName, testEvent)

	var projectionResult projections.ResultResponse

	require.Eventually(t, func() bool {
		resultOptions := projections.ResultOptionsRequest{
			ProjectionName: "MyContinuousProjection",
		}
		projectionResult, err = client.
			GetProjectionResult(context.Background(), resultOptions)
		return err == nil && projectionResult.GetType() == projections.ResultResponseStructType
	}, time.Second*5, time.Millisecond*500)

	type resultStructType struct {
		Count int64 `json:"count"`
	}

	var result resultStructType

	structResult := projectionResult.(*projections.ResultResponseStruct)

	stdErr := json.Unmarshal(structResult.Value(), &result)
	require.NoError(t, stdErr)
	require.EqualValues(t, 1, result.Count)
}

func Test_GetStateOfProjection(t *testing.T) {
	client, eventStreamsClient, closeFunc := initializeClientAndEventStreamsClient(t)
	defer closeFunc()

	streamName := "state_test_stream"

	const resultProjectionQuery = `
		fromStream('%s').when(
			{
				$init: function (state, ev) {
					return {
						count: 0
					};
				},
				$any: function(s, e) {
					s.count += 1;
				}
			}
		);
		`

	const eventData = `{
		message: "test",
		index: %d,
	}`

	createOptions := projections.CreateRequest{
		Mode: projections.ContinuousProjection{
			ProjectionName:      "MyContinuousProjection",
			TrackEmittedStreams: false,
		},
		Query: fmt.Sprintf(resultProjectionQuery, streamName),
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)

	testEvent := testCreateEvent()
	testEvent.EventType = "count_this"
	testEvent.Data = []byte(fmt.Sprintf(eventData, 0))

	pushEventsToStream(t, eventStreamsClient, streamName, testEvent)

	var projectionResult projections.StateResponse
	require.Eventually(t, func() bool {
		resultOptions := projections.StateOptionsRequest{
			ProjectionName: "MyContinuousProjection",
		}
		projectionResult, err = client.
			GetProjectionState(context.Background(), resultOptions)
		return err == nil && projectionResult.GetType() == projections.StateResponseStructType
	}, time.Second*5, time.Millisecond*500)

	type resultStructType struct {
		Count int64 `json:"count"`
	}

	var result resultStructType

	structResult := projectionResult.(*projections.StateResponseStruct)

	stdErr := json.Unmarshal(structResult.Value(), &result)
	require.NoError(t, stdErr)
	require.EqualValues(t, 1, result.Count)
}

func Test_GetStatusOfProjection(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	statisticsClient, err := client.GetProjectionStatistics(context.Background(), projections.StatisticsForProjectionByName{
		ProjectionName: StandardProjectionStreams,
	})
	require.NoError(t, err)
	require.NotNil(t, statisticsClient)

	result, stdErr := statisticsClient.Read()
	require.NoError(t, stdErr)
	require.EqualValues(t, StandardProjectionStreams, result.Name)
}

func Test_ResetProjection(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	resetOptions := projections.ResetOptionsRequest{
		ProjectionName: StandardProjectionStreams,
	}
	err := client.ResetProjection(context.Background(), resetOptions)
	require.NoError(t, err)

	statisticsClient, err := client.
		GetProjectionStatistics(context.Background(), projections.StatisticsForProjectionByName{
			ProjectionName: StandardProjectionStreams,
		})
	require.NoError(t, err)
	require.NotNil(t, statisticsClient)

	result, stdErr := statisticsClient.Read()
	require.NoError(t, stdErr)
	require.EqualValues(t, statistics.StatusRunning, result.Status)
}

func Test_RestartProjectionSubsystem(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	err := client.RestartProjectionsSubsystem(context.Background())
	require.NoError(t, err)
}

func Test_ListAllProjections(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	expectedStreamNames := []string{
		StandardProjectionStreams,
		StandardProjectionStreamByCategory,
		StandardProjectionByCategory,
		StandardProjectionByEventType,
		StandardProjectionByCorrelationId,
	}

	result, err := client.ListAllProjections(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, len(expectedStreamNames))

	var resultNames []string
	for _, resultItem := range result {
		resultNames = append(resultNames, resultItem.Name)
	}

	require.ElementsMatch(t, expectedStreamNames, resultNames)
}

func Test_ListContinuousProjections(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode: projections.ContinuousProjection{
			ProjectionName:      "MyContinuous_false",
			TrackEmittedStreams: false,
		},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)

	expectedStreamNames := []string{
		StandardProjectionStreams,
		StandardProjectionStreamByCategory,
		StandardProjectionByCategory,
		StandardProjectionByEventType,
		StandardProjectionByCorrelationId,
		"MyContinuous_false",
	}

	result, err := client.ListContinuousProjections(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, len(expectedStreamNames))

	var resultNames []string
	for _, resultItem := range result {
		resultNames = append(resultNames, resultItem.Name)
	}

	require.ElementsMatch(t, expectedStreamNames, resultNames)
}

func Test_ListOneTimeProjections(t *testing.T) {
	client, closeFunc := initializeContainerAndClient(t)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode:  projections.OneTimeProjection{},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.NoError(t, err)

	result, err := client.ListOneTimeProjections(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 1)

	require.Equal(t, statistics.ModeOneTime, result[0].Mode)
}

func Test_CreateProjection_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	createOptions := projections.CreateRequest{
		Mode:  projections.OneTimeProjection{},
		Query: "fromAll().when({$init: function (state, ev) {return {};}});",
	}

	err := client.CreateProjection(context.Background(), createOptions)
	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_UpdateProjection_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	updateOptions := projections.UpdateRequest{
		EmitOption:     projections.NoEmit{},
		Query:          "fromAll().when({$init: function (s, e) {return {};}});",
		ProjectionName: "MyTransient_no_emit",
	}

	err := client.UpdateProjection(context.Background(), updateOptions)
	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_AbortProjection_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	err := client.AbortProjection(context.Background(), StandardProjectionStreams)
	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_DisableProjection_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	err := client.DisableProjection(context.Background(), StandardProjectionStreams)
	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_EnableProjection_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	err := client.EnableProjection(context.Background(), StandardProjectionStreams)
	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_GetProjectionResult_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	resultOptions := projections.ResultOptionsRequest{
		ProjectionName: "MyContinuousProjection",
	}
	_, err := client.
		GetProjectionResult(context.Background(), resultOptions)
	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_GetProjectionState_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	resultOptions := projections.StateOptionsRequest{
		ProjectionName: "MyContinuousProjection",
	}
	_, err := client.
		GetProjectionState(context.Background(), resultOptions)
	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_GetProjectionStatistics_WithIncorrectCredentials_DoesNotFail(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	_, err := client.GetProjectionStatistics(context.Background(), projections.StatisticsForProjectionByName{
		ProjectionName: StandardProjectionStreams,
	})
	require.NoError(t, err)
}

func Test_ResetProjection_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	resetOptions := projections.ResetOptionsRequest{
		ProjectionName: StandardProjectionStreams,
	}
	err := client.ResetProjection(context.Background(), resetOptions)
	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_RestartProjectionsSubsystem_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	err := client.RestartProjectionsSubsystem(context.Background())

	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_ListAllProjections_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	_, err := client.ListAllProjections(context.Background())

	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_ListContinuousProjections_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	_, err := client.ListContinuousProjections(context.Background())

	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

func Test_ListOneTimeProjections_WithIncorrectCredentials(t *testing.T) {
	client, closeFunc := initializeContainerAndClientWithCredentials(t,
		"wrong_user_name", "wrong_password", nil)
	defer closeFunc()

	_, err := client.ListOneTimeProjections(context.Background())

	require.Equal(t, errors.UnauthenticatedErr, err.Code())
}

const (
	StandardProjectionStreams          = "$streams"
	StandardProjectionStreamByCategory = "$stream_by_category"
	StandardProjectionByCategory       = "$by_category"
	StandardProjectionByEventType      = "$by_event_type"
	StandardProjectionByCorrelationId  = "$by_correlation_id"
)
