export UDIR= .
export GOC = x86_64-xen-ethos-6g
export GOL = x86_64-xen-ethos-6l
export ETN2GO = etn2go
export ET2G   = et2g
export EG2GO  = eg2go

export GOARCH = amd64
export TARGET_ARCH = x86_64
export GOETHOSINCLUDE=/usr/lib64/go/pkg/ethos_$(GOARCH)
export GOLINUXINCLUDE=/usr/lib64/go/pkg/linux_$(GOARCH)

export ETHOSROOT=server/rootfs
export MINIMALTDROOT=server/minimaltdfs

.PHONY: all install
all: RpcServer

install: RpcServer RpcClient
	sudo rm -rf server/
	(ethosParams server && cd server/ && ethosMinimaltdBuilder)
	echo 7 > server/param/sleepTime
	ethosTypeInstall MyRpc
	ethosDirCreate $(ETHOSROOT)/services/MyRpc $(ETHOSROOT)/types/spec/MyRpc/MyRpc all
	ethosDirCreate $(ETHOSROOT)/services/MyRpc $(ETHOSROOT)/types/spec/User/User all
	cp RpcClient $(ETHOSROOT)/programs
	cp RpcServer $(ETHOSROOT)/programs
	ethosStringEncode /programs/RpcServer > $(ETHOSROOT)/etc/init/services/RpcServer
	ethosStringEncode /programs/RpcClient > $(ETHOSROOT)/etc/init/services/RpcClient

MyRpc.go: MyRpc.t
	$(ETN2GO) . MyRpc main $^

RpcServer: RpcServer.go MyRpc.go
	ethosGo $^

RpcClient: RpcClient.go MyRpc.go
	ethosGo $^

clean:
	sudo rm -rf server/
	rm -rf MyRpc/ MyRpcIndex/
	rm -f MyRpc.go
	rm -f RpcServer
	rm -f RpcClient
	rm -f RpcServer.goo.ethos
	rm -f RpcClient.goo.ethos
