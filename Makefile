# go binary path
EXE = $(shell which go)

# binary name
BIN = https2

# library path
LIB = github.com/lib/pq

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
$(BIN): clean fmt preprocess
	@$(EXE) build

# format
fmt:
	@$(EXE)fmt -w $(SRC)

# preprocess
preprocess:
	@$(EXE) get $(LIB) 

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
ifeq ("$(shell ls | grep -x $(BIN))", "$(BIN)")
	@rm $(BIN)
endif

#.PHONY
#.PHONY:	preprocess
