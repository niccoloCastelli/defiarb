const Arbitrage = artifacts.require("Arbitrage.sol");

module.exports = function (deployer) {
  deployer.deploy(
    Arbitrage,
    '0x05fF2B0DB69458A0750badebc4f9e13aDd608C7F', //PancakeSwap router
    "0xcde540d7eafe93ac5fe6233bee57e1270d3e330f",
    "0xb7e19a1188776f32e8c2b790d9ca578f2896da7c",
    "0xC0788A3aD43d79aa53B09c2EaCc313A787d1d607",
    "0xBCfCcbde45cE874adCB698cC183deBcF17952812" //PancakeSwap factory
  );
};
