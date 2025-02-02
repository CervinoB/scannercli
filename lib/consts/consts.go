// Package consts houses some constants needed across scannercli
package consts

import (
	"strings"
)

// Version contains the current semantic version of scannercli.
const Version = "0.0.0"

// Banner returns the ASCII-art banner with the scannercli logo
func Banner() string {
	banner := strings.Join([]string{
		`              _`,
		`             | |`,
		`             | |===( )   //////`,
		`             |_|   |||  | o o|`,
		`                    ||| ( c  )                  ____`,
		`                     ||| \= /                  ||   \_`,
		`                      ||||||                   ||     |`,
		`                      ||||||                ...||__/|-"`,
		`                      ||||||             __|________|__`,
		`                        |||             |______________|`,
		`                        |||             || ||      || ||`,
		`                        |||             || ||      || ||`,
		`------------------------|||-------------||-||------||-||-------`,
		`                        |__>            || ||      || ||`,
		`                                                        .CervinoB`,
	}, "\n")

	return banner
}
