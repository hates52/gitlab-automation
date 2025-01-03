package ldap

import (
	"fmt"

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

func (l *LDAPGroupSyncer) GetLdapGroups() (*ldap.SearchResult, error) {
	// Definice vyhledavaciho pozadavku
	searchRequest := ldap.NewSearchRequest(
		l.baseDN, // Zakladni DN
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		l.groupFilter, // Filtr pro vyhledani
		nil,
		nil,
	)

	// Provedeme vyhledani v LDAPu
	result, err := l.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (l *LDAPGroupSyncer) ListLdapGroupMembers(groupDN string) ([]string, error) {
	// Definice vyhledavaciho pozadavku
	searchRequest := ldap.NewSearchRequest(
		groupDN, // Zakladni DN
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		l.groupFilter, // Filtr pro vyhledani
		[]string{"member"},
		nil,
	)

	// Provedeme vyhledani v LDAPu
	result, err := l.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("not found members for group %s", groupDN)
	}

	members := result.Entries[0].GetAttributeValues("member")

	return members, nil
}
