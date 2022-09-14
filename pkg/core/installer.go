package core

import (
	"os"
	"os/exec"
	"path"
)

type Installer struct {
	binDir string
	logger *Logger
}

func NewInstaller(logger *Logger) (*Installer, error) {
	binDir, err := currentInstallDir()
	if err != nil {
		return nil, err
	}

	return &Installer{
		binDir: binDir,
		logger: logger,
	}, nil
}

func init() {
	postUpdateClean()
}

func postUpdateClean() {
	i, err := NewInstaller(NewNullLogger())
	if err != nil {
		return
	}
	i.clean()
}

func currentInstallDir() (string, error) {
	currentPath, err := exec.LookPath("kh")
	if err != nil {
		return "", err
	}

	binDir := path.Dir(currentPath)
	return binDir, nil
}

func (i *Installer) InstallNix(f *os.File) error {
	if _, err := os.Stat(i.binDir); os.IsNotExist(err) {
		return err
	}

	binPath := path.Join(i.binDir, "kh")
	i.clean()

	if err := os.Rename(f.Name(), binPath); err != nil {
		return err
	}

	if err := os.Chmod(binPath, 0775); err != nil {
		return err
	}

	i.logger.info("Installed binary at:", binPath)

	return nil
}

func (i *Installer) clean() {
	oldBin := path.Join(i.binDir, "kh.old")
	if _, err := os.Stat(oldBin); os.IsNotExist(err) {
		return
	}

	if err := os.Remove(oldBin); err != nil {
		i.logger.info("Unable to perform update clean up:", err)
	}
}
