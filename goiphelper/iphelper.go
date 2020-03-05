package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

var (
	modiphelper = windows.NewLazyDLL("Iphlpapi.dll")

	procGetExtendedTcpTable = modiphelper.NewProc("GetExtendedTcpTable")

)


type mibTcprowOwnerPid struct {
	/*  C declaration
		DWORD       dwState;
		DWORD       dwLocalAddr;
		DWORD       dwLocalPort;
		DWORD       dwRemoteAddr;
		DWORD       dwRemotePort;
		DWORD       dwOwningPid; */
	dwState			uint32
	dwLocalAddr		uint32			// network byte order
	dwLocalPort		uint32			// network byte order
	dwRemoteAddr	uint32			// network byte order
	dwRemotePort	uint32			// network byte order
	dwOwningPid		uint32			
}

const (
	TCP_TABLE_BASIC_LISTENER 			= uint32(0)
    TCP_TABLE_BASIC_CONNECTIONS     	= uint32(1)
    TCP_TABLE_BASIC_ALL					= uint32(2)
    TCP_TABLE_OWNER_PID_LISTENER    	= uint32(3)
    TCP_TABLE_OWNER_PID_CONNECTIONS 	= uint32(4)
    TCP_TABLE_OWNER_PID_ALL				= uint32(5)
    TCP_TABLE_OWNER_MODULE_LISTENER 	= uint32(6)
    TCP_TABLE_OWNER_MODULE_CONNECTIONS 	= uint32(7)
    TCP_TABLE_OWNER_MODULE_ALL			= uint32(8)
)
func getTcpTable() (table map[uint32][]mibTcprowOwnerPid, err error) {
	var size uint32
	var rawtableentry uintptr
	r, _, _ := procGetExtendedTcpTable.Call(rawtableentry, 
										   uintptr(unsafe.Pointer(&size)),
										   uintptr(0), // false, unsorted
										   uintptr(syscall.AF_INET),
										   uintptr(TCP_TABLE_OWNER_PID_ALL),
										   uintptr(0))

	if r != uintptr(windows.ERROR_INSUFFICIENT_BUFFER) {
		err = fmt.Errorf("Unexpected error %v", r)
		return
	}
	fmt.Printf("Size is %d\n", size)
	rawbuf := make([]byte, size)
	r, _, _ = procGetExtendedTcpTable.Call(uintptr(unsafe.Pointer(&rawbuf[0])),
		uintptr(unsafe.Pointer(&size)),
		uintptr(0), // false, unsorted
		uintptr(syscall.AF_INET),
		uintptr(TCP_TABLE_OWNER_PID_ALL),
		uintptr(0))
	if r != 0 {
		err = fmt.Errorf("Unexpected error %v", r)
		return
	}
	count := uint32(binary.LittleEndian.Uint32(rawbuf))
	fmt.Printf("Count is %d\n", count)
	table = make(map[uint32][]mibTcprowOwnerPid)

	entries := (*[1 << 30]mibTcprowOwnerPid)(unsafe.Pointer(&rawbuf[4]))[:count:count]
	for _, entry := range entries {
		pid := entry.dwOwningPid
		
		table[pid] = append(table[pid], entry)

	}
	return table, nil

}


func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

func Ntohs(i uint16) uint16 {
	return binary.BigEndian.Uint16((*(*[2]byte)(unsafe.Pointer(&i)))[:])
}

func Ntohl(i uint32) uint32 {
	return binary.BigEndian.Uint32((*(*[4]byte)(unsafe.Pointer(&i)))[:])
}
func Htonl(i uint32) uint32 {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return *(*uint32)(unsafe.Pointer(&b[0]))
}
func main() {
	t, err := getTcpTable()
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	for pid, list := range t {
		for _, entry := range list {
			fmt.Printf("Pid (%8v)    %16v:%5v  %16v:%5v\n",
						pid,
						int2ip(Ntohl(entry.dwLocalAddr)).String(),
						Ntohs(uint16(entry.dwLocalPort)),
						int2ip(Ntohl(entry.dwRemoteAddr)).String(),
						Ntohs(uint16(entry.dwRemotePort)))
		}

	}
	return
}