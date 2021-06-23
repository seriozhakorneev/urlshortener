package shortener

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/asaskevich/govalidator"
	"github.com/dgraph-io/badger"
	"golang.org/x/net/context"
)

type Server struct {
}

// Create method
func (s *Server) Create(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from Create: %s", message.Body)

	opt := badger.DefaultOptions("./data")
	opt.Logger = nil
	db, err := badger.Open(opt)
	defer db.Close()
	handle(err)

	// url validation
	if !govalidator.IsURL(message.Body) {
		return &Message{Body: "Not valid url."}, nil
	}
	// if original url already exist
	if isValInDB(db, message.Body) {
		return &Message{Body: "Already exist."}, nil
	}

	seq := randSeq()
	shortUrl := "ur.l/" + seq
	// while shortUrl is in db, make new
	for isKeyInDB(db, shortUrl) {
		seq = randSeq()
		shortUrl = "ur.l/" + seq
	}
	// insert new key: value
	insertKV(db, shortUrl, message.Body)

	return &Message{Body: shortUrl}, nil
}

// Get method
func (s *Server) Get(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from Get: %s", message.Body)

	opt := badger.DefaultOptions("./data")
	opt.Logger = nil
	db, err := badger.Open(opt)
	defer db.Close()

	handle(err)

	// if key not exist
	if !isKeyInDB(db, message.Body) {
		return &Message{Body: "Not found."}, nil
	}

	// get value by key
	origUrl := getValue(db, message.Body)

	return &Message{Body: origUrl}, nil
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// making unique sequence for url
func randSeq() string {
	var symbols = map[string]string{
		"letters":    "abcdefghijklmnopqrstuvwxyz",
		"capLetters": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"numbers":    "1234567890",
		"underscore": "_",
	}
	var keys = []string{"letters", "capLetters", "letters", "numbers", "underscore"}
	var seq string

	for i := 1; i <= 5; i++ {
		lenKeys := len(keys)
		keyIndex := rand.Intn(lenKeys)
		keyValue := keys[keyIndex]
		randIndex := rand.Intn(len(symbols[keyValue]))

		// add random symbol
		seq += string(symbols[keyValue][randIndex])

		// remove the element at index keyIndex from keys.
		keys[keyIndex] = keys[lenKeys-1]
		keys[lenKeys-1] = ""
		keys = keys[:lenKeys-1]
		// remove key from map if its not in keys
		_, found := Find(keys, keyValue)
		if !found {
			delete(symbols, keyValue)
		}
	}
	return seq
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

// insert key, value
func insertKV(db *badger.DB, k, v string) {
	key, val := []byte(k), []byte(v)

	txn := db.NewTransaction(true)
	defer txn.Discard()

	handle(txn.Set(key, val))
	handle(txn.Commit())
	fmt.Printf("Inserted key '%s' with value '%s' \n", key, val)
}

// get value by key
func getValue(db *badger.DB, k string) string {
	var valCopy []byte
	key := []byte(k)

	txn := db.NewTransaction(false)
	entry, err := txn.Get(key)
	handle(err)

	valCopy, err = entry.ValueCopy(nil)
	handle(err)

	fmt.Printf("Get key: '%s' with value: '%s' \n", entry.Key(), valCopy)
	return string(valCopy)
}

// check if key is in db, return bool
func isKeyInDB(db *badger.DB, urlKey string) bool {

	txn := db.NewTransaction(false)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false

	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.Key()
		if urlKey == string(k) {
			return true
		}
	}
	return false
}

// check if value is in db, return bool
func isValInDB(db *badger.DB, urlVal string) bool {
	var isIn bool

	txn := db.NewTransaction(false)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10

	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		err := item.Value(func(v []byte) error {
			if urlVal == string(v) {
				isIn = true
			}
			return nil
		})
		handle(err)
	}
	return isIn
}
