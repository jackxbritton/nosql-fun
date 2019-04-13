#ifndef SLAVE_H
#define SLAVE_H

#include <stdbool.h> //Getting compile error w/out this for some reason

typedef struct {
    //Slave socket file descriptor
    int slave_sockfd;
} Slave;

bool slave_init(Slave *s, const char *host, int port);

#endif
