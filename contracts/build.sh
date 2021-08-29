
mkdir -p bep20
solc --overwrite --abi --bin ./src/IBEP20.sol -o build
abigen --bin=./build/IBEP20.bin --abi=./build/IBEP20.abi --pkg=bep20 --out=bep20/IBEP20.go

rm -rf ./build
