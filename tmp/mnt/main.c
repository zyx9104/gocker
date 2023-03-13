#define _GNU_SOURCE

#include <errno.h>
#include <fcntl.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

__attribute__((constructor)) void set_namespace(void) {
  char *container_pid = NULL;
  container_pid = getenv("CONTAINER_PID");
  if (!container_pid) {
    fprintf(stderr, "CONTAINER_PID not set");
    return;
  }
  fprintf(stdout, "container_pid: %s", container_pid);
  char *container_cmd = NULL;

  container_cmd = getenv("CONTAINER_CMD");

  if (!container_cmd) {
    fprintf(stderr, "CONTAINER_CMD not set");
    return;
  }
  fprintf(stdout, "container_cmd: %s", container_cmd);

  int i;
  char nspath[1024];
  char *namespace[] = {"ipc", "uts", "pid", "user", "net", "mnt"};
  for (int i = 0; i < 6; i++) {
    sprintf(nspath, "/proc/%s/ns/%s", container_pid, namespace[i]);
    int fd = open(nspath, O_RDONLY);
    if (setns(fd, 0)) {
      fprintf(stderr, "set namespace: %s fail\n", namespace[i]);
    }
    fprintf(stdout, "set namespace: %s success", namespace[i]);
    close(fd);
  }
  exit(system(container_cmd));
}