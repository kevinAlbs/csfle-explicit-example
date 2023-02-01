An example using FLE2 Range Support for explicit encryption/decryption and automatic encryption/decryption.

# Run Instructions

At time of writing, this requires an unreleased branch of the Go driver and an unstable release of libmongocrypt.

First, make sure Brew had been updated & upgraded, which will bring along any required XCode updates as well.

Then, clone the branch with FLE2 Range Support:

```
git clone git@github.com:kevinAlbs/mongo-go-driver.git --branch DRIVERS-2505 mongo-go-driver-DRIVERS-2505
```

Update go.mod (for this csfle example repo) and change the path in the `replace` directive at the bottom of the file with the local checked-out go driver path. Something like this:
```
// Use local checkout of Go driver to get. TODO: replace this path.
replace go.mongodb.org/mongo-driver => /Users/kevin.albertson/code/csfle-explicit-example/mongo-go-driver-DRIVERS-2505
```

Download and install the latest libmongocrypt 1.7.0-alpha1. Might need to uninstall first (`brew uninstall --force libmongocrypt`).

```
brew install --HEAD libmongocrypt
```

Verify installation succeeded by running:
```
pkg-config --modversion libmongocrypt
# Expect to see "1.7.0-pre" printed
```

May require explicitly forcing go to read the go.mod; do this on the command line:
```
export GO111MODULE=auto
```

Start a MongoDB server locally on port 27017.

Run the example with:

```
./run.sh
```

Expect to see output resembling the following:
```
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
Get without decryption ... begin
Document is stored in database as: {"_id": {"$oid":"63da9953d3231298f54b1efa"},"encryptedInt": {"$binary":{"base64":"CSqVtMt9kEoZsJJNqSF9sOYQw6EhWDRHht6M3nzXnEoV2TDp1L5CiaU3/ORoCjV7OLIKO8SCBJUr8u5WWZZsD4QpI1lsgQreJemiJXceZdMqdSklWXm1OUMmaDPaToI0qCANJJ/tIeY5+tpPyp9G/3QMuny9ZYt6qPJt+ThenHuFrBMlZZhKBBPjbSqdz+Iyn9xP0mURVzFpG3SEK+enSZZ4P9MMOygbwstLeJkZAb+WhQYXfegd02mu/lCW1Hkf9pZpqf7XKINN8/sAwv2IlGXTi2SSoWKaqAOnLIF9KWryhRqNZJQqu1pX/0RRQVs2yarl1FOFnY0KKNlkNQXs5RAsrMiqHhr4d4RfKw9SiOUQvU0qv52F1sm0eceC9ifD9NpXcawDr8XufhpDvTI1tPJiufNXCXHYEnFahncWfO/6WRNIeMCqiekcvvqHeytAXBO4pEHSHABmFTx2aavPIitdOx3n14I4DDYv6SWAweLBkdLh0Jix7jZ1CkKIIn/nfuQ/WC/nqE+Zn7UO1Cbq4to9Kl1m3zYvYxdauLmCOhnHHfJrQCShg7NkAkMhosSQjyrG6oHWkbFd9o+/URoiZJWTGS/VDZ618X9JXRGLZSxe/MeEdhGPjy1jh6C2+f+vnU6OfsZwumBCMZQmq2PJFxXWRlx+wz2Vslg7Tbaxy+5DO/1tdgXxGWqe/PZzLeTg93MSKhamEMIg2Fi4Kqzh3/IoUlIo9QLSvYwoPT0l8mH+nJmcZd3Vyij7hjgds9rZ19wWS3q3WzLeShh8EAdNEQOy82V2SbW1PLP+l/NU/SXzoykO9Gl1Y+4t1ach3cg7Agwni4GNHPVNb1CXutZrqLjDzy1n4rbvyhCxGbBBluGi+dRGO0sUjnBUmmg5g0Wn2CEx+h9CwAna9nLiOE6d5RAeoUp6a1Fo2Yau59DBNvK4DMpSg6cY4VNCh3JaYYgEgqOyXQl+ROW3OakR28BeGXz4HFfg9tUZ6HDAGayhN6VxnNgupzjKg9cBPE9d0xQ9gxIg9iiyjJjDsfXiZdxgaqL8QooP5yAQpbGh09FgO+ejpLgAaE6WXBXrcOVXAEMlolje5Q8R/3EjBYVRNoD73iJnIK1lgzwnFl3elrkhT8xPhl7YZj6I05bJD0p+NQKp/jEzNvLabJIBnlFa34cC8PUDyUSPSger3Ua63B+TCCuur7tiIGMnJmpmGBl75sY3d1AxwnCMXh1F/qzv7VZeALo/z41jPOr9Q5/ZtlJlu3cdAMC4R86nANWY1Oa6yaLHduDuDSFXFOiLznvYmvOqP/Gyip4wrt2Lf7b4MHLjbfUZoM/DzaUrrmtz2oUezhh2crZ/TyNYdskphc+0iLbaeyWQ7zXN5++mN49aiFDjbdBrDzH7tXcH8INx","subType":"06"}},"__safeContent__": [{"$binary":{"base64":"+Vs+IbYTSVSIw2NN1R8PDlhE+u4Dr53698qiPUe9SFY=","subType":"00"}},{"$binary":{"base64":"MR+HISXvms8hf+WARl97XYfy41ETOaYJonGTrM+uMvY=","subType":"00"}},{"$binary":{"base64":"haxLQ4QPRHGo13Ri9/vYDFQbmzjIXErZUQOWJ/wFE3I=","subType":"00"}},{"$binary":{"base64":"TaJdPrvNqAOI07pHvDARFIuoKzX3pxprc/pn5Icqpzg=","subType":"00"}},{"$binary":{"base64":"lNKj0QT0GbRDfXef8+MxZN6fyqUs8wQTFdvGjjI9/9g=","subType":"00"}},{"$binary":{"base64":"YWV+Ra1AqvJ8ZFe5D0YLyT3Rpw2i++yrGku21Ur5wqU=","subType":"00"}},{"$binary":{"base64":"HJ5H4uyodGq0XKhD1HCTzjo3ehNW1kyii9ISKyW0r5c=","subType":"00"}},{"$binary":{"base64":"LRLQW9xtcrag+jEbTouUaL88hX6hhev0lFtZrXKoILQ=","subType":"00"}},{"$binary":{"base64":"Rt1FT0ZJ/XMFgFT2m9wDkcgRh0JVxQDhQKx7RdagqW0=","subType":"00"}}]}
Get without decryption ... end
```

Note, this uses an unstable release of libmongocrypt. Consider uninstalling after with `brew uninstall libmongocrypt`.
