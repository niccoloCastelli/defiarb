# Defi arbitrage bot (BSC)

__Flash loan arbitrage bot for BSC__

_The bot is not profitable as-is_

## Installation
The project relies on a PostgreSQL database to keep trusted tokens and LPs lists

### From binary
1. Download the binary
2. Run `chmod +x defiarb`
3. Copy the file in a _$PATH_ directory (eg. `/usr/bin`)
4. [Customize configuration](#Config file)

### From source
1. Clone the repo

   `git clone ...`

2. Install dependencies

   `go get ./...`

3. [Customize configuration](#Config file)


### Config file

Copy `config.example.json` and edit it according to your needs:
- `NodeUrl`: bsc node URL
- `ContractAddress`: bsc node URL
- `Db`: PostgreSQL db config
  - `Host`: DB host (default localhost)
  - `Port`: DB port
  - `DbName`: DB name
  - `User`: Username
  - `Password`: Password

### Smart contract
The arbitrage smart contract is in `contracts/flash_loans`, deploy it with `truffle migrate` and add the resulting address in `config.json`.

## Usage

Before starting the arbitrage bot the database should be updated:
1. Update trusted tokens from trustwallet list with command `updateTokens`
2. Update LP tokens list for desired exchanges with command `scan [EXCHANGE_1 ... EXCHANGE_n]` (this command may take a long time to complete)
3. Run the arbitrage bot with `run` command

### Help
Run `defiarb --help` for command list and `defiarb [COMMAND] --help` for detailed command help

### Write example config file
`defiarb writeconfig [CONFIG_FILE.json]`

### Update token list
`defiarb updateTokens`

### Update LP tokens
`defiarb scan [EXCHANGE_1 ... EXCHANGE_n]`


### Run arbitrage bot
`defiarb run`


## Contributing

The project was meant as a learning project, so it is not regularly maintained but any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request