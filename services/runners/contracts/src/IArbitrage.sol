pragma solidity ^0.8.3;
pragma experimental ABIEncoderV2;

interface IArbitrage {
    function startArbitrage(address router1,
        address router2,
        address token0,
        address token1,
        address token2,
        uint amount0,
        uint amount1,
        uint amount2
    ) external payable;
}