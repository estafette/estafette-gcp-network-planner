// Code generated for package network by go-bindata DO NOT EDIT. (@generated)
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

var _configJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x92\x41\x6b\xfa\x40\x10\xc5\xef\xf9\x14\x8f\x9c\xfe\x7f\xb0\xd1\x04\x9b\xaa\xf7\x1e\x0a\x3d\xf4\x58\x28\x25\xac\xc9\x24\x2e\x6d\x76\xc2\xec\x46\x11\xf1\xbb\x97\x8d\x4a\x2a\x6d\xea\x41\x5a\x02\x39\xcc\x0c\x6f\x7e\x6f\xde\xee\x02\x20\x14\x65\x2a\xca\x72\x36\xa5\xae\x6c\xb8\xc0\x4b\x00\x00\xbb\xee\x0f\x84\x6e\xdb\x50\xb8\x40\x68\xb8\xa0\x70\x74\xaa\x0a\x55\x9a\x8d\xaf\x53\x2b\xdc\xd0\xcd\x86\xac\x8b\xfb\xbe\x6e\xb2\x5c\x17\x92\x1d\xc4\x4f\x1a\x8d\xe8\x5a\xc9\xb6\x1f\xcb\xb9\xae\xc9\x38\xdf\xbb\x5f\x93\x6c\x61\xdb\xa5\x21\x87\x95\xb2\x28\xb9\x15\x08\x59\x92\x35\x15\x78\x78\x82\x2a\x0a\x21\x6b\xc9\x42\x1b\x68\x67\x71\x94\xf3\xbd\x6e\x4f\xaf\x6b\xc8\x6d\x58\xde\xbc\x6e\x7c\x97\x44\xc9\x2c\x9a\x44\x93\x71\x3c\xed\x27\x0e\x8b\xb2\x5a\x59\x3f\x95\xc4\x5d\x7d\x3f\xfa\xea\xbb\xe1\xe2\x5a\xdb\x96\x72\x36\xc5\x99\xf1\xcf\x80\x93\xa8\xfb\xc6\xf3\x6f\xef\xf2\xa8\xa4\x22\x90\xe1\xb6\x5a\xc1\x31\x4a\xed\xf0\xcf\xb1\x53\xef\x30\x6d\xbd\x24\x01\x97\xf0\xe1\x58\x3c\x23\xb9\x4d\xff\x9f\x9d\x6a\xc8\x70\x9c\x0e\x1a\xf6\x07\xd7\xf9\xd5\x59\x5f\x30\xed\x53\x99\x5e\x4a\x25\x19\x84\xac\x95\x75\x24\xbf\xcc\x38\x4f\xa2\x38\x3d\x3e\x9d\xd9\x20\xe4\x6c\x10\x92\xdd\xea\xaf\x18\xd3\xe9\xcf\x90\xc7\x4b\x06\xc0\x6b\xb0\xff\x08\x00\x00\xff\xff\x0c\x9f\x85\xcc\xf7\x03\x00\x00")

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

	info := bindataFileInfo{name: "config.json", size: 1015, mode: os.FileMode(420), modTime: time.Unix(1600091660, 0)}
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
