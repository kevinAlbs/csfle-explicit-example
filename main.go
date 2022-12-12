package main

// An example running explicit encrypt and decrypt with range index.
// Run with: ./run.sh
//
// Set the environment variable KMS_PROVIDERS_PATH to the path of a JSON file with KMS credentials.
// KMS_PROVIDERS_PATH defaults to ~/.csfle/kms_providers.json.
//
// Set the environment variable MONGODB_URI to set a custom URI. MONGODB_URI defaults to
// mongodb://localhost:27017.

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// readFile reads a file into a byte slice.
func readFile(path string) []byte {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Panicf("error in Open: %v", err)
	}
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicf("error in ReadAll: %v", err)
	}

	return contents
}

// getKMSProvidersFromFile reads a JSON file for use as the KmsProviders option.
func getKMSProvidersFromFile(path string) map[string]map[string]interface{} {
	var kmsProviders map[string]map[string]interface{}
	contents := readFile(path)
	err := bson.UnmarshalExtJSON(contents, false, &kmsProviders)
	if err != nil {
		log.Panicf("error in UnmarshalExtJSON: %v", err)
	}

	return kmsProviders
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}

	var kmsProvidersPath string
	if kmsProvidersPath = os.Getenv("KMS_PROVIDERS_PATH"); kmsProvidersPath == "" {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		kmsProvidersPath = dirname + "/.csfle/kms_providers.json"
	}

	keyvaultClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicf("Connect error: %v\n", err)
	}
	defer keyvaultClient.Disconnect(context.TODO())

	kmsProviders := getKMSProvidersFromFile(kmsProvidersPath)

	// A ClientEncryption struct provides admin helpers with three functions:
	// 1. create a data key
	// 2. explicit encrypt
	// 3. explicit decrypt
	ceopts := options.ClientEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace("keyvault.datakeys")

	ce, err := mongo.NewClientEncryption(keyvaultClient, ceopts)
	if err != nil {
		log.Panicf("NewClientEncryption error: %v\n", err)
	}

	var keyid primitive.Binary
	{
		fmt.Printf("CreateDataKey... begin\n")
		keyid, err = ce.CreateDataKey(context.TODO(), "local", options.DataKey())
		if err != nil {
			log.Panicf("CreateDataKey error: %v\n", err)
		}
		fmt.Printf("Created key with a UUID: %v\n", hex.EncodeToString(keyid.Data))
		fmt.Printf("CreateDataKey... end\n")
	}

	var ciphertext primitive.Binary
	{
		fmt.Printf("Encrypt with range... begin\n")
		ro := options.RangeOptions{
			Min:      &bson.RawValue{Type: bsontype.Int32, Value: bsoncore.AppendInt32(nil, 0)},
			Max:      &bson.RawValue{Type: bsontype.Int32, Value: bsoncore.AppendInt32(nil, 200)},
			Sparsity: 1,
		}
		plaintext := bson.RawValue{Type: bsontype.Int32, Value: bsoncore.AppendInt32(nil, 123)}
		eOpts := options.Encrypt().SetAlgorithm("RangePreview").SetKeyID(keyid).SetRangeOptions(ro).SetContentionFactor(0)
		ciphertext, err = ce.Encrypt(context.TODO(), plaintext, eOpts)
		if err != nil {
			log.Panicf("Encrypt error: %v\n", err)
		}
		fmt.Printf("Explicitly encrypted to ciphertext: %x\n", ciphertext.Data)
		fmt.Printf("Encrypt with range... end\n")
	}

	{
		fmt.Printf("Decrypt... begin\n")
		plaintext, err := ce.Decrypt(context.TODO(), ciphertext)
		if err != nil {
			log.Panicf("Decrypt error: %v\n", err)
		}
		fmt.Printf("Explicitly decrypted to plaintext: %v\n", plaintext)
		fmt.Printf("Decrypt... end\n")
	}

	var encryptedColl *mongo.Collection
	{
		fmt.Printf("Create encrypted collection db.coll ... begin\n")
		encryptedFields := bson.M{
			"fields": bson.A{
				bson.M{
					"keyId":    keyid,
					"path":     "encryptedInt",
					"bsonType": "int",
					"queries": bson.A{
						bson.M{
							"queryType":  "rangePreview",
							"contention": int64(0),
							"sparsity":   1,
							"min":        int32(0),
							"max":        int32(200),
						},
					},
				},
			},
		}
		encryptedFieldsMap := map[string]interface{}{"db.coll": encryptedFields}

		aeOpts := options.AutoEncryption().
			SetKmsProviders(kmsProviders).
			SetKeyVaultNamespace("keyvault.datakeys").
			SetEncryptedFieldsMap(encryptedFieldsMap)

		encryptedClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetAutoEncryptionOptions(aeOpts))
		if err != nil {
			log.Panicf("Connect error: %v\n", err)
		}
		defer encryptedClient.Disconnect(context.TODO())

		encryptedColl = encryptedClient.Database("db").Collection("coll")
		err = encryptedColl.Drop(context.TODO())
		if err != nil {
			log.Panicf("Drop error: %v\n", err)
		}

		err = encryptedClient.Database("db").CreateCollection(context.TODO(), "coll")
		if err != nil {
			log.Panicf("CreateCollection error: %v\n", err)
		}

		fmt.Printf("Create encrypted collection db.coll ... end\n")
	}

	{
		fmt.Printf("Automatic encryption begin...\n")
		_, err = encryptedColl.InsertOne(context.TODO(), bson.D{{"encryptedInt", 123}})
		if err != nil {
			log.Panicf("InsertOne error: %v\n", err)
		}
		fmt.Printf("Automatic encryption... end\n")
	}

	{
		fmt.Printf("Automatic decryption begin...\n")

		filter := bson.D{{"encryptedInt", bson.D{{"$lte", 123 }}}}

		res := encryptedColl.FindOne(context.TODO(), filter)
		if res.Err() != nil {
			log.Panicf("FindOne error: %v\n", res.Err())
		}
		var decoded bson.Raw
		if err = res.Decode(&decoded); err != nil {
			log.Panicf("Decode error: %v\n", err)
		}
		fmt.Printf("Decrypted document: %v\n", decoded)
		fmt.Printf("Automatic decryption... end\n")
	}
}

