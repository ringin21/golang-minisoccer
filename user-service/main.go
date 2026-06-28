// Package main menyediakan entry point executable user-service.
package main

import "user-service/cmd"

// main meneruskan proses startup ke command utama yang melakukan bootstrap
// dependency, route, middleware, dan HTTP server aplikasi.
func main() {
	cmd.Run()
}
