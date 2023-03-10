package main

import (
	"encoding/json"
	"bytes"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type EWallet struct {
}

//type account struct {
//	User    string `json:"user"`
//	Balance int `json:"balance"`
//}

const NO_CHANGE = "NO_CHANGE"

func (e *EWallet) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (e *EWallet) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "createAccount" { //create a new account
		return e.createAccount(stub, args)
	} else if function == "transferAccount" { //transfer one to another
		return e.transferAccount(stub, args)
	} else if function == "transferAccount_multiParty" { //transfer one to another
		return e.transferAccount_multiParty(stub, args)
	}else if function == "transferAccount_new" { //channge username to transfer one to another
		return e.transferAccount_new(stub, args)
	} else if function == "queryAccount" { //query an account
		return e.queryAccount(stub, args)
	} else if function == "deleteAccount" { //delete an account
		return e.deleteAccount(stub, args)
	} else if function == "saveMoney" { //save money
		return e.saveMoney(stub, args)
	} else if function == "drawMoney" { //draw money
		return e.drawMoney(stub, args)
	} else if function == "queryAllAccounts" { //query all accounts
		return e.queryAllAccounts(stub, args)
	} else if function == "getHistoryForAccount" { //query history for account
		return e.getHistoryForAccount(stub, args)
	} else if function == "beGoingToSave" { // change username to save money
		return e.beGoingToSave(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (e *EWallet) createAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user := args[0]
	balance := strconv.Itoa(0)

	// ==== Check if user already exists ====
	userAsBytes, err := stub.GetState(user)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if userAsBytes != nil {
		fmt.Println("This user already exists: " + user)
		return shim.Error("This user already exists: " + user)
	}

	err = stub.PutState(user, []byte(balance))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (e *EWallet) transferAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Auser, Buser string
	var Abalance, Bbalance int
	var X int
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	Auser = args[0]
	Buser = args[1]

	Abalancebytes, err := stub.GetState(Auser)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Abalancebytes == nil {
		return shim.Error("Account not found")
	}
	Abalance, _ = strconv.Atoi(string(Abalancebytes))

	Bbalancebytes, err := stub.GetState(Buser)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bbalancebytes == nil {
		return shim.Error("Account not found")
	}
	Bbalance, _ = strconv.Atoi(string(Bbalancebytes))

	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	if X > Abalance {
		return shim.Error("Not sufficient funds")
	}
	Abalance = Abalance - X
	Bbalance = Bbalance + X
	fmt.Printf("Abalance = %d, Bbalance = %d\n", Abalance, Bbalance)

	// Write the state back to the ledger
	err = stub.PutState(Auser, []byte(strconv.Itoa(Abalance)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(Buser, []byte(strconv.Itoa(Bbalance)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (e *EWallet) transferAccount_multiParty(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	strarg := args[0]
	//note = args[1]

	arg_map := make(map[string]string)
	_ = json.Unmarshal([]byte(strarg), &arg_map)

	var err error
	//暂时不考虑排序，如果排的话，就是把扣钱放在前面，加钱放在后面就行。但这个未必会导致错误哦。所以顺序未必重要，反而浪费性能。

	var oldValue int
	var changeValue int
	var newValue int


	for key, value := range arg_map {
		Bvalbytes, err := stub.GetState(key)
		if err != nil {
			return shim.Error("GetState error, key: " + key )
		}
		if Bvalbytes == nil {
			return shim.Error("no such key: " + key )
		}
		oldValue, _ = strconv.Atoi(string(Bvalbytes))
		changeValue, _ = strconv.Atoi(value)
		newValue = oldValue + changeValue

		// simu_state[key] = strconv.Itoa(newValue)
		if newValue < 0 {
			return shim.Error("not sufficient funds")
		}
		
		if newValue == oldValue {
			arg_map[key] = NO_CHANGE
		} else {
			arg_map[key] = strconv.Itoa(newValue)
		}
	}

	for key, value := range arg_map {
		if value == NO_CHANGE {
			continue
		}
		err = stub.PutState(key, []byte(value))
		if err != nil {
			return shim.Error("PutState failed :"+ key)
		}
	}

	return shim.Success(nil)
}

func (e *EWallet) transferAccount_new(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 付款方不变，收款方换名
	var Auser, Buser string
	var Abalance int
	var X int
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	Auser = args[0]
	Buser = args[1]

	Abalancebytes, err := stub.GetState(Auser)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Abalancebytes == nil {
		return shim.Error("Account not found")
	}
	Abalance, _ = strconv.Atoi(string(Abalancebytes))

	// Bbalancebytes, err := stub.GetState(Buser)
	// if err != nil {
	// 	return shim.Error("Failed to get state")
	// }
	// if Bbalancebytes == nil {
	// 	return shim.Error("Account not found")
	// }
	// Bbalance, _ = strconv.Atoi(string(Bbalancebytes))

	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	if X > Abalance {
		return shim.Error("Not sufficient funds")
	}
	Abalance = Abalance - X
	// Bbalance = Bbalance + X
	// fmt.Printf("Abalance = %d, Bbalance = %d\n", Abalance, Bbalance)

	// Write the state back to the ledger
	err = stub.PutState(Auser, []byte(strconv.Itoa(Abalance)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// err = stub.PutState(Buser, []byte(strconv.Itoa(Bbalance)))
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	/*
		换名：为Buser创建一个复合键新账户，余额为X（args[2]）
	*/
	// Retrieve info needed for the update procedure
	txid := stub.GetTxID()
	compositeIndexName := "varName~value~txID"

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{Buser, args[2], txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", Buser, compositeErr.Error()))
	}

	// Save the composite key index
	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", Buser, compositePutErr.Error()))
	}

	return shim.Success(nil)
}

func (e *EWallet) queryAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user := args[0]

	balancebytes, err := stub.GetState(user)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + user + "\"}"
		return shim.Error(jsonResp)
	}

	if balancebytes == nil {
		jsonResp := "{\"Error\":\"Nil balance for " + user + "\"}"
		return shim.Error(jsonResp)
	}
	
	// Get all deltas for the variable
	deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("varName~value~txID", []string{user})
	// 根据查询账中给定的部分复合键返回一个迭代器可用于遍历所有前缀匹配的复合键
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", user, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	//if !deltaResultsIterator.HasNext() {
		//return shim.Error(fmt.Sprintf("No variable by the name %s exists", user))
	//}

	// Iterate through result set and compute final value
	var finalVal float64
	finalVal = 0
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return shim.Error(splitKeyErr.Error())
		}

		// Retrieve the delta value and operation
		valueStr := keyParts[1]

		// Convert the value string and perform the operation
		value, convErr := strconv.ParseFloat(valueStr, 64)
		if convErr != nil {
			return shim.Error(convErr.Error())
		}

		finalVal += value
	}
	
	var balance, _ =  strconv.ParseFloat(string(balancebytes), 64)
	balance = finalVal + balance

	jsonResp := "{\"Name\":\"" + user + "\",\"Amount\":\"" + strconv.FormatFloat(balance, 'f', -1, 64) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(strconv.FormatFloat(balance, 'f', -1, 64)))
}

func (e *EWallet) deleteAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user := args[0]

	err := stub.DelState(user)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func (e *EWallet) saveMoney(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var balance int
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	user := args[0]

	balancebytes, err := stub.GetState(user)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if balancebytes == nil {
		e.createAccount(stub, []string{user})
	}
	balance, _ = strconv.Atoi(string(balancebytes))

	X, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}

	balance = balance + X
	fmt.Printf("balance = %d\n", balance)

	// Write the state back to the ledger
	err = stub.PutState(user, []byte(strconv.Itoa(balance)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (e *EWallet) drawMoney(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var balance int
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	user := args[0]

	balancebytes, err := stub.GetState(user)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if balancebytes == nil {
		return shim.Error("Account not found")
	}
	balance, _ = strconv.Atoi(string(balancebytes))

	X, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	if X > balance {
		return shim.Error("Not sufficient funds")
	}

	balance = balance - X
	fmt.Printf("balance = %d\n", balance)

	// Write the state back to the ledger
	err = stub.PutState(user, []byte(strconv.Itoa(balance)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (e *EWallet) queryAllAccounts(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

	startKey := "a"
	endKey := "zzzzzz"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"USER\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"BALANCE\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (e *EWallet) getHistoryForAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user := args[0]

	fmt.Printf("- start getHistoryForAccount: %s\n", user)

	resultsIterator, err := stub.GetHistoryForKey(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"USER\":")
		buffer.WriteString(user)

		buffer.WriteString(", \"BALANCE\":")
		buffer.WriteString(string(response.Value))
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		//		if response.IsDelete {
		//			buffer.WriteString("null")
		//		} else {
		//			buffer.WriteString(string(response.Value))
		//		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		//		buffer.WriteString(", \"IsDelete\":")
		//		buffer.WriteString("\"")
		//		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		//		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAccount returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (e *EWallet) beGoingToSave(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Extract the args
	user := args[0]
	_, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("Provided value was not a number")
	}

	// Retrieve info needed for the update procedure
	txid := stub.GetTxID()
	compositeIndexName := "varName~value~txID"

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{user, args[1], txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", user, compositeErr.Error()))
	}

	// Save the composite key index
	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", user, compositePutErr.Error()))
	}

	return shim.Success([]byte(fmt.Sprintf("Successfully added  %s to %s", args[1], user)))
}

func main() {
	err := shim.Start(new(EWallet))
	if err != nil {
		fmt.Printf("Error starting E-Wallet: %s", err)
	}
}
