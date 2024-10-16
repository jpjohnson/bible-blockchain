package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type BibleBlockData struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Text    string `json:"text"`
}

type BibleBlock struct {
	Data         BibleBlockData
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	Pow          int
}

type BibleBlockchain struct {
	GenesisBlock BibleBlock
	Chain        []BibleBlock
	Difficulty   int
}

// calculateHash calculates the hash of a BibleBlock by concatenating its data,
// previous hash, power, and timestamp, and then computing the SHA256 hash of the
// resulting string. It returns the hash as a hexadecimal string.
//
// Parameters:
// - b: a BibleBlock instance representing the block to calculate the hash for.
//
// Returns:
//   - string: the hexadecimal representation of the SHA256 hash of the block's data,
//     previous hash, power, and timestamp.
func (b BibleBlock) calculateHash() string {
	data, _ := json.Marshal(b.Data)
	bibleBlockData := b.PreviousHash + string(data) + strconv.Itoa(b.Pow) + b.Timestamp.String()
	bibleBlockHash := sha256.Sum256([]byte(bibleBlockData))
	return fmt.Sprintf("%x", bibleBlockHash)
}

// mineBlock mines a BibleBlock by incrementing its proof-of-work (pow) until its hash meets the specified difficulty.
//
// Parameters:
// - difficulty: the minimum number of leading zeros required in the block's hash.
//
// Returns:
// - none
func (b *BibleBlock) mineBlock(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Pow++
		b.Hash = b.calculateHash()
	}
}

// MakeBibleBlockchain creates a new Bible blockchain with the specified difficulty.
//
// Parameters:
// - difficulty: the minimum number of leading zeros required in the block's hash.
//
// Returns:
// - BibleBlockchain: a new Bible blockchain instance.
func MakeBibleBlockchain(difficulty int) BibleBlockchain {
	genesisBlock := BibleBlock{
		Hash:         "0",
		PreviousHash: "",
		Timestamp:    time.Now(),
		Pow:          0,
	}
	return BibleBlockchain{
		GenesisBlock: genesisBlock,
		Chain:        []BibleBlock{genesisBlock},
		Difficulty:   difficulty,
	}
}

