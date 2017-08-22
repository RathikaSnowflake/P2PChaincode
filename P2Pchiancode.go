package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Procure2Pay will implement the processes
type Procure2Pay struct {
}

type purchaseOrder struct {
	ID           string  `json:"po_id"`
	SupplierID   string  `json:"supplier_id"`
	Quantity     int     `json:"qty"`
	PriceUnit    float32 `json:"price_unit"`
	Item         string  `json:"item_details"`
	ItemID       string  `json:"item_id"`
	Date         string  `json:"date"`
	Status       string  `json:"status"`
	PoECustoms   string  `json:"poe_Customs"`
	BillOfLading string  `json:"billOflading"`
	Invoice      string  `json:"invoice"`
}

type item struct {
	ID        string  `json:"item_id"`
	Item      string  `json:"item"`
	Quantity  int     `json:"quantity"`
	PriceUnit float32 `json:"price_unit"`
}

type gnrMessage struct {
	Generated int    `json:"isgenerated"`
	Msg       string `json:"message"`
}

var supplierItemMap map[string][]item
var purchaseOrderMap map[string]purchaseOrder

//getPurchaseOrderMap
func getPurchaseOrderMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = stub.GetState("PurchaseOrderMap")
	if err != nil {
		fmt.Printf("Failed to get the PurchaseOrderMap for block chain :%v\n", err)
		return err
	}
	if len(bytesread) != 0 {
		fmt.Printf("PurchaseOrderMap map exists.\n")
		err = json.Unmarshal(bytesread, &purchaseOrderMap)
		if err != nil {
			fmt.Printf("Failed to initialize the PurchaseOrderMap for block chain :%v\n", err)
			return err
		}
	} else {
		fmt.Printf("PurchaseOrderMap map does not exist. To be created. \n")
		purchaseOrderMap = make(map[string]purchaseOrder)
		bytesread, err = json.Marshal(&purchaseOrderMap)
		if err != nil {
			fmt.Printf("Failed to initialize the PurchaseOrderMap for block chain :%v\n", err)
			return err
		}
		err = stub.PutState("PurchaseOrderMap", bytesread)
		if err != nil {
			fmt.Printf("Failed to initialize the PurchaseOrderMap for block chain :%v\n", err)
			return err
		}
	}
	return nil
}

//setPurchaseOrderMap
func setPurchaseOrderMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = json.Marshal(&purchaseOrderMap)
	if err != nil {
		fmt.Printf("Failed to set the PurchaseOrderMap for block chain :%v\n", err)
		return err
	}
	err = stub.PutState("PurchaseOrderMap", bytesread)
	if err != nil {
		fmt.Printf("Failed to set the PurchaseOrderMap %v\n", err)
		return errors.New("Failed to set the PurchaseOrderMap")
	}

	return nil
}

//getSupplierItemMap
func getSupplierItemMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = stub.GetState("SupplierItemMap")
	if err != nil {
		fmt.Printf("Failed to get the SupplierItemMap for block chain :%v\n", err)
		return err
	}
	if len(bytesread) != 0 {
		fmt.Printf("ItemMap map exists.\n")
		err = json.Unmarshal(bytesread, &supplierItemMap)
		if err != nil {
			fmt.Printf("Failed to initialize the SupplierItemMap for block chain :%v\n", err)
			return err
		}
	} else {
		fmt.Printf("ItemMap map does not exist. To be created. \n")
		supplierItemMap = make(map[string][]item)
		bytesread, err = json.Marshal(&supplierItemMap)
		if err != nil {
			fmt.Printf("Failed to initialize the SupplierItemMap for block chain :%v\n", err)
			return err
		}
		err = stub.PutState("SupplierItemMap", bytesread)
		if err != nil {
			fmt.Printf("Failed to initialize the SupplierItemMap for block chain :%v\n", err)
			return err
		}
	}
	return nil
}

//setSupplierItemMap
func setSupplierItemMap(stub shim.ChaincodeStubInterface) error {
	var err error
	var bytesread []byte

	bytesread, err = json.Marshal(&supplierItemMap)
	if err != nil {
		fmt.Printf("Failed to set the SupplierItemMap for block chain :%v\n", err)
		return err
	}
	err = stub.PutState("SupplierItemMap", bytesread)
	if err != nil {
		fmt.Printf("Failed to set the SupplierItemMap %v\n", err)
		return errors.New("Failed to set the SupplierItemMap")
	}

	return nil
}

