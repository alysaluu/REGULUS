package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"encoding/json"
	"strconv"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username string
	Password []byte
	Salt     []byte
	Sk       userlib.PKEDecKey
	SignKey  userlib.DSSignKey
	// maps file name to a list containing the iid and symmetric key

	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}

type File struct {
	Owner       string
	ContentHead LLNode
}

type LLNode struct {
	Next    []byte
	Content []byte
	Num     int
}

type TreeNode struct {
	Filename string
	Username string
	Children map[string]uuid.UUID
}

type PairShare struct {
	UUID []uuid.UUID // [fileUUID, treeUUID]
	Key  []byte
}

type MacEncVals struct {
	Ciphertext []byte
	MacVal     []byte
}

// NOTE: The following methods have toy (insecure!) implementations.

func InitUser(username string, password string) (userdataptr *User, err error) {
	// Error check #2: An empty username is provided.
	if username == "" {
		return nil, errors.New("An empty username is provided.")
	}

	// Create salt with randomly generated bytes
	salt := userlib.RandomBytes(16)

	// create uuid from username and store (username, salt) in datastore
	userUUID, err := uuid.FromBytes(userlib.Hash([]byte(username))[:16])
	// error check #1: A user with the same username exists.
	_, ok := userlib.DatastoreGet(userUUID)
	if ok {
		return nil, errors.New("A user with the same username exists.")
	}
	userlib.DatastoreSet(userUUID, salt)

	//init keys (pk/sk and signkey/vefifykey) and store them in keystore
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, err = userlib.PKEKeyGen()
	if err != nil {
		return nil, errors.New("Something went wrong with generating public keys")
	}

	var signKey userlib.DSSignKey
	var verifyKey userlib.DSVerifyKey
	signKey, verifyKey, err = userlib.DSKeyGen()
	if err != nil {
		return nil, errors.New("Something went wrong with generating digital signature keys")
	}

	userlib.KeystoreSet(username+"pk", pk)
	userlib.KeystoreSet(username+"vk", verifyKey)

	// rootkey for user struct = hash(username + password)
	key := userlib.Hash([]byte(username + password))[:16]
	encKey, macKey, err := GenerateMacEncKeys(key, "encryption", "mac")
	if err != nil {
		return nil, err
	}

	//initialize user struct
	newUser := User{
		username,
		[]byte(password),
		[]byte(salt),
		sk,
		signKey,
	}

	// storing user struct into datastore
	// key = uuid(hash(password + salt))
	hash := userlib.Hash(append(append([]byte(username), []byte(password)...), salt...))[:16]
	dsUserKey, err := uuid.FromBytes(hash)
	if err != nil {
		return nil, err
	}

	// marshal user struct and store bytes into userBytes
	userBytes, err := json.Marshal(newUser)
	if err != nil {
		return nil, err
	}
	err = StoreMEVDataStore(dsUserKey, macKey, encKey, userBytes)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User

	// check if username is already in datastore and error if not
	userUUID, err := uuid.FromBytes(userlib.Hash([]byte(username))[:16])
	if err != nil {
		return nil, err
	}
	saltByte, ok := userlib.DatastoreGet(userUUID)
	if !ok {
		return nil, errors.New("A user with this username does not exist.")
	}

	//generate keys to decrypt and verify mac of user struct
	key := userlib.Hash([]byte(username + password))[:16] // rootkey for user struct
	encKey, macKey, err := GenerateMacEncKeys(key, "encryption", "mac")
	if err != nil {
		return nil, err
	}

	//check user credentials (hash = where UserMEV is stored)
	hash := userlib.Hash(append(append([]byte(username), []byte(password)...), saltByte...))[:16]
	dsUserKey, err := uuid.FromBytes(hash)
	if err != nil {
		return nil, err
	}

	rawDecrypt, err := GetMEVDataStore(dsUserKey, macKey, encKey)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawDecrypt, &userdata)
	if err != nil {
		return nil, errors.New("Error converting json.Marshal to User struct")
	}

	return &userdata, nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {
	// key is used to mac + enc orginalFile Pairshare and uuid to store it
	key := userlib.Hash([]byte(filename + ".txt" + userdata.Username))[:16]
	keyUUID, err := uuid.FromBytes(key)
	if err != nil {
		return errors.New("Error generating UUID")
	}
	// root key and uuid of pairshare mev struct = username + filename
	encKey, macKey, err := GenerateMacEncKeys(key, "encryption", "mac")
	if err != nil {
		return err
	}

	// if file is already in userspace, want to access the file struct and replace contents with new contents
	_, isPresent := userlib.DatastoreGet(keyUUID)
	if isPresent {
		// decrypt pairshare struct and return as plaintext bytes
		pairsharePT, err := GetMEVDataStore(keyUUID, macKey, encKey)
		if err != nil {
			return err
		}
		var pairshare PairShare
		err = json.Unmarshal(pairsharePT, &pairshare)
		if err != nil {
			return err
		}

		// get filekeys from the key info in pairshare
		fileEncKey, fileMacKey, err := GenerateMacEncKeys(pairshare.Key, "encryption", "mac")
		if err != nil {
			return err
		}
		// decrypt file struct and return as plaintext bytes
		filePT, err := GetMEVDataStore(pairshare.UUID[0], fileMacKey, fileEncKey)
		if err != nil {
			return err
		}
		var file File
		err = json.Unmarshal(filePT, &file)
		if err != nil {
			return errors.New("Error converting json.Marshal to File struct")
		}

		// create new content node and set the existing head/tail to new node
		newContent := LLNode{
			nil,
			content,
			0,
		}
		file.ContentHead = newContent

		// marshal and re-enc/mac file struct, then store fileMEV in datastore
		fileBytes, err := json.Marshal(file)
		if err != nil {
			return err
		}
		err = StoreMEVDataStore(pairshare.UUID[0], fileMacKey, fileEncKey, fileBytes)
		if err != nil {
			return err
		}

		return nil
	}

	// new random bytes generated for every file uuid
	fileRandomBytes := userlib.RandomBytes(16)
	fileUUID, err := uuid.FromBytes(fileRandomBytes)
	if err != nil {
		return err
	}

	fileKey := userlib.RandomBytes(16)
	fileEncKey, fileMacKey, err := GenerateMacEncKeys(fileKey, "encryption", "mac")
	if err != nil {
		return err
	}

	treeUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
	if err != nil {
		return err
	}

	fileInfo := PairShare{
		[]uuid.UUID{fileUUID, treeUUID},
		fileKey,
	}

	// marshal and store fileInfo at uuid(filename + username)
	psBytes, err := json.Marshal(fileInfo)
	if err != nil {
		return err
	}

	err = StoreMEVDataStore(keyUUID, macKey, encKey, psBytes)
	if err != nil {
		return err
	}
	// if file does not exist in user file namespace, create new file struct and store fileMEV
	contentNode := LLNode{
		nil,
		content,
		0,
	}

	// treeHash := userlib.Hash(append([]byte(filename + ".txt" + "sharing tree" + userdata.Username), userdata.Password...))[:16]
	// newTree.Data = [filename, username]
	newTree := TreeNode{
		filename,
		userdata.Username,
		make(map[string]uuid.UUID),
	}

	treeEK, treeMk, err := GenerateMacEncKeys(fileKey, "tree encryption", "tree mac")
	if err != nil {
		return err
	}

	treeMarshaled, err := json.Marshal(newTree)
	if err != nil {
		return err
	}

	err = StoreMEVDataStore(treeUUID, treeMk, treeEK, treeMarshaled)
	if err != nil {
		return err
	}

	newFile := File{
		userdata.Username,
		contentNode,
	}

	// convert file struct into bytes, then store fileMEV into datastore (key = fileUUID)
	fileBytes, err := json.Marshal(newFile)
	if err != nil {
		return errors.New("Error with json.Marshal")
	}

	err = StoreMEVDataStore(fileUUID, fileMacKey, fileEncKey, fileBytes)
	if err != nil {
		return err
	}

	return nil
}

