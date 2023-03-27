package cli

import (
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/gui"
	"github.com/spf13/cobra"
)

func runGUI(t *core.Timetrace) *cobra.Command {
	gui := &cobra.Command{
		Use:   "gui",
		Short: "Show the GUI",
		Run: func(cmd *cobra.Command, args []string) {
			gui.RunGui(t)

		},
	}

	return gui
}
