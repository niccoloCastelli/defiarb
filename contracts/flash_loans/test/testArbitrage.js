const { accounts, contract } = require('@openzeppelin/test-environment');
const [ owner ] = accounts;

const { expect } = require('chai');

const ArbitrageContract = contract.fromArtifact('Arbitrage'); // Loads a compiled contract

describe('Arbitrage', function () {
    this.timeout(15000);

    beforeEach(async () => {
        this.arbContract = await ArbitrageContract.new({ from: owner });
    })

    it('call arbitrage', async () => {
        const env = {
            router1: '',
            router2: '',
            token0: '',
            token1: '',
            token2: '',
            amount0: '',
            amount1: '',
            amount2: '',
        }

        await this.arbContract.startArbitrage("0xBCfCcbde45cE874adCB698cC183deBcF17952812","0xBCfCcbde45cE874adCB698cC183deBcF17952812","0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c","0x55d398326f99059fF775485246999027B3197955","0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3","0","1000000000000000000000","997376987535316000000")
//        expect().to.equal(owner);
    });
});
