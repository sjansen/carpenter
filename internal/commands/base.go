package commands

import (
	"os"

	"github.com/sjansen/carpenter/internal/sys"
)

// Base contains dependencies and options common to all commands
type Base struct {
	sys.IO

	Debug     *os.File
	Verbosity int
}
