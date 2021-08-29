
mkdir -p build
mkdir -p arbitrage
solc --overwrite --abi --bin ./src/IArbitrage.sol -o build
abigen --bin=./build/IArbitrage.bin --abi=./build/IArbitrage.abi --pkg=arbitrage --out=arbitrage/IArbitrage.go
