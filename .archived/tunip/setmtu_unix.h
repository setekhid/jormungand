/* Copyright 2016 Huitse Tai. All rights reserved.
 * Use of this source code is governed by BSD 3-clause
 * license that can be found in the LICENSE file.
 */

#ifndef setmtu_unix_h_INCLUDE
#define setmtu_unix_h_INCLUDE

#include <sys/socket.h>
#include <net/if.h>
#include <sys/ioctl.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>

static int setmtu_unix(char const *ifnam, int mtu) {

  int net_fd = socket(PF_INET, SOCK_DGRAM, 0);
  if (net_fd < 0) {
    errno = EINVAL;
    return -1;
  }

  struct ifreq ifr;
  memset(&ifr, 0, sizeof(ifr));
  strncpy(ifr.ifr_name, ifnam, IFNAMSIZ);
  ifr.ifr_mtu = mtu;

  if (ioctl(net_fd, SIOCSIFMTU, &ifr) < 0) {

    int err = errno;
    close(net_fd);
    errno = err;
    return -1;
  }

  return close(net_fd);
}

#endif /*setmtu_unix_h_INCLUDE*/
