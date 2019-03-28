# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: bic android ios bic-cross swarm evm all test clean
.PHONY: bic-linux bic-linux-386 bic-linux-amd64 bic-linux-mips64 bic-linux-mips64le
.PHONY: bic-linux-arm bic-linux-arm-5 bic-linux-arm-6 bic-linux-arm-7 bic-linux-arm64
.PHONY: bic-darwin bic-darwin-386 bic-darwin-amd64
.PHONY: bic-windows bic-windows-386 bic-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

bic:
	build/env.sh go run build/ci.go install ./cmd/bic
	@echo "Done building."
	@echo "Run \"$(GOBIN)/bic\" to launch bic."

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

bic-cross: bic-linux bic-darwin bic-windows bic-android bic-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/bic-*

bic-linux: bic-linux-386 bic-linux-amd64 bic-linux-arm bic-linux-mips64 bic-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-*

bic-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/bic
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep 386

bic-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/bic
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep amd64

bic-linux-arm: bic-linux-arm-5 bic-linux-arm-6 bic-linux-arm-7 bic-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep arm

bic-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/bic
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep arm-5

bic-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/bic
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep arm-6

bic-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/bic
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep arm-7

bic-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/bic
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep arm64

bic-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/bic
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep mips

bic-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/bic
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep mipsle

bic-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/bic
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep mips64

bic-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/bic
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/bic-linux-* | grep mips64le

bic-darwin: bic-darwin-386 bic-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/bic-darwin-*

bic-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/bic
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/bic-darwin-* | grep 386

bic-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/bic
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/bic-darwin-* | grep amd64

bic-windows: bic-windows-386 bic-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/bic-windows-*

bic-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/bic
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/bic-windows-* | grep 386

bic-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/bic
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/bic-windows-* | grep amd64
