# binary path
EXE = $(shell which go)

# binary file name
BIN = https2

# sources
SRC = https2.go

# install destination
DST = /usr/local/bin/go/$(BIN)

# task
$(BIN):
	@$(EXE) build

fmt:
	@$(BIN) fmt -w $(SRC)

run:
	@$(BIN) run .

install:
	@cp $(BIN) $(DST)

clean:
	@rm $(BIN)
