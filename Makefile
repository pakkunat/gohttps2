# go binary path
EXE = $(shell which go)

# binary name
BIN = https2

# certification files
CRT = certificate.crt 
KEY = private.key

# sources
SRC = https2.go

# template files
TPL = *.html

# install
# base directory
BDR = /usr/local/bin

# base go directory
GDR = $(BDR)/go

# install destination
DST = $(GDR)/$(BIN)

# task
# compile
$(BIN):
	@$(EXE) build

# format
fmt:
	@$(EXE)fmt -w $(SRC)

# install
install:
ifeq ("$(shell ls $(BDR) | grep go)", "")
	@mkdir $(GDR)
endif
ifeq ("$(shell ls $(GDR) | grep $(BIN))", "")
	@mkdir $(DST)
endif
ifeq ("$(shell ls $(DST) | grep $(CRT))", "")
	@cp $(CRT) $(DST)
endif
ifeq ("$(shell ls $(DST) | grep $(KEY))", "")
	@cp $(KEY) $(DST)
endif
	@cp $(BIN) $(DST)
	@cp $(TPL) $(DST)

# uninstall
uninstall:
	@rm -r $(DST)

# run
run:
	@cd $(DST) && ./$(BIN)

# clean
clean:
	@rm $(BIN)
