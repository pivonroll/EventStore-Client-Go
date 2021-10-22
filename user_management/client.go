// Header

// Package user_management provides user management capabilities for EventStoreDb.
// A grpc connection needs to be established with EventStore through
// github.com/pivonroll/EventStore-Client-Go/core/connection package.
package user_management

import (
	"context"

	"github.com/pivonroll/EventStore-Client-Go/core/connection"
	"github.com/pivonroll/EventStore-Client-Go/core/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Client is used to manage users on EventStoreDB.
type Client struct {
	grpcClient            connection.GrpcClient
	grpcUserClientFactory grpcUserClientFactory
	detailsReaderFactory  detailsReaderFactory
}

// NewClient returns a new instance of a user management client.
func NewClient(grpcClient connection.GrpcClient) *Client {
	return &Client{
		grpcClient:            grpcClient,
		grpcUserClientFactory: grpcUserClientFactoryImpl{},
		detailsReaderFactory:  detailsReaderFactoryImpl{},
	}
}

// CreateUser creates a new user on EventStoreDB.
// You must have necessary access rights to do this operation.
func (client *Client) CreateUser(ctx context.Context, request CreateOrUpdateRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcUsersClient := client.grpcUserClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcUsersClient.Create(ctx, request.buildCreateRequest(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return err
	}

	return nil
}

// UpdateUser updates an existing user on EventStoreDB.
// You must have necessary access rights to do this operation.
func (client *Client) UpdateUser(ctx context.Context, request CreateOrUpdateRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcUsersClient := client.grpcUserClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcUsersClient.Update(ctx, request.buildUpdateRequest(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return err
	}

	return nil
}

// DeleteUser removes an existing user from EventStoreDB.
// You must have necessary access rights to do this operation.
func (client *Client) DeleteUser(ctx context.Context, loginName string) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcUsersClient := client.grpcUserClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcUsersClient.Delete(ctx, DeleteRequest(loginName).build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return err
	}

	return nil
}

// DisableUser disables an existing user on EventStoreDB.
// You must have necessary access rights to do this operation.
func (client *Client) DisableUser(ctx context.Context, loginName string) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcUsersClient := client.grpcUserClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcUsersClient.Disable(ctx, DisableRequest(loginName).build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return err
	}

	return nil
}

// EnableUser enables a previously disabled user on EventStoreDB.
// You must have necessary access rights to do this operation.
func (client *Client) EnableUser(ctx context.Context, loginName string) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcUsersClient := client.grpcUserClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcUsersClient.Enable(ctx, EnableRequest(loginName).build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return err
	}

	return nil
}

// GetUserDetails fetches details of the existing user from EventStoreDB.
// You must have necessary access rights to do this operation.
func (client *Client) GetUserDetails(ctx context.Context, loginName string) (DetailsResponse, errors.Error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return DetailsResponse{}, err
	}

	grpcUsersClient := client.grpcUserClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	ctx, cancel := context.WithCancel(ctx)
	protoStreamReader, protoErr := grpcUsersClient.Details(ctx, DetailsRequest(loginName).build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		defer cancel()
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return DetailsResponse{}, err
	}

	userDetailsReader := client.detailsReaderFactory.Create(protoStreamReader, cancel)

	return userDetailsReader.Recv()
}

// ListAllUsers lists details of all users on EventStoreDB.
// You must have necessary access rights to do this operation.
func (client *Client) ListAllUsers(ctx context.Context) ([]DetailsResponse, errors.Error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	grpcUsersClient := client.grpcUserClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	ctx, cancel := context.WithCancel(ctx)
	protoStreamReader, protoErr := grpcUsersClient.Details(ctx, DetailsRequest(AllUsers).build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		defer cancel()
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return nil, err
	}

	userDetailsReader := client.detailsReaderFactory.Create(protoStreamReader, cancel)

	var result []DetailsResponse
	for {
		response, readerErr := userDetailsReader.Recv()
		if readerErr != nil {
			if readerErr.Code() == errors.EndOfStream {
				break
			}
			return nil, readerErr
		}

		result = append(result, response)
	}

	return result, nil
}

// ChangeUserPassword changes a password of a user on EventStoreDB.
// You must have necessary access rights to do this operation.
func (client *Client) ChangeUserPassword(ctx context.Context, request ChangePasswordRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcUsersClient := client.grpcUserClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcUsersClient.ChangePassword(ctx, request.build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return err
	}

	return nil
}

// ResetUserPassword resets user's password on EventStoreDB.
// You must have necessary access rights to do this operation.
func (client *Client) ResetUserPassword(ctx context.Context, loginName string, newPassword string) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcUsersClient := client.grpcUserClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcUsersClient.ResetPassword(ctx, ResetPasswordRequest{
		LoginName:   loginName,
		NewPassword: newPassword,
	}.build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return err
	}

	return nil
}
