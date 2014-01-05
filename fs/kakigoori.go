package fs

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"github.com/fudanchii/kakigoori/event"
	"github.com/fudanchii/kakigoori/node"
)

type KakigooriFileSystem struct {
	pathfs.FileSystem
	Root string
}

func NewKakigooriFileSystem(root string) pathfs.FileSystem {
	return &KakigooriFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
		Root:       root,
	}
}

func (fs *KakigooriFileSystem) OnMount(nodeFs *pathfs.PathNodeFs) {}

func (fs *KakigooriFileSystem) OnUnmount() {}

func (fs *KakigooriFileSystem) GetPath(relPath string) string {
	return filepath.Join(fs.Root, relPath)
}

func (fs *KakigooriFileSystem) GetAttr(name string, context *fuse.Context) (a *fuse.Attr, code fuse.Status) {
	fullPath := fs.GetPath(name)
	var err error = nil
	st := syscall.Stat_t{}
	if name == "" {
		err = syscall.Stat(fullPath, &st)
	} else {
		err = syscall.Lstat(fullPath, &st)
	}

	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	a = &fuse.Attr{}
	a.FromStat(&st)
	return a, fuse.OK
}

func (fs *KakigooriFileSystem) OpenDir(name string, context *fuse.Context) (stream []fuse.DirEntry, status fuse.Status) {
	fullPath := fs.GetPath(name)
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	want := 500
	output := make([]fuse.DirEntry, 0, want)
	for {
		infos, err := f.Readdir(want)
		for i := range infos {
			// workaround for https://code.google.com/p/go/issues/detail?id=5960
			if infos[i] == nil {
				continue
			}
			n := infos[i].Name()
			d := fuse.DirEntry{
				Name: n,
			}
			if s := fuse.ToStatT(infos[i]); s != nil {
				d.Mode = uint32(s.Mode)
			} else {
				log.Printf("ReadDir entry %q for %q has no stat info", n, name)
			}
			output = append(output, d)
		}
		if len(infos) < want || err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Readdir() returned err:", err)
			break
		}
	}
	f.Close()
	go event.Notify(event.OpenDir, fullPath)
	return output, fuse.OK
}

func (fs *KakigooriFileSystem) Open(name string, flags uint32, context *fuse.Context) (fuseFile nodefs.File, status fuse.Status) {
	fullPath := fs.GetPath(name)
	f, err := os.OpenFile(fullPath, int(flags), 0)
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	go event.Notify(event.Open, fullPath)
	return node.NewAzukiFile(f), fuse.OK
}

func (fs *KakigooriFileSystem) Chmod(path string, mode uint32, context *fuse.Context) (code fuse.Status) {
	fullPath := fs.GetPath(path)
	err := os.Chmod(fullPath, os.FileMode(mode))
	go event.Notify(event.Chmod, fullPath)
	return fuse.ToStatus(err)
}

func (fs *KakigooriFileSystem) Chown(path string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	fullPath := fs.GetPath(path)
	go event.Notify(event.Chown, fullPath)
	return fuse.ToStatus(os.Chown(fullPath, int(uid), int(gid)))
}

func (fs *KakigooriFileSystem) Truncate(path string, offset uint64, context *fuse.Context) (code fuse.Status) {
	fullPath := fs.GetPath(path)
	go event.Notify(event.Trunc, fullPath)
	return fuse.ToStatus(os.Truncate(fullPath, int64(offset)))
}

func (fs *KakigooriFileSystem) Utimens(path string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	var a time.Time
	if Atime != nil {
		a = *Atime
	}
	var m time.Time
	if Mtime != nil {
		m = *Mtime
	}
	fullPath := fs.GetPath(path)
	return fuse.ToStatus(os.Chtimes(fullPath, a, m))
}

func (fs *KakigooriFileSystem) Readlink(name string, context *fuse.Context) (out string, code fuse.Status) {
	fullPath := fs.GetPath(name)
	f, err := os.Readlink(fullPath)
	go event.Notify(event.Readlink, fullPath)
	return f, fuse.ToStatus(err)
}

