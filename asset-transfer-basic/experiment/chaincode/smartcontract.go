package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Task
type SmartContract struct {
	contractapi.Contract
}

// Task describes basic details of what makes up a simple Task
//Insert struct field in alphabetic order => to achieve determinism accross languages
// golang keeps the order when marshal to json but doesn't order automatically


type Task struct {
	Publisher           string `json:"publisher"`
	PublishDateTime    string `json:"publishDateTime"`
	Cpu string `json:"cpu"`
	Mem string `json:"mem"`
	Storage string `json:"storage"`
	ConsumeTime string `json:"consumeTime"`
	Executor            string `json:"executor"`
}


// CreateTask issues a new Task to the world state with given details.
func (s *SmartContract) CreateTask(ctx contractapi.TransactionContextInterface, publisher_ string, executor_ string, publishDateTime_ string, consumeTime_ string, Cpu_ string, Mem_ string, Storage_ string) error {
	txID := ctx.GetStub().GetTxID()
	id := publisher_ + txID
	exists, err := s.TaskExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Task %s already exists", id)
	}

	task := Task{
		Publisher:  publisher_,
		Executor: executor_,
		PublishDateTime:           publishDateTime_,
		ConsumeTime:          consumeTime_,
		Cpu: Cpu_,
		Mem: Mem_,
		Storage: Storage_,
	}
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, taskJSON)
}

// Instantiate does nothing
func (c *SmartContract) Instantiate() {
	fmt.Println("Instantiated")
}

func (s *SmartContract) TaskExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	TaskJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return TaskJSON != nil, nil
}

func (s *SmartContract) ReadTask(ctx contractapi.TransactionContextInterface, id string) (*Task, error) {
	TaskJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if TaskJSON == nil {
		return nil, fmt.Errorf("the Task %s does not exist", id)
	}

	var task Task
	err = json.Unmarshal(TaskJSON, &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
