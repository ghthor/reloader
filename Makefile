TARGET		= $(shell pwd | xargs basename)
TARGET_APP	= $(TARGET)
TARGET_SRC	= $(shell find *.go)
SOCKET 		= /tmp/$(TARGET_APP)
ALIVE 		= $(shell ps aux | grep -v grep | grep "$(TARGET_APP) -c" -c)

lifeBegin 	= ./$(TARGET_APP) -c server &>> $(TARGET).log &
lifeRenew 	= ./$(TARGET_APP) -c rebuild
lifeRefresh = ./$(TARGET_APP) -c reload
lifeEnd 	= ./$(TARGET_APP) -c quit

ifeq (0, $(ALIVE))
	alive_cmd 	= $(lifeBegin)

	refresh_reqs 	= alive
	refresh_cmd		=

	kill_reqs	=
	kill_cmd 	=
else
	alive_cmd 	= $(lifeRefresh)

	refresh_reqs 	= $(TARGET_APP)
	refresh_cmd		= $(lifeRefresh)

	kill_reqs	= $(TARGET_APP)
	kill_cmd 	= $(lifeEnd)
endif

alive: $(TARGET_APP)
	$(alive_cmd)

refresh: $(refresh_reqs) ; $(refresh_cmd)
kill: $(kill_reqs) ; $(kill_cmd)

dead: kill
	go clean
	rm -f $(SOCKET)

$(TARGET_APP): $(TARGET_SRC)
	go build


.PHONY: alive refresh kill dead
