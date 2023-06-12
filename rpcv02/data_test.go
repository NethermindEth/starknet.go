package rpcv02

import (
	"github.com/smartcontractkit/caigo/types"
	ctypes "github.com/smartcontractkit/caigo/types"
)

var blockGoerli310370 = Block{
	BlockHeader: BlockHeader{
		BlockHash:        ctypes.StrToFelt("0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72"),
		ParentHash:       ctypes.StrToFelt("0x1ce6fa8ef59dfa1ad8f7ce7c3a4e6752e2d8ae6274f8257345f680e6ae0b5b5"),
		SequencerAddress: "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b",
		BlockNumber:      310370,
		NewRoot:          "0x5cd7a08312635206c0210b8c90e61ceac27cb09629064e12266fd79e4c05a3d",
		Timestamp:        1661450764,
	},
	Status: "ACCEPTED_ON_L1",
	Transactions: []Transaction{
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x28981b14353a28bc46758dff412ac544d16f2ffc8dde31867855592ea054ab1")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x41176c650076712f1618a141fc1cf9a8c39f0d9548a3458f29cf363310a1e72")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x43cd66f3ddbfbf681ab99bb57bf9d94c83d6e9b586bdbde78ab2deb0328ebd5")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x7602cfebe4f3cb3ef4c8b8c6d7dda2efaf4a500723020066f5db50acd5095cd")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x2612f3f870ee7e7617d4f9efdc41fa8fd571f9720b059b1aa14c1bf15d3a92a")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x1a7810a6c68adf0621ed384d915409c936efa0c9d436683ea0cf7ea171719b")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x26683aeef3e9d9bcc1f0d45a5f0b67d0aa1919726524b2a8dc59504dacfd1f4")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x1d374aa073435cdde1ec1caf972f7c175fd23438bb220848e71720e00fd7474")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0xfc13eabaa2f38981e68bb010370cad7a7d0b65a59101ec816042adca0d6841")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x672d007224128b99bcc145cd3dbd8930a944b6a5fff5c27e3b158a6ff701509")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x24795cbca6d2eba941082cea3f686bc86ef27dd46fdf84b32f9ba25bbeddb28")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x69281a4dd58c260a06b3266554c0cf1a4f19b79d8488efef2a1f003d67506ed")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x62211cc3c94d612b580eb729410e52277f838f962d91af91fb2b0526704c04d")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x5e4128b7680db32de4dff7bc57cb11c9f222752b1f875e84b29785b4c284e2a")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0xdb8ad2b7d008fd2ad7fba4315b193032dee85e17346c80276a2e08c7f09f80")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x67b9541ca879abc29fa24a0fa070285d1899fc044159521c827f6b6aa09bbd6")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x5d9c0ab1d4ed6e9376c8ab45ee02b25dd0adced12941aafe8ce37369d19d9c2")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x4e52da53e23d92d9818908aeb104b007ea24d3cd4a5aa43144d2db1011e314f")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x6cc05f5ab469a3675acb5885c274d5143dca75dd9835c582f59e85ab0642d39")},
		TransactionHash{TransactionHash: ctypes.StrToFelt("0x561ed983d1d9c37c964a96f80ccaf3de772e2b73106d6f49dd7c3f7ed8483d9")},
	},
}

