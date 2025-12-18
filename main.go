package main

import (
	"fmt"
	"log"
	"net"

	"github.com/hirochachacha/go-smb2"
)

func main() {
	// --- Configuration ---
	// Please change these values to match your SMB server details.
	serverAddress := "127.0.0.1:4789" // SMB server address and port
	shareName := "public"             // The name of the share you want to connect to.
	username := "user"                // The username for authentication.
	password := "password"            // The password for authentication.
	fileNameToRead := "hello.txt"     // The name of the file you want to read from the share.
	// ---------------------

	// Establish TCP connection
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("failed to establish TCP connection: %v", err)
	}
	defer conn.Close()

	// Create SMB dialer with NTLM authentication
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     username,
			Password: password,
		},
	}

	// Dial the SMB session
	session, err := d.Dial(conn)
	if err != nil {
		log.Fatalf("failed to dial SMB session: %v", err)
	}
	defer session.Logoff()

	fmt.Printf("Successfully connected to SMB server at %s.\n", serverAddress)

	// Mount the share
	fs, err := session.Mount(shareName)
	if err != nil {
		log.Fatalf("failed to mount share '%s': %v", shareName, err)
	}
	defer fs.Umount()

	fmt.Printf("Successfully mounted the '%s' share.\n", shareName)

	// List files in the root of the share
	files, err := fs.ReadDir(".")
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	fmt.Println("\nFiles in share:")
	for _, file := range files {
		fmt.Printf("- %s\n", file.Name())
	}

	// Read the content of the file
	content, err := fs.ReadFile(fileNameToRead)
	if err != nil {
		log.Fatalf("failed to read file '%s': %v", fileNameToRead, err)
	}

	fmt.Printf("\nContent of '%s':\n", fileNameToRead)
	fmt.Println(string(content))
}
