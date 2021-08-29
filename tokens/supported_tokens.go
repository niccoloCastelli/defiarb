package tokens

var supportedTokens []Erc20

// Default BSC tokens
func init() {
	supportedTokens = []Erc20{
		NewErc20("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c", "Wrapped BNB", "WBNB", "https://bscscan.com/token/images/binance_32.png"),

		NewErc20("0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c", "Binance-Peg BTC", "BTCB", "https://bscscan.com/token/images/btcb_32.png"),
		NewErc20("0x2170ed0880ac9a755fd29b2688956bd959f933f8", "Binance-Peg Ethereum", "ETH", "https://bscscan.com/token/images/ethereum_32.png"),
		NewErc20("0xe9e7cea3dedca5984780bafc599bd69add087d56", "Binance-Peg BUSD", "BUSD", "https://bscscan.com/token/images/busd_32.png"),
		NewErc20("0x55d398326f99059ff775485246999027b3197955", "Binance-Peg BUSD-T", "USDT", "https://bscscan.com/token/images/busdt_32.png"),
		NewErc20("0x1af3f329e8be154074d8769d1ffa4ee058b1dbc3", "Binance-Peg Dai Token", "DAI", "https://bscscan.com/token/images/dai_32.png"),
		NewErc20("0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d", "Binance-Peg USD Coin", "USDC", "https://bscscan.com/token/images/centre-usdc_28.png"),
		NewErc20("0x7083609fce4d1d8dc0c979aab8c869ea2c873402", "Binance-Peg Polkadot", "DOT", "https://bscscan.com/token/images/polkadot_32.png"),
		NewErc20("0xf8a0bf9cf54bb92f17374d9e9a321e6a111a51bd", "Binance-Peg ChainLink", "LINK", "https://bscscan.com/token/images/chainlink_32.png?v=2"),
		NewErc20("0x101d82428437127bf1608f699cd651e6abf9766e", "Binance-Peg Basic Attention Token", "BAT", "https://bscscan.com/token/images/bat_32.png"),
		NewErc20("0x947950BcC74888a40Ffa2593C5798F11Fc9124C4", "Binance-Peg SushiToken", "SUSHI", "https://bscscan.com/token/images/sushiswap_32.png"),
		NewErc20("0x4338665CBB7B2485A8855A139b75D5e34AB0DB94", "Binance-Peg Litecoin", "LTC", "https://bscscan.com/token/images/litecoin_32.png"),
		NewErc20("0x250632378e573c6be1ac2f97fcdf00515d0aa91b", "Binance Beacon ETH", "BETH", "https://bscscan.com/token/images/binance-beth_32.png"),

		NewErc20("0x0e09fabb73bd3ade0a17ecc321fd13a19e81ce82", "PancakeSwap Token", "CAKE", "https://bscscan.com/token/images/pancake_32.png?=v1"),
		NewErc20("0x4bd17003473389a42daf6a0a729f6fdb328bbbd7", "VAI Stablecoin", "VAI", "https://bscscan.com/token/images/venus-vai_32.png"),

		NewErc20("0xe0e514c71282b6f4e823703a39374cf58dc3ea4f", "BELT Token", "BELT", "https://bscscan.com/token/images/beltfinance_32.png"),
		NewErc20("0x603c7f932ed1fc6575303d8fb018fdcbb0f39a95", "ApeSwapFinance Banana", "BANANA", "https://bscscan.com/token/images/apeswap_32.png"),
		NewErc20("0xa184088a740c695e156f91f5cc086a06bb78b827", "AUTOv2", "AUTO", "https://bscscan.com/token/images/autofarm_32.png"),
		NewErc20("0xbcf39f0edda668c58371e519af37ca705f2bfcbd", "PolyCrowns", "pCWS", "https://bscscan.com/token/images/seascape_32.png"),
		NewErc20("0xCa3F508B8e4Dd382eE878A314789373D80A5190A", "beefy.finance", "BIFI", "https://bscscan.com/token/images/beefy_32.png?=v1"),
	}
}

func SupportedTokens() []Erc20 {
	cp := make([]Erc20, len(supportedTokens))
	copy(cp, supportedTokens)
	return cp
}
