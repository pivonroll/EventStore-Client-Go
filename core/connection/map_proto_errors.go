package connection

import (
	"github.com/pivonroll/EventStore-Client-Go/core/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	protoStreamDeleted                 = "stream-deleted"
	protoStreamNotFound                = "stream-not-found"
	protoMaximumAppendSizeExceeded     = "maximum-append-size-exceeded"
	protoWrongExpectedVersion          = "wrong-expected-version"
	protoNotLeader                     = "not-leader"
	protoUserNotFound                  = "user-not-found"
	protoMaximumSubscriberCountReached = "maximum-subscribers-reached"
	protoPersistentSubscriptionDropped = "persistent-subscription-dropped"
	protoScavengeNotFound              = "scavenge-not-found"
)

func ErrorFromStdErrorByStatus(err error) errors.Error {
	protoStatus, _ := status.FromError(err)
	if protoStatus.Code() == codes.PermissionDenied {
		return errors.NewError(errors.PermissionDeniedErr, err)
	} else if protoStatus.Code() == codes.Unauthenticated {
		return errors.NewError(errors.UnauthenticatedErr, err)
	} else if protoStatus.Code() == codes.DeadlineExceeded {
		return errors.NewError(errors.DeadlineExceededErr, err)
	} else if protoStatus.Code() == codes.Canceled {
		return errors.NewError(errors.CanceledErr, err)
	}
	return nil
}

func GetErrorFromProtoException(trailers metadata.MD, stdErr error) errors.Error {
	if isProtoException(trailers, protoStreamDeleted) {
		return errors.NewError(errors.StreamDeletedErr, stdErr)
	} else if isProtoException(trailers, protoMaximumAppendSizeExceeded) {
		return errors.NewError(errors.MaximumAppendSizeExceededErr, stdErr)
	} else if isProtoException(trailers, protoStreamNotFound) {
		return errors.NewError(errors.StreamNotFoundErr, stdErr)
	} else if isProtoException(trailers, protoWrongExpectedVersion) {
		return errors.NewError(errors.WrongExpectedStreamRevisionErr, stdErr)
	} else if isProtoException(trailers, protoNotLeader) {
		return errors.NewError(errors.NotLeaderErr, stdErr)
	} else if isProtoException(trailers, protoUserNotFound) {
		return errors.NewError(errors.UserNotFoundErr, stdErr)
	} else if isProtoException(trailers, protoMaximumSubscriberCountReached) {
		return errors.NewError(errors.MaximumSubscriberCountReached, stdErr)
	} else if isProtoException(trailers, protoPersistentSubscriptionDropped) {
		return errors.NewError(errors.PersistentSubscriptionDroppedErr, stdErr)
	} else if isProtoException(trailers, protoScavengeNotFound) {
		return errors.NewError(errors.ScavengeNotFoundErr, stdErr)
	}

	err := ErrorFromStdErrorByStatus(stdErr)
	if err != nil {
		return err
	}

	return nil
}

func isProtoException(trailers metadata.MD, protoException string) bool {
	values := trailers.Get("exception")
	return values != nil && values[0] == protoException
}
