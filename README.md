# Gridda
Gridda is a window positioning tool for Windows 10. Windows can be snapped to self defined grid lines and span one ore more grid tiles.

This tool does not utilize the mouse to position the windows but relies solely on keyboard shortcuts.

Gridda is configured via a simple YAML-formatted config file located inside the program directory and can be accessed via a notification icon in the system tray.

## Configuration
The project contains a default configuration file with all possible settings presented. Below you find a detailed description of all the configuration options.

The following excerpt shows a minimal configuration containing all the available options:
```
grids:               # defines the available grids to select from, up to 6
  -                  # are allowed
    name: Triptychon
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
    key: VK_LEFT        # the key to bind to, a list of available keys 
                        # can be found below
 
   
    
    
```


