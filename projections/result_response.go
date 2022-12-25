package projections

import (
	"encoding/json"

	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/projections"
	"google.golang.org/protobuf/types/known/structpb"
)

// ResultType is a type of the result received by reading a
// result of a projection of projection's state.
type ResultType string

const (
	// ResultTypeNull received value is null.
	ResultTypeNull ResultType = "ResultTypeNull"
	// ResultTypeFloat64 received value is a float64 type.
	ResultTypeFloat64 ResultType = "ResultTypeFloat64"
	// ResultTypeString received value is a string type.
	ResultTypeString ResultType = "ResultTypeString"
	// ResultTypeBool received value is a bool type.
	ResultTypeBool ResultType = "ResultTypeBool"
	// ResultTypeJSONObject received value is a JSON object type.
	ResultTypeJSONObject ResultType = "ResultTypeJSONObject"
	// ResultTypeJSONObjectList received value is a slice of JSON objects.
	ResultTypeJSONObjectList ResultType = "ResultTypeJSONObjectList"
)

// IsResult an interface of the result value.
type IsResult interface {
	// ResultType returns a type of the value.
	ResultType() ResultType
}

// BoolResult represents a boolean value.
type BoolResult bool

// ResultType returns ResultTypeBool.
func (s BoolResult) ResultType() ResultType {
	return ResultTypeBool
}

// NullResult represents a null.
type NullResult struct{}

// ResultType returns ResultTypeNull.
func (s NullResult) ResultType() ResultType {
	return ResultTypeNull
}

// Float64Result represents a float64 value.
type Float64Result float64

// ResultType returns ResultTypeFloat64.
func (s Float64Result) ResultType() ResultType {
	return ResultTypeFloat64
}

// StringResult represents a string value.
type StringResult string

// ResultType returns ResultTypeString.
func (s StringResult) ResultType() ResultType {
	return ResultTypeString
}

// JSONObjectResult represents a JSON object.
type JSONObjectResult JSONBytes

// ResultType returns ResultTypeJSONObject.
func (s JSONObjectResult) ResultType() ResultType {
	return ResultTypeJSONObject
}

// JSONObjectListResult represents a slice of JSON objects.
type JSONObjectListResult []JSONBytes

// ResultType returns ResultTypeJSONObjectList.
func (s JSONObjectListResult) ResultType() ResultType {
	return ResultTypeJSONObjectList
}

func newResultResponse(state *projections.ResultResp) IsResult {
	_, ok := state.Result.Kind.(*structpb.Value_NullValue)

	if ok {
		return &NullResult{}
	}

	numberValue, ok := state.Result.Kind.(*structpb.Value_NumberValue)

	if ok {
		return Float64Result(numberValue.NumberValue)
	}

	stringValue, ok := state.Result.Kind.(*structpb.Value_StringValue)

	if ok {
		return StringResult(stringValue.StringValue)
	}

	boolValue, ok := state.Result.Kind.(*structpb.Value_BoolValue)

	if ok {
		return BoolResult(boolValue.BoolValue)
	}

	structValue, ok := state.Result.Kind.(*structpb.Value_StructValue)

	if ok {
		return JSONObjectResult(newResultResponseStructMap(structValue.StructValue.Fields))
	}

	listValue, ok := state.Result.Kind.(*structpb.Value_ListValue)

	if ok {
		return JSONObjectListResult(newResultResponseList(listValue.ListValue.Values))
	}

	return nil
}

func newResultResponseStructMap(grpcValue map[string]*structpb.Value) JSONBytes {
	jsonValue, err := json.Marshal(grpcValue)
	if err != nil {
		panic(err)
	}

	return jsonValue
}

func newResultResponseList(grpcValue []*structpb.Value) []JSONBytes {
	var result []JSONBytes

	for _, value := range grpcValue {
		byteValue, err := value.MarshalJSON()
		if err != nil {
			panic(err)
		}
		result = append(result, byteValue)
	}

	return result
}

func newStateResponse(state *projections.StateResp) IsResult {
	_, ok := state.State.Kind.(*structpb.Value_NullValue)

	if ok {
		return NullResult{}
	}

	numberValue, ok := state.State.Kind.(*structpb.Value_NumberValue)

	if ok {
		return Float64Result(numberValue.NumberValue)
	}

	stringValue, ok := state.State.Kind.(*structpb.Value_StringValue)

	if ok {
		return StringResult(stringValue.StringValue)
	}

	boolValue, ok := state.State.Kind.(*structpb.Value_BoolValue)

	if ok {
		return BoolResult(boolValue.BoolValue)
	}

	structValue, ok := state.State.Kind.(*structpb.Value_StructValue)

	if ok {
		return JSONObjectResult(newStateResponseStructMap(structValue.StructValue.Fields))
	}

	listValue, ok := state.State.Kind.(*structpb.Value_ListValue)

	if ok {
		return JSONObjectListResult(newStateResponseList(listValue.ListValue.Values))
	}

	return nil
}

func newStateResponseStructMap(grpcValue map[string]*structpb.Value) JSONBytes {
	result, err := json.Marshal(grpcValue)
	if err != nil {
		panic(err)
	}

	return result
}

func newStateResponseList(grpcValue []*structpb.Value) []JSONBytes {
	var result []JSONBytes

	for _, value := range grpcValue {
		byteValue, err := value.MarshalJSON()
		if err != nil {
			panic(err)
		}
		result = append(result, byteValue)
	}

	return result
}
