.PHONY: all build sfs-provisioner docker clean

all:build

build:sfs-provisioner

package:
	mkdir -p  ./bin/

sfs-provisioner:package
	go build -o ./bin/sfs-provisioner ./cmd/sfs-provisioner

docker:sfs-provisioner
	cp ./bin/sfs-provisioner ./cmd/sfs-provisioner
	docker build cmd/sfs-provisioner -t swr.ap-southeast-1.myhuaweicloud.com/k8s-csi/sfs-provisioner:latest

clean:
	rm -rf ./bin/
	rm -rf ./cmd/sfs-provisioner/sfs-provisioner
