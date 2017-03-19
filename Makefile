PROGRAMS = dnsts/dnsts timeserializer/timeserializer
BUILD = go build -ldflags="-s -w" -v
.PHONY: all clean

all: $(PROGRAMS)

$(PROGRAMS):%:%.go
	@echo "===> Buiding $@ from $<"
	go fmt $^
	go vet $^
# prod is amd64 & linux. dev might be macos. go makes this easy
	@if [ -n "$$NATIVE" ]; then \
		echo $(BUILD) $^; \
		$(BUILD) $^; \
	else \
		echo GOOS=linux GOARCH=amd64 $(BUILD) $^; \
		GOOS=linux GOARCH=amd64 $(BUILD) $^; \
	fi

clean:
	@echo "===> Cleaning up"
	rm -f $(PROGRAMS)
