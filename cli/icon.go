package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-io/go-utils/pathutil"
)

func copyIconsToDir(icons models.Icons, outputDir string) error {
	if exist, err := pathutil.IsDirExists(outputDir); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("failed to copy icons, output dir does not exist")
	}

	for iconID, iconPath := range icons {
		if err := copyFile(iconPath, filepath.Join(outputDir, iconID)); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src string, dst string) (err error) {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		err = from.Close()
	}()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		err = to.Close()
	}()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}
