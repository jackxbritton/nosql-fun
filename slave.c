#include "slave.h"
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <errno.h>
#include <arpa/inet.h>


#define PORT 8080
#define TCP 0
#define MAXBUF 1024 //Max amount of data we can send

bool slave_init(Slave *s, const char *host, int port) {
	printf("host: %s\n", host);
	printf("port: %d\n", port);

	struct sockaddr_in server_addr;
	socklen_t addr_size;
	char buffer[MAXBUF];
	char *message = "#SlaveLyfe";


	//Open streaming socket
	if((s->slave_sockfd = socket(PF_INET, SOCK_STREAM, TCP)) < 0){
		perror("Socket Initialization");
		exit(errno);
	}

	//Configure server address
	server_addr.sin_family = AF_INET;
	server_addr.sin_port = htons(port);
	server_addr.sin_addr.s_addr = inet_addr(host);
	memset(server_addr.sin_zero, '\0', sizeof(server_addr.sin_zero));

	//Attempt to connect to the server
	addr_size = sizeof(server_addr);
	if(connect(s->slave_sockfd, (struct sockaddr *)&server_addr, addr_size) < 0){
		printf("\nClient unable to connect!\n");
		exit(errno);
	}

	send(s->slave_sockfd, message, strlen(message), 0);
	printf("Message sent!");
	return true;
}
