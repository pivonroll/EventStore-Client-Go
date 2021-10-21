package projections

import (
	"context"
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/errors"
	"github.com/pivonroll/EventStore-Client-Go/projections/internal/grpc_proto_client_factory"
	statistics_internal "github.com/pivonroll/EventStore-Client-Go/projections/internal/statistics"
	"github.com/pivonroll/EventStore-Client-Go/projections/statistics"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/pivonroll/EventStore-Client-Go/protos/persistent"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/pivonroll/EventStore-Client-Go/connection"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

func TestClientImpl_CreateProjection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	grpcClientConn := &grpc.ClientConn{}
	options := CreateOptionsRequest{}
	options.SetQuery("some query")
	options.SetMode(CreateConfigModeContinuousOption{
		Name:                "some mode",
		TrackEmittedStreams: true,
	})

	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD
		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Create(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(nil, nil),
		)
		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.CreateProjection(ctx, options)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		err := client.CreateProjection(ctx, options)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Create(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.CreateReq,
					options ...grpc.CallOption) (persistent.PersistentSubscriptions_ReadClient, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
			grpcClient:                   grpcClient,
		}

		err := client.CreateProjection(ctx, options)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_UpdateProjection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	grpcClientConn := &grpc.ClientConn{}
	options := UpdateOptionsRequest{}
	options.SetName("some name").
		SetQuery("some query").
		SetEmitOption(UpdateOptionsEmitOptionEnabled{
			EmitEnabled: true,
		})

	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Update(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(nil, nil),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.UpdateProjection(ctx, options)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		err := client.UpdateProjection(ctx, options)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Update(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.UpdateReq,
					options ...grpc.CallOption) (*persistent.UpdateResp, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.UpdateProjection(ctx, options)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_DeleteProjection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	grpcClientConn := &grpc.ClientConn{}
	options := DeleteOptionsRequest{}
	options.SetName("some name")

	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Delete(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(nil, nil),
		)
		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.DeleteProjection(ctx, options)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		err := client.DeleteProjection(ctx, options)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Delete(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.DeleteReq,
					options ...grpc.CallOption) (*persistent.DeleteResp, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.DeleteProjection(ctx, options)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_ProjectionStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	grpcClientConn := &grpc.ClientConn{}
	options := StatisticsOptionsRequest{}
	options.SetMode(StatisticsOptionsRequestModeAll{})

	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)
		statisticsClientSyncFactoryInstance := statistics_internal.NewMockClientSyncFactory(ctrl)
		statisticsClientSyncRead := statistics_internal.NewMockClientSync(ctrl)
		statisticsClient := projections.NewMockProjections_StatisticsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(statisticsClient, nil),
			statisticsClientSyncFactoryInstance.EXPECT().Create(statisticsClient).
				Return(statisticsClientSyncRead),
		)
		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
			statisticsClientSyncFactory:  statisticsClientSyncFactoryInstance,
		}

		result, err := client.GetProjectionStatistics(ctx, options)
		require.Equal(t, statisticsClientSyncRead, result)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		_, err := client.GetProjectionStatistics(ctx, options)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.StatisticsReq,
					options ...grpc.CallOption) (projections.Projections_StatisticsClient, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
			grpcClient:                   grpcClient,
		}

		statisticsClient, err := client.GetProjectionStatistics(ctx, options)
		require.Nil(t, statisticsClient)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_DisableProjection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	grpcClientConn := &grpc.ClientConn{}
	projectionName := "some name"
	options := disableOptionsRequest(projectionName)
	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Disable(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(nil, nil),
		)
		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.DisableProjection(ctx, projectionName)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		err := client.DisableProjection(ctx, projectionName)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Disable(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.DisableReq,
					options ...grpc.CallOption) (*projections.DisableResp, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
			grpcClient:                   grpcClient,
		}

		err := client.DisableProjection(ctx, projectionName)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_AbortProjection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	projectionName := "some name"
	grpcClientConn := &grpc.ClientConn{}
	options := abortOptionsRequest(projectionName)
	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Disable(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(nil, nil),
		)
		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.AbortProjection(ctx, projectionName)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		err := client.AbortProjection(ctx, projectionName)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Disable(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.DisableReq,
					options ...grpc.CallOption) (*projections.DisableResp, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.AbortProjection(ctx, projectionName)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_EnableProjection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	grpcClientConn := &grpc.ClientConn{}
	options := EnableOptionsRequest{}
	options.SetName("some name")
	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Enable(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(nil, nil),
		)
		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.EnableProjection(ctx, options)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		err := client.EnableProjection(ctx, options)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Enable(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.EnableReq,
					options ...grpc.CallOption) (*projections.EnableResp, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}
		err := client.EnableProjection(ctx, options)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_ResetProjection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	grpcClientConn := &grpc.ClientConn{}
	options := ResetOptionsRequest{}
	options.SetName("some name").SetWriteCheckpoint(true)
	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Reset(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(nil, nil),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.ResetProjection(ctx, options)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		err := client.ResetProjection(ctx, options)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Reset(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.ResetReq,
					options ...grpc.CallOption) (*projections.ResetResp, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.ResetProjection(ctx, options)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_ProjectionState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	grpcClientConn := &grpc.ClientConn{}
	options := StateOptionsRequest{}
	options.SetName("some name").SetPartition("some partition")
	grpcOptions := options.build()

	response := &projections.StateResp{
		State: &structpb.Value{
			Kind: &structpb.Value_NullValue{},
		},
	}

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().State(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(response, nil),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		expectedStateResponse := newStateResponse(response)

		state, err := client.GetProjectionState(ctx, options)
		require.Equal(t, expectedStateResponse, state)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		_, err := client.GetProjectionState(ctx, options)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().State(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.StateReq,
					options ...grpc.CallOption) (*projections.StateResp, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		state, err := client.GetProjectionState(ctx, options)
		require.Nil(t, state)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_ProjectionResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	grpcClientConn := &grpc.ClientConn{}
	options := ResultOptionsRequest{}
	options.SetName("some name").SetPartition("some partition")
	grpcOptions := options.build()

	response := &projections.ResultResp{
		Result: &structpb.Value{
			Kind: &structpb.Value_NullValue{},
		},
	}

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Result(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(response, nil),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		expectedStateResponse := newResultResponse(response)

		state, err := client.GetProjectionResult(ctx, options)
		require.Equal(t, expectedStateResponse, state)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		_, err := client.GetProjectionResult(ctx, options)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Result(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.ResultReq,
					options ...grpc.CallOption) (*projections.ResultResp, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		state, err := client.GetProjectionResult(ctx, options)
		require.Nil(t, state)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_RestartProjectionsSubsystem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpcClientConn := &grpc.ClientConn{}
	ctx := context.Background()

	t.Run("Success with RestartSubsystem return nil", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().RestartSubsystem(ctx, &shared.Empty{},
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(nil, nil),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.RestartProjectionsSubsystem(ctx)
		require.NoError(t, err)
	})

	t.Run("Success with RestartSubsystem return &shared.Empty{}", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().RestartSubsystem(ctx, &shared.Empty{},
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(&shared.Empty{}, nil),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.RestartProjectionsSubsystem(ctx)
		require.NoError(t, err)
	})

	t.Run("Error when obtaining connection", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		errorResult := errors.NewErrorCode("some error")
		grpcClient.EXPECT().GetConnectionHandle().Return(nil, errorResult)

		client := Client{
			grpcClient: grpcClient,
		}

		err := client.RestartProjectionsSubsystem(ctx)
		require.Equal(t, errorResult, err)
	})

	t.Run("Error returned from grpc", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().RestartSubsystem(ctx, &shared.Empty{},
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *shared.Empty,
					options ...grpc.CallOption) (*shared.Empty, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(
				errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		err := client.RestartProjectionsSubsystem(ctx)
		require.Equal(t, errors.FatalError, err.Code())
	})
}

func TestClientImpl_ListAllProjections(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	grpcClientConn := &grpc.ClientConn{}
	options := StatisticsOptionsRequest{}
	options.SetMode(StatisticsOptionsRequestModeAll{})

	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)
		statisticsClientSyncFactoryInstance := statistics_internal.NewMockClientSyncFactory(ctrl)
		statisticsClientSyncRead := statistics_internal.NewMockClientSync(ctrl)
		statisticsClient := projections.NewMockProjections_StatisticsClient(ctrl)

		responseList := []statistics.Response{
			{
				Name: "Some name 1",
			},
			{
				Name: "Some name 2",
			},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(statisticsClient, nil),
			statisticsClientSyncFactoryInstance.EXPECT().Create(statisticsClient).
				Return(statisticsClientSyncRead),
			statisticsClientSyncRead.EXPECT().Read().Return(responseList[0], nil),
			statisticsClientSyncRead.EXPECT().Read().Return(responseList[1], nil),
			statisticsClientSyncRead.EXPECT().Read().Return(statistics.Response{},
				errors.NewErrorCode(errors.EndOfStream)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
			statisticsClientSyncFactory:  statisticsClientSyncFactoryInstance,
		}

		allProjectionsResult, err := client.ListAllProjections(ctx)
		require.NoError(t, err)
		require.NotNil(t, allProjectionsResult)
		require.Len(t, allProjectionsResult, len(responseList))
	})

	t.Run("Error returned from Statistics", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.StatisticsReq,
					options ...grpc.CallOption) (projections.Projections_StatisticsClient, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(
				errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		allProjectionsResult, err := client.ListAllProjections(ctx)
		require.Error(t, err)
		require.Nil(t, allProjectionsResult)
		require.Equal(t, errors.FatalError, err.Code())
	})

	t.Run("Error returned from statistics client read", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)
		statisticsClient := projections.NewMockProjections_StatisticsClient(ctrl)
		statisticsClientSyncFactoryInstance := statistics_internal.NewMockClientSyncFactory(ctrl)
		statisticsClientSyncRead := statistics_internal.NewMockClientSync(ctrl)

		var headers, trailers metadata.MD
		errorCode := errors.FatalError

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(statisticsClient, nil),
			statisticsClientSyncFactoryInstance.EXPECT().Create(statisticsClient).
				Return(statisticsClientSyncRead),
			statisticsClientSyncRead.EXPECT().Read().Return(statistics.Response{},
				errors.NewErrorCode(errorCode)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
			statisticsClientSyncFactory:  statisticsClientSyncFactoryInstance,
		}

		allProjectionsResult, err := client.ListAllProjections(ctx)
		require.Error(t, err)
		require.Equal(t, errorCode, err.Code())
		require.Nil(t, allProjectionsResult)
	})
}

func TestClientImpl_ListContinuousProjections(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpcClientConn := &grpc.ClientConn{}
	ctx := context.Background()

	options := StatisticsOptionsRequest{}
	options.SetMode(StatisticsOptionsRequestModeContinuous{})

	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)
		statisticsClientSyncFactoryInstance := statistics_internal.NewMockClientSyncFactory(ctrl)
		statisticsClientSyncRead := statistics_internal.NewMockClientSync(ctrl)
		statisticsClient := projections.NewMockProjections_StatisticsClient(ctrl)

		var headers, trailers metadata.MD
		responseList := []statistics.Response{
			{
				Name: "Some name 1",
			},
			{
				Name: "Some name 2",
			},
		}

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(statisticsClient, nil),
			statisticsClientSyncFactoryInstance.EXPECT().Create(statisticsClient).
				Return(statisticsClientSyncRead),
			statisticsClientSyncRead.EXPECT().Read().Return(responseList[0], nil),
			statisticsClientSyncRead.EXPECT().Read().Return(responseList[1], nil),
			statisticsClientSyncRead.EXPECT().Read().Return(statistics.Response{},
				errors.NewErrorCode(errors.EndOfStream)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
			statisticsClientSyncFactory:  statisticsClientSyncFactoryInstance,
		}

		projectionsResult, err := client.ListContinuousProjections(ctx)
		require.NoError(t, err)
		require.NotNil(t, projectionsResult)
		require.Len(t, projectionsResult, len(responseList))
	})

	t.Run("Error returned from Statistics", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.StatisticsReq,
					options ...grpc.CallOption) (projections.Projections_StatisticsClient, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(
				errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		allProjectionsResult, err := client.ListContinuousProjections(ctx)
		require.Error(t, err)
		require.Nil(t, allProjectionsResult)
		require.Equal(t, errors.FatalError, err.Code())
	})

	t.Run("Error returned from statistics client read", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)
		statisticsClient := projections.NewMockProjections_StatisticsClient(ctrl)
		statisticsClientSyncFactoryInstance := statistics_internal.NewMockClientSyncFactory(ctrl)
		statisticsClientSyncRead := statistics_internal.NewMockClientSync(ctrl)

		var headers, trailers metadata.MD

		errorCode := errors.FatalError
		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(statisticsClient, nil),
			statisticsClientSyncFactoryInstance.EXPECT().Create(statisticsClient).
				Return(statisticsClientSyncRead),
			statisticsClientSyncRead.EXPECT().Read().Return(statistics.Response{},
				errors.NewErrorCode(errorCode)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
			statisticsClientSyncFactory:  statisticsClientSyncFactoryInstance,
		}

		allProjectionsResult, err := client.ListContinuousProjections(ctx)
		require.Error(t, err)
		require.Equal(t, errorCode, err.Code())
		require.Nil(t, allProjectionsResult)
	})
}

