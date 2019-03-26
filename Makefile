OBJS = main.o
CFLAGS = -Wall
LIBS =
TARGET = nosql-fun

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(OBJS) -o $(TARGET) $(LIBS)

clean:
	rm $(TARGET) $(OBJS)

