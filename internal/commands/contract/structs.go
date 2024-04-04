package contract

type ContractReadOnlyFunctionsSip10 struct {
	FunctionName string
	ResponseType string
}

type ContractReadOnlyFunctionsSip10Response struct {
	FunctionName string
	Result       string
}

type ContractIDRequiredError struct {
	message string
}

func (e *ContractIDRequiredError) Error() string {
	return e.message
}
