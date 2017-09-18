package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Money represents dollars and cents as an integer value
type Money int64

// Transaction represents a single transaction in an Account
type Transaction struct {
	Timestamp int64 `json:"timestamp"`
	Income    Money `json:"income"`
	Expense   Money `json:"expense"`
	Balance   Money `json:"balance"`
}

// Account represents a single ledger account
type Account struct {
	Balance      Money         `json:"balance"`
	StartBalance Money         `json:"startBalance"`
	Transactions []Transaction `json:"transactions"`
}

func fToMoney(f float64) Money {
	return Money(f * 100)
}

func moneyToF(money Money) float64 {
	return float64(money) * 0.01
}

func printMoney(money Money) {
	fmt.Printf("$%10.2f", moneyToF(money))
}

func formatMoney(money Money) string {
	return fmt.Sprintf("$%10.2f", moneyToF(money))
}

func createAccount(startBalance float64) Account {
	balance := fToMoney(startBalance)
	account := Account{Balance: balance, StartBalance: balance}
	account.Transactions = make([]Transaction, 0)
	return account
}

func (account *Account) makeTransaction(income Money, expense Money) Money {
	newBalance := income + account.Balance - expense
	account.Balance = newBalance
	timestamp := time.Now().Unix()
	transaction := Transaction{Balance: newBalance, Income: income, Expense: expense, Timestamp: timestamp}
	account.Transactions = append(account.Transactions, transaction)
	return newBalance
}

func (account *Account) deposit(amount Money) Money {
	fmt.Printf("Depositing\t\t%s\n", formatMoney(amount))
	return account.makeTransaction(amount, 0.0)
}

func (account *Account) withdraw(amount Money) Money {
	fmt.Printf("Withdrawing\t\t%s\n", formatMoney(amount))
	return account.makeTransaction(0.0, amount)
}

func (account *Account) saveToFile(name string) error {
	accountFileBytes, err := json.Marshal(account)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = ioutil.WriteFile(name, accountFileBytes, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (account *Account) readFromFile(name string) error {
	accountFileBytes, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = json.Unmarshal(accountFileBytes, &account)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func formatTimestamp(timestamp int64) string {
	// MM/DD/YYYY HH:MM:SS PM TZ
	return time.Unix(timestamp, 0).Format("01/02/2006 03:04:05 PM MST")
}

func getCommandLine() (int, []string) {
	argc := len(os.Args) - 1
	argv := os.Args[1:]
	return argc, argv
}

func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	os.Exit(0)
}

// WorkspaceConfig represents the configuration for the ledger workspace
type WorkspaceConfig struct {
	ActiveAccount string `json:"activeAccount"`
}

func createPath(userPath string) error {
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		err = os.MkdirAll(userPath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

func createFile(fileName string, fileBytes []byte) error {
	err := ioutil.WriteFile(fileName, fileBytes, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func getWorkspacePath() string {
	absolutePathToCurrentDirectory, _ := filepath.Abs("./")
	workspacePath := path.Join(absolutePathToCurrentDirectory, ".ledger/")
	return workspacePath
}

func getWorkspaceConfigPath() string {
	workspacePath := getWorkspacePath()
	configFilePath := path.Join(workspacePath, "config.json")
	return configFilePath
}

func createDefaultWorkspaceConfiguration() WorkspaceConfig {
	config := WorkspaceConfig{ActiveAccount: "default"}
	return config
}

func (config *WorkspaceConfig) save(fileName string) error {
	configFileBytes, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = createFile(fileName, configFileBytes)
	return err
}

func (config *WorkspaceConfig) load(fileName string) error {
	configFileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = json.Unmarshal(configFileBytes, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func initializeWorkspace() error {
	workspacePath := getWorkspacePath()
	createPath(workspacePath)

	configFilePath := path.Join(workspacePath, "config.json")
	config := createDefaultWorkspaceConfiguration()
	config.save(configFilePath)

	defaultAccountPath := path.Join(workspacePath, "default")
	createPath(defaultAccountPath)

	defaultAccount := createAccount(0)
	defaultAccountFilePath := path.Join(defaultAccountPath, "data.json")
	defaultAccount.saveToFile(defaultAccountFilePath)

	fmt.Println("\n*** initialized ledger workspace")
	return nil
}

func selectActiveAccount(accountName string) {
	configFilePath := getWorkspaceConfigPath()
	config := WorkspaceConfig{}
	config.load(configFilePath)
	config.ActiveAccount = accountName
	config.save(configFilePath)
}

func createNewAccount(accountName string) {
	workspacePath := getWorkspacePath()
	accountPath := path.Join(workspacePath, accountName)
	createPath(accountPath)
	account := createAccount(0)
	accountFilePath := path.Join(accountPath, "data.json")
	account.saveToFile(accountFilePath)
	selectActiveAccount(accountName)
}

func unknownAction(action string) {
	fmt.Printf("Unknown Action %s\n", action)
	os.Exit(-1)
}

func main() {
	fmt.Println("My Little Ledger v1.0")
	argc, argv := getCommandLine()

	if argc >= 1 {
		action := strings.ToLower(strings.Trim(argv[0], " \n"))
		switch action {
		case "init":
			initializeWorkspace()
		case "new":
			if argc >= 2 {
				accountName := strings.ToLower(strings.Trim(argv[1], " \n"))
				createNewAccount(accountName)
			} else {
				fmt.Println("Unable to create account. Missing Required Account Name")
				os.Exit(-1)
			}
		// case "ls": listTransactions()
		// case "shell": startInteractiveLedger()
		default:
			unknownAction(action)
		}
	} else {
		showHelp()
	}

	// usage
	// $ appname
	//    show help
	// $ appname init
	//    create .ledger/config.json
	// $ appname new accountname
	//    create .ledger/accountname/data.json
	//    write `"activeAccount":accountname`
	// $ appname ls
	//    list transactions in active account
	// $ appname shell
	//    start the interactive shell to create, update, and delete transactions

	// createAccountFlag := flag.Bool("createaccount", false, "specifies to create an account")
	// argv := flag.Args()
	// fmt.Println(*createAccountFlag)

	fmt.Print(argc, argv)

	// account := createAccount(0)

	// err := account.readFromFile("my-account.json")
	// if err != nil {
	// 	fmt.Println("There was a critical error reading the account database!")
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// // accountFileBytes := []byte(`{ "startBalance": 0.00, "balance": 0.00, "transactions": [] }`)

	// fmt.Printf("Starting Balance:\t%s\n", formatMoney(account.StartBalance))

	// // account.deposit(fToMoney(100.0))
	// // account.withdraw(fToMoney(30))

	// fmt.Printf("Balance:\t\t%s\n", formatMoney(account.Balance))
	// hr := strings.Repeat("-", 90)

	// fmt.Println("Transactions: ")
	// for i := 0; i < len(account.Transactions); i++ {
	// 	transaction := account.Transactions[i]
	// 	income := transaction.Income
	// 	expense := transaction.Expense
	// 	balance := transaction.Balance
	// 	ts := transaction.Timestamp

	// 	fmt.Printf("%04d: %s IN %s OUT %s BAL - %s\n%s\n", i, formatMoney(income), formatMoney(expense), formatMoney(balance), formatTimestamp(ts), hr)
	// }

	// // accountFileBytes, err := json.Marshal(account)
	// // if err != nil {
	// // 	log.Print(err)
	// // 	return
	// // }

	// // err = ioutil.WriteFile("account.json", accountFileBytes, 0644)
	// // if err != nil {
	// // 	log.Print(err)
	// // 	return
	// // }

}
