package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.

	_ "encoding/hex"
	_ "errors"
	"fmt"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"
const contentFour = "Sun stare, don't care with my head in my hands thinking of a simpler time Like Sun Ra, feel small, but I had it under control every time When I was a waitress wearing a white dress Look how I do this, look how I got this I was a waitress working the night shift You were my man, felt like I got this Down at the Men in Music Business Conference Down in Orlando, I was only 19 Down at the Men in Music Business Conference I only mention it 'cause it was such a scene And I felt seen Mmh, mmh-mmh-mmh-mmh, mhm Summer, sizzling Listening to jazz out on the lawn Listening to White Stripes when they were white-hot Listening to rock, all day long When I was a waitress wearing a tight dress, handling the heat I wasn't famous, just listening to Kings of Leon to the beat Like, look at how I got this Look how I got this, just singing in the street Down at the Men in Music Business Conference I felt free 'cause I was only 19 Such a scene Summer, summer's almost gone We were talking about life, we were sitting outside 'til dawn But I would still go back If I could do it all again, I'd fly Because it made me feel, made me feel like a god 'Cause it made me feel, made me feel like a god Somehow it made me feel, made me feel like a god"
const contentFive = "WHITE DRESS: Sun stare, don't care with my head in my hands thinking of a simpler time Like Sun Ra, feel small, but I had it under control every time When I was a waitress wearing a white dress Look how I do this, look how I got this I was a waitress working the night shift You were my man, felt like I got this Down at the Men in Music Business Conference Down in Orlando, I was only 19 Down at the Men in Music Business Conference I only mention it 'cause it was such a scene And I felt seen Mmh, mmh-mmh-mmh-mmh, mhm Summer, sizzling Listening to jazz out on the lawn Listening to White Stripes when they were white-hot Listening to rock, all day long When I was a waitress wearing a tight dress, handling the heat I wasn't famous, just listening to Kings of Leon to the beat Like, look at how I got this Look how I got this, just singing in the street Down at the Men in Music Business Conference I felt free 'cause I was only 19 Such a scene Summer, summer's almost gone We were talking about life, we were sitting outside 'til dawn But I would still go back If I could do it all again, I'd fly Because it made me feel, made me feel like a god 'Cause it made me feel, made me feel like a god Somehow it made me feel, made me feel like a god"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {
	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	dorisFile := "dorisFile.txt"
	//eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Efficiency Tests", func() {
		Specify("Testing effienceny with number of appends", func() {
			measureBandwidth := func(probe func()) (bandwidth int) {
				before := userlib.DatastoreGetBandwidth()
				probe()
				after := userlib.DatastoreGetBandwidth()
				return after - before
			}

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			fileData := contentOne
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			bw1stAppend := measureBandwidth(func() {
				err = alice.AppendToFile(aliceFile, []byte(contentOne))
				Expect(err).To(BeNil())
			})
			userlib.DebugMsg("1st Append to file bandwidth measure: %d", bw1stAppend)

			for i := 0; i < 1000; i++ {
				err = alice.AppendToFile(aliceFile, []byte(contentTwo))
				Expect(err).To(BeNil())
				fileData += contentTwo
			}

			bw1000thAppend := measureBandwidth(func() {
				err = alice.AppendToFile(aliceFile, []byte(contentOne))
				Expect(err).To(BeNil())
			})

			userlib.DebugMsg("~1000th Append to file bandwidth measure: %d", bw1000thAppend)

		})

		Specify("Testing effienceny with number of files", func() {
			//TODO: measure bandwidth of append on 1 file vs. on 4 files
			measureBandwidth := func(probe func()) (bandwidth int) {
				before := userlib.DatastoreGetBandwidth()
				probe()
				after := userlib.DatastoreGetBandwidth()
				return after - before
			}

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			bw1stAppend := measureBandwidth(func() {
				err = alice.AppendToFile(aliceFile, []byte(contentOne))
				Expect(err).To(BeNil())
			})
			userlib.DebugMsg("1st Append to file bandwidth measure: %d", bw1stAppend)

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(bobFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(charlesFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(dorisFile, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			err = alice.AppendToFile(bobFile, []byte(contentTwo))
			err = alice.AppendToFile(charlesFile, []byte(contentTwo))
			err = alice.AppendToFile(dorisFile, []byte(contentTwo))

			bw4thAppend := measureBandwidth(func() {
				err = alice.AppendToFile(aliceFile, []byte(contentOne))
				Expect(err).To(BeNil())
			})

			userlib.DebugMsg("Append to 4 files bandwidth measure: %d", bw4thAppend)

		})

		Specify("Testing effienceny with file size", func() {
			measureBandwidth := func(probe func()) (bandwidth int) {
				before := userlib.DatastoreGetBandwidth()
				probe()
				after := userlib.DatastoreGetBandwidth()
				return after - before
			}

			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			smallFileAppend := measureBandwidth(func() {
				err = alice.AppendToFile(aliceFile, []byte(contentOne))
				Expect(err).To(BeNil())
			})
			userlib.DebugMsg("1st Append to small file bandwidth measure: %d", smallFileAppend)

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(bobFile, []byte(contentFour))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte(contentOne))

			bigFileAppend := measureBandwidth(func() {
				err = alice.AppendToFile(bobFile, []byte(contentOne))
				Expect(err).To(BeNil())
			})

			userlib.DebugMsg("Append to big file bandwidth measure: %d", bigFileAppend)

		})
	})

	Describe("Client API Tests", func() {
		Specify("Init User: Empty/same username provided", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Testing for duplicate usernames.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Testing for empty username")
			alice, err = client.InitUser("", defaultPassword)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking for if user does not exist")
			_, ok := userlib.KeystoreGet("blahblahblah")
			Expect(ok).ToNot(BeTrue())

			userlib.DebugMsg("Checking for error if user credentials are invalid")
			aliceLaptop, err = client.GetUser("alice", "incorrect password")
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Store, Load, Append File Tests: ", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking for error if filename does not exist in user file namespace when appending")
			err = alice.AppendToFile("not aliceFile", []byte(contentTwo))
			Expect(err).NotTo(BeNil())

			userlib.DebugMsg("Checking for error if filename does not exist in user file namespace when loading")
			data, err := alice.LoadFile("not aliceFile")
			Expect(err).ToNot(BeNil())
			Expect(data).To(BeNil())
		})

		Specify("Sharing Tests", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Charles with a filename that does not exist")
			charlesinvite, err := alice.CreateInvitation("notalicefile", "charles")
			Expect(err).NotTo(BeNil())
			Expect(charlesinvite).NotTo(BeNil())

			userlib.DebugMsg("Alice creating invite for Doris, but should error since doris is not an initialized user")
			dorisinvite, err := alice.CreateInvitation(aliceFile, "doris")
			Expect(err).ToNot(BeNil())
			Expect(dorisinvite).NotTo(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			bobinvite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking for error if unable to verify inv was sent by %s", "alice")
			err = charles.AcceptInvitation("debera", bobinvite, "charlesFile2")
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Initializing user Doris.")
			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking for error if unable to verify inv was sent by %s", "alice")
			err = charles.AcceptInvitation("doris", bobinvite, "charlesFile2")
			Expect(err).ToNot(BeNil())

			err = bob.AcceptInvitation("alice", bobinvite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles storing file %s with content: %s", charlesFile, contentOne)
			charles.StoreFile(charlesFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Charles")
			charlesinvite, err = alice.CreateInvitation(aliceFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking for error if caller already has file with the given filename")
			err = charles.AcceptInvitation("alice", charlesinvite, charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking for error if unable to verify inv was sent by %s", "alice")
			err = charles.AcceptInvitation("doris", charlesinvite, "charlesFile2")
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Alice revokes charles, then charles tries to accept, should error")
			err = alice.RevokeAccess(aliceFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking for error when Charles tries to accept inv after being revoked")
			err = charles.AcceptInvitation("alice", charlesinvite, "charlesFile2")
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Revoke errors: alice revokes bob with incorrect filename")
			err = alice.RevokeAccess("Not alice's file", "bob")
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Bob should still have access to file (can load/append)")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Revoke errors: alice revokes doris, expecting error")
			err = alice.RevokeAccess(aliceFile, "doris")
			Expect(err).ToNot(BeNil())

			/*
				create invite:
				1. The given filename does not exist in the personal file namespace
					of the caller.
				2. The given recipientUsername does not exist.

				accept invite:
`				1. The caller already has a file with the given filename in their personal
					file namespace.
				2. The caller is unable to verify that the secure file share invitation
					pointed to by the given invitationPtr was created by senderUsername.
				3. The invitation is no longer valid due to revocation.`

				revoke access:
				1. The given filename does not exist in the caller’s personal file namespace.
				2. The given filename is not currently shared with recipientUsername.

			*/
		})

	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing shareTree.", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, Charles, and Doris")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice storing file %s with content: %s", aliceFile, contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice creating invite for Bob.")
			bobinvite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", bobinvite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("bob creating invite for charles.")
			charlesinvite, err := bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles accepting invite from Bob under filename %s.", charlesFile)
			err = charles.AcceptInvitation("bob", charlesinvite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles creating invite for doris.")
			dorisinvite, err := charles.CreateInvitation(charlesFile, "doris")
			Expect(err).To(BeNil())

			userlib.DebugMsg("doris accepting invite from charles under filename %s.", dorisFile)
			err = doris.AcceptInvitation("charles", dorisinvite, dorisFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Doris appending to file %s, content: %s", dorisFile, contentTwo)
			err = doris.AppendToFile(dorisFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice appending to file %s, content: %s", aliceFile, contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alice sees expected file data.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that charles sees expected file data.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that doris sees expected file data.")
			data, err = doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("alice revoking charles access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("checking that doris no longer has access")
			err = doris.AppendToFile(dorisFile, []byte(contentOne))
			Expect(err).NotTo(BeNil())
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Basic Test: Testing InitUser/GetUser on a single user that doesn't exist.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		Specify("Basic Test: File with no content.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", emptyString)
			err = alice.StoreFile(aliceFile, []byte(emptyString))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", emptyString)
			err = alice.AppendToFile(aliceFile, []byte(emptyString))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", emptyString)
			err = alice.AppendToFile(aliceFile, []byte(emptyString))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(emptyString + emptyString + emptyString)))
		})

		Specify("Basic Test: Testing More Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			bobinvite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", bobinvite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Charles and Charles accepting invite.")
			charlesinvite, err := alice.CreateInvitation(aliceFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("alice", charlesinvite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err := charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Charles can still load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that Bob cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Basic Test: Testing Revoke Functionality for Shared Users", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charles accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob revoking charles's access from %s.", aliceFile)
			err = bob.RevokeAccess(aliceFile, "charles")
			Expect(err).ToNot(BeNil())
		})

		Specify("Basic Test: Testing Same Username", func() {
			userlib.DebugMsg("Initializing users Alice twice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("alice", defaultPassword)
			Expect(err).NotTo(BeNil())
		})

		Specify("Basic Test: Testing Empty File Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(emptyString, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(emptyString, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(emptyString, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(emptyString)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing files of the same names.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentTwo)
			err = bob.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentThree)))

			userlib.DebugMsg("Bob Loading file...")
			data, err = bob.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)))
		})

		Specify("Basic Test: Testing Case Sensitivity in Username", func() {
			userlib.DebugMsg("Initializing users Alice twice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("Alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing a long input.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentFour)
			err = alice.StoreFile(aliceFile, []byte(contentFour))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentFive)
			err = alice.AppendToFile(aliceFile, []byte(contentFive))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentFour + contentFive)))
		})

		Specify("Basic Test: Testing multiple user sessions.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Creating another user session for Alice on laptop.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Creating another user session for Alice on phone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Creating another user session for Alice on desktop.")
			aliceDesktop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file in alice laptop.")
			err = aliceLaptop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("From alice's phone, appending file data: %s", contentTwo)
			err = alicePhone.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("From alice's desktop, appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing AcceptInvitation Errors.", func() {
			userlib.DebugMsg("Initializing users Alice andBob.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

		})

		Specify("Basic Test: Testing overriding files.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentTwo)
			err = alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)))
		})

		Specify("Basic Test: Test how many keys created.", func() {
			//number of keys in Keystore per user must be a small constant, must not depend on the number of files stored or length of any file, how many users a file has been shared with, or the number of users already in the system
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			num_keys_after_one_file := len(userlib.KeystoreGetMap())

			userlib.DebugMsg("Storing file 1 data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file 2 data: %s", contentOne)
			err = alice.StoreFile(bobFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file 3 data: %s", contentTwo)
			err = alice.StoreFile(charlesFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file 4 data: %s", contentTwo)
			err = alice.StoreFile(dorisFile, []byte(contentFive))
			Expect(err).To(BeNil())

			num_keys_after_4_files := len(userlib.KeystoreGetMap())

			Expect(num_keys_after_one_file).To(BeEquivalentTo(num_keys_after_4_files))
		})

		Specify("Basic Test: Sharing errors", func() {
			userlib.DebugMsg("Initializing users Alice and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Charles and Charles accepting invite.")
			invite, err := alice.CreateInvitation(aliceFile, "charles")
			Expect(err).To(BeNil())

			err = alice.AcceptInvitation("alice", invite, charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that Charles cannot load the file.")
			data, err := charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())
			Expect(data).ToNot(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob cannot load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())
			Expect(data).ToNot(Equal([]byte(contentOne)))
		})
	})

	Describe("Security Tests", func() {
		Specify("Security Test: Malicious content inserted to all entries", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Messing up the Data")
			temp_map := userlib.DatastoreGetMap()
			for temp_uuid := range temp_map {
				userlib.DatastoreSet(temp_uuid, []byte(" "))
			}

			userlib.DebugMsg("Loading file...")
			_, err := alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentFive)
			err = alice.AppendToFile(aliceFile, []byte(contentFive))
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)
			_, err = alice.CreateInvitation(aliceFile, "bob")
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())

		})

		Specify("Security Test: Datastore cleared and tampered with", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Messing up the Data")
			userlib.DatastoreClear()

			userlib.DebugMsg("Loading file...")
			_, err := alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())
			Expect(data).To(BeNil())
		})

		Specify("Swapping invitation ids", func ()  {
			userlib.DebugMsg("Initializing users Alice, Bob, Charles, and Doris")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentTwo)
			err = bob.StoreFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice creating invite for Charles.")
			charlesinvite, err := alice.CreateInvitation(aliceFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("bob creating invite for doris.")
			dorisinvite, err := bob.CreateInvitation(bobFile, "doris")
			Expect(err).To(BeNil())
			
			userlib.DebugMsg("swapping charles and doris inv ptr")
			datastoreMap := userlib.DatastoreGetMap()
			invPtr1 := datastoreMap[charlesinvite]
			invPtr2 := datastoreMap[dorisinvite]
			userlib.DatastoreSet(charlesinvite, invPtr2)
			userlib.DatastoreSet(dorisinvite, invPtr1)

			userlib.DebugMsg("charles tries accepting malicious ptr")
			err = charles.AcceptInvitation("bob", dorisinvite, bobFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("doris tries accepting malicious ptr")
			err = doris.AcceptInvitation("alice", charlesinvite, aliceFile)
			Expect(err).ToNot(BeNil())

		})

		Specify("Security Test: Append twice and swap", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			datastoreMap := userlib.DatastoreGetMap()
			beforeAppend := make(map[userlib.UUID][]byte)

			for k, v := range datastoreMap {
				beforeAppend[k] = v
			}

			userlib.DebugMsg("Appending file data once: %s", contentFour)
			err = alice.AppendToFile(aliceFile, []byte(contentFour))
			Expect(err).To(BeNil())

			var diffUUID1 userlib.UUID
			for key := range datastoreMap {
				if _, ok := beforeAppend[key]; !ok {
					diffUUID1 = key
				}
			}

			beforeAppend = make(map[userlib.UUID][]byte)

			for k, v := range datastoreMap {
				beforeAppend[k] = v
			}

			userlib.DebugMsg("Appending file data twice: %s", contentFive)
			err = alice.AppendToFile(aliceFile, []byte(contentFive))
			Expect(err).To(BeNil())

			append2_map := userlib.DatastoreGetMap()

			var diffUUID2 userlib.UUID

			for key := range append2_map {
				if _, ok := beforeAppend[key]; !ok {
					diffUUID2 = key
					fmt.Println("not in old map append 2")
				}
			}

			append1_data := datastoreMap[diffUUID1]
			append2_data := datastoreMap[diffUUID2]
			userlib.DatastoreSet(diffUUID1, append2_data)
			userlib.DatastoreSet(diffUUID2, append1_data)

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())
			Expect(data).To(BeNil())
		})

		Specify("Basic Test: Testing shareTree.", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, Charles, and Doris")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice storing file %s with content: %s", aliceFile, contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice creating invite for Bob.")
			bobinvite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", bobinvite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("bob creating invite for charles.")
			charlesinvite, err := bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles accepting invite from Bob under filename %s.", charlesFile)
			err = charles.AcceptInvitation("bob", charlesinvite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles creating invite for doris.")
			dorisinvite, err := charles.CreateInvitation(charlesFile, "doris")
			Expect(err).To(BeNil())

			userlib.DebugMsg("doris accepting invite from charles under filename %s.", dorisFile)
			err = doris.AcceptInvitation("charles", dorisinvite, dorisFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Doris appending to file %s, content: %s", dorisFile, contentTwo)
			err = doris.AppendToFile(dorisFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice appending to file %s, content: %s", aliceFile, contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alice sees expected file data.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that charles sees expected file data.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that doris sees expected file data.")
			data, err = doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("alice revoking charles access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("checking that doris no longer has access")
			err = doris.AppendToFile(dorisFile, []byte(contentOne))
			Expect(err).NotTo(BeNil())
		})

	})

})
