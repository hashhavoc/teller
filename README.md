# Teller

Teller is a sophisticated Command Line Interface (CLI) tool designed to facilitate developers' interactions with STX contracts. It streamlines the process of exploring source code and leveraging various analytics tools, thereby significantly enhancing the developer experience within the STX ecosystem. With its intuitive design, Teller ensures that integrating STX contracts into projects is both straightforward and efficient.

## Installing

Homebrew is the recommended way to install Teller:

```sh
brew install hashhavoc/tap/teller
```

## Source

### Dependencies

You will need to install some packages on linux to handle the build process: This is needed for the clipboard functionality.

```sh
sudo apt-get install libx11-dev
```

### Building

To use Teller, follow these simple steps:

1. Clone the repository:

    ```sh
    git clone https://github.com/hashhavoc/teller.git
    ```

2. Navigate to the project directory:

    ```sh
    cd teller
    ```

3. Build the executable:

    ```sh
    go build -o teller cmd/teller/main.go
    ```

4. Run Teller:

    ```sh
    ./teller
    ```

## Command Line Interface

Upon running Teller, you'll encounter the following command line interface:

```sh
➜  teller
NAME:
   teller - interact with the stx blockchain

USAGE:
   teller [global options] command [command options]

VERSION:
   v0.0.1

COMMANDS:
   init           Creates a new configuration file
   contracts      Provides interactions with contracts
   token          Provides interactions with tokens
   wallet         Provides interactions with wallets
   dex            Provides interactions with multiple dex
   transactions   Provides interactions with transactions
   ordinals, ord  Provides interactions with ordinals
   help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Commands

Teller offers the following commands:

- **contracts**: Provides interactions with contracts.
- **token**: Provides interactions with tokens.
- **wallet**: Provides interactions with wallets.
- **dex**: Provides interactions with multiple decentralized exchanges.
- **transactions**: Provides interactions with transactions.
- **ordinals**: Provides interactions with ordinals on bitcoin.
- **help**: Shows a list of commands or help for one command.

## Support

If you encounter any issues or have suggestions for improvement, please feel free to open an issue on [GitHub](https://github.com/hashhavoc/teller/issues). Your feedback is highly appreciated!

## License

This project is licensed under the GNU General Public License v3.0 License - see the [LICENSE](LICENSE) file for details.
