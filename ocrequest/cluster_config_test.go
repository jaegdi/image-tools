package ocrequest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadClusterConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{name: "valid", content: `{"config":{"test":{"name":"test","url":"https://example.invalid","token":"token"}}}`},
		{name: "invalid JSON", content: `{`, wantErr: "parse"},
		{name: "empty config", content: `{"config":{}}`, wantErr: "empty"},
		{name: "missing URL", content: `{"config":{"test":{"name":"test"}}}`, wantErr: "name and URL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "clusterconfig.json")
			if err := os.WriteFile(path, []byte(tt.content), 0600); err != nil {
				t.Fatal(err)
			}
			config, err := loadClusterConfig(path)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("got error %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if config.Config["test"].Token != "token" {
				t.Fatal("cluster configuration was not loaded")
			}
		})
	}
}

func TestLoadClusterConfigMissing(t *testing.T) {
	_, err := loadClusterConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil || !strings.Contains(err.Error(), "read") {
		t.Fatalf("got error %v, want read error", err)
	}
}
