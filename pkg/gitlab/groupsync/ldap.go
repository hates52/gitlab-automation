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

func (l *LDAPGroupSyncer) GetGroups() *ldap.SearchResult {
	// Definice vyhledavaciho pozadavku
	searchRequest := ldap.NewSearchRequest(
		l.baseDN, // Zakladni DN
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		l.groupFilter, // Filtr pro vyhledani skupin
		nil,
		nil,
	)

	// Provedeme vyhledani v LDAPu
	result, err := l.conn.Search(searchRequest)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	return result
}

func (l *LDAPGroupSyncer) GetMembers() *ldap.SearchResult {
	searchRequest := ldap.NewSearchRequest(
		l.baseDN, // Zakladni DN
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(cn=*)",            // Filtr pro vyhledani skupin
		[]string{"members"}, // Atributy, které chceme ziskat (napr. common name)
		nil,
	)
	// Provedeme vyhledani v LDAPu
	result, err := l.conn.Search(searchRequest)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	return result
}
