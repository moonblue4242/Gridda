//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package main

// Build with -ldflags "-H windowsgui" to create a gui element
func main() {

	// create business logic controller
	logic := NewLogic()
	logic.Loop()

}
