package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/junbin-yang/golib/bytesconv"
)

var (
	def           = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = def.Marshal
	Unmarshal     = def.Unmarshal
	MarshalIndent = def.MarshalIndent
	NewDecoder    = def.NewDecoder
	NewEncoder    = def.NewEncoder
	// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
	// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
	// interface{} as a Number instead of as a float64.
	EnableDecoderUseNumber = false
	// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
	// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
	// return an error when the destination is a struct and the input contains object
	// keys which do not match any non-ignored, exported fields in the destination.
	EnableDecoderDisallowUnknownFields = false
)

func initHandle(t ...string) jsoniter.API {
	tag := "json"
	if len(t) > 0 {
		tag = t[0]
	}
	return jsoniter.Config{TagKey: tag}.Froze()
}

func ObjectToJson(object interface{}, t ...string) (string, error) {
	json := initHandle(t...)
	str, e := json.Marshal(object)
	if e != nil {
		return "", e
	}
	return string(str), nil
}

func JsonToObject(jsonString string, object interface{}, t ...string) error {
	json := initHandle(t...)
	return json.Unmarshal(bytesconv.StringToBytes(jsonString), &object)
}
