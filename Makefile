TARGETS = reloader
SOCKET = /tmp/reloader
ALIVE = $(shell ps aux | grep -v grep | grep "reloader -c" -c)
ifeq (0, $(ALIVE))
	live = ./reloader -c server &> reloader.log &
else
	live = ./reloader -c reload
endif

alive: reloader
	$(live)

$(TARGETS): reloader.go
	go build

update:
	./reloader -c rebuild

dead:
	@if [ $(ALIVE) -gt 0 ] ; then \
		make reloader ; \
		./reloader -c quit ; \
	fi
	rm -f $(TARGETS) $(SOCKET) reloader.log

.PHONY: alive update dead
