// +build windows

package notify

import (
	"path/filepath"
	"syscall"
)

const (
	Create Event = 1 << (20 + iota)
	Delete
	Write
	Move
	Error

	// Internal events TOOD(rjeczalik): unexport
	//
	// Recursive is used to distinguish recursive eventsets from non-recursive ones.
	recursive

	// Inactive is used to bookkeep child watchpoints in the parent ones, which
	// have been set with an actual filesystem watch. This is to allow for
	// optimizing recursive watchpoint count - a single path can be watched
	// by at most 1 recursive filesystem watch.
	inactive

	// dirmarker TODO(pknap)
	dirmarker
)

// ReadDirectoryChangesW filters.
const (
	FileNotifyChangeFileName   = Event(syscall.FILE_NOTIFY_CHANGE_FILE_NAME)
	FileNotifyChangeDirName    = Event(syscall.FILE_NOTIFY_CHANGE_DIR_NAME)
	FileNotifyChangeAttributes = Event(syscall.FILE_NOTIFY_CHANGE_ATTRIBUTES)
	FileNotifyChangeSize       = Event(syscall.FILE_NOTIFY_CHANGE_SIZE)
	FileNotifyChangeLastWrite  = Event(syscall.FILE_NOTIFY_CHANGE_LAST_WRITE)
	FileNotifyChangeLastAccess = Event(syscall.FILE_NOTIFY_CHANGE_LAST_ACCESS)
	FileNotifyChangeCreation   = Event(syscall.FILE_NOTIFY_CHANGE_CREATION)
	FileNotifyChangeSecurity   = Event(syscallFileNotifyChangeSecurity)
)

// according to: http://msdn.microsoft.com/en-us/library/windows/desktop/aa365465(v=vs.85).aspx
// this flag should be declared in: http://golang.org/src/pkg/syscall/ztypes_windows.go
const syscallFileNotifyChangeSecurity = 0x00000100

// ReadDirectoryChangesW actions.
const (
	FileActionAdded          = Event(syscall.FILE_ACTION_ADDED) << 12
	FileActionRemoved        = Event(syscall.FILE_ACTION_REMOVED) << 12
	FileActionModified       = Event(syscall.FILE_ACTION_MODIFIED) << 14
	FileActionRenamedOldName = Event(syscall.FILE_ACTION_RENAMED_OLD_NAME) << 15
	FileActionRenamedNewName = Event(syscall.FILE_ACTION_RENAMED_NEW_NAME) << 16
)

var osestr = map[Event]string{
	FileNotifyChangeFileName:   "notify.FileNotifyChangeFileName",
	FileNotifyChangeDirName:    "notify.FileNotifyChangeDirName",
	FileNotifyChangeAttributes: "notify.FileNotifyChangeAttributes",
	FileNotifyChangeSize:       "notify.FileNotifyChangeSize",
	FileNotifyChangeLastWrite:  "notify.FileNotifyChangeLastWrite",
	FileNotifyChangeLastAccess: "notify.FileNotifyChangeLastAccess",
	FileNotifyChangeCreation:   "notify.FileNotifyChangeCreation",
	FileNotifyChangeSecurity:   "notify.FileNotifyChangeSecurity",

	FileActionAdded:          "notify.FileActionAdded",
	FileActionRemoved:        "notify.FileActionRemoved",
	FileActionModified:       "notify.FileActionModified",
	FileActionRenamedOldName: "notify.FileActionRenamedOldName",
	FileActionRenamedNewName: "notify.FileActionRenamedNewName",
}

var ekind = map[Event]Event{}

const (
	ObjectUnknown uint8 = iota
	ObjectFile
	ObjectDirectory
)

// TODO(ppknap) : doc.
type event struct {
	pathw   []uint16
	name    string
	objtype uint8
	action  uint32
	filter  uint32
	e       Event
}

func (e *event) Event() Event     { return e.e }
func (e *event) Path() string     { return filepath.Join(syscall.UTF16ToString(e.pathw), e.name) }
func (e *event) Sys() interface{} { return e.objtype }