func TestClientImpl_ListOneTimeProjections(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	grpcClientConn := &grpc.ClientConn{}
	ctx := context.Background()

	options := StatisticsOptionsRequest{}
	options.SetMode(StatisticsOptionsRequestModeOneTime{})

	grpcOptions := options.build()

	t.Run("Success", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)
		statisticsClientSyncFactoryInstance := statistics_internal.NewMockClientSyncFactory(ctrl)
		statisticsClientSyncRead := statistics_internal.NewMockClientSync(ctrl)
		statisticsClient := projections.NewMockProjections_StatisticsClient(ctrl)

		var headers, trailers metadata.MD
		responseList := []statistics.Response{
			{
				Name: "Some name 1",
			},
			{
				Name: "Some name 2",
			},
		}

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(statisticsClient, nil),
			statisticsClientSyncFactoryInstance.EXPECT().Create(statisticsClient).
				Return(statisticsClientSyncRead),
			statisticsClientSyncRead.EXPECT().Read().Return(responseList[0], nil),
			statisticsClientSyncRead.EXPECT().Read().Return(responseList[1], nil),
			statisticsClientSyncRead.EXPECT().Read().Return(statistics.Response{},
				errors.NewErrorCode(errors.EndOfStream)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
			statisticsClientSyncFactory:  statisticsClientSyncFactoryInstance,
		}

		projectionsResult, err := client.ListOneTimeProjections(ctx)
		require.NoError(t, err)
		require.NotNil(t, projectionsResult)
		require.Len(t, projectionsResult, len(responseList))
	})

	t.Run("Error returned from Statistics", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)

		errorResult := errors.NewErrorCode("some error")
		expectedHeader := metadata.MD{
			"header_key": []string{"header_value"},
		}

		expectedTrailer := metadata.MD{
			"trailer_key": []string{"trailer_value"},
		}

		var headers, trailers metadata.MD

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).
				DoAndReturn(func(
					_ctx context.Context,
					_protoRequest *projections.StatisticsReq,
					options ...grpc.CallOption) (projections.Projections_StatisticsClient, error) {

					*options[0].(grpc.HeaderCallOption).HeaderAddr = metadata.MD{
						"header_key": []string{"header_value"},
					}

					*options[1].(grpc.TrailerCallOption).TrailerAddr = metadata.MD{
						"trailer_key": []string{"trailer_value"},
					}
					return nil, errorResult
				}),
			grpcClient.EXPECT().HandleError(handle, expectedHeader, expectedTrailer, errorResult,
				errors.FatalError).Return(errors.NewErrorCode(errors.FatalError)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
		}

		allProjectionsResult, err := client.ListOneTimeProjections(ctx)
		require.Error(t, err)
		require.Nil(t, allProjectionsResult)
		require.Equal(t, errors.FatalError, err.Code())
	})

	t.Run("Error returned from statistics client read", func(t *testing.T) {
		grpcClient := connection.NewMockGrpcClient(ctrl)
		handle := connection.NewMockConnectionHandle(ctrl)
		grpcProjectionsClientFactoryInstance := grpc_proto_client_factory.NewMockFactory(ctrl)
		grpcProjectionsClientMock := projections.NewMockProjectionsClient(ctrl)
		statisticsClientSyncFactoryInstance := statistics_internal.NewMockClientSyncFactory(ctrl)
		statisticsClientSyncRead := statistics_internal.NewMockClientSync(ctrl)
		statisticsClient := projections.NewMockProjections_StatisticsClient(ctrl)

		var headers, trailers metadata.MD
		errorCode := errors.FatalError

		gomock.InOrder(
			grpcClient.EXPECT().GetConnectionHandle().Return(handle, nil),
			handle.EXPECT().Connection().Return(grpcClientConn),
			grpcProjectionsClientFactoryInstance.EXPECT().Create(grpcClientConn).
				Return(grpcProjectionsClientMock),
			grpcProjectionsClientMock.EXPECT().Statistics(ctx, grpcOptions,
				grpc.Header(&headers), grpc.Trailer(&trailers)).Times(1).Return(statisticsClient, nil),
			statisticsClientSyncFactoryInstance.EXPECT().Create(statisticsClient).
				Return(statisticsClientSyncRead),
			statisticsClientSyncRead.EXPECT().Read().Return(statistics.Response{},
				errors.NewErrorCode(errorCode)),
		)

		client := Client{
			grpcClient:                   grpcClient,
			grpcProjectionsClientFactory: grpcProjectionsClientFactoryInstance,
			statisticsClientSyncFactory:  statisticsClientSyncFactoryInstance,
		}

		allProjectionsResult, err := client.ListOneTimeProjections(ctx)
		require.Error(t, err)
		require.Equal(t, errorCode, err.Code())
		require.Nil(t, allProjectionsResult)
	})
}
