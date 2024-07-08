package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
)

const dirName = "./drive"

type File struct {
	Name         string `json:"name,omitempty"`
	Path         string `json:"path,omitempty"`
	IsFolder     bool   `json:"is_folder,omitempty"`
	Size         string `json:"size"`
	LastModified string `json:"last_modified,omitempty"`
	Items        []File `json:"items,omitempty"`
}

var Files []File

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

func GetRootItems() ([]File, error) {
	var files []File
	err := filepath.WalkDir(dirName, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() {
			files = append(files, File{
				Name:         d.Name(),
				Path:         filepath.ToSlash(path),
				IsFolder:     true,
				Items:        nil,
				Size:         convertSizeToMB(fileInfo.Size()),
				LastModified: fileInfo.ModTime().String(),
			})
		} else {
			parentDir := filepath.ToSlash(filepath.Dir(path))
			for i := range files {
				if files[i].Path == parentDir {
					files[i].Items = append(files[i].Items, File{
						Name:         d.Name(),
						Path:         filepath.ToSlash(path),
						IsFolder:     false,
						Items:        nil,
						Size:         convertSizeToMB(fileInfo.Size()),
						LastModified: fileInfo.ModTime().String(),
					})
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func CreateFolder(folderName string) error {
	return os.Mkdir(filepath.Join(dirName, folderName), 0755)
}
