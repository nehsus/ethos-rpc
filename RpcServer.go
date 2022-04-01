package main

import (
	"ethos/altEthos"
	"ethos/syscall"
	"log"
	"strconv"
)

var userListServer []User
var PATH = "/home/me/" + syscall.GetUser() + "/"

func init() {
	SetupMyRpcmakeAccounts(makeAccounts)
	SetupMyRpcgetBalance(getBalance)
	SetupMyRpctransfer(transfer)
	SetupMyRpcwriteUser(writeUser)
}

func generateID() int64 {
	r := syscall.GetTime()
	return int64(((r * 7621) + 1) % 32768)
}

func generateBalance() int64 {
	r := syscall.GetTime()
	return int64(((r * 7621) + 1) % 32768)
}

func makeAccounts(count int64) MyRpcProcedure {

	for i := 0; i < int(count); i++ {

		user := User{
			UserID:      generateID(),
			UserBalance: generateBalance(),
		}

		userListServer = append(userListServer, user)
	}

	return &MyRpcmakeAccountsReply{userListServer}
}

func getBalance(user User) MyRpcProcedure {
	log.Printf("MyRpcService : getBalance called\n")

	return &MyRpcgetBalanceReply{user.UserID, user.UserBalance}
}

func transfer(user1 User, user2 User, amt int64) MyRpcProcedure {
	if user1.UserBalance >= amt {
		user1.UserBalance -= amt
		user2.UserBalance += amt
		log.Printf("Amount transferred between %v and %v\n", user1.UserID, user2.UserID)
	} else {
		log.Printf("Could not transfer, negative\n")
	}

	return &MyRpctransferReply{user1.UserBalance, user2.UserBalance}
}

func writeUser(user User) MyRpcProcedure {
	status := altEthos.Write(PATH+"user-"+strconv.Itoa(int(user.UserID)), &user)
	if status != syscall.StatusOk {
		log.Fatalf("Failed to write for ID: %v/%v\n", user.UserID, status)
		return nil
	}

	return &MyRpcwriteUserReply{user.UserID}
}

func main() {

	altEthos.LogToDirectory("test/rpcServer")

	listeningFd, status := altEthos.Advertise("MyRpc")
	if status != syscall.StatusOk {
		log.Printf("Advertising_service_failed:_%s\n", status)
		altEthos.Exit(status)
	}

	for {
		_, fd, status := altEthos.Import(listeningFd)
		if status != syscall.StatusOk {
			log.Printf("Error_calling_Import:_%v\n", status)
			altEthos.Exit(status)
		}
		log.Printf("MyRpcService:_new_connection_accepted\n")

		t := MyRpc{}
		altEthos.Handle(fd, &t)
	}

}
