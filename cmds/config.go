package cmds

import (
	"errors"
	"io/ioutil"

	"sonnenfroh.de/test/actions"
	"sonnenfroh.de/test/winapi"

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
		"VK_LEFT":  winapi.VK_LEFT,
		"VK_RIGHT": winapi.VK_RIGHT,
		"VK_UP":    winapi.VK_UP,
		"VK_DOWN":  winapi.VK_DOWN,
		"VK_PRIOR": winapi.VK_PRIOR,
		"VK_NEXT":  winapi.VK_NEXT,
		"VK_A":     winapi.CharToVK('a'),
		"VK_B":     winapi.CharToVK('b'),
		"VK_C":     winapi.CharToVK('c'),
		"VK_D":     winapi.CharToVK('d'),
		"VK_E":     winapi.CharToVK('e'),
		"VK_F":     winapi.CharToVK('f'),
		"VK_G":     winapi.CharToVK('g'),
		"VK_H":     winapi.CharToVK('h'),
		"VK_I":     winapi.CharToVK('i'),
		"VK_J":     winapi.CharToVK('j'),
		"VK_K":     winapi.CharToVK('k'),
		"VK_L":     winapi.CharToVK('l'),
		"VK_M":     winapi.CharToVK('m'),
		"VK_N":     winapi.CharToVK('n'),
		"VK_O":     winapi.CharToVK('o'),
		"VK_P":     winapi.CharToVK('p'),
		"VK_Q":     winapi.CharToVK('q'),
		"VK_R":     winapi.CharToVK('r'),
		"VK_S":     winapi.CharToVK('s'),
		"VK_T":     winapi.CharToVK('t'),
		"VK_U":     winapi.CharToVK('u'),
		"VK_V":     winapi.CharToVK('v'),
		"VK_W":     winapi.CharToVK('w'),
		"VK_X":     winapi.CharToVK('x'),
		"VK_Y":     winapi.CharToVK('y'),
		"VK_Z":     winapi.CharToVK('z'),
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
