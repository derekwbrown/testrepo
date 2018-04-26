package main

/*
#cgo LDFLAGS: -l wevtapi
#include "event.h"
*/
import "C"
import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var (
	modWinEvtApi = syscall.NewLazyDLL("wevtapi.dll")

	procEvtSubscribe       = modWinEvtApi.NewProc("EvtSubscribe")
	procEvtClose           = modWinEvtApi.NewProc("EvtClose")
	procEvtRender          = modWinEvtApi.NewProc("EvtRender")
	procEvtOpenChannelEnum = modWinEvtApi.NewProc("EvtOpenChannelEnum")
	procEvtNextChannelPath = modWinEvtApi.NewProc("EvtNextChannelPath")
	procEvtNext            = modWinEvtApi.NewProc("EvtNext")
)

// EvtRender takes an event handle and reders it to XML
func EvtRender(h C.ULONGLONG) (xml string, err error) {
	var bufSize uint32
	var bufUsed uint32
	var propCount uint32

	ret, _, err := procEvtRender.Call(uintptr(0), // this handle is always null for XML renders
		uintptr(h),                 // handle of event we're rendering
		uintptr(EvtRenderEventXml), // for now, always render in xml
		uintptr(bufSize),
		uintptr(0),                          //no buffer for now, just getting necessary size
		uintptr(unsafe.Pointer(&bufUsed)),   // filled in with necessary buffer size
		uintptr(unsafe.Pointer(&propCount))) // not used but must be provided
	if err != error(syscall.ERROR_INSUFFICIENT_BUFFER) {
		return
	}
	bufSize = bufUsed
	buf := make([]uint8, bufSize)
	ret, _, err = procEvtRender.Call(uintptr(0), // this handle is always null for XML renders
		uintptr(h),                 // handle of event we're rendering
		uintptr(EvtRenderEventXml), // for now, always render in xml
		uintptr(bufSize),
		uintptr(unsafe.Pointer(&buf[0])),    //no buffer for now, just getting necessary size
		uintptr(unsafe.Pointer(&bufUsed)),   // filled in with necessary buffer size
		uintptr(unsafe.Pointer(&propCount))) // not used but must be provided
	if ret == 0 {
		return
	}
	// Call will set error anyway.  Clear it so we don't return an error
	err = nil
	xml = ConvertWindowsString(buf)
	return

}

type eventContext struct {
	name  string
	index int
	// could put the channel here
}

/* These are entry points for the callback to hand the pointer to Go-land.
   Note: handles are only valid within the callback. Don't pass them out. */

//export goStaleCallback
func goStaleCallback(errCode C.ULONGLONG, ctx C.PVOID) {
	fmt.Printf("Stale callback\n")
}

//export goErrorCallback
func goErrorCallback(errCode C.ULONGLONG, ctx C.PVOID) {
	fmt.Printf("Error callback %v\n", errCode)
}

//export goNotificationCallback
func goNotificationCallback(handle C.ULONGLONG, ctx C.PVOID) {
	fmt.Printf("Notification Callback\n")
	var goctx eventContext
	goctx = *(*eventContext)(unsafe.Pointer(uintptr(ctx)))
	fmt.Printf("Callback from %s\n", goctx.name)
	xml, err := EvtRender(handle)
	if err == nil {
		fmt.Printf("Rendered XML: %s\n", xml)
	} else {
		fmt.Printf("Error rendering xml %v\n", err)
	}
	return
}

type EvtSubscribeNotifyAction int32
type EvtSubscribeFlags int32

const (
	EvtSubscribeActionError   EvtSubscribeNotifyAction = 0
	EvtSubscribeActionDeliver EvtSubscribeNotifyAction = 1

	EvtSubscribeOriginMask          EvtSubscribeFlags = 0x3
	EvtSubscribeTolerateQueryErrors EvtSubscribeFlags = 0x1000
	EvtSubscribeStrict              EvtSubscribeFlags = 0x10000

	EvtRenderEventValues = 0 // Variants
	EvtRenderEventXml    = 1 // XML
	EvtRenderBookmark    = 2 // Bookmark

	ERROR_NO_MORE_ITEMS syscall.Errno = 259
)

type EVT_SUBSCRIBE_FLAGS int

const (
	_ = iota
	EvtSubscribeToFutureEvents
	EvtSubscribeStartAtOldestRecord
	EvtSubscribeStartAfterBookmark
)

func main() {
	fmt.Printf("starting event log watcher\n")
	/*
		channels, err := EnumerateChannels()
		if err != nil {
			fmt.Printf("Error enumerating channels %v\n", err)
			return
		}
		for _, ch := range channels {
			fmt.Printf("Channel:  %s\n", ch)
		}
		return
	*/
	ctx := eventContext{
		name:  "Application",
		index: 2,
	}

	C.startEventSubscribe(C.CString("Application"),
		C.CString("*"),
		C.ULONGLONG(0),
		C.int(EvtSubscribeToFutureEvents),
		C.PVOID(uintptr(unsafe.Pointer(&ctx))))
	for {
		time.Sleep(2 * time.Second)
	}

}

// ConvertWindowsString converts a windows c-string
// into a go string.  Even though the input is array
// of uint8, the underlying data is expected to be
// uint16 (unicode)
func ConvertWindowsString(winput []uint8) string {
	var retstring string
	for i := 0; i < len(winput); i += 2 {
		dbyte := (uint16(winput[i+1]) << 8) + uint16(winput[i])
		if dbyte == 0 {
			break
		}
		retstring += string(rune(dbyte))
	}
	return retstring
}

type EvtEnumHandle uintptr

// EnumerateChannels() enumerates available log channels
func EnumerateChannels() (chans []string, err error) {
	err = nil

	ret, _, err := procEvtOpenChannelEnum.Call(uintptr(0), // local computer
		uintptr(0)) // must be zero

	hEnum := EvtEnumHandle(ret)
	if hEnum == EvtEnumHandle(0) {
		return
	}
	defer procEvtClose.Call(uintptr(hEnum))

	for {
		var str string
		str, err = evtNextChannel(hEnum)
		if err == nil {
			chans = append(chans, str)
		} else if err == error(ERROR_NO_MORE_ITEMS) {
			fmt.Printf("setting err to nil\n")
			err = nil
			break
		} else {
			break
		}
	}
	fmt.Printf("returning error %v\n", err)
	return
}

func evtNextChannel(h EvtEnumHandle) (ch string, err error) {

	var bufSize uint32
	var bufUsed uint32

	ret, _, err := procEvtNextChannelPath.Call(uintptr(h), // this handle is always null for XML renders
		uintptr(bufSize),
		uintptr(0),                        //no buffer for now, just getting necessary size
		uintptr(unsafe.Pointer(&bufUsed))) // filled in with necessary buffer size
	if err != error(syscall.ERROR_INSUFFICIENT_BUFFER) {
		fmt.Printf("Next: %v %v", ret, err)
		return
	}
	bufSize = bufUsed
	buf := make([]uint8, bufSize*2)
	ret, _, err = procEvtNextChannelPath.Call(uintptr(h), // handle of event we're rendering
		uintptr(bufSize),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufUsed))) // filled in with necessary buffer size
	if ret == 0 {
		fmt.Printf("Next: %v %v", ret, err)
		return
	}
	err = nil
	// Call will set error anyway.  Clear it so we don't return an error
	ch = ConvertWindowsString(buf)
	return
}
