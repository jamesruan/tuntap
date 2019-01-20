/*
 * Simple TUN/TAP wrapper for golang
 * Copyright (C) 2019  James Ruan <ruanbeihong@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 **/

// Package tuntap provides a simple wrapper for TUN/TAP device.
package tuntap

/*
#include <stdlib.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <fcntl.h>
#include <string.h>
#include <linux/if.h>
#include <linux/if_tun.h>

static int alloc_internal(int fd, const char* dev, char *dev_o, short flags) {
	struct ifreq ifr;
	int err;

	memset(&ifr, 0, sizeof(ifr));
	if (*dev) strncpy(ifr.ifr_name, dev, IFNAMSIZ);
	ifr.ifr_flags = flags;

	if ((err = ioctl(fd, TUNSETIFF, (void *) &ifr)) < 0) {
		close(fd);
		return err;
	}

	strcpy(dev_o, ifr.ifr_name);
	return 0;
}

int alloc_tuntap(_GoString_ dev, char* dev_o, int tap) {
	const char * c_dev = _GoStringPtr(dev);
	int fd, err;

	if((fd = open("/dev/net/tun", O_RDWR|O_NONBLOCK|O_DSYNC)) < 0) {
		return -1;
	}

	if ((err = alloc_internal(fd, c_dev, dev_o, (tap?IFF_TAP:IFF_TUN) | IFF_NO_PI)) < 0) {
		return -1;
	}
	return fd;
}
*/
import "C"
import (
	"github.com/pkg/errors"
	"os"
	"unsafe"
)

// TunTap provides a File interface for TUN/TAP device
type TunTap struct {
	*os.File
	Tap bool
}

// NewTunTap creates a TUN/TAP device.
func NewTunTap(name string, tap bool) (*TunTap, error) {
	var err error
	var fd C.int

	name_o := (*C.char)(C.malloc(C.IFNAMSIZ + 1))
	defer C.free(unsafe.Pointer(name_o))

	if tap {
		fd, err = C.alloc_tuntap(name, name_o, 1)
	} else {
		fd, err = C.alloc_tuntap(name, name_o, 0)
	}

	if err != nil {
		return nil, errors.Errorf("can't allocate device %s: %s", name, err)
	}
	return &TunTap{
		File: os.NewFile(uintptr(fd), C.GoString(name_o)),
		Tap:  tap,
	}, nil
}
