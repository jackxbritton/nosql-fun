#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <arpa/inet.h>

typedef struct {
	bool occupied;
} SlaveNode;

int main(int argc, char *argv[]) {

	// Parse port from CLI args.
	int port;
	if (argc < 2 || sscanf(argv[1], "%d", &port) != 1) {
		fprintf(stderr, "Usage: %s PORT\n", argv[0]);
		return 1;
	}

	int num_slaves = 64;
	SlaveNode *slaves = calloc(num_slaves, sizeof(SlaveNode));
	if (slaves == NULL) {
		perror("calloc");
		return 1;
	}

	// Start two TCP servers:
	// One to receive queries from outside,
	// and one to accept slave connections.

	int sock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
	if (sock == -1) {
		perror("socket");
		return 1;
	}

	struct sockaddr_in addr;
	memset(&addr, 0, sizeof(addr));
	if (bind(sock, (struct sockaddr *) &addr, sizeof(addr)) == -1) {
		perror("bind");
		return 1;
	}

	while (true) {


		// Wait for a connection.
		struct sockaddr_in client_addr;
		memset(&client_addr, 0, sizeof(client_addr));
		unsigned int len = sizeof(client_addr);
		int client_sock = accept(sock, (struct sockaddr *) &client_addr, &len);
		if (client_sock == -1) {
			perror("accept");
			return 1;
		}

		// TODO Do stuff.

	}

	return 0;

}

