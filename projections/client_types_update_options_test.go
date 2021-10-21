package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"

	"github.com/stretchr/testify/require"
)

func TestUpdateOptionsRequest_Build(t *testing.T) {
	t.Run("Build with name, query and emit option set to no emit", func(t *testing.T) {
		options := UpdateOptionsRequest{
			EmitOption: UpdateOptionsEmitOptionNoEmit{},
			Query:      "some query",
			Name:       "some name",
		}

		result := options.build()

		expectedResult := &projections.UpdateReq{
			Options: &projections.UpdateReq_Options{
				Name:  "some name",
				Query: "some query",
				EmitOption: &projections.UpdateReq_Options_NoEmitOptions{
					NoEmitOptions: &shared.Empty{},
				},
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Build with name, query and emit option set to EmitOptionEnabled with EmitEnabled set to false",
		func(t *testing.T) {
			options := UpdateOptionsRequest{
				EmitOption: UpdateOptionsEmitOptionEnabled{
					EmitEnabled: false,
				},
				Query: "some query",
				Name:  "some name",
			}
			result := options.build()

			expectedResult := &projections.UpdateReq{
				Options: &projections.UpdateReq_Options{
					Name:  "some name",
					Query: "some query",
					EmitOption: &projections.UpdateReq_Options_EmitEnabled{
						EmitEnabled: false,
					},
				},
			}

			require.Equal(t, expectedResult, result)
		})

	t.Run("Build with name, query and emit option set to EmitOptionEnabled with EmitEnabled set to true",
		func(t *testing.T) {
			options := UpdateOptionsRequest{
				EmitOption: UpdateOptionsEmitOptionEnabled{
					EmitEnabled: true,
				},
				Query: "some query",
				Name:  "some name",
			}
			result := options.build()

			expectedResult := &projections.UpdateReq{
				Options: &projections.UpdateReq_Options{
					Name:  "some name",
					Query: "some query",
					EmitOption: &projections.UpdateReq_Options_EmitEnabled{
						EmitEnabled: true,
					},
				},
			}

			require.Equal(t, expectedResult, result)
		})

	t.Run("Build with name with trailing spaces", func(t *testing.T) {
		options := UpdateOptionsRequest{
			EmitOption: UpdateOptionsEmitOptionNoEmit{},
			Query:      "some query",
			Name:       " some name ",
		}
		result := options.build()

		expectedResult := &projections.UpdateReq{
			Options: &projections.UpdateReq_Options{
				Name:  " some name ",
				Query: "some query",
				EmitOption: &projections.UpdateReq_Options_NoEmitOptions{
					NoEmitOptions: &shared.Empty{},
				},
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Build with query with trailing spaces", func(t *testing.T) {
		options := UpdateOptionsRequest{
			EmitOption: UpdateOptionsEmitOptionNoEmit{},
			Query:      " some query ",
			Name:       "some name",
		}
		result := options.build()

		expectedResult := &projections.UpdateReq{
			Options: &projections.UpdateReq_Options{
				Name:  "some name",
				Query: " some query ",
				EmitOption: &projections.UpdateReq_Options_NoEmitOptions{
					NoEmitOptions: &shared.Empty{},
				},
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Build with empty name panics", func(t *testing.T) {
		options := UpdateOptionsRequest{
			Name:  "",
			Query: "some query",
		}

		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics with name consisting of spaces only", func(t *testing.T) {
		options := UpdateOptionsRequest{
			Name:  "     ",
			Query: "some query",
		}

		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics with empty query", func(t *testing.T) {
		options := UpdateOptionsRequest{
			Name:  "some name",
			Query: "",
		}

		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics with query consisting of spaces only", func(t *testing.T) {
		options := UpdateOptionsRequest{
			Name:  "some name",
			Query: "    ",
		}

		require.Panics(t, func() {
			options.build()
		})
	})
}
