BIN_DIR = ./bin/
BINARY_NAME=seacrane
ARMBIN = ${BINARY_NAME}-arm
MIPSBIN = ${BINARY_NAME}-mips
MIPSBINSF = ${BINARY_NAME}-mips-sf
MIPSBINLE = ${BINARY_NAME}-mips-le
WINDOWS = ${BINARY_NAME}.exe
ANDROID = ${BINARY_NAME}-android
OSX = ${BINARY_NAME}-osx
OSXARM = ${BINARY_NAME}-osxarm
WASM = ${BINARY_NAME}-wasm

x86:
	go build -o ${BIN_DIR}${BINARY_NAME} ${BINARY_NAME}
 
arm:
	env GOARCH=arm GOOS=linux go build -o ${BIN_DIR}${ARMBIN} main.go

mips:
	env GOARCH=mips GOOS=linux go build -o ${BIN_DIR}${MIPSBIN} main.go

mipssf:
	env GOARCH=mips GOMIPS=softfloat GOOS=linux go build -o ${BIN_DIR}${MIPSBINSF} main.go

mipsle:
	env GOARCH=mipsle GOOS=linux go build -o ${BIN_DIR}${MIPSBINLE} main.go

windows:
	env GOARCH=amd64 GOOS=windows go build -o ${BIN_DIR}${WINDOWS} main.go

osx:
	env GOHOSTOS=linux GOOS=darwin go build -o ${BIN_DIR}${OSX} main.go

wasm:
	env GOHOSTOS=linux GOOS=js GOARCH=wasm go build -o ${BIN_DIR}${WASM} main.go

arm-osx:
	GOHOSTOS=linux GOARCH=arm64 GOOS=darwin go build -o ${BIN_DIR}${OSXARM} main.go
# golang 1.18 will atleast be needed for this. darwin arm support came in this release

android:
	env GOHOSTOS=linux GOARCH=arm64 GOOS=android go build -o ${BIN_DIR}${ANDROID} main.go
# requires gomobile

# all: x86 arm windows osx arm-osx mips mipssf mipsle android
all: all-32bit all-non32bit

all-32bit: arm mips mipssf mipsle

all-non32bit: x86 windows osx arm-osx android

run:
	go run ${BINARY_NAME}

clean:
	echo "all clean lol"

clean-all:
	rm ${BIN_DIR}${BINARY_NAME} ${BIN_DIR}${ARMBIN} ${BIN_DIR}${WINDOWS} ${BIN_DIR}${OSX} ${BIN_DIR}${OSX-ARM} ${BIN_DIR}${MIPSBIN} ${BIN_DIR}${MIPSBINLE} ${BIN_DIR}${MIPSBINSF} ${BIN_DIR}${ANDROID}
	file ${BIN_DIR}seacrane*
