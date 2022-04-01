MyRpc interface { 
	makeAccounts(count int64) ([]User)
	getBalance(user User) (int64, int64)
	transfer(user1 User, user2 User, amt int64) (User, User)
	writeUser(user User) (int64)
}

User struct {
	UserID int64
	UserBalance int64
}
