# **CS161 Proj2: Design Doc**
## **Data Structures**
*What data structures are you going to use? List any struct definitions that you plan on including, along with the attributes these structs will contain. We’d recommend starting with a few core data structures (e.g. `struct user`, `struct file`, etc.), and adding additional attributes and structs as you need them.*
### **`User struct`**
### Attributes
- `string username`
- `string password`
- `string salt`
- `map files`
    - key = filename, value = invID
### **`File struct`**
### Attributes
- `string filename`
- `bytes[] contents`
- `User owner`
- `FileShareMap`
- `PairShare Map`
- `sharingTree`
### **`FileShare struct`**
### Attributes
- `PairShare shareFileInfo`
- `[]bytes symmKey`
### **`PairShare struct`**
### Attributes
- `[]bytes macKey`
- `[]bytes encKey`
- `UUID filePtr`
___
## **User Authentication**
*How will you authenticate users? What information will you store in Datastore/Keystore for each user? How will you ensure that a user can have multiple client instances (e.g. laptop, phone, etc.) running simultaneously?*

We authenticate users with their username, password, and a unique salt corresponding to each user. We will store their public keys in Keystore, and their private and symmetric keys in the User struct, which is stored in Datastore. A user can have multiple client instances running simultaneously since we are storing the files in the Datastore and the different client instances would just be accessing the same file.

### `InitUser`
**Input:** `username string, password string` \
**Output:** `userdataptr *User, error err` \
**Error Checks:**
1. Username already exists
    - Check datastore by getting the `UUID(username)` and error if true
2. Username is empty

Actions:
1. Generate random **salt** and store into DataStore 
    - key = `uuid(username)`, value(aka salt) = `RandomBytes()`
2. Initialize a user struct and store as (key, value) in DataStore
    - key = `mac(enc(UUID(Hash(password || salt))))`
    - value = `mac(enc(json.Marshal(User struct)))`
3. Generate set of public keys and store into KeyStore
    -  Generate a set of public and private keys (`PK, SK = PKEKeyGen()`)
    -  Generate a set of Digital Signature keys (`SignKey, VerifyKey = DSKeyGen()`)
### `GetUser`
**Input:** `username string, password string` \
**Output:** `userdataptr *User, error err` \
**Error Checks:**
1. No initialized user given username
    - Check datastore by getting the `UUID(username)` and error if false
2. Invalid user credentials = incorrect password
    - Check datastore by getting the intended key and error if false
3. User struct cannout be obtained due to malicious actions
    - Since we mac the user struct stored in DataStore if the mac is invalid, error

**Actions:** \
Return `json.Unmarshal(userlib.DatastoreGet(key))`

___
## **File Storage and Retrieval**
*How will a user store their files? How will a user retrieve their files from the server? How will your design support efficient file append?*

Files will be stored in the Datastore as key-value pairs of `(UUID(invID), mac(enc(struct File)))`. To keep track of them within the personal file directory of a user, there are also key-value pairs in the `fileMap` of the User struct. The key will be `fileName`, and the value will be `invID` which we use to look the file up in Datastore. Since the File struct in Datastore is MAC'ed and encrypted, it will have confidentiality and integrity against the attacker. 

### `User.StoreFile`
This function instantiates a new File struct with the appropriate attributes (file name, file contents, file owner, an empty `FileShareMap`, an empty `PairShareMap`, and a `sharingTree` with the owner as the root node) and adds the key-value pair of `(UUID(invID), mac(enc(struct File)))` to user's `fileMap`.

If the file name already exists in the user's directory, the function finds the File struct with the same file name in Datastore and replaces the contents.

**Input:** `fileName string, content []bytes` \
**Output:** `error err` \
**Error Checks:**
1. Given that the user has proven their identity with the correct username and password to log in, no other malicious action  

**Actions:** 
- Generate 16 random bytes for the file's' `invID` using `UUID.New()`. 
- Add `(fileName, invID)` to user's `fileMap`
- Add `(invID, mac(enc(struct File)))` to Datastore-- MAC-ing with user's `K1` and encrypting user's `K2`. 

Return `error` or `None`

### `User.LoadFile`
This function loads a file if the user has the appropriate login, invitation ID, and the file exists in their directory (`fileMap`).

**Input:** `fileName string` \
**Output:** `(content []byte, err error)` \
**Error Checks:**
1. If the file name does not exist in the personal file namespace of the user (there are no keys in `fileMap` with the file name), it means that the file does not exist and it errors.
2. If the file contents have been tampered with (if `mac(enc(struct File))` does not match the one stored in Datastore) then the function call also errors.  

**Actions:** 
- Get `invID` from user's `fileMap` using `fileName`. Return `error` Error Check 1 fails. 
- Get file struct from Datastore with `invID`. Return error if Error Check 2 fails.
- Get `contents` of file struct and load into `content` to be returned to user.

