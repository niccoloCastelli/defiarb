pragma solidity ^0.6.6;
pragma experimental ABIEncoderV2;

import './PancakeLibrary.sol';

import './interfaces/IERC20.sol';
import './interfaces/IPancakeRouter.sol';
import './interfaces/IValueLiquidRouter.sol';
import './interfaces/BakerySwap.sol';


interface IArbitrage {
    function startArbitrage(
        address router1,
        address router2,
        address token0,
        address token1,
        address token2,
        uint amount0,
        uint amount1,
        uint amount2
    ) external payable;
}

// token0 -> token1 (token0 == 0)
// token1 -> token2 (router1)
// token2 -> token0 (router2)


// token1 -> token0 (token1 == 0)
// token0 -> token2 (router1)
// token2 -> token1 (router2)

// Router1, Router2, token0, token1, token2
contract Arbitrage is IArbitrage {
    uint constant deadline = 1000000 days;

    IPancakeRouter02 public pancakeRouter;
    IPancakeRouter02 public bakeryRouter;
    IPancakeRouter02 public apeRouter;
    IValueLiquidRouter public valueRouter;

    address pancakeFactory;

    struct RouteArgs {
        address router1;
        address router2;
        address token2;
        uint amount2;
    }

    struct ArbInfo {
        address[] path;
        address[] path1;
        address[] path2;

        address token0;
        address token1;
        address token2;

        uint amount2;

        address router1;
        address router2;
    }

    event Called(address indexed from, bytes _data);

    constructor(address _pancakeRouter, address _bakeryRouter, address _valueRouter, address _apeRouter, address _pancakeFactory) public payable  {
        pancakeRouter = IPancakeRouter02(_pancakeRouter);
        bakeryRouter = IPancakeRouter02(_bakeryRouter);
        valueRouter = IValueLiquidRouter(_valueRouter);
        apeRouter = IPancakeRouter02(_apeRouter);
        //Arbitrage.deployed().then((instance) => {return instance.startArbitrage("0xC0788A3aD43d79aa53B09c2EaCc313A787d1d607","0xC0788A3aD43d79aa53B09c2EaCc313A787d1d607","0xC0788A3aD43d79aa53B09c2EaCc313A787d1d607","0x603c7f932ed1fc6575303d8fb018fdcbb0f39a95","0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c","0x947950BcC74888a40Ffa2593C5798F11Fc9124C4","0","5000000000000000000","997376987535316000000")})

        pancakeFactory = _pancakeFactory;
    }

    receive() external payable {}

    //TODO: only owner
    function startArbitrage(
        address router1,
        address router2,
        address token0,
        address token1,
        address token2,
        uint amount0,
        uint amount1,
        uint amount2
    ) override external payable {

        address pairAddress = IPancakeFactory(pancakeFactory).getPair(token0, token1);
        require(pairAddress != address(0), 'This pool does not exist');

        IPancakePair(pairAddress).swap(
            IPancakePair(pairAddress).token0() == token0 ? 0 : amount1,
            IPancakePair(pairAddress).token0() == token0 ? amount1 : 0,
            address(this),
            abi.encode(RouteArgs(router1, router2, token2, amount2))
        );
    }

    function swapTokens(address router, address[] memory path, uint amountToken) private returns (uint) {
        uint amountReceived = 0;
        if (router == address(pancakeRouter) || router == address(apeRouter)) {

            address factoryAddr = IPancakeRouter02(router).factory();
            address pairAddress = IPancakeFactory(factoryAddr).getPair(path[0], path[1]);
            require(pairAddress != address(0), 'This pool does not exist');

            address[] memory sortedPath = new address[](2);
            sortedPath[0] = IPancakePair(pairAddress).token0();
            sortedPath[1] = IPancakePair(pairAddress).token1();

            if (sortedPath[0] == path[0]) {
                uint amountOut = IPancakeRouter02(router).getAmountsOut(amountToken, sortedPath)[1];
                //IERC20(sortedPath[0]).approve(pairAddress, 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF);
                IERC20(sortedPath[0]).transfer(pairAddress, amountToken);
                IPancakePair(pairAddress).swap(0, amountOut, address(this), "");
                amountReceived = amountOut;

            } else {
                address[] memory swapPath = new address[](2);
                swapPath[0] = sortedPath[1];
                swapPath[1] = sortedPath[0];

                uint amountOut = IPancakeRouter02(router).getAmountsOut(amountToken, swapPath)[1];
                //IERC20(sortedPath[1]).approve(pairAddress, 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF);
                IERC20(sortedPath[1]).transfer(pairAddress, amountToken);
                IPancakePair(pairAddress).swap(amountOut, 0, address(this), "");
                amountReceived = amountOut;
            }
        } else if (router == address(bakeryRouter)) {
            address factoryAddr = IPancakeRouter02(router).factory();
            address pairAddress = IPancakeFactory(factoryAddr).getPair(path[0], path[1]);
            require(pairAddress != address(0), 'This pool does not exist');

            address[] memory sortedPath = new address[](2);
            sortedPath[0] = IBakerySwapPair(pairAddress).token0();
            sortedPath[1] = IBakerySwapPair(pairAddress).token1();

            (uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast) = IBakerySwapPair(pairAddress).getReserves();


            if (sortedPath[0] == path[0]) {
                uint amountOut = IPancakeRouter02(router).getAmountOut(amountToken, reserve0, reserve1);
                IERC20(sortedPath[0]).approve(pairAddress, 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF);
                IERC20(sortedPath[0]).transfer(pairAddress, amountToken);
                IBakerySwapPair(pairAddress).swap(0, amountOut, address(this));
                amountReceived = amountOut;

            } else {
                address[] memory swapPath = new address[](2);
                swapPath[0] = sortedPath[1];
                swapPath[1] = sortedPath[0];

                uint amountOut = IPancakeRouter02(router).getAmountOut(amountToken, reserve1, reserve0);
                IERC20(sortedPath[1]).approve(pairAddress, 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF);
                IERC20(sortedPath[1]).transfer(pairAddress, amountToken);
                IBakerySwapPair(pairAddress).swap(amountOut, 0, address(this));
                amountReceived = amountOut;
            }

        } else if (router == address(valueRouter)) {
            address pairAddress = IValueLiquidFactory(valueRouter.factory()).getPair(path[0], path[1], 50, 3);
            require(pairAddress != address(0), 'Value: pair not found');
            address formula = valueRouter.formula();

            address[] memory sortedPath = new address[](2);
            sortedPath[0] = IValueLiquidPair(pairAddress).token0();
            sortedPath[1] = IValueLiquidPair(pairAddress).token1();

            (uint112 reserve0,uint112 reserve1, uint32 _) = IValueLiquidPair(pairAddress).getReserves();

            /*
            uint256 amountIn,
            uint256 reserveIn,
            uint256 reserveOut,
            uint32 tokenWeightIn,
            uint32 tokenWeightOut,
            uint32 swapFee
            */


            if (sortedPath[0] == path[0]) {
                uint amountOut = IValueLiquidFormula(formula).getPairAmountOut(pairAddress, sortedPath[0], amountToken);
                IERC20(sortedPath[0]).approve(pairAddress, 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF);
                IERC20(sortedPath[0]).transfer(pairAddress, amountToken);
                IValueLiquidPair(pairAddress).swap(0, amountOut, address(this), "");
                amountReceived = amountOut;
            } else {
                uint amountOut = IValueLiquidFormula(formula).getPairAmountOut(pairAddress, sortedPath[1], amountToken);
                IERC20(sortedPath[1]).approve(pairAddress, 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF);
                IERC20(sortedPath[1]).transfer(pairAddress, amountToken);
                IValueLiquidPair(pairAddress).swap(amountOut, 0, address(this), "");
                amountReceived = amountOut;
            }
        } else {
            revert('Unknown router');
        }
        require(amountReceived > 0, 'Amount received 0');
        return amountReceived;
    }

    function pancakeCall(address _sender, uint _amount0, uint _amount1, bytes calldata _data) external {

        ArbInfo memory arbInfo = ArbInfo({
        path : new address[](2),
        path1 : new address[](2),
        path2 : new address[](2),
        token0 : IPancakePair(msg.sender).token0(),
        token1 : IPancakePair(msg.sender).token1(),
        token2 : address(0),
        amount2 : 0,
        router1 : address(0),
        router2 : address(0)
        });

        (arbInfo.router1, arbInfo.router2, arbInfo.token2, arbInfo.amount2) = abi.decode(_data, (address, address, address, uint));

        require(
            msg.sender == PancakeLibrary.pairFor(pancakeFactory, arbInfo.token0, arbInfo.token1),
            'Unauthorized'
        );
        require(_amount0 == 0 || _amount1 == 0, 'Amount zero');

        uint amountToken = _amount0 == 0 ? _amount1 : _amount0;
        uint amountReceived = 0;
        uint amountReceivedFinal = 0;

        arbInfo.path[0] = _amount0 == 0 ? arbInfo.token1 : arbInfo.token0;
        arbInfo.path[1] = _amount0 == 0 ? arbInfo.token0 : arbInfo.token1;


        IERC20 otherToken;

        if (_amount0 == 0) {
            _amount1 = IERC20(arbInfo.token1).balanceOf(address(this));
            // token1 -> token2 -> token0
            // paths: [fl 0->1] [1-> 2] [2 -> 0]
            if (arbInfo.token2 == address(0)) {
                // Swap 1
                arbInfo.path1[0] = arbInfo.token1;
                arbInfo.path1[1] = arbInfo.token0;
                // Selezione router e scambio (amount1 -> amount2)
                amountReceivedFinal = swapTokens(arbInfo.router1, arbInfo.path1, _amount1);
            } else {
                // Swap 1
                arbInfo.path1[0] = arbInfo.token1;
                arbInfo.path1[1] = arbInfo.token2;
                // Selezione router e scambio (amount1 -> amount2)
                amountReceived = swapTokens(arbInfo.router1, arbInfo.path1, _amount1);

                // Swap 2
                arbInfo.path2[0] = arbInfo.token2;
                arbInfo.path2[1] = arbInfo.token0;

                //uint _amount2 = IERC20(arbInfo.token2).balanceOf(address(this));
                // Selezione router e scambio (out_swap_1 -> amountRequired)
                amountReceivedFinal = swapTokens(arbInfo.router2, arbInfo.path2, amountReceived);
            }


            otherToken = IERC20(arbInfo.token0);
            arbInfo.path[0] = arbInfo.token0;
            arbInfo.path[1] = arbInfo.token1;
        } else {
            // token0 -> token2 -> token1
            // paths: [fl 1->0] [0-> 2] [2 -> 1]
            _amount0 = IERC20(arbInfo.token0).balanceOf(address(this));

            if (arbInfo.token2 == address(0)) {
                // Swap 1
                arbInfo.path1[0] = arbInfo.token0;
                arbInfo.path1[1] = arbInfo.token1;

                // Selezione router e scambio (amount1 -> amount2)
                amountReceivedFinal = swapTokens(arbInfo.router1, arbInfo.path1, _amount0);
            } else {
                // Swap 1
                arbInfo.path1[0] = arbInfo.token0;
                arbInfo.path1[1] = arbInfo.token2;

                // Selezione router e scambio (amount1 -> amount2)
                amountReceived = swapTokens(arbInfo.router1, arbInfo.path1, _amount0);

                // Swap 2
                arbInfo.path2[0] = arbInfo.token2;
                arbInfo.path2[1] = arbInfo.token1;

                //uint _amount2 = IERC20(arbInfo.token2).balanceOf(address(this));
                // Selezione router e scambio (out_swap_1 -> amountRequired)
                amountReceivedFinal = swapTokens(arbInfo.router2, arbInfo.path2, amountReceived);
            }
            otherToken = IERC20(arbInfo.token1);
            arbInfo.path[0] = arbInfo.token1;
            arbInfo.path[1] = arbInfo.token0;
        }

        // Importo da rimborsare
        uint amountRequired = PancakeLibrary.getAmountsIn(
            pancakeFactory,
            amountToken,
            arbInfo.path
        )[0];

        require(amountReceivedFinal > amountRequired, 'amount required exceeds amount received');


        // Restituzione flash loan
        otherToken.transfer(msg.sender, amountRequired);


        // Trasferimento profitto
        uint earning = otherToken.balanceOf(address(this));
        otherToken.transfer(tx.origin, earning);
    }

}