// GetAllPurchaseOrders gets the details of the All PO
func GetAllPurchaseOrders(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var list []purchaseOrder
	var err error
	var bytesRead []byte

	err = getPurchaseOrderMap(stub)
	if err != nil {
		fmt.Printf("Unable to read the list of purchase orders : %s\n", err)
		return nil, err
	}

	for _, value := range purchaseOrderMap {
		// add any checks if required
		list = append(list, value)
	}
	fmt.Printf("list of POs : %v\n", list)
	bytesRead, err = json.Marshal(&list)
	if err != nil {
		fmt.Printf("Unable to return the list of purchase orders : %s\n", err)
		return nil, err
	}

	return bytesRead, nil
}

// GetPOForSupplier gets the list of the PO for a Supplier
func GetPOForSupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var object []purchaseOrder
	var bytes []byte

	fmt.Println("Entering GetPOForSupplier")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing Supplier ID")
	}
	fmt.Printf("Entering GetPurchaseOrder : %v\n", args[0])
	var supplierID = args[0]

	for _, value := range purchaseOrderMap {
		if value.SupplierID == supplierID {
			object = append(object, value)
		}
	}
	bytes, err := json.Marshal(&object)
	if err != nil {
		fmt.Printf("Unable to marshal the object array %s\n", err)
		return nil, err
	}

	fmt.Printf(" List of POs for Supplier : %v\n", object)
	return bytes, nil
}

// GetPurchaseOrder gets the details of the PO
func GetPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var bytes []byte
	fmt.Println("Entering GetPurchaseOrder : ", args[0])

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing purchase order ID")
	}

	getPurchaseOrderMap(stub)
	if object, ok := purchaseOrderMap[args[0]]; ok {
		bytes, err = json.Marshal(&object)
		if err != nil {
			fmt.Printf(" Error : %v\n", err)
			return nil, errors.New("PO not found")
		}
		fmt.Printf("PO found : %v\n", object)
		return bytes, nil
	}
	return nil, errors.New("PO not found")

}

// CreatePurchaseOrder creates/updates the PO
func CreatePurchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var object purchaseOrder
	var list []item
	var ok bool

	var err error
	fmt.Println("Entering CreatePurchaseOrder")

	if len(args) < 2 {
		fmt.Println("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for purchase order creation")
	}

	err = json.Unmarshal([]byte(args[1]), &object)
	if err != nil {
		fmt.Printf("Unable to marshal the purchase order input : %s\n", err)
		return nil, nil
	}

	// saving the purchase order in the map
	getPurchaseOrderMap(stub)
	getSupplierItemMap(stub)

	purchaseOrderMap[object.ID] = object

	if list, ok = supplierItemMap[object.SupplierID]; ok {
		fmt.Printf(" Item : %v\n", list)
		for i, id := range list {
			fmt.Printf(" id : %v\n", id)
			fmt.Printf(" i : %v\n", i)
			if id.ID == object.ItemID {
				fmt.Printf(" id.ID : %v\n", id.ID)
				fmt.Printf(" object.ItemID : %v\n", object.ItemID)
				id.Quantity = id.Quantity - object.Quantity
				supplierItemMap[object.SupplierID][i] = id

			}
		}
	}

	setSupplierItemMap(stub)
	setPurchaseOrderMap(stub)
	fmt.Printf("purchaseOrderMap : %v \n", purchaseOrderMap)
	fmt.Printf("supplierItemMap : %v \n", supplierItemMap)

	fmt.Println("Successfully saved purchase order")
	return nil, nil
}

// GetItemForSupplier gets the details of the Item
func GetItemForSupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var object []item
	var ok bool
	var bytes []byte
	var err error
	fmt.Println("Entering GetItemForSupplier")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing Supplier ID")
	}
	fmt.Printf("Supplier ID : %v\n", args[0])
	getSupplierItemMap(stub)
	if object, ok = supplierItemMap[args[0]]; ok {
		bytes, err = json.Marshal(&object)
		if err != nil {
			fmt.Printf(" error : %v\n", err)
		}
		fmt.Printf(" Item : %v\n", object)
		return bytes, nil
	}

	return nil, errors.New("Didnt find any items for the supplier")

}

