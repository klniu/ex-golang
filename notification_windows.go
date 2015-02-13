// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://msdn.microsoft.com/en-us/library/aa365261(v=vs.85).aspx
package main

import (
	"log"
	"path/filepath"
	"syscall"
	"unsafe"
)

var INFINITE int32 = -1

// maximum number of jobs to run at once
const (
	MAXBG int32 = 6
)

// indicates the event that caused the function to return
var (
	WAIT_OBJECT_0    uint32 = 0
	WAIT_ABANDONED_0 uint32 = 0x00000080
	WAIT_TIMEOUT     uint32 = 0x00000102
	WAIT_FAILED      uint32 = 0xFFFFFFFF
)

// The filter conditions that satisfy a change notification wait. This parameter can be one or more of the following values.
var (
	FILE_NOTIFY_CHANGE_FILE_NAME  uint32 = 0x00000001
	FILE_NOTIFY_CHANGE_DIR_NAME   uint32 = 0x00000002
	FILE_NOTIFY_CHANGE_ATTRIBUTES uint32 = 0x00000004
	FILE_NOTIFY_CHANGE_SIZE       uint32 = 0x00000008
	FILE_NOTIFY_CHANGE_LAST_WRITE uint32 = 0x00000010
	FILE_NOTIFY_CHANGE_SECURITY   uint32 = 0x00000100
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procFindFirstChangeNotification = kernel32.NewProc("FindFirstChangeNotificationW")
	procFindNextChangeNotification  = kernel32.NewProc("FindNextChangeNotification")
	procFindCloseChangeNotification = kernel32.NewProc("FindCloseChangeNotification")
	procReadDirectoryChangesW       = kernel32.NewProc("ReadDirectoryChangesW")
	procWaitForMultipleObjects      = kernel32.NewProc("WaitForMultipleObjects")
)

func boolToUint32(b bool) uint32 {
	if b {
		return 1
	} else {
		return 0
	}
}

