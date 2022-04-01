package main

import (
	"ethos/altEthos"
	"ethos/syscall"
	"log"
)

/**
 * RPC Client for Users.
 *
 * @author Sushen Kumar <skumar88@uic.edu>
 */

var usersList []User

func init() {
	SetupMyRpcmakeAccountsReply(makeAccountsReply)
	SetupMyRpcgetBalanceReply(getBalanceReply)
	SetupMyRpctransferReply(transferReply)
	SetupMyRpcwriteUserReply(writeUserReply)
}

// Reply for making accounts
func makeAccountsReply(users []User) MyRpcProcedure {
	log.Printf("Created %v accounts!\n ", len(users))
	usersList = users
	return nil
}

// Reply for getting balance
func getBalanceReply(userID int64, userBalance int64) MyRpcProcedure {
	log.Printf("Client_Received_Balance_Reply --- ID:_%v---Balance:_%v\n", userID, userBalance)
	return nil
}

// Reply for transferring balance
func transferReply(amt1 int64, amt2 int64) MyRpcProcedure {
	log.Printf("First user has balance : %v and second user has balance: %v\n", amt1, amt2)
	return nil
}

// Reply for writing users to FS
func writeUserReply(userID int64) MyRpcProcedure {
	log.Printf("User written to fileSystem: %v\n", userID)
	return nil
}

// Main function
func main() {
	altEthos.LogToDirectory("test/rpcClient")
	log.Printf("rpcClient:_before_call\n")

	fd, status := altEthos.IpcRepeat("MyRpc", "", nil)
	if status != syscall.StatusOk {
		log.Printf("Ipc_failed:_%v\n", status)
		altEthos.Exit(status)
	}

	// An initial call to make 3 users
	initialCall := MyRpcmakeAccounts{int64(3)}
	status = altEthos.ClientCall(fd, &initialCall)

	// Calls
	if len(usersList) > 0 {

		// Get balance for all users
		for _, user := range usersList {
			fd, status = altEthos.IpcRepeat("MyRpc", "", nil)
			if status != syscall.StatusOk {
				log.Printf("Ipc_failed:_%v\n", status)
				altEthos.Exit(status)
			}

			call := MyRpcgetBalance{user}
			status = altEthos.ClientCall(fd, &call)
		}

		// Transferring balance of one user to another
		for i := 0; i < len(usersList)-1; i++ {
			fd, status = altEthos.IpcRepeat("MyRpc", "", nil)
			if status != syscall.StatusOk {
				log.Printf("Ipc_failed:_%v\n", status)
				altEthos.Exit(status)
			}

			call := MyRpctransfer{usersList[i], usersList[i+1], int64(100)}
			status = altEthos.ClientCall(fd, &call)
		}
	}

	// end
	log.Printf("RpcClient:done\n")
}