func (userdata *User) AppendToFile(filename string, content []byte) (err error) {
	//generate UUID --hash(filename + username)--to access pairshare info in datastore
	psHash := userlib.Hash([]byte(filename + ".txt" + userdata.Username))[:16]
	psUUID, err := uuid.FromBytes(psHash)
	if err != nil {
		return err
	}
	_, isPresent := userlib.DatastoreGet(psUUID)
	if !isPresent {
		return errors.New(filename + " not in " + userdata.Username + "'s file namespace.")
	}

	//generate pairshare rootKey, then macKey & encKey
	encKey, macKey, err := GenerateMacEncKeys(psHash, "encryption", "mac")
	if err != nil {
		return err
	}

	// get uuid/key of file by decrypting pairshare
	var pairshare PairShare
	psMEVMarshaled, err := GetMEVDataStore(psUUID, macKey, encKey)
	if err != nil {
		return err
	}
	err = json.Unmarshal(psMEVMarshaled, &pairshare)
	if err != nil {
		return errors.New("error unmarshaling pairshare in append")
	}

	// generate file keys with pairshare info
	fileEK, fileMK, err := GenerateMacEncKeys(pairshare.Key, "encryption", "mac")
	if err != nil {
		return err
	}

	var file File
	fileMEVMarshaled, err := GetMEVDataStore(pairshare.UUID[0], fileMK, fileEK)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileMEVMarshaled, &file)
	if err != nil {
		return errors.New("Error converting json.Marshal to File struct.")
	}

	newBytes := userlib.RandomBytes(16)
	newUUID, err := uuid.FromBytes(newBytes)
	if err != nil {
		return err
	}

	newHead := LLNode{
		newBytes,
		content,
		file.ContentHead.Num + 1,
	}

	// store current head at new node.next
	contentEK, contentMK, err := GenerateMacEncKeys(pairshare.Key, strconv.Itoa(file.ContentHead.Num)+"encryption", strconv.Itoa(file.ContentHead.Num)+"mac")
	marshaledHead, err := json.Marshal(file.ContentHead)
	if err != nil {
		return err
	}
	err = StoreMEVDataStore(newUUID, contentMK, contentEK, marshaledHead)
	if err != nil {
		return err
	}

	// set cur head to new node
	file.ContentHead = newHead

	// re-marshal content and store back into datastore
	fileBytes, err := json.Marshal(file)
	if err != nil {
		return errors.New("Error with json.Marshal")
	}
	err = StoreMEVDataStore(pairshare.UUID[0], fileMK, fileEK, fileBytes)
	if err != nil {
		return err
	}

	return nil
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	//generate UUID --hash(filename + username)--to access pairshare info in datastore
	psHash := userlib.Hash([]byte(filename + ".txt" + userdata.Username))[:16]
	psUUID, err := uuid.FromBytes(psHash)
	if err != nil {
		return nil, errors.New("Error loading pairshare uuid")
	}
	_, isPresent := userlib.DatastoreGet(psUUID)
	if isPresent == false {
		return nil, errors.New("File not in user's file namespace.")
	}

	//generate pairshare rootKey, then macKey & encKey
	encKey, macKey, err := GenerateMacEncKeys(psHash, "encryption", "mac")
	if err != nil {
		return nil, errors.New("error generating enc/mac keys")
	}

	// get uuid/key of file by decrypting pairshare
	var pairshare PairShare
	psMEVMarshaled, err := GetMEVDataStore(psUUID, macKey, encKey)
	if err != nil {
		return nil, err
	}
	// fmt.Println("success getting pairshare in load")

	err = json.Unmarshal(psMEVMarshaled, &pairshare)
	if err != nil {
		return nil, errors.New("error converting marshal value to pairshare")
	}

	// generate file keys with pairshare info
	fileEK, fileMK, err := GenerateMacEncKeys(pairshare.Key, "encryption", "mac")
	if err != nil {
		return nil, errors.New("file mac/enc keys")
	}
	var file File
	fileMEVMarshaled, err := GetMEVDataStore(pairshare.UUID[0], fileMK, fileEK)
	if err != nil {
		return nil, err // something wrong with getting file
	}
	// fmt.Println("success in getting file mev")
	err = json.Unmarshal(fileMEVMarshaled, &file)
	if err != nil {
		return nil, errors.New("Error converting json.Marshal to File struct.")
	}

	node := file.ContentHead
	res := node.Content
	nextUUIDBytes := node.Next

	for nextUUIDBytes != nil {
		nextUUID, err := uuid.FromBytes(nextUUIDBytes)
		if err != nil {
			return nil, errors.New("Error getting next uuid")
		}
		// unmarshal and store node value into node
		contentEK, contentMK, err := GenerateMacEncKeys(pairshare.Key, strconv.Itoa(node.Num-1)+"encryption", strconv.Itoa(node.Num-1)+"mac")
		nodeMarshalMEV, err := GetMEVDataStore(nextUUID, contentMK, contentEK)
		if err != nil {
			return nil, errors.New("Error getting marshal node")
		}
		err = json.Unmarshal(nodeMarshalMEV, &node)
		if err != nil {
			return nil, errors.New("error converting marshal value to node")
		}

		res = append(node.Content, res...)
		nextUUIDBytes = node.Next
	}

	return res, nil
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (invitationPtr uuid.UUID, err error) {
	// check if recipientUser exists
	userUUID, err := uuid.FromBytes(userlib.Hash([]byte(recipientUsername))[:16])
	if err != nil {
		return uuid.Nil, err
	}
	_, ok := userlib.DatastoreGet(userUUID)
	if !ok {
		return uuid.Nil, errors.New("The recipient user does not exist.")
	}

	psHash := userlib.Hash([]byte(filename + ".txt" + userdata.Username))[:16]
	psUUID, err := uuid.FromBytes(psHash)
	if err != nil {
		return uuid.Nil, err
	}
	_, isPresent := userlib.DatastoreGet(psUUID)
	if !isPresent {
		return uuid.Nil, errors.New("File not in user's file namespace.")
	}

	//generate pairshare rootKey, then macKey & encKey
	encKey, macKey, err := GenerateMacEncKeys(psHash, "encryption", "mac")
	if err != nil {
		return uuid.Nil, err
	}

	// get uuid/key of file by decrypting pairshare
	var pairshare PairShare
	psMEVMarshaled, err := GetMEVDataStore(psUUID, macKey, encKey)
	if err != nil {
		return uuid.Nil, err
	}
	err = json.Unmarshal(psMEVMarshaled, &pairshare)
	if err != nil {
		return uuid.Nil, errors.New("error unmarshalling pairshare")
	}

	// dsKey := userlib.Hash([]byte(userdata.Username + recipientUsername + filename))[:16]
	// store copyPSMEV at hash(filename + sender username + recipient username) and mac/enc with symmetric key and new tree uuid
	// create a copy pairshare that can be sent to the recipient
	treeEK, treeMK, err := GenerateMacEncKeys(pairshare.Key, "tree encryption", "tree mac")
	if err != nil {
		return uuid.Nil, err
	}

	// treeUUID, err := uuid.FromBytes(pairshare.UUID[1][:])
	// if err != nil {
	// 	return uuid.Nil, err
	// }
	treeMarshaled, err := GetMEVDataStore(pairshare.UUID[1], treeMK, treeEK)
	if err != nil {
		return uuid.Nil, err
	}

	var parentTree TreeNode
	err = json.Unmarshal(treeMarshaled, &parentTree)
	if err != nil {
		return uuid.Nil, err
	}

	newTreeUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
	if err != nil {
		return uuid.Nil, err
	}

	parentTree.Children[recipientUsername] = newTreeUUID

	parentTreeBytes, err := json.Marshal(parentTree)
	if err != nil {
		return uuid.Nil, err
	}
	err = StoreMEVDataStore(pairshare.UUID[1], treeMK, treeEK, parentTreeBytes)
	if err != nil {
		return uuid.Nil, err
	}

	copyPS := PairShare{
		[]uuid.UUID{pairshare.UUID[0], newTreeUUID},
		pairshare.Key,
	}

	recipTree := TreeNode{
		"",
		recipientUsername,
		make(map[string]uuid.UUID),
	}

	recipTreeMarshaled, err := json.Marshal(recipTree)
	if err != nil {
		return uuid.Nil, err
	}

	err = StoreMEVDataStore(newTreeUUID, treeMK, treeEK, recipTreeMarshaled)
	if err != nil {
		return uuid.Nil, err
	}

	copyPSMarshaled, err := json.Marshal(copyPS)
	if err != nil {
		return uuid.Nil, err
	}
	symKeyBytes := userlib.Hash([]byte(filename + ".txt" + userdata.Username + recipientUsername))[:16] // might need to change this up later
	copyPSUUID, err := uuid.FromBytes(symKeyBytes)
	if err != nil {
		return uuid.Nil, err
	}
	symmKey, psMacKey, err := GenerateMacEncKeys(symKeyBytes, "symmetric enc key for sharing", "symmetric mac key for sharing")
	if err != nil {
		return uuid.Nil, err
	}
	err = StoreMEVDataStore(copyPSUUID, psMacKey, symmKey, copyPSMarshaled)
	if err != nil {
		return uuid.Nil, err
	}

	cpsMEVPairshare := PairShare{
		[]uuid.UUID{copyPSUUID, uuid.Nil},
		symKeyBytes,
	}

	invBytes := userlib.RandomBytes(16)
	invitationPtr, err = uuid.FromBytes(invBytes)
	if err != nil {
		return uuid.Nil, err
	}
	cpsMEVPSBytes, err := json.Marshal(cpsMEVPairshare)
	if err != nil {
		return uuid.Nil, err
	}
	recipKey, isPresent := userlib.KeystoreGet(recipientUsername + "pk")
	if !isPresent {
		return uuid.Nil, err
	}
	fmt.Printf("size of marshled cpsmevms: %d\n", len(cpsMEVPSBytes))
	encMEVPairshare, err := userlib.PKEEnc(recipKey, cpsMEVPSBytes)
	if err != nil {
		return uuid.Nil, err
	}

	signedEncMEVPairshare, err := userlib.DSSign(userdata.SignKey, encMEVPairshare)
	if err != nil {
		return uuid.Nil, err
	}

	invIDMEV := MacEncVals{
		encMEVPairshare,
		signedEncMEVPairshare,
	}

	invIDMEVmarshaled, err := json.Marshal(invIDMEV)
	if err != nil {
		return uuid.Nil, err
	}
	userlib.DatastoreSet(invitationPtr, invIDMEVmarshaled)
	return invitationPtr, nil
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	// check if filename already exists
	psHash := userlib.Hash([]byte(filename + ".txt" + userdata.Username))[:16]
	psUUID, err := uuid.FromBytes(psHash)
	if err != nil {
		return err
	}
	_, isPresent := userlib.DatastoreGet(psUUID)
	if isPresent {
		return errors.New(filename + " already in " + userdata.Username + "'s file namespace.")
	}

	invIDMEVmarshaled, isPresent := userlib.DatastoreGet(invitationPtr)
	if !isPresent {
		return errors.New("cannot find invitationPtr")
	}

	var copyPSMEVPSMEV MacEncVals //A; asymmetrically encrypted PSMEV
	err = json.Unmarshal(invIDMEVmarshaled, &copyPSMEVPSMEV)
	if err != nil {
		return err
	}
	//get sender's verification key from keystore then verify signature
	senderVerifyKey, isPresent := userlib.KeystoreGet(senderUsername + "vk")
	if !isPresent {
		return errors.New("cannot find verification key of " + senderUsername)
	}
	err = userlib.DSVerify(senderVerifyKey, copyPSMEVPSMEV.Ciphertext, copyPSMEVPSMEV.MacVal)
	if err != nil {
		return errors.New("Error verifying symmKey") //A; changed to error verfiying symmKey
	}
	// get copyPSMEVPS marshal by decrypting the ciphertext with personal secret key
	copyPSMEVPSMarshaled, err := userlib.PKEDec(userdata.Sk, copyPSMEVPSMEV.Ciphertext)
	if err != nil {
		return err
	}
	// unmarshal pairshare and store into copypsmevps
	var copyPSMEVPS PairShare //A; recipient's own copy of the pairshare struct, containing UUID to actual pairshare w/ address of file & fileKey
	err = json.Unmarshal(copyPSMEVPSMarshaled, &copyPSMEVPS)
	if err != nil {
		return err
	}

	// generating symmkey and mackey with the rootkey from copypsmevps
	symmKey, macKey, err := GenerateMacEncKeys(copyPSMEVPS.Key, "symmetric enc key for sharing", "symmetric mac key for sharing")
	if err != nil {
		return err
	}
	// get marshaled copyps as plaintext
	copyPSMarshaled, err := GetMEVDataStore(copyPSMEVPS.UUID[0], macKey, symmKey)
	if err != nil {
		return err
	}

	// unmarshal plaintext into a pairshare struct
	var copyPS PairShare
	err = json.Unmarshal(copyPSMarshaled, &copyPS)
	if err != nil {
		return err
	}

	// check if uuid is valid
	_, isPresent = userlib.DatastoreGet(copyPS.UUID[0])
	if !isPresent {
		return errors.New("user was possibly revoked")
	}

	// create personal copy of pairshare to store in personal file map
	personalPS := PairShare{
		[]uuid.UUID{copyPS.UUID[0], copyPS.UUID[1]},
		copyPS.Key,
	}

	//delete old hash(sender, recipient) entry from datastore
	// userlib.DatastoreDelete(invitationPtr)
	// userlib.DatastoreDelete(invitationPtr)

	//add personal hash(filename, username) entry into datastore with PairShare struct-> will this work once a user tries to loadfile?
	//rootKey for this entry is hash(username, password, filename) -> will this work once a user tries to loadfile?
	rootKey := userlib.Hash([]byte(filename + ".txt" + userdata.Username))[:16]
	encKey, macKey, err := GenerateMacEncKeys(rootKey, "encryption", "mac")
	if err != nil {
		return err
	}
	newUUID, err := uuid.FromBytes(rootKey)
	if err != nil {
		return err
	}

	psCopyMarshal, err := json.Marshal(personalPS)
	if err != nil {
		return err
	}
	err = StoreMEVDataStore(newUUID, macKey, encKey, psCopyMarshal)
	if err != nil {
		return err
	}

	treeEK, treeMk, err := GenerateMacEncKeys(personalPS.Key, "tree encryption", "tree mac")
	if err != nil {
		return err
	}

	treeMarshal, err := GetMEVDataStore(personalPS.UUID[1], treeMk, treeEK)
	if err != nil {
		return errors.New("tree mev cannot be retrieved")
	}

	var userTree TreeNode
	err = json.Unmarshal(treeMarshal, &userTree)
	if err != nil {
		return err
	}

	//change node attributes
	userTree.Filename = filename

	//remarshal and store in datastore
	treeBytes, err := json.Marshal(userTree)
	if err != nil {
		return errors.New("Error with json.Marshal")
	}

	treeEK, treeMK, err := GenerateMacEncKeys(personalPS.Key, "tree encryption", "tree mac")
	if err != nil {
		return err
	}
	err = StoreMEVDataStore(personalPS.UUID[1], treeMK, treeEK, treeBytes)
	if err != nil {
		return err
	}

	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) (err error) {
	//load file struct
	psHash := userlib.Hash([]byte(filename + ".txt" + userdata.Username))[:16]
	psUUID, err := uuid.FromBytes(psHash)
	if err != nil {
		return errors.New("Error loading pairshare uuid")
	}
	_, isPresent := userlib.DatastoreGet(psUUID)
	if isPresent == false {
		return errors.New("File not in user's file namespace.")
	}

	//generate pairshare rootKey, then macKey & encKey
	encKey, macKey, err := GenerateMacEncKeys(psHash, "encryption", "mac")
	if err != nil {
		return errors.New("error generating enc/mac keys")
	}

	// get uuid/key of file by decrypting pairshare
	var pairshare PairShare
	psMEVMarshaled, err := GetMEVDataStore(psUUID, macKey, encKey)
	if err != nil {
		return err
	}
	err = json.Unmarshal(psMEVMarshaled, &pairshare)
	if err != nil {
		return errors.New("error converting marshal value to pairshare")
	}

	// generate tree node mac/enc keys
	treeEK, treeMK, err := GenerateMacEncKeys(pairshare.Key, "tree encryption", "tree mac")
	if err != nil {
		return err
	}

	fileEK, fileMK, err := GenerateMacEncKeys(pairshare.Key, "encryption", "mac")
	if err != nil {
		return errors.New("file mac/enc keys")
	}
	var file File
	fileMEVMarshaled, err := GetMEVDataStore(pairshare.UUID[0], fileMK, fileEK)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileMEVMarshaled, &file)
	if err != nil {
		return errors.New("Error converting json.Marshal to File struct.")
	}

	if file.Owner != userdata.Username { //might not be necessary but doesn't hurt to have
		return errors.New("Not owner of file, cannot revoke other users")
	}

	// check if file is shared with recipientUser
	found := false
	queue := []uuid.UUID{}
	// treeUUID, err := uuid.FromBytes(pairshare.TreeID)
	// if err != nil {
	// 	return err
	// }
	queue = append(queue, pairshare.UUID[1])
	for len(queue) > 0 {
		curTreeUUID := queue[0]
		queue = queue[1:]

		treeMarshal, err := GetMEVDataStore(curTreeUUID, treeMK, treeEK)
		if err != nil {
			return err
		}

		var tree TreeNode
		err = json.Unmarshal(treeMarshal, &tree)
		if err != nil {
			return err
		}

		for name, treeUUID := range tree.Children {
			if name == recipientUsername {
				// fmt.Println(recipientUsername + " found")
				found = true
				break
			} else {
				queue = append(queue, treeUUID)
			}
		}
		if found {
			break
		}
	}
	fmt.Printf("%s found in tree: %t\n", recipientUsername, found)
	if !found {
		return errors.New(filename + " is not shared with " + recipientUsername)
	}

	// fmt.Println("REVOKE - file is shared with " + recipientUsername)
	/*
		move file:
		- do not need to change keys
		[Q:]how should we generate a new uuid for file if previously it was based on the root user's username + filename?
		- generate new uuid, in datastore set new uuid -> fileMEV, delete old datastore entry
	*/
	// get file ShareTree, start from root node and traverse until revoked user is reached (bfs)
	newFileUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
	if err != nil {
		return err
	}
	queue = []uuid.UUID{}

	queue = append(queue, pairshare.UUID[1])
	for len(queue) > 0 {
		curTreeUUID := queue[0]
		queue = queue[1:]

		treeMarshal, err := GetMEVDataStore(curTreeUUID, treeMK, treeEK)
		if err != nil {
			return err
		}
		var tree TreeNode
		err = json.Unmarshal(treeMarshal, &tree)
		if err != nil {
			return err
		}
		for name, nodeUUID := range tree.Children {
			// fmt.Println(name)
			if name == recipientUsername {
				// delete this node, then proceed to share new iid with remain children
				// do not add this node to queue
				delete(tree.Children, name)
			} else {
				/*
					- need to get pairshare of curUser.Data[1] and name, stored at Hash(sender + recipent)
					- decrypt and verify pairshare with symmetric key
				*/
				childTreeBytes, err := GetMEVDataStore(nodeUUID, treeMK, treeEK)
				if err != nil {
					return err
				}
				var childTree TreeNode
				err = json.Unmarshal(childTreeBytes, &childTree)
				if err != nil {
					return err
				}

				// fmt.Println(childTree.Filename + ".txt" + childTree.Username)
				userPSHash := userlib.Hash([]byte(childTree.Filename + ".txt" + childTree.Username))[:16]
				userEK, userMK, err := GenerateMacEncKeys(userPSHash, "encryption", "mac")
				if err != nil {
					return err
				}
				userPSUUID, err := uuid.FromBytes(userPSHash)
				if err != nil {
					return err
				}

				var userPS PairShare
				psMEVMarshaled, err := GetMEVDataStore(userPSUUID, userMK, userEK)
				if err != nil {
					return err
				}
				err = json.Unmarshal(psMEVMarshaled, &userPS)
				if err != nil {
					return errors.New("error unmarshaling pairshare in revoke")
				}

				userPS.UUID[0] = newFileUUID
				userPSMarshaled, err := json.Marshal(userPS)
				if err != nil {
					return err
				}
				err = StoreMEVDataStore(userPSUUID, userMK, userEK, userPSMarshaled)
				if err != nil {
					return err
				}

				queue = append(queue, nodeUUID)
			}
		}
	}

	fileBytes, err := json.Marshal(file)
	if err != nil {
		return errors.New("Error with json.Marshal")
	}
	err = StoreMEVDataStore(newFileUUID, fileMK, fileEK, fileBytes)
	if err != nil {
		return err
	}
	userlib.DatastoreDelete(pairshare.UUID[0])

	pairshare.UUID[0] = newFileUUID
	psMarshaled, err := json.Marshal(pairshare)
	if err != nil {
		return err
	}
	err = StoreMEVDataStore(psUUID, macKey, encKey, psMarshaled)
	if err != nil {
		return err
	}

	return nil
}

