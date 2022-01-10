## Golang Library for StarkNet

<!-- :exclamation::exclamation::exclamation: Dr. Spacemn is not a cryptographer and this library has not been audited by Starkware Ltd. :exclamation::exclamation::exclamation: -->

### Caigo is predominately a transcription of the following libraries:
- https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature
- https://github.com/seanjameshan/starknet.js
- https://github.com/software-mansion/starknet.py




#### Test
```go
$ go test -v
=== RUN   TestPedersenHash
--- PASS: TestPedersenHash (0.02s)
=== RUN   TestInitCurveWithConstants
--- PASS: TestInitCurveWithConstants (0.01s)
=== RUN   TestDivMod
--- PASS: TestDivMod (0.00s)
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
=== RUN   TestMultAir
--- PASS: TestMultAir (0.00s)
=== RUN   TestGetY
--- PASS: TestGetY (0.00s)
=== RUN   TestVerifySignature
--- PASS: TestVerifySignature (0.01s)
=== RUN   TestUIVerifySignature
--- PASS: TestUIVerifySignature (0.02s)
PASS
ok      github.com/dontpanicdao/caigo   0.605s
```