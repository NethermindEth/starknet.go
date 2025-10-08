package account_test

import (
	"context"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/internal/tests"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestVerify tests the account.Verify method.

func TestVerify(t *testing.T) {
	tests.RunTestOn(t, tests.MockEnv)
	t.Parallel()

	// setup mock account
	mockCtrl := gomock.NewController(t)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)
	mockRpcProvider.EXPECT().ChainID(context.Background()).Return(gomock.Any().String(), nil)
	// TODO: remove this once the braavos bug is fixed. Ref: https://github.com/NethermindEth/starknet.go/pull/691
	mockRpcProvider.EXPECT().
		ClassHashAt(context.Background(), gomock.Any(), gomock.Any()).
		Return(internalUtils.RANDOM_FELT, nil)

	ks := account.NewMemKeystore()
	accAddress := internalUtils.TestHexToFelt(
		t,
		"0x2d54b7dc47eafa80f8e451cf39e7601f51fef6f1bfe5cea44ff12fa563e5457",
	)
	acc, err := account.NewAccount(
		mockRpcProvider,
		accAddress,
		"0x3904dda2cdd58e15dd8667b51a49deec6ce9c53e17b28fffb28fe9ccfddda92",
		ks,
		2,
	)
	require.NoError(t, err)

	// test cases
	type TestSet struct {
		txnHash string
		r       string
		s       string
		result  bool
	}

	testSet := map[tests.TestEnv][]TestSet{
		tests.MockEnv: { // all values were taken from Sepolia Voyager
			{
				txnHash: "0x3211b659c22723984cc6082f48fc0a1f1ef5de635ae50a4faa6dffc2eb902b0",
				r:       "0x3fe1dd281704514cb0dfefba6219d7e51104dd1a16f41612e0ab48274c8fb60",
				s:       "0x15b86524cab1fda7b3cb01d523ca9ff93d35e0fa6f68b13ec618da90c90c60c",
				result:  true,
			},
			{
				txnHash: "0x6e32e1a84a37cfa5dec99b58a104da98a763fb11993ddc7b1d3b6a7b1bcfcd7",
				r:       "0x4d2578f8269eb9ccd4555adf616bedbccec2bc5e76c4bf7e16331f50c25bbfa",
				s:       "0x661ff093034030d7c127c8bc14921559b555cab9dbce76d1019983e4e92ae9b",
				result:  true,
			},
			{
				txnHash: "0x6a52e4846693ae216927264e2ad6ebfc20de8041fb07e35db65e183097789a1",
				r:       "0x3a027b8e495b0e7825f065ab74400a8b62118bd28ae01daa3b50731a403287a",
				s:       "0x16aaaea02a74920e303747036037e9c8a64b00b0389da2b7fb4b4a81c173ec6",
				result:  true,
			},
			{
				txnHash: "0x3a71bea5bb50f004cb3d46cd62e5d1cb8938261be20d55bd1be2612b99db211",
				r:       "0x40473d6aeefcd875d79c7cf429c7f81a4a3a697d591b914c16ec9fc86d6676f",
				s:       "0x6a1a0c6121e9fc1bdc4e593b73e02d67be307c1e05b9542634acb44565ac3f",
				result:  true,
			},
			{
				txnHash: "0x3a80c85fc9717c629559e5fdc472f748068b4de98ccab92f0464dc25c530bda",
				r:       "0x4faef6fe54191e1de71136f54eb42b310f6fc8bdaa05d57dd6c5df8a519d48f",
				s:       "0x32dfee9c2d08bc5cfe3344ab341398dfff120dad0a71e38b0a4754302daa1b0",
				result:  true,
			},
			{
				txnHash: "0x6c5a532788bf32900c77d3d3d70bc295e251be482245d16e057783e5c7abd",
				r:       "0x39205f7574047ee8291eeca37f50ec1c13f25d2f61d321cde4fd319fb6afe8",
				s:       "0x4ac4b429e8d15683ab94490172a76f9f961330adb7cec033aa593d6c87dcbbf",
				result:  true,
			},
			{
				txnHash: "0x2091bcfb7ed291d85bba2dffddc18168b5e213627b74122371850af2ba1fd2a",
				r:       "0xb65136913adbd84bc4b71e1ecd537692b6e844c4fc89c91ffe091cf6e40ef9",
				s:       "0x612fdfd8a1738cceee73225d2fda8e93147316e39b22b20a825ecc9d49729ff",
				result:  true,
			},
			{
				txnHash: "0x2580427fc234409ff342804523dc6b328d1afd8a3194bef10f9052c7e22710",
				r:       "0x7ba0a7e939a2c89f9b69c19fe200b3b11cbd21085ef53740ece48c341b969e8",
				s:       "0x5b5ebd0ea83fc16820b94983375ddb42091f7cac77cb7785b47e19dc3bb149d",
				result:  true,
			},
			{
				txnHash: "0x260badc37c9734af6e5ec3be5a2e33343357f36a4047d6e8050e878af36de60",
				r:       "0x2a91e8c00252dd641ccf844687864c7c1e7c7c3515862d9384f8dab81afe3cb",
				s:       "0x3adde07e9a7201a8222b3d52ec46ad1ff409888e37758e9e1750fe0260005fd",
				result:  true,
			},
			{
				txnHash: "0x22432bdbf1f5376e6f3df154f172c3ca340d1649ad3c497490732ee33d9eec3",
				r:       "0x7ff1b184a38757cbfa5d8a44460e4fa73211e8a7f3ea6be8465d93ef055baf6",
				s:       "0x67d58da2fc852e0c204e3eb01b3479f9b3c4cbd3fbec5de9eb32ef4295f058",
				result:  true,
			},
			{
				txnHash: "0x11fafa59f202971d2d937a65d7db0c999c11b429b9821ae0a4e547246c41ee7",
				r:       "0x37726818b9768a390c711bd697535e0277d3c32b535b79e143aeb9cacb6f17b",
				s:       "0x14306c2afc62bdbf9d66ca3015249bf1b1c42dcd2bfea45e1642e16961ea9e6",
				result:  true,
			},
			{
				txnHash: "0x42294e78570e837b75551eab761e35d301d03e0d0fb2bbac4973f185ba5023",
				r:       "0x172dee5e3aae15e4f47606d69f3c2bf36745b41c58589979bc5818714afe9a6",
				s:       "0x69b71ad993b11b2a3d664cdadf79a01fcfa18cdb94391a2857c5adf9174e75b",
				result:  true,
			},
			// false signatures (the three ones above, but with a different final number)
			{
				txnHash: "0x22432bdbf1f5376e6f3df154f172c3ca340d1649ad3c497490732ee33d9eec3",
				r:       "0x7ff1b184a38757cbfa5d8a44460e4fa73211e8a7f3ea6be8465d93ef055baf6",
				s:       "0x67d58da2fc852e0c204e3eb01b3479f9b3c4cbd3fbec5de9eb32ef4295f051",
				result:  false,
			},
			{
				txnHash: "0x11fafa59f202971d2d937a65d7db0c999c11b429b9821ae0a4e547246c41ee7",
				r:       "0x37726818b9768a390c711bd697535e0277d3c32b535b79e143aeb9cacb6f17b",
				s:       "0x14306c2afc62bdbf9d66ca3015249bf1b1c42dcd2bfea45e1642e16961ea9e1",
				result:  false,
			},
			{
				txnHash: "0x42294e78570e837b75551eab761e35d301d03e0d0fb2bbac4973f185ba5023",
				r:       "0x172dee5e3aae15e4f47606d69f3c2bf36745b41c58589979bc5818714afe9a6",
				s:       "0x69b71ad993b11b2a3d664cdadf79a01fcfa18cdb94391a2857c5adf9174e751",
				result:  false,
			},
		},
	}[tests.TEST_ENV]

	// tests
	for _, test := range testSet {
		t.Run(test.txnHash, func(t *testing.T) {
			t.Parallel()

			txnHash := internalUtils.TestHexToFelt(t, test.txnHash)
			r := internalUtils.TestHexToFelt(t, test.r)
			s := internalUtils.TestHexToFelt(t, test.s)

			result, err := acc.Verify(txnHash, []*felt.Felt{r, s})
			require.NoError(t, err)
			assert.Equal(t, test.result, result)
		})
	}
}
