# ![](./logo.png) Gridda
Gridda is a window positioning tool for Windows 10 64-bit. Windows can be snapped to self defined grid lines and span one ore more grid tiles.

This tool does not utilize the mouse to position the windows but relies solely on keyboard shortcuts.

Gridda is configured via a simple YAML-formatted config file located inside the program directory and can be accessed via a notification icon in the system tray.

## Configuration
The project contains a default configuration file with all possible settings presented. Below you find a detailed description of all the configuration options.

In case of an invalid configuration the application will exit and write the error to a file called gridda.log inside the application directory.

The following excerpt shows a minimal configuration containing all the available options:
```
grids:               # defines the available grids to select from, up to 6
  -                  # are allowed
    name: Triptych
    columns: [1,2,1] # array of weights defining the amount of  
                     # columns along with their respective weights
    rows: [1,1]      # array of rows defining the amount of
                     # rows along with their respective weights
    presets:         # presets define default values for application windows 
                     # per Grid
      -
        executable: firefox.exe  # the name of the executable file the window 
                                 # belongs, case insensitive
        span:                    # span defines the initial span of the window
          columns: 1             # the columns to span
          rows: 2                # the rows to span
bindings:            # defines the key bindings
  -
    action: SNAP_LEFT   # the Gridda action to perform on hot key
    alt: true           # alt defines if the alt key must be pressed
                        # if omitted reverts to false
    ctrl: true          # ctrl defines if the alt key must be pressed
                        # if omitted reverts to false
    shift: true         # shift defines if the shift key must be pressed
                        # if omitted reverts to false
    win: true           # win defines if the windows key must be pressed
                        # if omitted reverts to false
    key: VK_LEFT        # the key to bind to, a list of available keys 
                        # can be found below

    
```

### Actions
The following actions are available and can be bind to a hotkey:

* SNAP_LEFT   
   Snaps the focused window to the next grid tile on the left
* SNAP_RIGHT   
   Snaps the focused window to the next grid tile on the right
* SNAP_TOP   
   Snaps the focused window to the next grid tile above
* SNAP_BOTTOM   
   Snaps the focused window to the next grid tile below
* SPAN_HORIZONTAL_ADD   
   Increase the horizontal span of the focused window
* SPAN_HORIZONTAL_DEC   
   Decrease the horizontal span of the focused window
* SPAN_VERTICAL_ADD   
   Increase the vertical span of the focused window
* SPAN_VERTICAL_DEC   
   Decrease the vertical span of the focused window
* TO_PREVIOUS_GRID   
   Cycle through the available grids in reverse definition order
* TO_NEXT_GRID   
   Cycle through the available grids in definition order


### Keys

The following key shortcuts are defined and can be used directly

* VK_SPACE   
   The space key
* VK_LEFT    
   The left cursor key
* VK_RIGHT   
   The right cursor key
* VK_UP   
   The up cursor key
* VK_DOWN   
   The down cursor key
* VK_PRIOR   
   The page up key
* VK_NEXT   
   The page down key
* VK_INSERT   
   The insert key
* VK_DELETE   
   The delete key
* VK_HOME   
   The home key
* VK_END   
   The end key
*	VK_MULTIPLY   
  The multiply key
*	VK_ADD   
   The plus key
*	VK_SUBTRACT   
   The subtract key
*	VK_DIVIDE   
  The divide key
* VK_A - VK_Z   
   Represents a single char a-z
* VK_0 - VK_9   
   Represents a single char 0-9
* VK_F1 - VK_F12   
   The functions keys from 1 to 12   
* VK_NUMPAD0 - VK_NUMPAD9  
  The numpad keys from 0 to 9   

