package groupsync

import (
	"log"

	"github.com/go-ldap/ldap/v3"
)

type LDAPGroupSyncer struct {
	conn        *ldap.Conn
	baseDN      string
	groupFilter string
}

type LDAPConfig struct {
	Host        string
	BindDN      string
	Password    string
	BaseDN      string
	GroupFilter string
}

func NewLDAPGroupSyncer(ldapHost, ldapBindDN, ldapPassword, ldapSearchBase, ldapFilter string) (*LDAPGroupSyncer, error) {
	// Pripojeni k LDAPu
	conn, err := ldap.DialURL(ldapHost)
	if err != nil {
		return nil, err
	}

	// Autentizace uzivatele
	err = conn.Bind(ldapBindDN, ldapPassword)
	if err != nil {
		return nil, err
	}

	// Vytvoreni instance
	return &LDAPGroupSyncer{
		conn:        conn,
		baseDN:      ldapSearchBase,
		groupFilter: ldapFilter,
	}, nil
}

func (l *LDAPGroupSyncer) Close() {
	if l.conn != nil {
		l.conn.Close()
	}
}

func (l *LDAPGroupSyncer) GetGroups() []string {
	// Definice vyhledávacího požadavku
	searchRequest := ldap.NewSearchRequest(
		l.baseDN, // Zakladni DN
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		l.groupFilter,  // Filtr pro vyhledani skupin
		[]string{"cn"}, // Atributy, které chceme ziskat (napr. common name)
		nil,
	)

	// Provest hledani
	result, err := l.conn.Search(searchRequest)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	// Extrahovani nazvu skupin
	var groups []string
	for _, entry := range result.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}

	return groups
}

func (l *LDAPGroupSyncer) GetMembers(group string) []string {
	// Implementujte kód pro získání členů skupiny z LDAP.
	return []string{"member1", "member2"}
}
