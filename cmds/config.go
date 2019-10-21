//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package cmds

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/moonblue4242/Gridda/actions"
	"github.com/moonblue4242/Gridda/winapi"

	"gopkg.in/yaml.v2"
)

var (
	// activeConfig     *actions.ActiveConfig
	snapActions      = actions.NewSnapActions()
	availableActions = map[string]actions.ActionHandler{
		"MOVE_LEFT":           actions.MoveLeft,
		"MOVE_RIGHT":          actions.MoveRight,
		"TO_PREVIOUS_GRID":    actions.ToPreviousGrid,
		"TO_NEXT_GRID":        actions.ToNextGrid,
		"SNAP_LEFT":           snapActions.ToLeft(),
		"SNAP_RIGHT":          snapActions.ToRight(),
		"SNAP_TOP":            snapActions.ToTop(),
		"SNAP_BOTTOM":         snapActions.ToBottom(),
		"SPAN_HORIZONTAL_ADD": snapActions.SpanHorizontal(true),
		"SPAN_HORIZONTAL_DEC": snapActions.SpanHorizontal(false),
		"SPAN_VERTICAL_ADD":   snapActions.SpanVertical(true),
		"SPAN_VERTICAL_DEC":   snapActions.SpanVertical(false),
	}
	keys = map[string]int32{
		"VK_LEFT":   winapi.VK_LEFT,
		"VK_RIGHT":  winapi.VK_RIGHT,
		"VK_UP":     winapi.VK_UP,
		"VK_DOWN":   winapi.VK_DOWN,
		"VK_PRIOR":  winapi.VK_PRIOR,
		"VK_NEXT":   winapi.VK_NEXT,
		"VK_INSERT": winapi.VK_INSERT,
		"VK_DELETE": winapi.VK_DELETE,
		"VK_HOME":   winapi.VK_HOME,
		"VK_END":    winapi.VK_END,
		"VK_SPACE":  winapi.VK_SPACE,

		"VK_NUMPAD0":   winapi.VK_NUMPAD0,
		"VK_MULTIPLY":  winapi.VK_MULTIPLY,
		"VK_ADD":       winapi.VK_ADD,
		"VK_SEPARATOR": winapi.VK_SEPARATOR,
		"VK_SUBTRACT":  winapi.VK_SUBTRACT,
		"VK_DECIMAL":   winapi.VK_DECIMAL,
		"VK_DIVIDE":    winapi.VK_DIVIDE,
	}
)

// Config defines the configuration information for moving and sizing the windows
// via hotkey
type Config struct {
	Grids    []actions.Grid
	Bindings []Binding
}

// Binding defines a key binding along with an action
type Binding struct {
	Alt    bool
	Ctrl   bool
	Shift  bool
	Win    bool
	Key    string
	Action string
}

func (binding *Binding) validate() error {
	_, ok := availableActions[binding.Action]
	if !ok {
		return errors.New("Unknown action: " + binding.Action)
	}
	_, ok = keys[binding.Key]
	if !ok {
		return errors.New("Unknown key: " + binding.Key)
	}
	return nil
}

// LoadConfig will load a configuration from file
func LoadConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	for _, binding := range config.Bindings {
		if err := binding.validate(); err != nil {
			return nil, err
		}
	}
	for _, grid := range config.Grids {
		if err := grid.Validate(); err != nil {
			return nil, err
		}
	}
	return &config, nil
}

func init() {
	// numbers
	for i := 0; i < 10; i++ {
		keys["VK_"+strconv.Itoa(i)] = winapi.DigitToVK(i)
	}
	// chars
	for i := 0; i < 26; i++ {
		var ch rune = rune('a' + i)
		keys["VK_"+strings.ToUpper(fmt.Sprintf("%c", ch))] = winapi.CharToVK(ch)
	}
	// numpad
	for i := 0; i < 10; i++ {
		keys["VK_NUMPAD"+strconv.Itoa(i)] = winapi.NumToVK(i)
	}
	// function keys
	for i := 1; i <= 12; i++ {
		keys["VK_F"+strconv.Itoa(i)] = winapi.FunctionToVK(i)
	}
}
