package parse

import (
	"regexp"
	"strconv"
	"strings"

	govaluate "github.com/blackskygg/cc/third/govaluate_modified"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func makeParameterMap(stub shim.ChaincodeStubInterface, exp, id string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	re, err := regexp.Compile(`a(_\w+)+`)
	if err != nil {
		return result, err
	}

	l := re.FindAllString(exp, -1)
	for _, w := range l {
		if w == "a_id" {
			result[w] = id
			continue
		}

		key := strings.Join([]string{id, w[2:]}, "_")
		re, err := stub.GetState(key)
		val, err := strconv.Atoi(string(re))
		if err != nil {
			return result, err
		}
		result[w] = interface{}(val)
	}

	return result, nil

}

func Eval(exp string, stub shim.ChaincodeStubInterface, id string) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(exp)
	var params map[string]interface{}
	if err != nil {
		return false, err
	}

	params, err = makeParameterMap(stub, exp, id)
	if err != nil {
		return false, err
	}

	result, err := expression.Evaluate(params)
	return result.(bool), err
}