var fullBlockGoerli310370 = Block{
	BlockHeader: BlockHeader{
		BlockHash:        ctypes.StrToFelt("0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72"),
		ParentHash:       ctypes.StrToFelt("0x1ce6fa8ef59dfa1ad8f7ce7c3a4e6752e2d8ae6274f8257345f680e6ae0b5b5"),
		SequencerAddress: "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b",
		BlockNumber:      310370,
		NewRoot:          "0x5cd7a08312635206c0210b8c90e61ceac27cb09629064e12266fd79e4c05a3d",
		Timestamp:        1661450764,
	},
	Status: "ACCEPTED_ON_L1",
	Transactions: []Transaction{
		InvokeTxnV0{
			CommonTransaction: CommonTransaction{
				TransactionHash: ctypes.StrToFelt("0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99"),
				Type:            "INVOKE",
				MaxFee:          "0xde0b6b3a7640000",
				Version:         "0x0",
				Signature: []string{
					"0x7bc0a22005a54ec6a005c1e89ab0201cbd0819621edd9fe4d5ef177a4ff33dd",
					"0x13089e5f38de4ea98e9275be7fadc915946be15c14a8fed7c55202818527bea",
				},
				Nonce: "0x0",
			},
			ContractAddress:    ctypes.StrToFelt("0x2e28403d7ee5e337b7d456327433f003aa875c29631906908900058c83d8cb6"),
			EntryPointSelector: "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
			Calldata: []string{
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
		},
	},
}

var receiptTxn310370_0 = InvokeTransactionReceipt(CommonTransactionReceipt{
	TransactionHash: ctypes.StrToFelt("0x40c82f79dd2bc1953fc9b347a3e7ab40fe218ed5740bf4e120f74e8a3c9ac99"),
	ActualFee:       "0x1709a2f3a2",
	Type:            "INVOKE",
	Status:          types.TransactionAcceptedOnL1,
	BlockHash:       ctypes.StrToFelt("0x6c2fe3db009a2e008c2d65fca14204f3405cb74742fcf685f02473acaf70c72"),
	BlockNumber:     310370,
	MessagesSent:    []MsgToL1{},
	Events: []Event{
		{
			FromAddress: ctypes.StrToFelt("0x37de00fb1416936b3074fc78bcc811d83046009b162c4a822ce84dabedd0ea9"),
			Data: []string{
				"0x0",
				"0x35b32bb4a1969175fb14b6c09838d1b3200724cc4d2b0891be319764021f5ac",
				"0xe9",
				"0x0",
			},
			Keys: []string{"0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9"},
		},
		{
			FromAddress: ctypes.StrToFelt("0x33830ce413e4c096eef81b5e6ffa9b9f5d963f57b8cd63c9ae4c839c383c1a6"),
			Data: []string{
				"0x61c6e7484657e5dc8b21677ffa33e4406c0600bba06d12cf1048fdaa55bdbc3",
				"0x2e28403d7ee5e337b7d456327433f003aa875c29631906908900058c83d8cb6",
			},
			Keys: []string{"0xf806f71b19e4744968b37e3fb288e61309ab33a782ea9d11e18f67a1fbb110"},
		},
	},
})

var fullBlockGoerli310843 = Block{
	BlockHeader: BlockHeader{
		BlockHash:        ctypes.StrToFelt("0x424fba26a7760b63895abe0c366c2d254cb47090c6f9e91ba2b3fa0824d4fc9"),
		ParentHash:       ctypes.StrToFelt("0x30e34dedf00bb35a9076b2b0f50a5a74fd2501f62094b6e687277be6ef3d444"),
		SequencerAddress: "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b",
		BlockNumber:      310843,
		NewRoot:          "0x32bd4ff21288c898d4d3b6a7aea4ebdb3f1c7089cd52bde98316b4ecb8a50be",
		Timestamp:        1661486036,
	},
	Status: "ACCEPTED_ON_L1",
	Transactions: []Transaction{
		DeployTxn{
			ConstructorCalldata: []string{
				"0x31ad196615d50956d98be085eb1774624106a6936c7c38696e730d2a6df735a",
				"0x736affc32af71f8d361c855b38ffef58ec151bd8361a3b160017b90ada1068e",
			},
			ContractAddressSalt: "0x4241e90ee6a33a1e2e70b088f7e4acfb3d6630964c1a85e96fa96acd56dcfcf",
			ClassHash:           "0x1ca349f9721a2bf05012bb475b404313c497ca7d6d5f80c03e98ff31e9867f5",
			Type:                "DEPLOY",
			Version:             "0x0",
			TransactionHash:     ctypes.StrToFelt("0x35bd2978d2061b3463498f83c09322ed6a82e4b2a188506525e272a7adcdf6a"),
		},
	},
}

var receiptTxn310843_14 = DeployTransactionReceipt{
	CommonTransactionReceipt: CommonTransactionReceipt{
		TransactionHash: types.StrToFelt("0x035bd2978d2061b3463498f83c09322ed6a82e4b2a188506525e272a7adcdf6a"),
		ActualFee:       "0x0",
		Status:          "ACCEPTED_ON_L1",
		BlockHash:       types.StrToFelt("0x0424fba26a7760b63895abe0c366c2d254cb47090c6f9e91ba2b3fa0824d4fc9"),
		BlockNumber:     310843,
		Type:            "DEPLOY",
		MessagesSent:    []MsgToL1{},
		Events:          []Event{},
	},
	ContractAddress: "0x21c40b1377353924e185c9536469787dbe0cdb77b6877fa3a9946b795c71ec7",
}

var fullBlockGoerli300114 = Block{
	BlockHeader: BlockHeader{
		BlockHash:        ctypes.StrToFelt("0x184268bfbce24766fa53b65c9c8b30b295e145e8281d543a015b46308e27fdf"),
		ParentHash:       ctypes.StrToFelt("0x7307cb0d7fa65c111e71cdfb6209bdc90d2454d4c0f34d8bf5a3fe477826c3c"),
		SequencerAddress: "0x46a89ae102987331d369645031b49c27738ed096f2789c24449966da4c6de6b",
		BlockNumber:      300114,
		NewRoot:          "0x239a44410e78665f41f7a65ef3b5ed244ce411965498a83f80f904e22df1045",
		Timestamp:        1660701246,
	},
	Status: "ACCEPTED_ON_L1",
	Transactions: []Transaction{
		DeclareTxn{
			CommonTransaction: CommonTransaction{
				TransactionHash: ctypes.StrToFelt("0x46a9f52a96b2d226407929e04cb02507e531f7c78b9196fc8c910351d8c33f3"),
				Type:            "DECLARE",
				MaxFee:          "0x0",
				Version:         "0x0",
				Signature:       []string{},
				Nonce:           "0x0",
			},
			ClassHash:     "0x6feb117d1c3032b0ae7bd3b50cd8ec4a78c621dca0d63ddc17890b78a6c3b49",
			SenderAddress: "0x1",
		},
	},
}

var receiptTxn300114_3 = DeclareTransactionReceipt(
	CommonTransactionReceipt{
		TransactionHash: ctypes.StrToFelt("0x46a9f52a96b2d226407929e04cb02507e531f7c78b9196fc8c910351d8c33f3"),
		ActualFee:       "0x0",
		Status:          types.TransactionAcceptedOnL1,
		BlockHash:       ctypes.StrToFelt("0x184268bfbce24766fa53b65c9c8b30b295e145e8281d543a015b46308e27fdf"),
		BlockNumber:     300114,
		Type:            "DECLARE",
		MessagesSent:    []MsgToL1{},
		Events:          []Event{},
	})

var InvokeTxnV00x705547f8f2f8f = InvokeTxnV0{
	CommonTransaction: CommonTransaction{
		TransactionHash: ctypes.StrToFelt("0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0"),
		Type:            TransactionType_Invoke,
		MaxFee:          "0x53685de02fa5",
		Version:         "0x0",
		Nonce:           "0x0",
		Signature: []string{
			"0x4a7849de7b91e52cd0cdaf4f40aa67f54a58e25a15c60e034d2be819c1ecda4",
			"0x227fcad2a0007348e64384649365e06d41287b1887999b406389ee73c1d8c4c",
		},
	},
	ContractAddress:    ctypes.StrToFelt("0x315e364b162653e5c7b23efd34f8da27ba9c069b68e3042b7d76ce1df890313"),
	EntryPointSelector: "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
	Calldata: []string{
		"0x1",
		"0x13befe6eda920ce4af05a50a67bd808d67eee6ba47bb0892bef2d630eaf1bba",
		"0x3451875d57805682e40d0ad8e604fc4cc5f949d14ca8228d4a2eaeee7f48688",
		"0x0",
		"0x19",
		"0x19",
		"0x6",
		"0x6274632f757364",
		"0x46c9e02299ca9c00000",
		"0x62aba221",
		"0x657175696c69627269756d2d636578",
		"0x6574682f757364",
		"0x3c7c6b5765ad980000",
		"0x62aba221",
		"0x657175696c69627269756d2d636578",
		"0x736f6c2f757364",
		"0x1b6009e149e038000",
		"0x62aba221",
		"0x657175696c69627269756d2d636578",
		"0x617661782f757364",
		"0xe167e2d85ad98000",
		"0x62aba221",
		"0x657175696c69627269756d2d636578",
		"0x646f67652f757364",
		"0xc94bf844f94000",
		"0x62aba222",
		"0x657175696c69627269756d2d636578",
		"0x736869622f757364",
		"0x764e9c31400",
		"0x62aba222",
		"0x657175696c69627269756d2d636578",
		"0x4c83",
	},
}

var InvokeTxnV0_300000_0 = InvokeTxnV0{
	CommonTransaction: CommonTransaction{
		TransactionHash: ctypes.StrToFelt("0x680b0616e65633dfaf06d5a5ee5f1d1d4b641396009f00a67c0d18dc0f9638"),
		Type:            TransactionType_Invoke,
		MaxFee:          "0x17e817dc96bf",
		Version:         "0x0",
		Signature: []string{
			"0x28fd7fdff06bb65438e10bb9af23d7daf354ec7c75173056c0bb5a55690bf42",
			"0x6eb86c1d6b8efcfd9cefe2d068e9f34874a71e0583abdb38513f92addfb36ca",
		},
		Nonce: "0x0",
	},
	ContractAddress:    ctypes.StrToFelt("0x661638e27a00f65819559a0612874d1e3865b2372c539a85b5f0fb47c8ec683"),
	EntryPointSelector: "0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad",
	Calldata: []string{
		"0x1",
		"0x12fadd18ec1a23a160cc46981400160fbf4a7a5eed156c4669e39807265bcd4",
		"0x1a3af362287b0d09060bb79b7eeb6da636a9cd720a9b805c081181c8bcc58f",
		"0x0",
		"0x10",
		"0x10",
		"0x3",
		"0x6274632f757364",
		"0x50de402c68b9dc00000",
		"0x62fc247c",
		"0x636f696e62617365",
		"0x656d7069726963",
		"0x6574682f757364",
		"0x65a2e8c490e4e00000",
		"0x62fc247c",
		"0x636f696e62617365",
		"0x656d7069726963",
		"0x6461692f757364",
		"0xde000cd866f8000",
		"0x62fc247c",
		"0x636f696e62617365",
		"0x656d7069726963",
		"0x62fc2512",
	},
}
