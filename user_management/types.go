package user_management

import (
	"strings"
	"time"

	"github.com/pivonroll/EventStore-Client-Go/protos/users"
)

// CreateOrUpdateRequest represents the data necessary to create a user on EventStoreDB.
type CreateOrUpdateRequest struct {
	LoginName string
	Password  string
	FullName  string
	Groups    []string
}

func (request CreateOrUpdateRequest) buildCreateRequest() *users.CreateReq {
	return &users.CreateReq{
		Options: &users.CreateReq_Options{
			LoginName: request.LoginName,
			Password:  request.Password,
			FullName:  request.FullName,
			Groups:    request.Groups,
		},
	}
}

func (request CreateOrUpdateRequest) buildUpdateRequest() *users.UpdateReq {
	return &users.UpdateReq{
		Options: &users.UpdateReq_Options{
			LoginName: request.LoginName,
			Password:  request.Password,
			FullName:  request.FullName,
			Groups:    request.Groups,
		},
	}
}

// DeleteRequest represents a login name of the user we want to remove from EventStoreDB.
type DeleteRequest string

func (request DeleteRequest) build() *users.DeleteReq {
	return &users.DeleteReq{
		Options: &users.DeleteReq_Options{
			LoginName: string(request),
		},
	}
}

// DisableRequest represents a login name of the enabled user which we want to disable on EventStoreDB.
type DisableRequest string

func (request DisableRequest) build() *users.DisableReq {
	return &users.DisableReq{
		Options: &users.DisableReq_Options{
			LoginName: string(request),
		},
	}
}

// EnableRequest represents a login name of the disabled user which we want to enable on EventStoreDB.
type EnableRequest string

func (request EnableRequest) build() *users.EnableReq {
	return &users.EnableReq{
		Options: &users.EnableReq_Options{
			LoginName: string(request),
		},
	}
}

// DetailsRequest represents a login name of the user which details we want to fetch from EventStoreDB.
type DetailsRequest string

const AllUsers = ""

func (request DetailsRequest) build() *users.DetailsReq {
	if (strings.TrimSpace(string(request))) == "" {
		return &users.DetailsReq{}
	}
	return &users.DetailsReq{
		Options: &users.DetailsReq_Options{
			LoginName: string(request),
		},
	}
}

// DetailsResponse are fetched details of one user from EventStoreDB.
type DetailsResponse struct {
	LoginName   string
	FullName    string
	Groups      []string
	LastUpdated time.Time
	Disabled    bool
}

type detailsResponseAdapter interface {
	Create(proto *users.DetailsResp) DetailsResponse
}

type detailsResponseAdapterImpl struct{}

func (adapter detailsResponseAdapterImpl) Create(proto *users.DetailsResp) DetailsResponse {
	return DetailsResponse{
		LoginName:   proto.UserDetails.LoginName,
		FullName:    proto.UserDetails.FullName,
		Groups:      proto.UserDetails.Groups,
		LastUpdated: time.Unix(0, proto.UserDetails.LastUpdated.TicksSinceEpoch*100).UTC(),
		Disabled:    proto.UserDetails.Disabled,
	}
}

// ChangePasswordRequest are data required for password change of a user.
type ChangePasswordRequest struct {
	LoginName       string
	CurrentPassword string
	NewPassword     string
}

func (request ChangePasswordRequest) build() *users.ChangePasswordReq {
	return &users.ChangePasswordReq{
		Options: &users.ChangePasswordReq_Options{
			LoginName:       request.LoginName,
			CurrentPassword: request.CurrentPassword,
			NewPassword:     request.NewPassword,
		},
	}
}

// ResetPasswordRequest represent a set of data required to reset user's password.
type ResetPasswordRequest struct {
	LoginName   string
	NewPassword string
}

func (request ResetPasswordRequest) build() *users.ResetPasswordReq {
	return &users.ResetPasswordReq{
		Options: &users.ResetPasswordReq_Options{
			LoginName:   request.LoginName,
			NewPassword: request.NewPassword,
		},
	}
}