// addBlock adds a new block to the Bible blockchain.
//
// Parameters:
// - from: the sender of the transaction.
// - to: the recipient of the transaction.
// - amount: the amount of the transaction.
//
// Returns:
// - none
func (b *BibleBlockchain) addBlock(book string, chapter int, verse int, text string) {
	blockData := BibleBlockData{
		Book:    book,
		Chapter: chapter,
		Verse:   verse,
		Text:    text,
	}
	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := BibleBlock{
		Data:         blockData,
		Hash:         "",
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}
	newBlock.mineBlock(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}

// String returns a string representation of the BibleBlock.
//
// Parameters:
// - none
// Returns:
// - string: a string representation of the BibleBlock.
func (b BibleBlock) String() string {
	return fmt.Sprintf("data: %v\nhash: %v\npreviousHash: %v\ntimestamp: %v\npow: %v\n", b.Data, b.Hash, b.PreviousHash, b.Timestamp, b.Pow)
}

// String returns a string representation of the BibleBlockData.
//
// Parameters:
// - none
// Returns:
// - string: a string representation of the BibleBlockData.
func (b BibleBlockData) String() string {
	return fmt.Sprintf("book: %v\nchapter: %v\nverse: %v\ntext: %v\n", b.Book, b.Chapter, b.Verse, b.Text)
}

// isValid checks if the Bible blockchain is valid by verifying the hashes and previous hashes of its blocks.
//
// Parameters:
// - b: a BibleBlockchain instance representing the blockchain to check.
//
// Returns:
// - bool: true if the blockchain is valid, false otherwise.
func (b BibleBlockchain) isValid() bool {
	for i := range b.Chain[1:] {
		previousBlock := b.Chain[i]
		currentBlock := b.Chain[i+1]
		if currentBlock.Hash != currentBlock.calculateHash() || currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}
	return true
}

// toFile saves the BibleBlockchain instance to a file named "bible.bin" using gob encoding.
//
// Parameters:
// - none
//
// Returns:
// - none
func (b BibleBlockchain) toFile() {
	f, err := os.Create("bible.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	err = encoder.Encode(b)
	if err != nil {
		log.Fatal(err)
	}

}

// fromFile loads a BibleBlockchain instance from a file.
//
// Parameters:
// - file: The path to the file to load from.
//
// Return: None.
func (b *BibleBlockchain) fromFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(b)
	if err != nil {
		log.Fatal(err)
	}
}

// splitBibleReference parses a Bible reference string into its components: book, chapter, verse, and text.
//
// The input string is expected to be in the format "Book Chapter:Verse Text", e.g., "2 Timothy 3:14 But".
// The function extracts and returns the book name, chapter number, verse number, and the text following the reference.
//
// Parameters:
// - reference: a string representing the Bible reference to be split.
//
// Returns:
// - string: the name of the book.
// - int: the chapter number.
// - int: the verse number.
// - string: the text associated with the verse.
func splitBibleReference(reference string) (string, int, int, string) {
	// 2 Timothy 3:14	But
	partsA := strings.Split(reference, ":")
	partsA1 := strings.Split(partsA[0], " ")
	book := strings.Join(partsA1[0:len(partsA1)-1], " ")

	chapterVerse := partsA1[len(partsA1)-1]
	chapter, _ := strconv.Atoi(chapterVerse)

	re := regexp.MustCompile(`\d+`)
	match := re.FindString(partsA[1])

	verse, _ := strconv.Atoi(match)

	re = regexp.MustCompile(`^\d+\s*`)
	match = re.ReplaceAllString(partsA[1], "")
	text := match
	return book, chapter, verse, text
}

// CreateBibleBlockchain creates a Bible blockchain from a file and sets the difficulty.
//
// Parameters:
// - readFile: The path to the file to read.
// - difficulty: The difficulty level for the blockchain.
//
// Return: None.
func CreateBibleBlockchain(readFile string, difficulty int) {
	file, err := os.ReadFile(readFile)
	if err != nil {
		log.Fatal(err)
	}

	//log.Println(string(file))
	log.Println("Creating Bible Blockchain...")
	// open file
	scanner := bufio.NewScanner(strings.NewReader(string(file)))
	index := 0

	// create Bible Blockchain and set difficulty
	BibleBlockchain := MakeBibleBlockchain(difficulty)
	var bibleVerse BibleBlockData

	// read file line by line
	// skip first 3 lines of file
	for scanner.Scan() {
		line := scanner.Text()

		index++
		if index < 3 {
			// skips first 2 line of file
			continue
		}
		//log.Println("line: ", line)

		// parse BibleVerse
		bibleVerse.Book, bibleVerse.Chapter, bibleVerse.Verse, bibleVerse.Text = splitBibleReference(line)
		//fmt.Println("bibleVerse", bibleVerse)

		// add BibleVerse to Bible Blockchain
		BibleBlockchain.addBlock(bibleVerse.Book, bibleVerse.Chapter, bibleVerse.Verse, bibleVerse.Text)
	}

	// for _, block := range BibleBlockchain.Chain {
	// 	fmt.Println(block.Data.String())
	// }

	// validate Bible Blockchain
	log.Println("Validating Bible Blockchain...")
	if !BibleBlockchain.isValid() {
		log.Fatal("Bible Blockchain is not valid")
	}

	fmt.Println("Saving Bible Blockchain...")
	BibleBlockchain.toFile()
	log.Println("Bible Blockchain created and saved to bible.bin")
}

// menu displays the main menu of the program.
//
// No parameters.
// No return values.
func menu() {
	fmt.Println("\n-------- Menu ------")
	fmt.Println("1. Create Bible Blockchain")
	fmt.Println("2. Load Bible Blockchain")
	fmt.Println("3. Exit")
	fmt.Println("---------------------")
}

// main is the entry point of the program.
//
// It continuously prompts the user to choose an option from the menu and performs the corresponding action.
// No parameters.
// No return values.
func main() {

	var input string

	for {
		menu()
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&input)
		switch input {
		case "1":
			fmt.Print("Enter difficulty level: ")
			var difficulty int
			fmt.Scanln(&difficulty)
			CreateBibleBlockchain("./bibles/kjv.txt", difficulty)
		case "2":
			fmt.Print("Enter blockchain file: ")
			var file string
			fmt.Scanln(&file)
			LoadBibleBlockchain(file)
		case "3":
			return
		}
	}
}

// LoadBibleBlockchain loads a Bible blockchain from a file.
//
// Parameters:
// - readFile: The path to the file to read.
//
// Return: None.
func LoadBibleBlockchain(readFile string) {

	// parse file into blockchain
	BibleBlockchain := BibleBlockchain{}
	BibleBlockchain.fromFile(readFile)

	// for _, block := range BibleBlockchain.Chain {
	// 	fmt.Println(block.Data.String())
	// }

	// // validate Bible Blockchain
	// log.Println("Validating Bible Blockchain...")
	// if !BibleBlockchain.isValid() {
	// 	log.Fatal("Bible Blockchain is not valid")
	// }

	fmt.Println("Bible Blockchain loaded from ", readFile)

	// sub menut to search for bible verse given book, chapter, and verse
	subMenu(BibleBlockchain)
}

// subMenu is a sub menu to search for bible verse given book, chapter, and verse.
//
// Parameters:
// - BibleBlockchain: a BibleBlockchain instance representing the blockchain to search in.
//
// Return: None.
func subMenu(BibleBlockchain BibleBlockchain) {

	var book string
	var chapter int
	var verse int
	fmt.Print("Enter book, chapter, and verse: ")
	fmt.Scanln(&book, &chapter, &verse)
	fmt.Println("You searched for: ", book, " ", chapter, ":", verse)

	// search for BibleVerse
	bibleVerse := BibleBlockchain.searchBibleVerse(book, chapter, verse)
	fmt.Println("Found BibleVerse: ", bibleVerse.Data.String())
}

// searchBibleVerse searches for a BibleVerse in the BibleBlockchain.
//
// Parameters:
// - book: The book of the BibleVerse.
// - chapter: The chapter of the BibleVerse.
// - verse: The verse of the BibleVerse.
//
// Return: The BibleVerse if found, nil otherwise.
func (b *BibleBlockchain) searchBibleVerse(book string, chapter int, verse int) BibleBlock {
	for _, block := range b.Chain {
		if block.Data.Book == book && block.Data.Chapter == chapter && block.Data.Verse == verse {
			return block
		}
	}
	return BibleBlock{}
}