/* Sample output
CreateDataKey... begin
Created key with a UUID: 8900362f9853455698370c91120a600c
CreateDataKey... end
Encrypt with range... begin
Explicitly encrypted to ciphertext: 04080900000564002000000000378f47763ec8758687124d9c8c10fe062d7e9f855a4dc9a8e0984596b11501ac0573002000000000359478ef8ddf3d64bbdb3c1151b9add47629ad535bc798b3cc4de888ba709f3e056300200000000055255dfabe586be364dc48253878881d207c3cac5840e089f708fa53c4953cd505700050000000009f16c1b6da8e177f388594042926f7b36a9539638ad194f781f00bd399c818b0f95c25aab78aef2016ac7a74b1b3f084a7ac794ec559c822d3da670922088cc1494fe1d94d43c56c5edd7796ed36a7cf05750010000000048900362f9853455698370c91120a600c1074001000000005760044000000008900362f9853455698370c91120a600c0421e6788642252f87ebdfc49d76e9a6e399ee34326fa5fe5fda08cb44dea30314868d20856d42278c07130ea03232b2e068595d056500200000000071ab491cbff1c61cdcec12ef1120f569a802d9309caa223b513edba8b209c8e00467009d070000033000d500000005640020000000006c5a1ffbc3fbb062c8cbfcd50e7ee2c3d36ccddb31847c700edf8e69acacebb7057300200000000090659ac186e2553c4036c0f22a4de7678435e7956ded090011a066ba7d5efb6e0563002000000000c94fdb5f1e634e7d51767d2cbde669cdecacc3420628708e340456231d6aeb3905700050000000002e9c80f66666883bc6ff7d1b6057b338de3f392a6334e7aaf5fd5975bae0fdfe6299a8c7d44b7a379fe802a200280b9113d25dbfad3bfc6885aadf7f9fe7e6e140b20f03e9b3143b48fb17698c075b7100033100d500000005640020000000009c7f3e15e14ec1482ecabe55f2765f9fdd190689a369f72d6eedf3e41f1950950573002000000000b5a5a124e8d736a90279a22ca1ad46a8632dc39ba595e276f35f2899009833c80563002000000000e19ba8a7acf7f9eca53d4bc043597e6f14aba96aa153964185c5edc7f1e58e78057000500000000001b335c6c1a516af663dd874ccb415430c8f3bc3b3de46b504e14a6c36a2e2467dfc29b323d533d7aac4df08d4bff394175acd566ad8d379c1869ee199866e608f7669d25c5b5ba1d6fdb39ea4c0da8100033200d500000005640020000000007d99cfc8b93ab784cbfee8d7923b1469a1cdc393c99cb7af47f2c0ca3f7550d205730020000000000662d755bd6829bbf9e035e33a30f8b5b944a7d7eba7c31e169f1039120eb92d05630020000000009af90e1f3d9043469be7d55234881b75dee40c6f90a5ae257a6ee1e0c4732f7c05700050000000005a6eeff6c027fd2e0399df158b3a8635b326e4dcc9254095f3afbf8c9bc07e4b16930c1ca8ce5a665903392b9d408a8cbbfa72bbd5a33c823518fdd10b9231869a4bfc7ede1ba45f8f4b6b7afd57efe500033300d50000000564002000000000b02cb69b0f4a439117545242a03e8b593dbc24318839f84bc894d78dde1860280573002000000000a70666769ea0180de7615faf441929b00505cb96978cdeeb8af1c4d7f53c36d7056300200000000039220459ebf0d67cc4206c04ba7b941f8ee116ac55ba2478b12fcc91a2347bc505700050000000009761c27e37a09c3c225213d812acd56ff7c8fbfe325d99b9fb6788cb3de267c3a496c8601b6b9da453581cb8dcbbc985a8250b070a563bf82f375d18ddd6d80a460d28d2662444c647bf2c89f3dcf5ee00033400d500000005640020000000000e38cd11a0877a285b648847444e72064b881ee446edf68120663cc43e47c11f0573002000000000f7c9fea0415b533caca1db8365fe7e048433c62ecaf551f523e2baa342f0cc5a0563002000000000e9ae4a3ed5e443e1432918effc48d741632e4295943387e3a9e577dcb67bd22605700050000000007e7ffe61dcfbe503ef97eb92b1fc93af828ee1b66f5b0e718b6d1e0ee7616c04505557f84169e262cecfdf6d079fd6e445efa390521ce3ab1c6e3ff2fb6f890a629aefc0cdaf34882c7f08930819aea400033500d50000000564002000000000f8f4986b34f3405fa1991053ee849c4ef3fb7a76640c58fa484a0dbf1d6ee6560573002000000000ccf0f72536a095c360f8921042c8a9980c5aba7f2b0400f1f79ea199186be6170563002000000000ab4620844ff14d18150558fac661462ccdc4641991326fffa6306c0cf515232f0570005000000000051c1a2430bcaefadaeb7ec1e62e45e03ad748b44d70ec593d27cff906a261f8136791a99e397a1cb44d212355ae22bba0391566e12b2694cafba74168ff0c3df93b18826d8d956663f73a46213fcfcb00033600d500000005640020000000003a8869ace0b03a0567e50221f5b843c22f34581abb5731409232d8b2391b0c8e0573002000000000b2eeca40ec82f0a68a19cef72eb02c87e37dce5c6b2d115b0f9a72f7fb1225b00563002000000000e71d8f64e01330d82759e97821701027156e393ca6637927c9268f79b5ec2afd057000500000000067374c50a7fad62573beba81e3f8666be73f33afac8c5bbbbe0b858c6727ee6cc6284248abf6ac73d171b05fe01d0afa7805b3c794ccd5cba503ecece026ea074d38406ffde041fcea9fc8feb71d4c3100033700d500000005640020000000009738dc36d8bed628010eba8ff09dfdcb619adf0ddbd52235b1239b155e50620205730020000000004b08db00d136d1ba32c7e9112b926098c576a4feb25ad2f6eab47a8a9a47d9fa05630020000000004546612471650e2072a8e7fc3bcc9fa86e48b5727cc482958912d28fa3f68da8057000500000000016fa68946e1fb10b55f837d9b25404ac18cb338b96c8ebdae3a10c123c7bc85b47f7e669876a93c9a64bb0246f19860ece852eeb02c293993c3ffbfc512cd1db457479e42f139f60ead35ae20316db2100033800d50000000564002000000000db4ad746405e58501c82259bf106ff55c2471d9933a361151651c1249cae9de305730020000000009a6535ad054c087a47bb1cb9ca0cf1fb847975e5fe9ac8c35ab1782c032c43b2056300200000000074fdd5463d0b183d8787eebbb0258ee2daa7f8acf9f69ef70e87ad33ed1b4782057000500000000092137fc7078dd697668b039a672fb01f8612aa1dc03642e2b5c47f3e7499cccb7bedef875ffadc33644f451bd60dd840b76f45b23b37560e4a09dbfc17cd5590b9364fb9a10ca8e9e34ee3bb0194e931000000
Encrypt with range... end
Decrypt... begin
Explicitly decrypted to plaintext: {"$numberInt":"123"}
Decrypt... end
Create encrypted collection db.coll ... begin
Create encrypted collection db.coll ... end
Automatic encryption begin...
Automatic encryption... end
Automatic decryption begin...
Decrypted document: {"_id": {"$oid":"63949f407d48e5b2bc6b57cc"},"encryptedInt": {"$numberInt":"123"},"__safeContent__": [{"$binary":{"base64":"PJb5JVhm/l8uIZuV4fZFzcL5u0EMEsbcxOnc1h62jaQ=","subType":"00"}},{"$binary":{"base64":"3OJ5Gw0TMDp2YqdCU5DSHk7cviYTfutkpulyj1kxL18=","subType":"00"}},{"$binary":{"base64":"pmKytP+7M1WLkHdDZ9SVY4MmpnI88JQntUi3sb6zSbI=","subType":"00"}},{"$binary":{"base64":"JpS5n/U5aIyjCUF1k/mfmmtHqj8Txw884l0wAu4mfcs=","subType":"00"}},{"$binary":{"base64":"ENO8YgCr0f4b9GHy/fVzVo8jyuMn6emnV2dWHa9u/6k=","subType":"00"}},{"$binary":{"base64":"+9EQWsw/g7dc5Te78eUFAOw0GpdWAfaRxvAsM/nJRvk=","subType":"00"}},{"$binary":{"base64":"e4tKJjMJ/NPO68zAXk7pU/CfdcgvbY51GLSJwOTqnlA=","subType":"00"}},{"$binary":{"base64":"pTLeSYTT+YP+dXDFcXGOuPgo/6e0AJ1amU2qXJlfdD8=","subType":"00"}},{"$binary":{"base64":"SBPSKi7Da60MaeiK1cH3z7P8LA8W9dU85sY9C/xrTD4=","subType":"00"}}]}
Automatic decryption... end
*/
