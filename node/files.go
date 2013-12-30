package node

import (
	"fmt"
	"os"
	"sync"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

    "github.com/fudanchii/kakigoori/event"
)

type AzukiFile struct {
	File *os.File
	lock sync.Mutex
}

func NewAzukiFile(f *os.File) File {
	return &AzukiFile{
		File:       f,
		Event_chan: make(chan *event.Intent, 128),
	}
}

func (f *AzukiFile) InnerFile() File {
	return nil
}

func (f *AzukiFile) SetInode(n *Inode) {
}

func (f *AzukiFile) String() string {
	return fmt.Sprintf("AzukiFile(%s)", f.File.Name())
}

func (f *AzukiFile) Read(buf []byte, off int64) (res fuse.ReadResult, code fuse.Status) {
	f.lock.Lock()
	r := fuse.ReadResultFd(f.File.Fd(), off, len(buf))
	f.lock.Unlock()
	go event.Notify("read", f.File.Name())
	return r, fuse.OK
}

func (f *AzukiFile) Write(data []byte, off int64) (uint32, fuse.Status) {
	f.lock.Lock()
	n, err := f.File.WriteAt(data, off)
	f.lock.Unlock()
	go event.Notify("write", f.File.Name())
	return uint32(n), fuse.ToStatus(err)
}

func (f *AzukiFile) Release() {
	go event.Notify("close", f.File.Name())
	f.lock.Lock()
	f.File.Close()
	f.lock.Unlock()
}

func (f *AzukiFile) Flush() fuse.Status {
	f.lock.Lock()
	newFd, err := syscall.Dup(int(f.File.Fd()))
	f.lock.Unlock()

	if err != nil {
		return fuse.ToStatus(err)
	}

	err = syscall.Close(newFd)

	go event.Notify("flush", f.File.Name())
	return fuse.ToStatus(err)
}

func (f *AzukiFile) Fsync(flags int) (code fuse.Status) {
	f.lock.Lock()
	r := fuse.ToStatus(syscall.Fsync(int(f.File.Fd())))
	f.lock.Unlock()
	go event.Notify("fsync", f.File.Name())
	return r
}

func (f *AzukiFile) Truncate(size uint64) fuse.Status {
	f.lock.Lock()
	r := fuse.ToStatus(syscall.Ftruncate(int(f.File.Fd()), int64(size)))
	f.lock.Unlock()
	go event.Notify("trunc", f.File.Name())
	return r
}

func (f *AzukiFile) Chmod(mode uint32) fuse.Status {
	f.lock.Lock()
	r := fuse.ToStatus(f.File.Chmod(os.FileMode(mode)))
	f.lock.Unlock()
	go event.Notify("chmod", f.File.Name())
	return r
}

func (f *AzukiFile) Chown(uid uint32, gid uint32) fuse.Status {
	f.lock.Lock()
	r := fuse.ToStatus(f.File.Chown(int(uid), int(gid)))
	f.lock.Unlock()
	go event.Notify("chown", f.File.Name())
	return r
}

func (f *AzukiFile) GetAttr(a *fuse.Attr) fuse.Status {
	st := syscall.Stat_t{}
	f.lock.Lock()
	err := syscall.Fstat(int(f.File.Fd()), &st)
	f.lock.Unlock()
	if err != nil {
		return fuse.ToStatus(err)
	}
	a.FromStat(&st)
	go event.Notify("getattr", f.File.Name())
	return fuse.OK
}
