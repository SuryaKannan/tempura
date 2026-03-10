package cli

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var port int = 8732

func init() {
	serve.Flags().IntVarP(&port, "port", "p", port, "port to serve tempura")
	rootCmd.AddCommand(serve)
}

func HandleFry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "tempura is running!")
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "start your tempura server",
	Run: func(cmd *cobra.Command, args []string) {

		mux := http.NewServeMux()

		mux.HandleFunc("/", HandleFry)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		}

		color.Green("\nserving Tempura on port %d!\n", port)

		server.ListenAndServe()
	},
}
