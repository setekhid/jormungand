/* Copyright 2016 Huitse Tai. All rights reserved.
 * Use of this source code is governed by BSD 3-clause
 * license that can be found in the LICENSE file.
 */

#ifndef tun_fd_h_INCLUDE
#define tun_fd_h_INCLUDE

static int const IFNAME_LEN = 32;

typedef struct tun_fd_info {
  int fd;
  char nam[IFNAME_LEN];
} tun_fd_info;

extern tun_fd_info create_tun_fd(int mtu);
extern int const IFPKT_OFFSET;

#endif /* tun_fd_h_INCLUDE */
