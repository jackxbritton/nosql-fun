#ifndef SLAVE_H
#define SLAVE_H

typedef struct {
    //Slave socket file descriptor
    int slave_sockfd;
} Slave;

bool slave_init(Slave *s, const char *host, int port);

#endif
