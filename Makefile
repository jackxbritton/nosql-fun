MASTER_OBJS = master.o
SLAVE_OBJS = slave_main.o slave.o
CFLAGS = -Wall
LIBS =

all: master slave

master: $(MASTER_OBJS)
	$(CC) $(MASTER_OBJS) -o master $(LIBS)

slave: $(SLAVE_OBJS)
	$(CC) $(SLAVE_OBJS) -o slave $(LIBS)

clean:
	rm master $(MASTER_OBJS) slave $(SLAVE_OBJS)

