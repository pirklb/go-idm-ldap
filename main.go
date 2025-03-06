// Version: 0.0.2503061
// 0.0.2503061 ... added support for commandline parameters

package main

import (
	"flag"
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
		log.Printf("Error loading .env file: %s, but as this is not the only option to get the needed parameters the program continues", err)
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

	flag.StringVar(&url, "url", url, "LDAP URL")
	flag.StringVar(&user, "user", user, "Bind User DN")
	flag.StringVar(&pawd, "pawd", pawd, "Bind User Password")
	flag.Parse()

	if pawd == "" || user == "" || url == "" {
		log.Fatal("You have not provided all required parameters, so the program ends now")
	}
	fmt.Println("Parameters:")
	fmt.Printf("Url: %v\nUser: %v\nPassword: ***\n", url, user)
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

	baseDN := "cn=RoleDefs,cn=RoleConfig,cn=AppConfig,cn=UserApplication,cn=DriverSet,o=System"
	//filter := fmt.Sprintf("(objectClass=%s)", ldap.EscapeFilter("*"))
	filter := "(objectClass=*)"
	filter = "(nrfRoleCategoryKey=owner-managed)"
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
