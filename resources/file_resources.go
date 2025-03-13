package resources

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

var _ Resources = &FileResources{}

type FileResources struct {
	root   string
	dbRoot string
}

func NewFileResources(root string, mkdir bool) (*FileResources, error) {
	fileInfo, err := os.Stat(root)
	if err != nil {
		if _, ok := err.(*fs.PathError); !ok && !mkdir {
			return nil, err
		}
		err = os.MkdirAll(root, 0750)
		if err != nil {
			return nil, err
		}
	} else if !fileInfo.IsDir() {
		return nil, errors.New("path is not a directory")
	}

	return &FileResources{
		root: root,
	}, nil
}

func (res *FileResources) SetDatabaseRoot(rootDb string) {
	res.dbRoot = rootDb
}

func (res *FileResources) Read(name string) ([]byte, error) {
	bytes, err := os.ReadFile(path.Join(res.root, name))

	switch {
	case err == nil:

	case strings.Contains(err.Error(), "no such file or directory"):
		err = ErrNotFound
	}

	return bytes, err
}

func (res *FileResources) Write(name string, data []byte) error {
	if s, _ := os.Stat(res.root); s == nil || !s.IsDir() {
		if err := os.MkdirAll(res.root, 0700); err != nil {
			return fmt.Errorf("cannot create config directory: %w", err)
		}
	}

	return os.WriteFile(path.Join(res.root, name), data, 0600)
}

func (res *FileResources) Root() string {
	return res.root
}

func (res *FileResources) DatabaseRoot() string {
	if len(res.dbRoot) > 0 {
		return res.dbRoot
	}

	return res.Root()
}
