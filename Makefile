TARGETS = reloader

$(TARGETS): reloader.go
	go build

clean:
	rm -f $(TARGETS)

.PHONY: clean
