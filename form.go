package requests

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

// File represents a form file
type File struct {
	Name     string        // Form field name
	FileName string        // File name
	Content  io.ReadCloser // File content
}

func (f *File) SetContent(content io.ReadCloser) {
	f.Content = content
}

func (f *File) SetFileName(fileName string) {
	f.FileName = fileName
}

func (f *File) SetName(name string) {
	f.Name = name
}

func parseFormFields(fields any) (url.Values, error) {
	switch data := fields.(type) {
	case url.Values:
		// Directly return url.Values data.
		return data, nil
	case map[string][]string:
		// Convert and return map[string][]string data as url.Values.
		return url.Values(data), nil
	case map[string]string:
		// Convert and return map[string]string data as url.Values.
		values := make(url.Values)
		for key, value := range data {
			values.Set(key, value)
		}
		return values, nil
	default:
		// Attempt to use query.Values for encoding struct types.
		if values, err := query.Values(fields); err == nil {
			return values, nil
		} else {
			// Return an error if encoding fails or type is unsupported.
			return nil, fmt.Errorf("%w: %v", ErrUnsupportedFormFieldsType, err)
		}
	}
}

func parseForm(v any) (url.Values, []*File, error) {
	switch data := v.(type) {
	case url.Values:
		// Directly return url.Values data.
		return data, nil, nil
	case map[string][]string:
		// Convert and return map[string][]string data as url.Values.
		return url.Values(data), nil, nil
	case map[string]string:
		// Convert and return map[string]string data as url.Values.
		values := make(url.Values)
		for key, value := range data {
			values.Set(key, value)
		}
		return values, nil, nil
	case map[string]any:
		// Convert and return map[string]any data as url.Values and []*File.
		values := make(url.Values)
		files := make([]*File, 0)
		for key, value := range data {
			switch v := value.(type) {
			case string:
				values.Set(key, v)
			case []string:
				for _, v := range v {
					values.Add(key, v)
				}
			case *File:
				v.SetName(key)
				files = append(files, v)
			default:
				// Return an error if type is unsupported.
				return nil, nil, fmt.Errorf("%w: %T", ErrUnsupportedDataType, value)
			}
		}
		return values, files, nil
	default:
		// Attempt to use query.Values for encoding struct types.
		if values, err := query.Values(v); err == nil {
			return values, nil, nil
		} else {
			// Return an error if encoding fails or type is unsupported.
			return nil, nil, fmt.Errorf("%w: %v", ErrUnsupportedFormFieldsType, err)
		}
	}
}

// FormEncoder handles encoding of form data.
type FormEncoder struct{}

// Encode encodes the given value into URL-encoded form data.
func (e *FormEncoder) Encode(v any) (io.Reader, error) {
	switch data := v.(type) {
	case url.Values:
		// Directly encode url.Values data.
		return strings.NewReader(data.Encode()), nil
	case map[string][]string:
		// Convert and encode map[string][]string data as url.Values.
		values := url.Values(data)
		return strings.NewReader(values.Encode()), nil
	case map[string]string:
		// Convert and encode map[string]string data as url.Values.
		values := make(url.Values)
		for key, value := range data {
			values.Set(key, value)
		}
		return strings.NewReader(values.Encode()), nil
	default:
		// Attempt to use query.Values for encoding struct types.
		if values, err := query.Values(v); err == nil {
			return strings.NewReader(values.Encode()), nil
		} else {
			// Return an error if encoding fails or type is unsupported.
			return nil, fmt.Errorf("%w: %v", ErrEncodingFailed, err)
		}
	}
}

var DefaultFormEncoder = &FormEncoder{}
