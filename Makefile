TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=serko.com
NAMESPACE=serko
NAME=registry
BINARY=terraform-provider-${NAME}
OS_ARCH=windows_amd64
VERSION=0.1.0

default: install


build:
	go build -o ${BINARY}.exe


release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p C:\Users\MasonChen\AppData\Roaming\terraform.d\plugins\${HOSTNAME}\${NAMESPACE}\${NAME}\${VERSION}\${OS_ARCH}
	mv ${BINARY} C:\Users\MasonChen\AppData\Roaming\terraform.d\plugins\${HOSTNAME}\${NAMESPACE}\${NAME}\${VERSION}\${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

# go build -gcflags="all=-N -l"
# go build -o terraform-provider-registry.exe
# mkdir -p C:\Users\MasonChen\AppData\Roaming\terraform.d\plugins\serko.com\serko\registry\0.1.0\windows_amd64
# rm C:\Users\MasonChen\AppData\Roaming\terraform.d\plugins\serko.com\serko\registry\0.1.0\windows_amd64\terraform-provider-registry.exe
# move terraform-provider-registry.exe C:\Users\MasonChen\AppData\Roaming\terraform.d\plugins\serko.com\serko\registry\0.1.0\windows_amd64
