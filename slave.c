#include "slave.h"
#include <stdio.h>
#include <stdbool.h>
#include <sys/socket.h>
#include <stdlib.h>
#include <netinet/in.h>
#include <string.h>
#include <errno.h>


#define PORT 8080
#define MAXBUF 1024 //Max amount of data we can send

bool slave_init(Slave *s, const char *host, int port) {
	printf("host: %s\n", host);
	printf("port: %d\n", port);

	struct sockaddr_in
	char buffer[MAXBUF];
	char *message = "#SlaveLyfe";
	int sock = 0, server_return;

	//Open streaming socket
	if((Slave->slave_sockfd = socket(AF_NET, SOCK_STREAM, 0)) < 0){
		perror("Socket Initialization");
		exit(errno);
	}

	memset(&host, '0', sizeof(host));
	host.sin_family = AF_INET;
	host.sin_port = htons(port);

	//Normally would check to make sure the host:port address is valid here, but
	//this is currently handled by slave_main.c

	//Attempt to connect to the server
	if(connect(sock, (struct sockaddr *)&host, sizeof(host)) < 0){
		printf("\nClient unable to connect!\n");
		exit(errno);
	}

	send(sock, message, strlen(message), 0);
	printf("Message sent!");
	server_return = read(sock, buffer, MAXBUF);
	printf("%s\n", buffer);
	return true;
}
