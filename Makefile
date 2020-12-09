.PHONY: all clean

ifeq ($(OS), Windows_NT)
    OUTPUT=seata-benchmark.exe
else
    OUTPUT=seata-benchmark
endif

all: clean
	go build -o ./bin/${OUTPUT} main.go

clean:
	rm -f bin/${OUTPUT}
