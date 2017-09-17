package main

import "fmt"
import "strconv"

// Account represents a single ledger account
type Account struct {
	balance      float64
	startBalance float64
}

func (account *Account) makeTransaction(income float64, expense float64) float64 {
	newBalance := income + account.balance - expense
	account.balance = newBalance
	return newBalance
}

func (account *Account) deposit(amount float64) float64 {
	return account.makeTransaction(amount, 0.0)
}

func (account *Account) withdraw(amount float64) float64 {
	return account.makeTransaction(0.0, amount)
}

func main() {
	fmt.Println("My Little Ledger")

	account := Account{balance: 0.0, startBalance: 0.0}

	account.deposit(100.0)
	account.withdraw(30.0)

	fmt.Println("Balance: " + strconv.FormatFloat(account.balance, 'f', 2, 64))
	fmt.Println("Starting Balance: " + strconv.FormatFloat(account.startBalance, 'f', 2, 64))
}
