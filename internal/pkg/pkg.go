package pkg

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Package struct {
	Version   string    `json:"version"`
	URI       string    `json:"uri"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(pkg string) (*Package, error) {
	pkgList := strings.Split(pkg, "@")
	name := getBinaryNameFromURI(pkgList[0])
	if name == "" {
		return nil, fmt.Errorf("could not determine package name from URI: %s", pkgList[0])
	}

	return &Package{
		Version:   pkgList[1],
		URI:       pkgList[0],
		Name:      name,
		UpdatedAt: time.Now(),
	}, nil
}

func (p *Package) String() string {
	return "Name: " + p.Name + "\n" +
		"URI: " + p.URI + "@" + p.Version + "\n" +
		"Updated: " + p.UpdatedAt.String()
}

func (p *Package) ID() string {
	return p.Name
}

func (p *Package) URIWithVersion() string {
	return fmt.Sprintf("%s@%s", p.URI, p.Version)
}

func (p *Package) Install() (string, error) {
	cmd := exec.Command("go", "install", p.URIWithVersion())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to install package: %v, stderr: %s", err, stderr.String())
	}

	if stderr.Len() > 0 {
		return "", fmt.Errorf("failed to install package, stderr: %s", stderr.String())
	}

	return stdout.String(), nil
}

func (p *Package) UpdateVersion(version string) {
	p.Version = version
	p.UpdatedAt = time.Now()
}

func getBinaryNameFromURI(uri string) string {
	splitSlash := strings.Split(uri, "/")

	// example.com/cmd/tool/v2 becomes tool
	matched, err := regexp.MatchString(`v\d+`, splitSlash[len(splitSlash)-1])
	if err != nil {
		slog.Error("failed to match regex", "error", err)
		return ""
	}

	if matched {
		return splitSlash[len(splitSlash)-2]
	}

	return splitSlash[len(splitSlash)-1]
}
