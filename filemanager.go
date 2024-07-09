package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
)

var DefaultRoot = "./drive"

type File struct {
	Name         string `json:"name,omitempty"`
	Path         string `json:"path,omitempty"`
	IsFolder     bool   `json:"is_folder"`
	Size         string `json:"size"`
	DataModified string `json:"data_modified,omitempty"`
	Items        []File `json:"items,omitempty"`
}

func convertSizeToMB(size int64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	}

	return fmt.Sprintf("%d bytes", size)
}

func GetRootFolders() ([]File, error) {
	var folders []File

	err := filepath.WalkDir(DefaultRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == DefaultRoot {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() && filepath.Dir(path) == DefaultRoot {
			folders = append(folders, File{
				Name:         d.Name(),
				Path:         filepath.ToSlash(path),
				IsFolder:     true,
				Items:        nil,
				Size:         convertSizeToMB(fileInfo.Size()),
				DataModified: fileInfo.ModTime().String(),
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return folders, nil
}

func GetFolderItems(folderPath string) ([]File, error) {
	var items []File

	err := filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == folderPath {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		items = append(items, File{
			Name:         d.Name(),
			Path:         filepath.ToSlash(path),
			IsFolder:     d.IsDir(),
			Size:         convertSizeToMB(fileInfo.Size()),
			DataModified: fileInfo.ModTime().String(),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func CreateFolder(folderName string) error {
	return os.Mkdir(filepath.Join(DefaultRoot, folderName), 0755)
}

func DeleteItem(files []File, path string) error {
	exeDir, err := os.Executable()
	for i := range files {
		if files[i].Path == path {
			if files[i].IsFolder {
				err := os.RemoveAll(path)
				if err != nil {
					return err
				}
			} else {
				if err != nil {
					return err
				}

				err = os.Remove(filepath.Join(filepath.Base(exeDir), path))
				if err != nil {
					return err
				}
			}
			return nil
		}

		for j := range files[i].Items {
			if files[i].Items[j].Path == path {
				if files[i].IsFolder {
					err := os.RemoveAll(path)
					if err != nil {
						return err
					}
				} else {
					if err != nil {
						return err
					}

					err = os.Remove(filepath.Join(filepath.Base(exeDir), path))
					if err != nil {
						return err
					}
				}
				return nil
			}
		}
	}
	return fmt.Errorf("file or directory not found")
}