func (fs *KakigooriFileSystem) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) (code fuse.Status) {
	fullPath := fs.GetPath(name)
	go event.Notify(event.Mknod, fullPath)
	return fuse.ToStatus(syscall.Mknod(fullPath, mode, int(dev)))
}

func (fs *KakigooriFileSystem) Mkdir(path string, mode uint32, context *fuse.Context) (code fuse.Status) {
	fullPath := fs.GetPath(path)
	go event.Notify(event.Mkdir, fullPath)
	return fuse.ToStatus(os.Mkdir(fullPath, os.FileMode(mode)))
}

func (fs *KakigooriFileSystem) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	fullPath := fs.GetPath(name)
	go event.Notify(event.Unlink, fullPath)
	return fuse.ToStatus(syscall.Unlink(fullPath))
}

func (fs *KakigooriFileSystem) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	fullPath := fs.GetPath(name)
	go event.Notify(event.Rmdir, fullPath)
	return fuse.ToStatus(syscall.Rmdir(fullPath))
}

func (fs *KakigooriFileSystem) Symlink(pointedTo string, linkName string, context *fuse.Context) (code fuse.Status) {
	fullPath := fs.GetPath(pointedTo)
	linkPath := fs.GetPath(linkName)
	go event.Notify(event.Symlink, fmt.Sprintf("%s -> %s", fullPath, linkPath))
	return fuse.ToStatus(os.Symlink(pointedTo, linkPath))
}

func (fs *KakigooriFileSystem) Rename(oldPath string, newPath string, context *fuse.Context) (codee fuse.Status) {
	fullOldPath := fs.GetPath(oldPath)
	fullNewPath := fs.GetPath(newPath)
	err := os.Rename(fullOldPath, fullNewPath)
	go event.Notify(event.Rename, fmt.Sprintf("%s -> %s", fullOldPath, fullNewPath))
	return fuse.ToStatus(err)
}

func (fs *KakigooriFileSystem) Link(orig string, newName string, context *fuse.Context) (code fuse.Status) {
	fullOrig := fs.GetPath(orig)
	fullNewName := fs.GetPath(newName)
	go event.Notify(event.Link, fmt.Sprintf("%s -> %s", fullOrig, fullNewName))
	return fuse.ToStatus(os.Link(fullOrig, fullNewName))
}

func (fs *KakigooriFileSystem) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	fullPath := fs.GetPath(name)
	go event.Notify(event.Access, fullPath)
	return fuse.ToStatus(syscall.Access(fullPath, mode))
}

func (fs *KakigooriFileSystem) Create(path string, flags uint32, mode uint32, context *fuse.Context) (fuseFile nodefs.File, code fuse.Status) {
	fullPath := fs.GetPath(path)
	f, err := os.OpenFile(fullPath, int(flags)|os.O_CREATE, os.FileMode(mode))
	go event.Notify(event.Create, fullPath)
	return node.NewAzukiFile(f), fuse.ToStatus(err)
}

func (fs *KakigooriFileSystem) StatFs(name string) *fuse.StatfsOut {
	s := syscall.Statfs_t{}
	err := syscall.Statfs(fs.GetPath(name), &s)
	if err == nil {
		return &fuse.StatfsOut{
			Blocks:  s.Blocks,
			Bsize:   uint32(s.Bsize),
			Bfree:   s.Bfree,
			Bavail:  s.Bavail,
			Files:   s.Files,
			Ffree:   s.Ffree,
			Frsize:  uint32(s.Frsize),
			NameLen: uint32(s.Namelen),
		}
	}
	return nil
}

func (fs *KakigooriFileSystem) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {
	data, err := listXAttr(fs.GetPath(name))
	return data, fuse.ToStatus(err)
}

func (fs *KakigooriFileSystem) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {
	err := sysRemovexattr(fs.GetPath(name), attr)
	return fuse.ToStatus(err)
}

func (fs *KakigooriFileSystem) String() string {
	return fmt.Sprintf("KakigooriFs(%s)", fs.Root)
}

func (fs *KakigooriFileSystem) GetXAttr(name string, attr string, context *fuse.Context) ([]byte, fuse.Status) {
	data := make([]byte, 1024)
	data, err := getXAttr(fs.GetPath(name), attr, data)
	return data, fuse.ToStatus(err)
}
