package jsonrpcws

import "encoding/json"

// Convert func
func Convert(param interface{}, output interface{}) error {
	r, err := json.Marshal(param)
	if err != nil {
		return err
	}
	err = json.Unmarshal(r, &output)
	if err != nil {
		switch t := err.(type) {
		default:
			println("Type default: ", t)
		case *json.SyntaxError:
			println("Type SyntaxError: ", t)
		case *json.UnmarshalFieldError:
			println("Type UnmarshalFieldError: ", t)
		case *json.UnmarshalTypeError:
			println("Type UnmarshalTypeError: ", t)
		case *json.UnsupportedTypeError:
			println("Type UnsupportedTypeError: ", t)
		case *json.UnsupportedValueError:
			println("Type UnsupportedValueError: ", t)
		}

		return err
	}
	return nil
}
