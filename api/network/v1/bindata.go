// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// config.json
package network

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _configJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x91\x41\x6b\xa3\x50\x14\x85\xf7\xfe\x8a\x83\xab\x19\xc8\x64\x46\x71\x6c\x92\x7d\x17\x85\x2e\xba\x2c\x94\x22\x2f\x7a\x35\x8f\xd6\x77\xe5\xbe\x67\x42\x08\xf9\xef\xe5\x19\x83\x4d\xa8\xcd\x22\xb4\x45\x70\x71\x8f\xdc\xfb\x7d\xc7\x5d\x00\x84\xa2\x4c\x45\x59\xce\xa6\xd4\x95\x0d\x17\x78\x0a\x00\x60\xd7\xbd\x81\xd0\x6d\x1b\x0a\x17\x08\x1b\x2e\xc2\xc9\x71\x28\x54\x69\x36\x7e\x4c\xad\x70\x43\x7f\x36\x64\x5d\x34\xe4\xba\xc9\x72\x5d\x48\x76\xd8\x7d\x5c\x61\x29\x67\x53\x28\xd9\x0e\x1f\x1a\x72\x1b\x96\x17\x9f\x46\xff\xa6\xdd\xf3\x77\x3e\xc4\x39\xd7\x35\x19\xe7\xe3\x7b\x25\x15\x81\x0c\xb7\xd5\x0a\x8e\x51\x6a\x87\x5f\x8e\x9d\x7a\x85\x69\xeb\x25\x09\xb8\x84\xe1\x82\x2c\x1e\x11\xff\x4f\x7f\xe3\xee\x01\xaa\x28\x84\xac\x25\x3b\xac\xb4\xed\xd2\x90\xcb\x6a\x65\xfd\xd5\x28\xed\xe6\xfb\xc9\xcf\x08\x47\xf1\xec\x5c\xf9\x8c\x2f\x19\xe5\xf3\xae\xd7\x02\x36\xa2\xeb\x13\xbc\x77\x85\xdf\xae\x49\xb6\x38\xe0\x60\xa5\x2c\x4a\x6e\x05\x42\x96\x64\x4d\xc5\x49\xbb\xd0\x06\xda\x59\xf4\xeb\x7c\xd6\xdd\xf9\x58\xfb\x26\x9e\xf6\xda\x51\x32\xe6\x1d\x47\xa3\xde\xfe\xbe\xce\xaf\x56\xbf\xf0\x6f\x3c\x64\x72\x09\x32\x1e\x85\xac\x95\x75\x24\x5f\xcc\x38\x8f\xa7\x51\xda\x37\x39\x1b\x85\x9c\x8d\x42\xb2\x5b\x7d\x17\x63\x9a\x7c\x0e\xd9\x37\x19\x00\xcf\xc1\xfe\x2d\x00\x00\xff\xff\x2e\x95\xa8\x88\x97\x04\x00\x00")

func configJsonBytes() ([]byte, error) {
	return bindataRead(
		_configJson,
		"config.json",
	)
}

func configJson() (*asset, error) {
	bytes, err := configJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "config.json", size: 1175, mode: os.FileMode(420), modTime: time.Unix(1600086973, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"config.json": configJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"config.json": &bintree{configJson, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
