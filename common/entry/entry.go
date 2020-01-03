package entry

import (
	"log"
	"os"
)

// HandleMainE calls the passed mainE function
// and if that returns an error it prints it with full stack trace (when
// "github.com/pkg/errors" is used for creating the error)
// and then exits with exit code 1.
// If mainE returns with a nil error then nothing will be printed and it'll
// exit with os.Exit(0).
//
// The purpose of this function is to have a common way of handling
// the entry point of a go program, the main() function, which should
// be as minimal as possible. Requiring a mainE function enforces
// the code to be more testable as mainE itself will be easier to test.
//
// Example usage:
//
//   package main
//
//   import (
//     "github.com/bitrise-io/bitrise-devops-microservice-common/common/entry"
//   )
//
//   func mainE() error {
//     return nil
//   }
//
//   func main() {
//     entry.HandleMainE(mainE)
//   }
func HandleMainE(mainE func() error) {
	if err := mainE(); err != nil {
		log.Printf("[!] Exception: %+v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
