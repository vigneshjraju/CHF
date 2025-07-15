package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type RiceBatch struct {
	BatchID       string `json:"batchID"`
	Farmer        string `json:"farmer"`
	HarvestDate   string `json:"harvestDate"`
	Quantity      int    `json:"quantity"`
	Location      string `json:"location"`
	QualityGrade  string `json:"qualityGrade"`
	CurrentHolder string `json:"currentHolder"`
}

// 🔐 Private Data Structure
type PrivateDetails struct {
	PricePerKg float64 `json:"pricePerKg"`
	GradeNote  string  `json:"gradeNote"`
}

// 🚜 Add Rice Batch (Farmer)
func (s *SmartContract) AddRiceBatch(ctx contractapi.TransactionContextInterface, batchID, harvestDate, quantityStr, location, qualityGrade string) error {
	exists, err := s.RiceBatchExists(ctx, batchID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("batch %s already exists", batchID)
	}

	quantity, _ := strconv.Atoi(quantityStr)
	clientID, _ := ctx.GetClientIdentity().GetID()

	batch := RiceBatch{
		BatchID:       batchID,
		Farmer:        clientID,
		HarvestDate:   harvestDate,
		Quantity:      quantity,
		Location:      location,
		QualityGrade:  qualityGrade,
		CurrentHolder: "Farmer",
	}

	batchJSON, _ := json.Marshal(batch)
	return ctx.GetStub().PutState(batchID, batchJSON)
}

// 🔐 Add Private Details (Farmer)
func (s *SmartContract) AddPrivateDetails(ctx contractapi.TransactionContextInterface, batchID string) error {
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return err
	}

	privateJSON, ok := transientMap["privateDetails"]
	if !ok {
		return fmt.Errorf("privateDetails key not found in transient map")
	}

	var details PrivateDetails
	err = json.Unmarshal(privateJSON, &details)
	if err != nil {
		return fmt.Errorf("unmarshal error: %v", err)
	}

	privateDetailsJSON, _ := json.Marshal(details)
	return ctx.GetStub().PutPrivateData("ricePrivateData", batchID, privateDetailsJSON)
}

// 🏭 Transfer to Miller
func (s *SmartContract) TransferToMiller(ctx contractapi.TransactionContextInterface, batchID string) error {
	batch, err := s.ReadRiceBatch(ctx, batchID)
	if err != nil {
		return err
	}
	batch.CurrentHolder = "Miller"
	batchJSON, _ := json.Marshal(batch)
	return ctx.GetStub().PutState(batchID, batchJSON)
}

// 🛒 Transfer to Retailer
func (s *SmartContract) TransferToRetailer(ctx contractapi.TransactionContextInterface, batchID string) error {
	batch, err := s.ReadRiceBatch(ctx, batchID)
	if err != nil {
		return err
	}
	batch.CurrentHolder = "Retailer"
	batchJSON, _ := json.Marshal(batch)
	return ctx.GetStub().PutState(batchID, batchJSON)
}

// 📖 Read Rice Batch (Public)
func (s *SmartContract) ReadRiceBatch(ctx contractapi.TransactionContextInterface, batchID string) (*RiceBatch, error) {
	batchJSON, err := ctx.GetStub().GetState(batchID)
	if err != nil {
		return nil, err
	}
	if batchJSON == nil {
		return nil, fmt.Errorf("batch %s does not exist", batchID)
	}

	var batch RiceBatch
	_ = json.Unmarshal(batchJSON, &batch)
	return &batch, nil
}

// 🔎 Read Private Data (Farmer/Miller only)
func (s *SmartContract) ReadPrivateDetails(ctx contractapi.TransactionContextInterface, batchID string) (*PrivateDetails, error) {
	data, err := ctx.GetStub().GetPrivateData("ricePrivateData", batchID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("no private data found for batch %s", batchID)
	}

	var details PrivateDetails
	_ = json.Unmarshal(data, &details)
	return &details, nil
}

// 📊 Rich Query by Location
func (s *SmartContract) QueryByLocation(ctx contractapi.TransactionContextInterface, location string) ([]*RiceBatch, error) {
	query := fmt.Sprintf(`{"selector":{"location":"%s"}}`, location)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []*RiceBatch
	for resultsIterator.HasNext() {
		result, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var batch RiceBatch
		_ = json.Unmarshal(result.Value, &batch)
		results = append(results, &batch)
	}
	return results, nil
}

// 🔍 Query by QualityGrade
func (s *SmartContract) QueryByQuality(ctx contractapi.TransactionContextInterface, quality string) ([]*RiceBatch, error) {
	query := fmt.Sprintf(`{"selector":{"qualityGrade":"%s"}}`, quality)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []*RiceBatch
	for resultsIterator.HasNext() {
		result, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var batch RiceBatch
		_ = json.Unmarshal(result.Value, &batch)
		results = append(results, &batch)
	}
	return results, nil
}

// ✅ Check if RiceBatch exists
func (s *SmartContract) RiceBatchExists(ctx contractapi.TransactionContextInterface, batchID string) (bool, error) {
	data, err := ctx.GetStub().GetState(batchID)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(fmt.Sprintf("Error creating rice supply chain chaincode: %v", err))
	}
	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Error starting chaincode: %v", err))
	}
}
