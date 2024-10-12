package log

import (
	"fmt"
	"strings"
)

// Format determines what kind of output is written by the logger.
type Format int

// Common log types.
//
//go:generate stringer -type Format -trimprefix Format
const (
	FormatText Format = iota
	FormatJSON
)

// UnmarshalText implements [encoding.TextUnmarshaler].
func (lt *Format) UnmarshalText(data []byte) error {
	// Convert the input data to uppercase
	str := strings.ToUpper(string(data))

	switch str {
	case "TEXT":
		*lt = FormatText
	case "JSON":
		*lt = FormatJSON
	default:
		return fmt.Errorf("invalid LogType: %s", data)
	}

	return nil
}
