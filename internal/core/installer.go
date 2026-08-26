package core

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/config"
)

const singBoxVersion = "1.11.4"

// EnsureSingBoxBinary finds sing-box or automatically downloads it into ~/.vify/bin/sing-box
func EnsureSingBoxBinary() (string, error) {
	// 1. Check if ~/.vify/bin/sing-box exists
	localBin := filepath.Join(config.GetVifyDir(), "bin", "sing-box")
	if runtime.GOOS == "windows" {
		localBin += ".exe"
	}

	if _, err := os.Stat(localBin); err == nil {
		return localBin, nil
	}

	// 2. Check if sing-box is in PATH
	if p, err := exec.LookPath("sing-box"); err == nil {
		return p, nil
	}

	// 3. Auto-download from GitHub releases
	binDir := filepath.Join(config.GetVifyDir(), "bin")
	_ = os.MkdirAll(binDir, 0755)

	err := downloadSingBox(binDir)
	if err != nil {
		return "", fmt.Errorf("failed to download sing-box core: %w", err)
	}

	if _, err := os.Stat(localBin); err == nil {
		_ = os.Chmod(localBin, 0755)
		return localBin, nil
	}

	return "", fmt.Errorf("sing-box binary not found after download")
}

func downloadSingBox(destDir string) error {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	var archiveName, fileExt string
	if osName == "windows" {
		archiveName = fmt.Sprintf("sing-box-%s-%s-%s.zip", singBoxVersion, osName, arch)
		fileExt = ".zip"
	} else {
		archiveName = fmt.Sprintf("sing-box-%s-%s-%s.tar.gz", singBoxVersion, osName, arch)
		fileExt = ".tar.gz"
	}

	downloadURL := fmt.Sprintf("https://github.com/SagerNet/sing-box/releases/download/v%s/%s", singBoxVersion, archiveName)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		if strings.Contains(err.Error(), "certificate signed by unknown authority") || strings.Contains(err.Error(), "x509") {
			return fmt.Errorf("%w\n\nyour system is missing root CA certificates, which are required to download sing-box securely. "+
				"On minimal/headless Linux images this package is often not installed by default — install it and retry:\n"+
				"  Debian/Ubuntu: apt install ca-certificates\n"+
				"  Alpine:        apk add ca-certificates\n"+
				"  RHEL/Fedora:   dnf install ca-certificates", err)
		}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download from %s (status %d)", downloadURL, resp.StatusCode)
	}

	tmpFile := filepath.Join(destDir, "singbox_tmp"+fileExt)
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	if fileExt == ".tar.gz" {
		return extractTarGz(tmpFile, destDir)
	}
	return extractZip(tmpFile, destDir)
}

func extractTarGz(tarGzPath, destDir string) error {
	f, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg && strings.HasSuffix(header.Name, "sing-box") {
			target := filepath.Join(destDir, "sing-box")
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, tr)
			outFile.Close()
			return err
		}
	}
	return fmt.Errorf("sing-box binary not found inside archive")
}

func extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "sing-box.exe") {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			target := filepath.Join(destDir, "sing-box.exe")
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				rc.Close()
				return err
			}
			_, err = io.Copy(outFile, rc)
			rc.Close()
			outFile.Close()
			return err
		}
	}
	return fmt.Errorf("sing-box.exe not found inside zip")
}
