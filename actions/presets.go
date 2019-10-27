//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package actions

import "strings"

// Presets define preconditions for windows
type Presets []Preset

// Preset defines values for a window which are set as defaults
type Preset struct {
	Executable string
	Span       *Span
	Expandable *Expandable
}

// FindFirst will retrieve the first preset with the defined properties
func (presets *Presets) FindFirst(executable string) (preset *Preset) {
	for _, item := range *presets {
		if strings.ToUpper(item.Executable) == strings.ToUpper(executable) {
			preset = &item
			break
		}
	}
	return
}
