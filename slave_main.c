#include <stdio.h>
#include "slave.h"

int main(int argc, char *argv[]) {

	// Parse CLI args for hostname and port.
	static char host[128];
	int port;
	if (argc < 2 || sscanf(argv[1], "%[^:]:%d", host, &port) != 2) {
		fprintf(stderr, "Usage: %s HOST:PORT\n", argv[0]);
		return 1;
	}

	Slave slave;
	slave_init(&slave, host, port);

	return 0;

}
