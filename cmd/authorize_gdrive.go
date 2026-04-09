package cmd

import "github.com/spf13/cobra"

func newAuthorizeGDriveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "authorize-gdrive",
		Short: "Authorize Google Drive access (required once before first use)",
		Long: `Opens a browser for Google Drive OAuth2 authorization and saves the token
to token_file (configured in config.yaml).

Run this command once from the terminal before launching the tray app.
The tray app has no terminal, so it cannot prompt for authorization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.GDrive.Authorize()
		},
	}
}
