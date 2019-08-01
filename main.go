package statistics_client_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
)

type Statistics struct {
	AppName    string
	TargetPath string
	URL        string
	UUID       string
}

type Request struct {
	AppName string `json:"app_name"`
	UUID    string `json:"uuid"`
}

func Client(appName, url string) *Statistics {
	return &Statistics{
		AppName:    appName,
		URL:        url,
		TargetPath: getTargetPath(appName),
	}
}

func (s *Statistics) generateUUID() error {
	u, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	s.UUID = u.String()

	return nil
}

func getTargetPath(appName string) string {
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		home = os.Getenv("APPDATA")
	}

	configDirectoryPath := filepath.Join(home, ".config")
	if _, err := os.Stat(configDirectoryPath); os.IsNotExist(err) {
		_ = os.Mkdir(configDirectoryPath, 0700)
	}

	configPath := filepath.Join(configDirectoryPath, appName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		_ = os.Mkdir(configPath, 0700)
	}

	return filepath.Join(configPath, "uuid")
}

func (s *Statistics) saveUUID() error {
	return ioutil.WriteFile(s.TargetPath, []byte(s.UUID), 0644)
}

// Check files exists. don't check body.
func (s *Statistics) checkUUID() bool {
	_, err := os.Stat(s.TargetPath)
	return err == nil
}

func (s *Statistics) request() error {
	r := Request{
		AppName: s.AppName,
		UUID:    s.UUID,
	}

	body, err := json.Marshal(r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.URL+"/uuid", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to request: %s, status code: %d", s.URL, resp.StatusCode)
	}

	return err
}

func (s *Statistics) Register() error {
	if !s.checkUUID() {
		err := s.generateUUID()
		if err != nil {
			return err
		}

		err = s.request()
		if err != nil {
			return err
		}

		err = s.saveUUID()
		if err != nil {
			return err
		}
	}

	return nil
}
