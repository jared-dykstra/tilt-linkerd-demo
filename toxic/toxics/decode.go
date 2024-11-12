package toxics

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

func StringToDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		// Check if the data is a string and the target type is time.Duration
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Duration(0)) {
			return data, nil
		}

		// Convert it by parsing
		return time.ParseDuration(data.(string))
	}
}

func mapAnyToStruct(src map[string]any, dst any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: StringToDurationHookFunc(),
		Result:     &dst,
	})
	if err != nil {
		return err
	}
	if err := decoder.Decode(src); err != nil {
		return err
	}
	return nil
}
