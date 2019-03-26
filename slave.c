#include "slave.h"
#include <stdio.h>

bool slave_init(Slave *s, const char *host, int port) {
	printf("host: %s\n", host);
	printf("port: %d\n", port);
	return true;
}

