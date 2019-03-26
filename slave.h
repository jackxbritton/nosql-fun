#ifndef SLAVE_H
#define SLAVE_H

#include <stdbool.h>

typedef struct {

} Slave;

bool slave_init(Slave *s, const char *host, int port);

#endif

