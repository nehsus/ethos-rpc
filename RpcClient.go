package main

import (
	"ethos/altEthos"
	"ethos/syscall"
	"log"
)

var usersList []User

func init() {
	SetupMyRpcmakeAccountsReply(makeAccountsReply)
	SetupMyRpcgetBalanceReply(getBalanceReply)
	SetupMyRpctransferReply(transferReply)
	SetupMyRpcwriteUserReply(writeUserReply)
}

func makeAccountsReply(users []User) MyRpcProcedure {
	log.Printf("Created %v accounts!\n ", len(users))
	usersList = users
	return nil
}

func getBalanceReply(userID int64, userBalance int64) MyRpcProcedure {
	log.Printf("Client_Received_Balance_Reply --- ID:_%v---Balance:_%v\n", userID, userBalance)
	return nil
}

func transferReply(amt1 int64, amt2 int64) MyRpcProcedure {
	log.Printf("First user has balance : %v and second user has balance: %v\n", amt1, amt2)
	return nil
}

func writeUserReply(userID int64) MyRpcProcedure {
	log.Printf("User written to fileSystem: %v\n", userID)
	return nil
}

func main() {
	altEthos.LogToDirectory("test/rpcClient")
	log.Printf("rpcClient:_before_call\n")

	fd, status := altEthos.IpcRepeat("MyRpc", "", nil)
	if status != syscall.StatusOk {
		log.Printf("Ipc_failed:_%v\n", status)
		altEthos.Exit(status)
	}

	initialCall := MyRpcmakeAccounts{int64(3)}
	status = altEthos.ClientCall(fd, &initialCall)

	if len(usersList) > 0 {
		for _, user := range usersList {
			fd, status = altEthos.IpcRepeat("MyRpc", "", nil)
			if status != syscall.StatusOk {
				log.Printf("Ipc_failed:_%v\n", status)
				altEthos.Exit(status)
			}

			call := MyRpcgetBalance{user}
			status = altEthos.ClientCall(fd, &call)
		}
	}
	log.Printf("RpcClient:done\n")
}
