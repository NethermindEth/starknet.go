package rpc

var blockGoerli310370 = BlockWithTxHashes{
	BlockHeader: BlockHeader{
		BlockHash:        "0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72",
		ParentHash:       "0x1ce6fa8ef59dfa1ad8f7ce7c3a4e6752e2d8ae6274f8257345f680e6ae0b5b5",
		SequencerAddress: "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b",
		BlockNumber:      310370,
		NewRoot:          "0x5cd7a08312635206c0210b8c90e61ceac27cb09629064e12266fd79e4c05a3d",
		Timestamp:        1661450764,
	},
	Status: "ACCEPTED_ON_L1",
	BlockBodyWithTxHashes: BlockBodyWithTxHashes{
		Transactions: []TxnHash{
			"0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99",
			"0x28981b14353a28bc46758dff412ac544d16f2ffc8dde31867855592ea054ab1",
			"0x41176c650076712f1618a141fc1cf9a8c39f0d9548a3458f29cf363310a1e72",
			"0x43cd66f3ddbfbf681ab99bb57bf9d94c83d6e9b586bdbde78ab2deb0328ebd5",
			"0x7602cfebe4f3cb3ef4c8b8c6d7dda2efaf4a500723020066f5db50acd5095cd",
			"0x2612f3f870ee7e7617d4f9efdc41fa8fd571f9720b059b1aa14c1bf15d3a92a",
			"0x1a7810a6c68adf0621ed384d915409c936efa0c9d436683ea0cf7ea171719b",
			"0x26683aeef3e9d9bcc1f0d45a5f0b67d0aa1919726524b2a8dc59504dacfd1f4",
			"0x1d374aa073435cdde1ec1caf972f7c175fd23438bb220848e71720e00fd7474",
			"0xfc13eabaa2f38981e68bb010370cad7a7d0b65a59101ec816042adca0d6841",
			"0x672d007224128b99bcc145cd3dbd8930a944b6a5fff5c27e3b158a6ff701509",
			"0x24795cbca6d2eba941082cea3f686bc86ef27dd46fdf84b32f9ba25bbeddb28",
			"0x69281a4dd58c260a06b3266554c0cf1a4f19b79d8488efef2a1f003d67506ed",
			"0x62211cc3c94d612b580eb729410e52277f838f962d91af91fb2b0526704c04d",
			"0x5e4128b7680db32de4dff7bc57cb11c9f222752b1f875e84b29785b4c284e2a",
			"0xdb8ad2b7d008fd2ad7fba4315b193032dee85e17346c80276a2e08c7f09f80",
			"0x67b9541ca879abc29fa24a0fa070285d1899fc044159521c827f6b6aa09bbd6",
			"0x5d9c0ab1d4ed6e9376c8ab45ee02b25dd0adced12941aafe8ce37369d19d9c2",
			"0x4e52da53e23d92d9818908aeb104b007ea24d3cd4a5aa43144d2db1011e314f",
			"0x6cc05f5ab469a3675acb5885c274d5143dca75dd9835c582f59e85ab0642d39",
			"0x561ed983d1d9c37c964a96f80ccaf3de772e2b73106d6f49dd7c3f7ed8483d9",
		},
	},
}

var fullBlockGoerli310370 = BlockWithTxs{
	BlockHeader: BlockHeader{
		BlockHash:        "0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72",
		ParentHash:       "0x1ce6fa8ef59dfa1ad8f7ce7c3a4e6752e2d8ae6274f8257345f680e6ae0b5b5",
		SequencerAddress: "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b",
		BlockNumber:      310370,
		NewRoot:          "0x5cd7a08312635206c0210b8c90e61ceac27cb09629064e12266fd79e4c05a3d",
		Timestamp:        1661450764,
	},
	Status: "ACCEPTED_ON_L1",
	BlockBodyWithTxs: BlockBodyWithTxs{
		Transactions: []Txn{
			InvokeTxnV0{
				CommonTxnProperties{
					TransactionHash: "0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99",
					BroadcastedCommonTxnProperties: BroadcastedCommonTxnProperties{
						Type:    "INVOKE",
						MaxFee:  "0xde0b6b3a7640000",
						Version: "0x0",
						Signature: []string{
							"0x7bc0a22005a54ec6a005c1e89ab0201cbd0819621edd9fe4d5ef177a4ff33dd",
							"0x13089e5f38de4ea98e9275be7fadc915946be15c14a8fed7c55202818527bea",
						},
						Nonce: "0x0",
					},
				},
				InvokeV0(FunctionCall{
					ContractAddress:    "0x2e28403d7ee5e337b7d456327433f003aa875c29631906908900058c83d8cb6",
					EntryPointSelector: "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
					CallData: []string{
						"0x1",
						"0x33830ce413e4c096eef81b5e6ffa9b9f5d963f57b8cd63c9ae4c839c383c1a6",
						"0x2db698626ed7f60212e1ce6e99afb796b6b423d239c3f0ecef23e840685e866",
						"0x0",
						"0x2",
						"0x2",
						"0x61c6e7484657e5dc8b21677ffa33e4406c0600bba06d12cf1048fdaa55bdbc3",
						"0x6307b990",
						"0x2b81",
					},
				}),
			},
		},
	},
}

var fullBlockGoerli310843 = BlockWithTxs{
	BlockHeader: BlockHeader{
		BlockHash:        "0x424fba26a7760b63895abe0c366c2d254cb47090c6f9e91ba2b3fa0824d4fc9",
		ParentHash:       "0x30e34dedf00bb35a9076b2b0f50a5a74fd2501f62094b6e687277be6ef3d444",
		SequencerAddress: "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b",
		BlockNumber:      310843,
		NewRoot:          "0x32bd4ff21288c898d4d3b6a7aea4ebdb3f1c7089cd52bde98316b4ecb8a50be",
		Timestamp:        1661486036,
	},
	Status: "ACCEPTED_ON_L1",
	BlockBodyWithTxs: BlockBodyWithTxs{
		Transactions: []Txn{
			DeployTxn{
				ContractAddress:     "0x21c40b1377353924e185c9536469787dbe0cdb77b6877fa3a9946b795c71ec7",
				ContractAddressSalt: "0x4241e90ee6a33a1e2e70b088f7e4acfb3d6630964c1a85e96fa96acd56dcfcf",
				ConstructorCalldata: []string{
					"0x31ad196615d50956d98be085eb1774624106a6936c7c38696e730d2a6df735a",
					"0x736affc32af71f8d361c855b38ffef58ec151bd8361a3b160017b90ada1068e",
				},
				Type:            "DEPLOY",
				Version:         "0x0",
				ClassHash:       "0x1ca349f9721a2bf05012bb475b404313c497ca7d6d5f80c03e98ff31e9867f5",
				TransactionHash: "0x35bd2978d2061b3463498f83c09322ed6a82e4b2a188506525e272a7adcdf6a",
			},
		},
	},
}

//
