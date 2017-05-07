package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blackskygg/cc/parse"
	"github.com/blackskygg/chaincode/config"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	conf, err := config.FromFile("init.conf")
	conf.ApplyConfig(stub)
	return nil, err
}

func (t *SimpleChaincode) checkPermission(table_name, role string) error {
	switch {
	case table_name == "student" && role == "ABoss":
	case table_name == "pay" && role == "FBoss":
	case table_name == "staff" && role == "PBoss":
	case table_name == "netusr" && role == "NBoss":
	default:
		return errors.New("Permission denied!")
	}

	return nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	role := string(args[0])
	switch function {
	case "add":
		id := string(args[1])
		table_name := string(args[2])
		var_name := string(args[3])
		if err := t.checkPermission(table_name, role); err != nil {
			return "Permission Denied!", err
		}

		key := strings.Join([]string{id, table_name, var_name}, "_")
		stub.PutState(key, args[4])
		return nil, nil

	case "del":
		id := strings(args[1])
		table_name := string(args[2])
		var_name := string(args[3])
		if err := t.checkPermission(table_name, role); err != nil {
			return "Permission Denied!", err
		}

		key := strings.Join([]string{id, table_name, var_name}, "_")
		stub.DelState(key)

		return nil, nil
	default:
		return nil, errors.New("")
	}

	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "query":
		table_name := string(args[0])
		var_name := string(args[1])
		key := strings.Join([]string{table_name, var_name}, "_")
		if ret, err := stub.GetState(key); err != nil {
			return nil, errors.New("no such key")
		}

		return ret, nil

	case "evaluate":
		var result string
		var val bool
		var err error
		if val, err = parse.Eval(expression, stub, id); err != nil {
			return nil, err
		}

		if val {
			result = "true"
		} else {
			result = "false"
		}

		return []byte(result), nil
	default:
		return []byte{}, nil
	}

	return []byte{}, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
