/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

const keyPrefix = "key"

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Demo Init")
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// Initialize the chaincode
	//keyPrefix = copy(keyPrefix, args[0])

	keyFrom := keyPrefix + args[1]
	keyTo := keyPrefix + args[2]

	fmt.Printf("keyPrefix = %s, key range from [%s] to [%s]\n", keyPrefix, keyFrom, keyTo)

	from, errfrom := strconv.Atoi(args[1])
	if errfrom != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	to, errTo := strconv.Atoi(args[2])
	if errTo != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	// Write the state to the ledger
	for i := from; i <= to; i++ {
		err := stub.PutState(generalKey(i), []byte(strconv.Itoa(i)))
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Demo Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "increase" {
		// Increase x from a to b
		return t.increase(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// increase increase x from a to b
func (t *SimpleChaincode) increase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	from, errfrom := strconv.Atoi(args[0])
	if errfrom != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	to, errTo := strconv.Atoi(args[1])
	if errTo != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	addVal, errAddVal := strconv.Atoi(args[2])
	if errAddVal != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	for i := from; i <= to; i++ {
		originNum, err := stub.GetState(generalKey(i))
		if err != nil {
			return shim.Error("Failed to get state")
		}
		err = stub.PutState(generalKey(i), []byte(addSI(string(originNum), addVal)) )
		if err != nil {
			return shim.Error("Failed to put state")
		}
	}

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	// Delete the key from the state in ledger
	err = stub.DelState(generalKey(A))
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	qKey := args[0]
	//fmt.Printf("Query Response:%s\n", keyPrefix + qKey)
	// Get the state from the ledger
	result, err := stub.GetState(keyPrefix + qKey)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + qKey + "\"}"
		return shim.Error(jsonResp)
	}

	if result == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + qKey + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + qKey + "\",\"Amount\":\"" + string(result) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(result)
}


func generalKey (x int) string { 
	return keyPrefix + strconv.Itoa(x)
}

func addSI(snum string, anum int) string {
	bnum, _ := strconv.Atoi(snum)
	sum := strconv.Itoa(anum + bnum)
	return sum
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
