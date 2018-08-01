.PHONY: all build sfs-provisioner clean

all:build

build:sfs-provisioner

package:
	mkdir -p  ./bin/

sfs-provisioner:package
	go build -o ./bin/sfs-provisioner ./cmd/sfs-provisioner

clean:
	rm -rf ./bin/ 
