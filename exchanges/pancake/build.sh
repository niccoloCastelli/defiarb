WD=$PWD
cd contracts || exit
mkdir -p build
mkdir -p factory
mkdir -p router
mkdir -p pair
solc --overwrite --abi --bin ./src/PancakeRouter.sol -o build
abigen --bin=./build/IPancakeFactory.bin --abi=./build/IPancakeFactory.abi --pkg=factory --out=./factory/IPancakeFactory.go
abigen --bin=./build/IPancakeRouter01.bin --abi=./build/IPancakeRouter01.abi --pkg=router --out=./router/IPancakeRouter01.go
abigen --bin=./build/IPancakePair.bin --abi=./build/IPancakePair.abi --pkg=pair --out=./pair/IPancakePair.go
#abigen --bin=./build/IPancakeRouter02.bin --abi=./build/IPancakeRouter02.abi --pkg=contracts --out=./IPancakeRouter02.go

#abigen --bin=./build/IERC20.bin --abi=./build/IERC20.abi --pkg=contracts --out=./IERC20.go
#abigen --bin=./build/IWETH.bin --abi=./build/IWETH.abi --pkg=contracts --out=./IWETH.go
cd $WD || exit