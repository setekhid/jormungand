/* Copyright 2016 Huitse Tai. All rights reserved.
 * Use of this source code is governed by BSD 3-clause
 * license that can be found in the LICENSE file.
 */

#if defined(macintosh) || defined(Macintosh) ||                                \
    (defined(__APPLE__) && defined(__MACH__))

#include "tun_fd.h"
#include "setmtu_unix.h"

#include <sys/sys_domain.h>
#include <sys/kern_control.h>
#include <net/if_utun.h>
#include <sys/socket.h>
#include <sys/ioctl.h>
#include <unistd.h>
#include <errno.h>
#include <string.h>
#include <stdio.h>

/* from openvpn tun.c */
/* Helper functions that tries to open utun device
   return -2 on early initialization failures (utun not supported
   at all (old OS X) and -1 on initlization failure of utun
   device (utun works but utunX is already used */
static int utun_open_helper(struct ctl_info ctlInfo, int utunnum) {

  struct sockaddr_ctl sc;
  int fd;

  fd = socket(PF_SYSTEM, SOCK_DGRAM, SYSPROTO_CONTROL);
  if (fd < 0) {
    errno = errno;
    return -2;
  }

  if (ioctl(fd, CTLIOCGINFO, &ctlInfo) < 0) {

    int err = errno;
    close(fd);
    errno = err;
    return -2;
  }

  sc.sc_id = ctlInfo.ctl_id;
  sc.sc_len = sizeof(sc);
  sc.sc_family = AF_SYSTEM;
  sc.ss_sysaddr = AF_SYS_CONTROL;
  sc.sc_unit = utunnum + 1;

  /* If the connect is successful, a utun%d device will be created, where "%d"
   * is (sc.sc_unit - 1) */

  if (connect(fd, (struct sockaddr *)&sc, sizeof(sc)) < 0) {

    int err = errno;
    close(fd);
    errno = err;
    return -1;
  }

  errno = 0;
  return fd;
}

/* from openvpn tun.c */
static int open_darwin_utun(tun_fd_info *info) {

  struct ctl_info ctlInfo;
  memset(&ctlInfo, 0, sizeof(ctlInfo));
  strlcpy(ctlInfo.ctl_name, UTUN_CONTROL_NAME, sizeof(ctlInfo.ctl_name));

  int fd;

  /* try to open first available utun device if no specific utun is requested */
  int utunnum;
  for (utunnum = 0; utunnum < 255; utunnum++) {

    fd = utun_open_helper(ctlInfo, utunnum);
    /* Break if the fd is valid,
     * or if early initalization failed (-2) */
    if (fd != -1)
      break;
  }

  if (fd < 0) {
    errno = errno;
    return -1;
  }

  info->fd = fd;

  /* Retrieve the assigned interface name. */
  socklen_t utunname_len = sizeof(info->nam);
  if (getsockopt(fd, SYSPROTO_CONTROL, UTUN_OPT_IFNAME, info->nam,
                 &utunname_len) < 0) {
    errno = errno;
    return -1;
  }

  errno = 0;
  return 0;
}

int const IFPKT_OFFSET = 4;

tun_fd_info create_tun_fd(int mtu) {

  tun_fd_info info;
  memset(&info, 0, sizeof(info));

  if (open_darwin_utun(&info) < 0) {
    errno = errno;
    return info;
  }

  if (setmtu_unix(info.nam, mtu) < 0) {

    int err = errno;
    close(info.fd);
    errno = err;
    return info;
  }

  errno = 0;
  return info;
}

int const built_macos = 1;
#else

int const built_macos = 0;
#endif /*mac os guard*/
