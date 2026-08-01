package tun

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	ifNameSize = 16

	// Linux ioctl request to create/configure a TUN/TAP interface.
	tunSetIFF = 0x400454ca

	// Create a TUN device (Layer 3).
	iffTUN = 0x0001

	// Do not prepend packet information.
	iffNoPI = 0x1000
)

type ifreq struct {
	Name  [ifNameSize]byte
	Flags uint16
	Pad   [22]byte
}

// Device represents a Linux TUN interface.
type Device struct {
	Name string
	File *os.File
}

// Open creates (or opens) a Linux TUN interface.
func Open(name string) (*Device, error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("open /dev/net/tun: %w", err)
	}

	var req ifreq

	copy(req.Name[:], name)
	req.Flags = iffTUN | iffNoPI

	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		file.Fd(),
		uintptr(tunSetIFF),
		uintptr(unsafe.Pointer(&req)),
	)

	if errno != 0 {
		file.Close()
		return nil, fmt.Errorf("TUNSETIFF: %v", errno)
	}

	dev := &Device{
		Name: string(req.Name[:]),
		File: file,
	}

	return dev, nil
}

// Read reads an IP packet from the TUN interface.
func (d *Device) Read(buf []byte) (int, error) {
	return d.File.Read(buf)
}

// Write writes an IP packet to the TUN interface.
func (d *Device) Write(buf []byte) (int, error) {
	return d.File.Write(buf)
}

// Close closes the TUN interface.
func (d *Device) Close() error {
	if d.File == nil {
		return nil
	}

	return d.File.Close()
}
