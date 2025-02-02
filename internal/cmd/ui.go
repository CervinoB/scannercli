package cmd

import (
	"fmt"
	"os"

	"github.com/CervinoB/sonarcli/cmd/state"
	"github.com/CervinoB/sonarcli/internal/ui/pb"
	"github.com/CervinoB/sonarcli/lib/consts"
	"github.com/fatih/color"
)

const (
	// Max length of left-side progress bar text before trimming is forced
	maxLeftLength = 30
	// Amount of padding in chars between rendered progress
	// bar text and right-side terminal window edge.
	termPadding      = 1
	defaultTermWidth = 80
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

func printBanner(gs *state.GlobalState) {
	// if gs.Flags.Quiet {
	// 	return // do not print banner when --quiet is enabled
	// }

	banner := getBanner(true)
	_, err := fmt.Fprintf(os.Stdout, "\n%s\n\n", banner)
	if err != nil {
		gs.Logger.Warnf("could not print k6 banner message to stdout: %s", err.Error())
	}
}

func printBar(gs *state.GlobalState, bar *pb.ProgressBar) {
	// if gs.Flags.Quiet {
	// 	return
	// }
	end := "\n"
	// TODO: refactor widthDelta away? make the progressbar rendering a bit more
	// stateless... basically first render the left and right parts, so we know
	// how long the longest line is, and how much space we have for the progress
	widthDelta := -defaultTermWidth
	rendered := bar.Render(0, widthDelta)
	// Only output the left and middle part of the progress bar
	printToStdout(gs, rendered.String()+end)
}
