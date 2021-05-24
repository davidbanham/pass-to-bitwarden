package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	output := csv.NewWriter(os.Stdout)
	output.Write([]string{"folder", "favorite", "type", "name", "notes", "fields", "login_uri", "login_username", "login_password", "login_totp"})

	homeDir := os.Getenv("HOME")
	storeDir := fmt.Sprintf("%s/.password-store", homeDir)

	err := filepath.Walk(storeDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) != ".gpg" {
				return nil
			}

			fileName := info.Name()

			username := strings.TrimSuffix(fileName, filepath.Ext(fileName))

			dir, _ := filepath.Split(path)

			site := strings.TrimPrefix(dir, storeDir+"/")

			passID := fmt.Sprintf("%s/%s", site, username)
			if site == "" {
				passID = username
			}

			cmd := exec.Command("pass", passID)
			var out bytes.Buffer
			cmd.Stdout = &out
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}

			notes := out.String()
			lines := strings.Split(notes, "\n")
			var password string
			if len(lines) > 0 {
				password = lines[0]
			}
			//log.Printf("DEBUG out.String(): %+v \n", out.String())
			//log.Printf("DEBUG site: %+v \n", site)
			//log.Printf("DEBUG username: %+v \n", username)
			//log.Printf("DEBUG password: %+v \n", password)
			//log.Printf("DEBUG notes: %+v \n", notes)

			line := []string{
				"import", // folder
				"",       //favorite
				"",       //type
				site,     //name
				notes,    //notes
				"",       //fields
				"",       //login_uri
				username, //login_username
				password, //login_password
				"",       //login_totp
			}

			if err := output.Write(line); err != nil {
				log.Fatal(err)
			}

			output.Flush()

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	output.Flush()
}
