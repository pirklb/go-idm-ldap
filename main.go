package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-ldap/ldap/v3"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("eDirectory LDAP")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	url := os.Getenv("LDAP_URL")
	if url == "" {
		url = "ldap://wipidm03.lkw-walter.com"
	}
	user := os.Getenv("LDAP_USER")
	if user == "" {
		user = "cn=LDAPReader,o=System"
	}
	pawd := os.Getenv("LDAP_PASSWORD")
	if pawd == "" {
		log.Fatal("LDAP_PASSWORD is not set")
	}

	l, err := ldap.DialURL(url)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// First bind with a read only user
	fmt.Println("l.Bind")
	err = l.Bind(user, pawd)
	if err != nil {
		log.Fatal(err)
	}

	baseDN := "cn=WURM,cn=Level10,cn=RoleDefs,cn=RoleConfig,cn=AppConfig,cn=UserApplication,cn=DriverSet,o=System"
	//filter := fmt.Sprintf("(objectClass=%s)", ldap.EscapeFilter("*"))
	filter := "(objectClass=*)"
	fmt.Println("filter: ", filter)
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"modifyTimestamp"}, []ldap.Control{})

	result, err := l.Search(searchReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %v entries\n", len(result.Entries))
	for _, entry := range result.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("modifyTimestamp"))
	}
}