// CreateItems creates the Items
func CreateItems(stub shim.ChaincodeStubInterface) error {
	var object item
	var err error

	getSupplierItemMap(stub)

	fmt.Println("Entering createItems")
	object.ID = "CFBEAN"
	object.Item = "Coffee Beans"
	object.PriceUnit = 150
	object.Quantity = 5000
	supplierItemMap["8888"] = append(supplierItemMap["8888"], object)

	object.ID = "CFBEAN"
	object.Item = "Coffee Beans"
	object.PriceUnit = 150
	object.Quantity = 6000
	supplierItemMap["9999"] = append(supplierItemMap["9999"], object)

	object.ID = "CFBEAN"
	object.Item = "Coffee Beans"
	object.PriceUnit = 150
	object.Quantity = 7000
	supplierItemMap["7777"] = append(supplierItemMap["7777"], object)

	if err != nil {
		fmt.Println("error : ", err)
	}
	setSupplierItemMap(stub)
	fmt.Printf("Successfully saved items : %v\n", supplierItemMap)
	return nil

}

// UpdateItemForSupplier updates the Item
func UpdateItemForSupplier(stub shim.ChaincodeStubInterface, args []string) error {
	var list []item
	var inputItem item
	var ok bool
	fmt.Println("Entering UpdateItemForSupplier")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return errors.New("Missing Supplier ID")
	}

	getSupplierItemMap(stub)
	fmt.Printf("item map : %v\n", supplierItemMap)
	var err = json.Unmarshal([]byte(args[1]), &inputItem)
	fmt.Printf(" inputItem: %v\n", inputItem)
	fmt.Printf("Supplier ID : %v\n", args[0])
	if list, ok = supplierItemMap[args[0]]; ok {
		fmt.Printf(" Item : %v\n", list)
		for i, id := range list {
			if id.ID == inputItem.ID {
				id.PriceUnit = inputItem.PriceUnit
				id.Quantity = inputItem.Quantity
				supplierItemMap[args[0]][i] = id

			}
		}
	}
	if err != nil {
		fmt.Println("Error : ", err)
	}
	setSupplierItemMap(stub)
	fmt.Printf("item map : %v\n", supplierItemMap)
	fmt.Println("Successfully saved item")
	return nil

}

// UpdateBillOfLading adds the BillOfLading to PO
func UpdateBillOfLading(stub shim.ChaincodeStubInterface, args []string) error {
	var object purchaseOrder
	fmt.Println("Entering UpdateBillOfLading")

	if len(args) < 2 {
		fmt.Println("Invalid number of arguments")
		return errors.New("Missing details for purchase order ID " + args[0])
	}
	fmt.Printf("PO id : %v\n", args[0])
	fmt.Printf("BillOfLading : %v\n", args[1])

	getPurchaseOrderMap(stub)
	object = purchaseOrderMap[args[0]]
	object.BillOfLading = args[1]
	purchaseOrderMap[args[0]] = object
	setPurchaseOrderMap(stub)

	//fmt.Printf("Successfully saved purchase order :%v\n", purchaseOrderMap)
	return nil
}

// UpdatePoECustoms adds the POE Docs to PO
func UpdatePoECustoms(stub shim.ChaincodeStubInterface, args []string) error {

	var object purchaseOrder
	fmt.Println("Entering UpdatePoECustoms")

	if len(args) < 2 {
		fmt.Println("Invalid number of arguments")
		return errors.New("Missing details for purchase order ID " + args[0])
	}
	fmt.Printf("PO id : %v\n", args[0])
	fmt.Printf("PoECustoms : %v\n", args[1])

	getPurchaseOrderMap(stub)
	object = purchaseOrderMap[args[0]]
	object.PoECustoms = args[1]
	purchaseOrderMap[args[0]] = object
	setPurchaseOrderMap(stub)

	//fmt.Printf("Successfully saved purchase order :%v\n", purchaseOrderMap)
	return nil
}

