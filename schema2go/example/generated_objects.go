package objects

/*
autogenerated from:
schema files:
core.schema nis.schema

object description:
[
{
	"Name": "User",
	"Desc": "A user with several extensions",
	"ObjectClasses": ["posixAccount","person"],
	"FilterObjectClass": "posixAccount"
}
]

*/

import (
	"errors"
	"fmt"
	"github.com/bytemine/ldap-crud/crud"
	"github.com/rbns/ldap"
	"strings"
)

// User: A user with several extensions
type User struct {
	dn           string
	Person       bool
	PosixAccount bool

	//MUST attributes
	// attribute definition missing
	Cn []string `json:",omitempty"`
	// An integer uniquely identifying a group in an administrative domain
	GidNumber string `json:",omitempty"`
	// The absolute path to the home directory
	HomeDirectory string `json:",omitempty"`
	// RFC2256: last (family) name(s) for which the entity is known by
	Sn []string `json:",omitempty"`
	// attribute definition missing
	Uid []string `json:",omitempty"`
	// An integer uniquely identifying a group in an administrative domain
	UidNumber string `json:",omitempty"`

	// MAY attributes
	// attribute definition missing
	Description []string `json:",omitempty"`
	// The GECOS field; the common name
	Gecos string `json:",omitempty"`
	// The path to the login shell
	LoginShell string `json:",omitempty"`
	// attribute definition missing
	SeeAlso []string `json:",omitempty"`
	// RFC2256: Telephone Number
	TelephoneNumber []string `json:",omitempty"`
	// attribute definition missing
	UserPassword []string `json:",omitempty"`
}

func NewUser(dn string) *User {
	o := new(User)
	o.dn = dn
	return o
}

func (o *User) FilterObjectClass() string {
	return "posixAccount"
}

func (o *User) Copy() crud.Item {
	c := NewUser(o.dn)

	c.Cn = make([]string, len(o.Cn))
	copy(c.Cn, o.Cn)
	c.GidNumber = o.GidNumber
	c.HomeDirectory = o.HomeDirectory
	c.Sn = make([]string, len(o.Sn))
	copy(c.Sn, o.Sn)
	c.Uid = make([]string, len(o.Uid))
	copy(c.Uid, o.Uid)
	c.UidNumber = o.UidNumber

	c.Description = make([]string, len(o.Description))
	copy(c.Description, o.Description)
	c.Gecos = o.Gecos
	c.LoginShell = o.LoginShell
	c.SeeAlso = make([]string, len(o.SeeAlso))
	copy(c.SeeAlso, o.SeeAlso)
	c.TelephoneNumber = make([]string, len(o.TelephoneNumber))
	copy(c.TelephoneNumber, o.TelephoneNumber)
	c.UserPassword = make([]string, len(o.UserPassword))
	copy(c.UserPassword, o.UserPassword)
	return c
}

func (o *User) Dn() string {
	return o.dn
}

func (o *User) MarshalLDAP() (*ldap.Entry, error) {
	e := ldap.NewEntry(o.dn)

	if o.Person {
		e.AddAttributeValue("objectClass", "person")
		if len(o.Sn) == 0 {
			return nil, errors.New(fmt.Sprintf("Marshalling %v: Attribute %v is empty", "User", "Sn"))
		}

		e.AddAttributeValues("sn", o.Sn)

		e.AddAttributeValues("seeAlso", o.SeeAlso)
		e.AddAttributeValues("telephoneNumber", o.TelephoneNumber)
	}
	if o.PosixAccount {
		e.AddAttributeValue("objectClass", "posixAccount")
		if len(o.Cn) == 0 {
			return nil, errors.New(fmt.Sprintf("Marshalling %v: Attribute %v is empty", "User", "Cn"))
		}
		if len(o.GidNumber) == 0 {
			return nil, errors.New(fmt.Sprintf("Marshalling %v: Attribute %v is empty", "User", "GidNumber"))
		}
		if len(o.HomeDirectory) == 0 {
			return nil, errors.New(fmt.Sprintf("Marshalling %v: Attribute %v is empty", "User", "HomeDirectory"))
		}
		if len(o.Uid) == 0 {
			return nil, errors.New(fmt.Sprintf("Marshalling %v: Attribute %v is empty", "User", "Uid"))
		}
		if len(o.UidNumber) == 0 {
			return nil, errors.New(fmt.Sprintf("Marshalling %v: Attribute %v is empty", "User", "UidNumber"))
		}

		e.AddAttributeValues("cn", o.Cn)
		e.AddAttributeValue("gidNumber", o.GidNumber)
		e.AddAttributeValue("homeDirectory", o.HomeDirectory)
		e.AddAttributeValues("uid", o.Uid)
		e.AddAttributeValue("uidNumber", o.UidNumber)

		e.AddAttributeValues("description", o.Description)
		e.AddAttributeValue("gecos", o.Gecos)
		e.AddAttributeValue("loginShell", o.LoginShell)
		e.AddAttributeValues("userPassword", o.UserPassword)
	}
	return e, nil
}

func (o *User) UnmarshalLDAP(e *ldap.Entry) error {
	o.dn = e.DN

	for _, v := range e.GetAttributeValues("objectClass") {
		switch strings.ToLower(v) {
		case "person":
			o.Person = true
		case "posixaccount":
			o.PosixAccount = true

		}
	}

	o.Cn = e.GetAttributeValues("cn")
	o.GidNumber = e.GetAttributeValue("gidNumber")
	o.HomeDirectory = e.GetAttributeValue("homeDirectory")
	o.Sn = e.GetAttributeValues("sn")
	o.Uid = e.GetAttributeValues("uid")
	o.UidNumber = e.GetAttributeValue("uidNumber")

	o.Description = e.GetAttributeValues("description")
	o.Gecos = e.GetAttributeValue("gecos")
	o.LoginShell = e.GetAttributeValue("loginShell")
	o.SeeAlso = e.GetAttributeValues("seeAlso")
	o.TelephoneNumber = e.GetAttributeValues("telephoneNumber")
	o.UserPassword = e.GetAttributeValues("userPassword")
	return nil
}
