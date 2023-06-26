module deploy

go 1.18

replace github.com/smartcontractkit/caigo => ../../

require github.com/smartcontractkit/caigo v0.3.1-0.20220909184134-51c4e68080bd

require (
	github.com/google/go-querystring v1.1.0 // indirect
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
	golang.org/x/sys v0.0.0-20220909162455-aba9fc2a8ff2 // indirect
)
