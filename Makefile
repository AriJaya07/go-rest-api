.PHONY: all build run clean

#  Go parameters
GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=${GOCMD} clean
BINARY_NAME=main

all: build

build:
	${GOBUILD} -o ${BINARY_NAME} -v

run: build 
	./${BINARY_NAME}

clean: 
	${GOCLEAN}
	rm -f ${BINARY_NAME}