/*
Helper Function: MacAndEnc
Inputs: macKey, encKey, and value
Return: ciphertext and mac(ciphertext)
*/
func MacAndEnc(macKey []byte, encKey []byte, value []byte) (ciphertext []byte, mac []byte, err error) {
	ciphertext = userlib.SymEnc(encKey, userlib.RandomBytes(16), value)
	mac, err = userlib.HMACEval(macKey, ciphertext)
	if err != nil {
		return nil, nil, errors.New("Error generating MAC value")
	}
	return ciphertext, mac, nil
}

/*
Helper Function: StoreMEVDataStore
Inputs: uuid of datastore key, macKey, encKey, and value
Return: err
*/
func StoreMEVDataStore(cipherKey uuid.UUID, macKey []byte, encKey []byte, value []byte) (err error) {
	var ciphertext []byte
	ciphertext = userlib.SymEnc(encKey, userlib.RandomBytes(16), value)
	macVal, err := userlib.HMACEval(macKey, ciphertext)
	if err != nil {
		return errors.New("Error evaluating mac")
	}
	newMEV := MacEncVals{
		ciphertext,
		macVal,
	}
	newMEVMarshaled, err := json.Marshal(newMEV)
	if err != nil {
		return errors.New("Error with json.Marshal")
	}
	userlib.DatastoreSet(cipherKey, newMEVMarshaled)
	return nil
}

