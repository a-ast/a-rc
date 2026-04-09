package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"a-rc/internal/core"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Drive implements core.Uploader using the Google Drive API with OAuth2.
// On first use it opens a browser for authorization and caches the token in token_file.
type Drive struct {
	configRepo core.ConfigRepository
	configPath *string
}

func New(repo core.ConfigRepository, configPath *string) *Drive {
	return &Drive{configRepo: repo, configPath: configPath}
}

func (a *Drive) Upload(localPath string) error {
	cfg, err := a.configRepo.Load(*a.configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	g := cfg.GDrive
	if g.CredentialsFile == "" {
		return fmt.Errorf("gdrive.credentials_file is not set in config")
	}
	if g.TokenFile == "" {
		return fmt.Errorf("gdrive.token_file is not set in config")
	}
	if g.Folder == "" {
		return fmt.Errorf("gdrive.folder is not set in config")
	}

	svc, err := newDriveService(g.CredentialsFile, g.TokenFile)
	if err != nil {
		return fmt.Errorf("creating Drive client: %w", err)
	}

	folderID, err := findFolder(svc, g.Folder)
	if err != nil {
		return fmt.Errorf("finding folder %q: %w", g.Folder, err)
	}

	return uploadFile(svc, localPath, folderID)
}

func newDriveService(credFile, tokenFile string) (*drive.Service, error) {
	data, err := os.ReadFile(credFile)
	if err != nil {
		return nil, fmt.Errorf("reading credentials file: %w", err)
	}
	oauthCfg, err := google.ConfigFromJSON(data, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("parsing credentials: %w", err)
	}

	tok, err := loadToken(tokenFile)
	if err != nil {
		// No cached token — run the interactive auth flow.
		tok, err = authorizeInteractive(oauthCfg, tokenFile)
		if err != nil {
			return nil, fmt.Errorf("authorization failed: %w", err)
		}
	}

	ctx := context.Background()
	client := oauthCfg.Client(ctx, tok)
	return drive.NewService(ctx, option.WithHTTPClient(client))
}

// authorizeInteractive runs the OAuth2 browser flow, saves the token, and returns it.
func authorizeInteractive(cfg *oauth2.Config, tokenFile string) (*oauth2.Token, error) {
	authURL := cfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("\nOpening browser for Google Drive authorization...")
	fmt.Printf("If the browser does not open, visit:\n%s\n\n", authURL)
	_ = exec.Command("open", authURL).Start() // macOS

	fmt.Print("Paste the authorization code here: ")
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("reading auth code: %w", err)
	}

	tok, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("exchanging auth code: %w", err)
	}

	if err := saveToken(tokenFile, tok); err != nil {
		return nil, fmt.Errorf("saving token: %w", err)
	}
	fmt.Println("Authorization successful, token saved.")
	return tok, nil
}

func loadToken(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tok oauth2.Token
	return &tok, json.NewDecoder(f).Decode(&tok)
}

func saveToken(path string, tok *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(tok)
}

// findFolder searches for a folder by name in all drives visible to the user.
func findFolder(svc *drive.Service, name string) (string, error) {
	q := fmt.Sprintf("name=%q and mimeType='application/vnd.google-apps.folder' and trashed=false", name)
	res, err := svc.Files.List().
		Q(q).
		Fields("files(id, name)").
		PageSize(1).
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		Corpora("allDrives").
		Do()
	if err != nil {
		return "", err
	}
	if len(res.Files) == 0 {
		return "", fmt.Errorf("folder %q not found in Drive", name)
	}
	return res.Files[0].Id, nil
}

func uploadFile(svc *drive.Service, localPath, folderID string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	name := filepath.Base(localPath)

	// Check if a file with the same name already exists in the folder.
	existingID, err := findExistingFile(svc, name, folderID)
	if err != nil {
		return fmt.Errorf("checking for existing file: %w", err)
	}

	if existingID != "" {
		// Overwrite: update content, keep metadata.
		_, err = svc.Files.Update(existingID, &drive.File{}).
			Media(f).
			SupportsAllDrives(true).
			Do()
	} else {
		_, err = svc.Files.Create(&drive.File{Name: name, Parents: []string{folderID}}).
			Media(f).
			SupportsAllDrives(true).
			Do()
	}
	if err != nil {
		return fmt.Errorf("uploading to Drive: %w", err)
	}
	return nil
}

func findExistingFile(svc *drive.Service, name, folderID string) (string, error) {
	q := fmt.Sprintf("name=%q and %q in parents and mimeType!='application/vnd.google-apps.folder' and trashed=false", name, folderID)
	res, err := svc.Files.List().
		Q(q).
		Fields("files(id)").
		PageSize(1).
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		Do()
	if err != nil {
		return "", err
	}
	if len(res.Files) == 0 {
		return "", nil
	}
	return res.Files[0].Id, nil
}