// UpdateInvoice adds the POE Docs to PO
func UpdateInvoice(stub shim.ChaincodeStubInterface, args []string) error {
	var object purchaseOrder
	fmt.Println("Entering UpdateInvoice")

	if len(args) < 2 {
		fmt.Println("Invalid number of arguments")
		return errors.New("Missing details for purchase order ID " + args[0])
	}
	fmt.Printf("PO id : %v\n", args[0])
	fmt.Printf("Invoice : %v\n", args[1])

	getPurchaseOrderMap(stub)
	object = purchaseOrderMap[args[0]]
	object.Invoice = args[1]
	object.Status = "sent"
	purchaseOrderMap[args[0]] = object
	setPurchaseOrderMap(stub)

	//fmt.Printf("Successfully saved purchase order :%v\n", purchaseOrderMap)
	return nil
}

// CreateGNR creates the GNR for PO
func CreateGNR(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var err error
	var object purchaseOrder
	var response gnrMessage
	fmt.Println("Entering createGNR")

	if len(args) < 2 {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing details for purchase order ID " + args[0])
	}
	fmt.Printf("PO id : %v\n", args[0])
	fmt.Printf("Qty Received : %v\n", args[1])

	getPurchaseOrderMap(stub)
	object = purchaseOrderMap[args[0]]

	var qtyint, _ = strconv.Atoi(args[1])

	if qtyint < object.Quantity {
		var qty = object.Quantity * 3 / 4
		if qtyint >= qty {
			response.Msg = "GNR generated successfully as quantity received is more than 70% of what was requested"
			response.Generated = 1
		} else {
			response.Msg = "GNR failed to generate as quantity received is less than 70% of what was requested"
			response.Generated = 0
		}
	}
	fmt.Printf("GNR status : %v\n", response)
	bytes, err := json.Marshal(&response)
	return bytes, err
}

// DeleteAllPOs will create a new PurchaseOrderMap
func DeleteAllPOs(stub shim.ChaincodeStubInterface) ([]byte, error) {
	var err error
	var byteArray []byte

	fmt.Println("DeleteAllPOs will create a new instance of PurchaseOrderMap")
	err = getPurchaseOrderMap(stub)
	purchaseOrderMap = make(map[string]purchaseOrder)
	byteArray, err = json.Marshal(&purchaseOrderMap)
	err = setPurchaseOrderMap(stub)
	err = getPurchaseOrderMap(stub)
	fmt.Printf("PurchaseOrderMap created : %v\n", purchaseOrderMap)

	return byteArray, err

}

// Init sets up the chaincode
func (t *Procure2Pay) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Inside INIT for test chaincode")
	return nil, nil
}

// Query the chaincode
func (t *Procure2Pay) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetPurchaseOrder" {
		return GetPurchaseOrder(stub, args)
	} else if function == "GetItemForSupplier" {
		return GetItemForSupplier(stub, args)
	} else if function == "GetPOForSupplier" {
		return GetPOForSupplier(stub, args)
	} else if function == "GetAllPurchaseOrders" {
		return GetAllPurchaseOrders(stub, args)
	} else if function == "CreateGNR" {
		return CreateGNR(stub, args)
	}
	return nil, nil
}

// Invoke the function in the chaincode
func (t *Procure2Pay) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "CreatePurchaseOrder" {
		return CreatePurchaseOrder(stub, args)
	} else if function == "UpdateItemForSupplier" {
		return nil, UpdateItemForSupplier(stub, args)
	} else if function == "UpdateBillOfLading" {
		return nil, UpdateBillOfLading(stub, args)
	} else if function == "UpdatePoECustoms" {
		return nil, UpdatePoECustoms(stub, args)
	} else if function == "UpdateInvoice" {
		return nil, UpdateInvoice(stub, args)
	} else if function == "CreateItems" {
		return nil, CreateItems(stub)
	} else if function == "DeleteAllPOs" {
		return DeleteAllPOs(stub)
	}

	fmt.Println("Function not found")
	return nil, nil
}

func main() {
	err := shim.Start(new(Procure2Pay))
	if err != nil {
		fmt.Println("Could not start Procure2Pay")
	} else {
		fmt.Println("Procure2Pay successfully started")
	}

}
