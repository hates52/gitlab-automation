package ldap

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type LDAPConfig struct {
	Host     string
	BindDN   string
	Password string
	BaseDN   string
}

type LDAPConnector struct {
	conn   *ldap.Conn
	baseDN string
}

type LDAPGroupSyncer struct {
	connector   *LDAPConnector
	groupFilter string
}

func NewLDAPConnector(config LDAPConfig) (*LDAPConnector, error) {
	// Pripojeni k LDAPu
	conn, err := ldap.DialURL(config.Host)
	if err != nil {
		return nil, err
	}

	// Autentizace uzivatele
	err = conn.Bind(config.BindDN, config.Password)
	if err != nil {
		return nil, err
	}

	// Vytvoreni instance
	return &LDAPConnector{
		conn:   conn,
		baseDN: config.BaseDN,
	}, nil
}

func (l *LDAPConnector) Close() {
	if l.conn != nil {
		l.conn.Close()
	}
}

func NewLDAPGroupSyncer(connector *LDAPConnector, groupFilter string) *LDAPGroupSyncer {
	return &LDAPGroupSyncer{
		connector:   connector,
		groupFilter: groupFilter,
	}
}

func (s *LDAPGroupSyncer) GetLdapGroups() (*ldap.SearchResult, error) {
	// Definice vyhledavaciho pozadavku
	searchRequest := ldap.NewSearchRequest(
		s.connector.baseDN, // Zakladni DN
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		s.groupFilter, // Filtr pro vyhledani
		nil,
		nil,
	)

	// Provedeme vyhledani v LDAPu
	result, err := s.connector.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *LDAPGroupSyncer) ListLdapGroupMemberDNs(groupDN string) ([]string, error) {
	// Definice vyhledavaciho pozadavku
	searchRequest := ldap.NewSearchRequest(
		groupDN, // Zakladni DN
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		s.groupFilter, // Filtr pro vyhledani
		[]string{"member"},
		nil,
	)

	// Provedeme vyhledani v LDAPu
	result, err := s.connector.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("not found members for group %s", groupDN)
	}

	members := result.Entries[0].GetAttributeValues("member")

	return members, nil
}
