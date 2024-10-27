package cmd

import (
	"github.com/smoqadam/kindle-send/config"
	"github.com/smoqadam/kindle-send/server"
	"github.com/smoqadam/kindle-send/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 8080, "Port to run the server on")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server with send and download endpoints",
	Long: `Starts a web server that provides two endpoints:
- GET /send?urls=url1,url2 : Processes URLs and send email
- GET /download?urls=url1,url2 : Downloads content from provided URLs`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")

		_, err := config.Load()
		if err != nil {
			util.Red.Println("Error loading config:", err)
			return
		}

		server.Start(port)
	},
}