// HANDLE WINAPI FindFirstChangeNotification(
//   _In_  LPCTSTR lpPathName,
//   _In_  BOOL bWatchSubtree,
//   _In_  DWORD dwNotifyFilter
// );
func FindFirstChangeNotification(pathName string, watchSubTree bool, mask uint32) (handle syscall.Handle, err error) {
	r1, _, e1 := syscall.Syscall(
		procFindFirstChangeNotification.Addr(),
		3,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(pathName))),
		uintptr(boolToUint32(watchSubTree)),
		uintptr(mask))

	handle = syscall.Handle(r1)
	if handle == syscall.InvalidHandle {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// BOOL WINAPI FindNextChangeNotification(
//   _In_  HANDLE hChangeHandle
// );
func FindNextChangeNotification(handle syscall.Handle) (b bool) {
	r1, _, _ := syscall.Syscall(
		procFindNextChangeNotification.Addr(),
		1,
		uintptr(handle),
		0,
		0)

	b = r1 != 0
	return
}

// BOOL WINAPI FindCloseChangeNotification(
//   _In_  HANDLE hChangeHandle
// );
func FindCloseChangeNotification(handle syscall.Handle) (b bool) {
	// call 1
	// r1, _, _ := syscall.Syscall(procFindCloseChangeNotification.Addr(), 1, uintptr(handle), 0, 0)

	// call 2
	r1, _, _ := procFindCloseChangeNotification.Call(uintptr(handle))

	b = r1 != 0
	return

}

// BOOL WINAPI ReadDirectoryChangesW(
//   _In_         HANDLE hDirectory,
//   _Out_        LPVOID lpBuffer,
//   _In_         DWORD nBufferLength,
//   _In_         BOOL bWatchSubtree,
//   _In_         DWORD dwNotifyFilter,
//   _Out_opt_    LPDWORD lpBytesReturned,
//   _Inout_opt_  LPOVERLAPPED lpOverlapped,
//   _In_opt_     LPOVERLAPPED_COMPLETION_ROUTINE lpCompletionRoutine
// );
func ReadDirectoryChanges(handle syscall.Handle, buf *byte, buflen uint32, watchSubTree bool, mask uint32, retlen *uint32, overlapped *syscall.Overlapped, completionRoutine uintptr) (err error) {
	r1, _, e1 := syscall.Syscall9(
		procReadDirectoryChangesW.Addr(),
		8,
		uintptr(handle),
		uintptr(unsafe.Pointer(buf)),
		uintptr(buflen),
		uintptr(boolToUint32(watchSubTree)),
		uintptr(mask),
		uintptr(unsafe.Pointer(retlen)),
		uintptr(unsafe.Pointer(overlapped)),
		uintptr(completionRoutine),
		0)

	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// DWORD WINAPI WaitForMultipleObjects(
//   _In_  DWORD nCount,
//   _In_  const HANDLE *lpHandles,
//   _In_  BOOL bWaitAll,
//   _In_  DWORD dwMilliseconds
// );
func WaitForMultipleObjects(count uint32, handles *[MAXBG]syscall.Handle, waitAll bool, milliseconds int32) (status uint32, err error) {
	r1, _, e1 := syscall.Syscall6(
		procWaitForMultipleObjects.Addr(),
		4,
		uintptr(count),
		uintptr(unsafe.Pointer(handles)),
		uintptr(boolToUint32(waitAll)),
		uintptr(milliseconds),
		0,
		0)

	status = uint32(r1)
	if status == WAIT_FAILED {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WatchDirectory(dir string) {
	var drive string = filepath.VolumeName("d:/mplus") + "/"
	var err error
	var dwChangeHandles [MAXBG]syscall.Handle
	var flags uint32 = FILE_NOTIFY_CHANGE_FILE_NAME | FILE_NOTIFY_CHANGE_DIR_NAME | FILE_NOTIFY_CHANGE_ATTRIBUTES | FILE_NOTIFY_CHANGE_SIZE | FILE_NOTIFY_CHANGE_LAST_WRITE | FILE_NOTIFY_CHANGE_SECURITY

	// Watch the directory for file creation and deletion.
	dwChangeHandles[0], err = FindFirstChangeNotification(dir, false, flags)
	if err != nil {
		log.Fatal(err)
	}

	// Watch the subtree for directory creation and deletion.
	dwChangeHandles[1], err = FindFirstChangeNotification(drive, true, flags)
	if err != nil {
		log.Fatal(err)
	}

	// Change notification is set. Now wait on both notification
	// handles and refresh accordingly.
	for {
		// Wait for notification.
		log.Printf("Waiting for notification...\n")

		dwWaitStatus, err := WaitForMultipleObjects(2, &dwChangeHandles, false, INFINITE)
		if err != nil {
			log.Fatal(err)
		}

		switch dwWaitStatus {
		case WAIT_OBJECT_0:

			// A file was created, renamed, or deleted in the directory.
			// Refresh this directory and restart the notification.

			RefreshDirectory(dir)
			if FindNextChangeNotification(dwChangeHandles[0]) == false {
				log.Fatal("ERROR: FindNextChangeNotification function failed.\n")
			}

		case WAIT_OBJECT_0 + 1:

			// A directory was created, renamed, or deleted.
			// Refresh the tree and restart the notification.

			RefreshTree(drive)
			if FindNextChangeNotification(dwChangeHandles[1]) == false {
				log.Fatal("ERROR: FindNextChangeNotification function failed.\n")
			}

		case WAIT_TIMEOUT:

			// A timeout occurred, this would happen if some value other
			// than INFINITE is used in the Wait call and no changes occur.
			// In a single-threaded environment you might not want an
			// INFINITE wait.

			log.Printf("No changes in the timeout period.\n")

		default:
			log.Fatal("ERROR: Unhandled dwWaitStatus.\n")
		}
	}
}

func RefreshDirectory(dir string) {
	// This is where you might place code to refresh your
	// directory listing, but not the subtree because it
	// would not be necessary.

	log.Printf("Directory (%s) changed.\n", dir)
}

func RefreshTree(drive string) {
	// This is where you might place code to refresh your
	// directory listing, including the subtree.

	log.Printf("Directory tree (%s) changed.\n", drive)
}

func main() {
	log.SetFlags(log.Lshortfile)
	WatchDirectory("E:/home/LD/zhgo/src/github.com/liudng/examples")
}
