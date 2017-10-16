/* Copyright 2016 Huitse Tai. All rights reserved.
 * Use of this source code is governed by BSD 3-clause
 * license that can be found in the LICENSE file.
 */

#if defined(__linux__) && __linux__

#include "tun_fd.h"
#include "setmtu_unix.h"

#include <linux/if.h>
#include <linux/if_tun.h>
#include <sys/ioctl.h>
#include <stddef.h>
#include <string.h>
#include <errno.h>

int const IFPKT_OFFSET = 0;

tun_fd_info create_tun_fd(int mtu) {

  tun_fd_info info;
  memset(&info, 0, sizeof(info));

  struct ifreq ifr;
  int fd, err;

  if ((fd = open("/dev/net/tun", O_RDWR)) < 0) {
    errno = errno;
    return info;
  }

  memset(&ifr, 0, sizeof(ifr));

  ifr.ifr_flags = IFF_TUN | IFF_NO_PI;

  if ((ioctl(fd, TUNSETIFF, &ifr)) < 0) {

    int err = errno;
    close(fd);
    errno = err;
    return info;
  }

  if (setmtu_unix(ifr.ifr_name, mtu) < 0) {

    int err = errno;
    close(fd);
    errno = err;
    return info;
  }

  info.fd = fd;
  strncpy(info.nam, ifr.ifr_name, IFNAMSIZ);

  errno = 0;
  return info;
}

int const built_linux = 0;
#else

int const built_linux = 0;
#endif /*linux guard*/
