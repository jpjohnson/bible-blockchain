# Bible Blockchain

Bible Blockchain is a simple implementation of a blockchain written in Go. It creates a blockchain of the King James Version of the Bible and validates it using proof-of-work.

## Installation

1. Clone the repository: 

## Usage

The program provides a simple menu-driven interface to create and load Bible blockchains.
The bible text being used is the King James Version from https://openbible.com/textfiles/, in bibles/kjv.txt.
It stores the blockchain in bible.bin. 
### Menu Options

*   1 - **Create Bible Blockchain**: Creates a new Bible blockchain from a file.
*   2 - **Load Bible Blockchain**: Loads an existing Bible blockchain from a file.
*   3 - **Exit**: Exits the program.

#### Load Bible
After loading the Bible blockchain, you are prompted to search for a verse
ex `Galatians 5 14` (note no `:` between chapter and verse)

output:
```
You searched for:  Galatians   5 : 14
Found BibleVerse:  book: Galatians
chapter: 5
verse: 14
text: For all the law is fulfilled in one word, [even] in this; Thou shalt love thy neighbour as thyself.
```

## Implementation

The program uses a [BibleBlockchain](/bible.go#L32) struct to represent the blockchain, which contains a `GenesisBlock` and a slice of [BibleBlock](/bible.go#L24)s. The `addBlock` method is used to add new blocks to the blockchain, and the `isValid` method checks the validity of the blockchain by verifying the hashes and previous hashes of its blocks.

The `searchBibleVerse` method is used to search for a specific Bible verse in the blockchain.