/*
Helper Function: GenerateMacEncKeys
Inputs: root key, encPurpose, macPurpose
Return: enc and mac key
*/
func GenerateMacEncKeys(rootKey []byte, encPurpose string, macPurpose string) (encKey []byte, macKey []byte, err error) {
	encKey, err = userlib.HashKDF(rootKey, []byte(encPurpose))
	if err != nil {
		return nil, nil, err
	}
	encKey = encKey[:16]

	macKey, err = userlib.HashKDF(rootKey, []byte(macPurpose))
	if err != nil {
		return nil, nil, err
	}
	macKey = macKey[:16]

	return encKey, macKey, nil
}

/*
Helper Function: GetMEVDataStore
Inputs: datastore key, macKey, encKey
Return: plaintext, err
*/
func GetMEVDataStore(DSkey uuid.UUID, macKey []byte, encKey []byte) (value []byte, err error) {
	var MEV MacEncVals
	MEVmarshaled, ok := userlib.DatastoreGet(DSkey)
	if !ok {
		return nil, errors.New("Error retrieving MEV.")
	}
	err = json.Unmarshal(MEVmarshaled, &MEV)
	if err != nil {
		return nil, errors.New("Error with json.Unmarshal(MEV)")
	}
	tempMacVal, err := userlib.HMACEval(macKey, MEV.Ciphertext)
	if err != nil {
		return nil, errors.New("Error calculating new MAC value.")
	}
	if !userlib.HMACEqual(MEV.MacVal, tempMacVal) {
		return nil, errors.New("The integrity has been compromised.")
	}
	// need to check if mac is the same, then decrypt with enckey
	plaintext := userlib.SymDec(encKey, MEV.Ciphertext)
	return plaintext, nil
}
