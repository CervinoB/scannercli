package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"go.k6.io/k6/cmd/state"
	"go.k6.io/k6/lib/consts"
	"go.k6.io/k6/ui/pb"
)

func setColor(noColor bool, c *color.Color) *color.Color {
	if noColor {
		c.DisableColor()
	} else {
		c.EnableColor()
	}
	return c
}

// getColor returns the requested color, or an uncolored object, depending on
// the value of noColor. The explicit EnableColor() and DisableColor() are
// needed because the library checks os.Stdout itself otherwise...
func getColor(noColor bool, attributes ...color.Attribute) *color.Color {
	return setColor(noColor, color.New(attributes...))
}

func getBanner(noColor bool) string {
	c := setColor(noColor, color.RGB(0xFF, 0x67, 0x1d).Add(color.Bold))
	return c.Sprint(consts.Banner())
}

func printBanner() {

	banner := getBanner(true)
	_, err := fmt.Fprintf(os.Stdout, "\n%s\n\n", banner)
	if err != nil {
		gs.Logger.Warnf("could not print k6 banner message to stdout: %s", err.Error())
	}
}

func printBar(gs *state.GlobalState, bar *pb.ProgressBar) {
	if gs.Flags.Quiet {
		return
	}
	end := "\n"
	// TODO: refactor widthDelta away? make the progressbar rendering a bit more
	// stateless... basically first render the left and right parts, so we know
	// how long the longest line is, and how much space we have for the progress
	widthDelta := -defaultTermWidth
	if gs.Stdout.IsTTY {
		// If we're in a TTY, instead of printing the bar and going to the next
		// line, erase everything till the end of the line and return to the
		// start, so that the next print will overwrite the same line.
		//
		// TODO: check for cross platform support
		end = "\x1b[0K\r"
		widthDelta = 0
	}
	rendered := bar.Render(0, widthDelta)
	// Only output the left and middle part of the progress bar
	printToStdout(gs, rendered.String()+end)
}
