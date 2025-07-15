package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	ccName      = "rice"
	channelName = "mychannel"
	walletPath  = "wallet"
	ccpPath     = "../../network/organizations/peerOrganizations/farmerorg.example.com/connection-farmerorg.yaml"
)

func main() {
	fmt.Println("\n========== RICE SUPPLY CHAIN ==========")

	gw, err := connectGateway()
	if err != nil {
		log.Fatalf("❌ Gateway connection failed: %v", err)
	}
	defer gw.Close()

	contract := getContract(gw)

	// 👇 Uncomment whichever function you want to test
	AddRiceBatch(contract)
	AddPrivateDetails(contract)
	ReadRiceBatch(contract)
	ReadPrivateDetails(contract)
	QueryByLocation(contract)
	QueryByQuality(contract)
	TransferToMiller(contract)
	TransferToRetailer(contract)
}

// 🔌 Setup Gateway and Contract

func connectGateway() (*gateway.Gateway, error) {
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return nil, err
	}
	if !wallet.Exists("appUser") {
		return nil, fmt.Errorf("appUser not enrolled in wallet")
	}
	return gateway.Connect(
		gateway.WithConfig(gateway.ConfigFromYAMLFile(ccpPath)),
		gateway.WithIdentity(wallet, "appUser"),
	)
}

func getContract(gw *gateway.Gateway) *gateway.Contract {
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	return network.GetContract(ccName)
}

////////////////////////////////////////////////////
// 🧩 INDIVIDUAL FUNCTION WRAPPERS
////////////////////////////////////////////////////

// 🚜 Farmer adds new rice batch
func AddRiceBatch(contract *gateway.Contract) {
	fmt.Println("\n➡️ AddRiceBatch")
	_, err := contract.SubmitTransaction("AddRiceBatch", "BATCH101", "2025-07-01", "100", "Kuttanad", "GradeA")
	if err != nil {
		log.Fatalf("❌ AddRiceBatch failed: %v", err)
	}
	fmt.Println("✅ Rice batch BATCH101 added successfully.")
}

// 🔐 Add private details via transient map
func AddPrivateDetails(contract *gateway.Contract) {
	fmt.Println("\n➡️ AddPrivateDetails (Transient Map)")
	privateData := map[string]interface{}{
		"pricePerKg": 45.0,
		"gradeNote":  "Premium Quality",
	}
	transientJSON, _ := json.Marshal(privateData)
	transient := map[string][]byte{
		"privateDetails": transientJSON,
	}

	_, err := contract.SubmitWithOpts("AddPrivateDetails",
		gateway.WithArguments("BATCH101"),
		gateway.WithTransient(transient),
	)
	if err != nil {
		log.Fatalf("❌ AddPrivateDetails failed: %v", err)
	}
	fmt.Println("✅ Private details added.")
}

// 🧐 Read public data of batch
func ReadRiceBatch(contract *gateway.Contract) {
	fmt.Println("\n🔍 ReadRiceBatch")
	result, err := contract.EvaluateTransaction("ReadRiceBatch", "BATCH101")
	if err != nil {
		log.Fatalf("❌ ReadRiceBatch failed: %v", err)
	}
	fmt.Printf("📦 Batch Info: %s\n", result)
}

// 🔒 Read private data
func ReadPrivateDetails(contract *gateway.Contract) {
	fmt.Println("\n🔐 ReadPrivateDetails")
	result, err := contract.EvaluateTransaction("ReadPrivateDetails", "BATCH101")
	if err != nil {
		log.Printf("⚠️ Failed to read private data (maybe unauthorized): %v\n", err)
		return
	}
	fmt.Printf("🔐 Private Info: %s\n", result)
}

// 📍 Query all batches by location
func QueryByLocation(contract *gateway.Contract) {
	fmt.Println("\n📍 QueryByLocation")
	result, err := contract.EvaluateTransaction("QueryByLocation", "Kuttanad")
	if err != nil {
		log.Fatalf("❌ QueryByLocation failed: %v", err)
	}
	fmt.Printf("📍 Result: %s\n", result)
}

// 🏷️ Query all batches by grade
func QueryByQuality(contract *gateway.Contract) {
	fmt.Println("\n🏷️ QueryByQuality")
	result, err := contract.EvaluateTransaction("QueryByQuality", "GradeA")
	if err != nil {
		log.Fatalf("❌ QueryByQuality failed: %v", err)
	}
	fmt.Printf("🏷️ Result: %s\n", result)
}

// 🚚 Transfer to Miller
func TransferToMiller(contract *gateway.Contract) {
	fmt.Println("\n➡️ TransferToMiller")
	_, err := contract.SubmitTransaction("TransferToMiller", "BATCH101")
	if err != nil {
		log.Fatalf("❌ TransferToMiller failed: %v", err)
	}
	fmt.Println("✅ Transferred to Miller.")
}

// 🏪 Transfer to Retailer
func TransferToRetailer(contract *gateway.Contract) {
	fmt.Println("\n➡️ TransferToRetailer")
	_, err := contract.SubmitTransaction("TransferToRetailer", "BATCH101")
	if err != nil {
		log.Fatalf("❌ TransferToRetailer failed: %v", err)
	}
	fmt.Println("✅ Transferred to Retailer.")
}