Return `content` or `error`

### `User.AppendToFile`
The function appends to a file if the user has the appropriate login, invitation ID, and the file exists in their directory (`fileMap`).

This implementation is not dependent on the size of the file/how many users it is shared with since every user is referencing the same File struct in Datastore (would not have to edit every person's copy of the file).

**Input:** `fileName string` \
**Output:** `(filename string, content []byte)` \
**Error Checks:**
1. If the file name does not exist in the personal file namespace of the user (there are no keys in `fileMap` with the file name), it means that the file does not exist and it errors.
2. If the file contents have been tampered with (if `mac(enc(struct File))` does not match the one stored in Datastore) then the function call also errors.  

**Actions:** 
- Get `invID` from user's `fileMap` using `fileName`. Return `error` Error Check 1 fails. 
- Get file struct from Datastore with `UUID(invID)`. Return error if Error Check 2 fails.
- Get `contents` of file struct and append `content` to `contents`.

Return `error` or `None`

___
## **File Sharing and Revocation**
*How will a user share files with another user? How does this shared user access the shared file after accepting the invitation? How will a user revoke a different user’s access to a file? How will you ensure a revoked user cannot take any malicious actions on a file?*

Using the digital signing key generated when we instantiated the User struct, we will sign the new `UUID` of the FileShare struct that we put in Datastore to share it with the recipient to validate the identity of the sender. This FileShare struct will be asymmetrically encrypted with the recipient's public key. The FileShare struct will also contain a symmetrically encrypted PairShare object to access the File struct in Datastore (the sender and recipient shares the same key), as well as the symmetric key to decrypt the PairShare object.

 The recipient will verify the signature on the UUID and access the shared file by using the symmetric key they got from the sender in the FileShare object to decrypt the Pairshare object, which will contain the UUID of the File struct in Datastore, key to decrypt and verify the MAC of the File struct to access the File struct in Datastore.

 After the file is shared, the File struct's FileShareMap and PairshareMap are updated with the recipient's username mapped to the UUID of the new FileShare and Pairshare objects. The recipient will be added as a node to the `sharingTree` which is stored in the file.
 
 If the user's access to a file is revoked, we call the helper method `removeUser` to go recursively traverse the `sharingTree` and update the MAC verification and decryption keys for all its children except for the revoked user. The revoked user and its children are pruned from the tree. We then delete the file name from the user's `fileMap`.

### `User.CreateInvitation`
This function shares the UUID which points to a FileShare struct that has been encrypted with the recipient's public key.

**Input:** `filename string`, `recipientUsername string` \
**Output:** `(invitationPtr UUID, err error)` \
**Error Checks:**
1. If the file name does not exist in the personal file namespace of the user (there are no keys in `fileMap` with the file name), it means that the file does not exist and it errors.
2. If the file contents have been tampered with (if `mac(enc(struct File))` does not match the one stored in Datastore) then the function call also errors.  

**Actions:** 
-  Get `invID` from user's `fileMap` using `fileName`. Return `error` Error Check 1 fails. 
- Get file struct from Datastore with `UUID(invID)`. Return error if Error Check 2 fails.
- Instantiate a FileShare struct that contains a symmetrically encrypted PairShare struct and the symmetric key to decrypt that. The PairShare struct is instantiated to contain the UUID that points to the File struct in Datastore, and the MAC and decryption to verify the confidentiality and integrity of the file. 
- Save the new FileShare struct to Datastore with a randomly generated UUID.
- Return the randomly generated UUID as `invitationPtr`.

Return `invitationPtr (UUID)`, `err (error)`
### `User.AcceptInvitation`

**Input:** `senderUsername string, invitationPtr UUID, filename string` \
**Output:** `error err` \
**Error Checks:**
1. Caller already has a file with `filename` in personal file nampespace
    - If `filename` in `files` map, error
2. Caller cannot verify that the `invitationPtr` was given by `senderUsername`
    - If `DSVerify(vk DSVerifyKey, invitationPtr, sig []byte)` is not valid, this means we cannot verify that the invitation was given by the provider user. The DSVerifyKey can be obtained from the keystore.\
3. Invitation is no longer valid due to revocation
    - Call `is_valid_user(root of sharingTree, user's username)` to check.

**Actions:**
- Verify the signature of the UUID with the DSVerify
- Decrypt the `fileShare` struct with their private key and inside the `fileShare` struct, the sender will have the symmetric key stored. The recipient adds the symmetric key to their `fileMap` and can use that to decrypt the `PairShare` struct in datastore. We then update the file's `FileShareMap` and `PairShareMap`.
- Add recipient as a node to the `sharingTree` which is stored in the file.

Return `error err` or `None`

### `User.RevokeAccess`
If the user's access to a file is revoked, we call the helper method `removeUser` to go recursively traverse the `sharingTree` and update the MAC verification and decryption keys for all its children except for the revoked user. The revoked user and its children are pruned from the tree. We then delete the file name from the user's `fileMap`.

**Input:** `filename string, recipientUsername string`\
**Output:** `error err` \
**Error Checks:**
1. Caller already has a file with `filename` in personal file nampespace
    - If `filename` in `files` map, error
2. Caller cannot verify that the `invitationPtr` was given by `senderUsername`
    - If `DSVerify(vk DSVerifyKey, invitationPtr, sig []byte)` is not valid, this means we cannot verify that the invitation was given by the provider user. The DSVerifyKey can be obtained from the keystore.\
    **CHECK: (how do we access the signature to compare)**
3. Given `filename` is not currently shared with the user. 
    - Call `is_valid_user(root of sharingTree, recipientUsername)` to check.

**Actions:**
- Generate a new set of MAC (`macKey`) and decryption key(`decKey`) for the corresponding `filename`. 
From the user's `fileMap`, get the `FileShareTree` and call `removeUser(root of FileShareTree, macKey, decKey, recipientUsername)`.

Return `error err` or `None`
___
## **Helper Methods**
Are there helper methods that you’ll create? There are a few in particular that may help simplify your implementation.
- `store_datastore(key, macKey, encKey, value)`
    - Use userlib.DatastoreSet(key UUID, value[] byte) to store something in Datastore while ensuring confidentiality and integrity against attackers
    - value = MAC(enc(value))
- `get_datastore(key, macKey, encKey)`
    - Use userlib.DatastoreGet(key UUID) to verify the MAC and decrypt the value that we want to get from Datastore
    - Returns the decrypted value if the MAC is valid, else error.
- `is_valid_user(Fileshare Node, string username)`
    - Traverses the `sharingTree` and returns True if user is in tree, False otherwise.
- `removeUser(FileShare node, macKey, decKey, revokedUser)`
    -   removes the revoked user from the `sharingTree` and updates everyone else's `mackey` and `decKey` for the file
___
# Test Proposals
1. **Design Requirement 3.5.7**: The client MUST NOT assume that filenames are globally unique. For example, user bob can have a file named foo.txt and user alice can have a file named foo.txt. The client MUST keep each user’s file namespace independent from one another.
    - Create 2 users.
    - Each user creates a file named “foo” in their personal file namespace with different contents – User1.StoreFile(“foo”, “abc”) and User2.StoreFile(“foo”, “def”)
    - Expect that User1.LoadFile(“foo”) != User2.LoadFile(“foo”)

2. **Design Requirement 3.1.1**: Client should have unique usernames, be case-sensitive, and support usernames of any length > 0.
    - Create 2 users
    - If the user2 has the same username as user1, userlib.
    - DatastoreGet(key UUID) should return true, where the UUID = UUID(username).
    - Expect that user1.username != user2.username.

3. **Design Requirement 3.3.3**: The following SHOULD be avoided because they are dangerous design patterns that often lead to subtle vulnerabilities (Reusing the same key for multiple purposes, Authenticate-then-encrypt, Decrypt-then-verify)
    - Create 1 user.
    - Create function called key_count(key) contains a dictionary that counts how many times a key has been used, and for what purpose. 
    - MAC_count, encrypt_count = key_count(key).get[‘MAC’], key_count(key).get[‘encrypt’]
    - Expect MAC_count == encrypt_count
4. **Design Requirement 3.5.3**: The client MUST prevent adversaries from learning filenames and the length of filenames. The client MAY use and store filenames in a deterministic manner.
    - Create 1 user.
    - Create a file “file1.txt” with some content “foo” and user calls User.storeFile
    - File is stored into datastore → when user logs back in the can load file from datastore
5. **Design Requirement 3.4**: The client MUST NOT save any data to the local file system. If the client is restarted, it must be able to pick up where it left off given only a username and password. Any data requiring persistent storage MUST be stored in either Keystore or Datastore.
    - Create 1 user.
    - Create file “foo” with contents “abc”
    - Get user from two different devices (aka, get user twice)
    - In the second login, append “def” to file.
    - In the first login, X = load file contents. 
    - In the second login, Y = load file contents.
    - Expect X == Y.
6. **Design Requirement 3.5.6**: Filenames MAY be any length, including zero (empty string).
    - Create 1 user.
    - User stores file with name: “foo” * 1000 and contents = “abc”
    - User stores file with name: “f” and contents = “def”
    - Expect User.LoadFile(“foo”*1000) == “abc and User.LoadFile(“f”) == “def”
