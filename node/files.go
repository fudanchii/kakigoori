package node

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	"github.com/fudanchii/kakigoori/event"
)

type AzukiFile struct {
	File *os.File
	lock sync.Mutex
}

func NewAzukiFile(f *os.File) nodefs.File {
	return &AzukiFile{
		File: f,
	}
}

func (f *AzukiFile) InnerFile() nodefs.File {
	return nil
}

func (f *AzukiFile) SetInode(n *nodefs.Inode) {
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
	return fuse.OK
}

func (f *AzukiFile) Allocate(off uint64, sz uint64, mode uint32) fuse.Status {
	f.lock.Lock()
	err := syscall.Fallocate(int(f.File.Fd()), mode, int64(off), int64(sz))
	f.lock.Unlock()
	if err != nil {
		return fuse.ToStatus(err)
	}
	go event.Notify("fallocate", f.File.Name())
	return fuse.OK
}

const _UTIME_NOW = ((1 << 30) - 1)
const _UTIME_OMIT = ((1 << 30) - 2)

func (f *AzukiFile) Utimens(a *time.Time, m *time.Time) fuse.Status {
	tv := make([]syscall.Timeval, 2)
	if a == nil {
		tv[0].Usec = _UTIME_OMIT
	} else {
		n := a.UnixNano()
		tv[0] = syscall.NsecToTimeval(n)
	}

	if m == nil {
		tv[1].Usec = _UTIME_OMIT
	} else {
		n := a.UnixNano()
		tv[1] = syscall.NsecToTimeval(n)
	}

	err := syscall.Futimes(int(f.File.Fd()), tv)
	return fuse.ToStatus(err)
}
