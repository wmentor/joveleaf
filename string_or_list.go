package joveleaf

import (
	"encoding/json"
	"strings"
)

// StringOrList handles the notebook convention where multi-line strings
// may be stored as either a single string or a list of strings.
// When unmarshaling, a list is joined into a single string.
type StringOrList string

// UnmarshalJSON implements json.Unmarshaler.
func (s *StringOrList) UnmarshalJSON(data []byte) error {
	// Try string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = StringOrList(str)
		return nil
	}

	// Try list of strings
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	*s = StringOrList(strings.Join(list, ""))
	return nil
}

// MarshalJSON implements json.Marshaler.
func (s StringOrList) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

// String returns the underlying string value.
func (s StringOrList) String() string {
	return string(s)
}
