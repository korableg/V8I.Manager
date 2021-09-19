package lstparser

import (
	"errors"
	"regexp"
	"strings"

	"github.com/korableg/V8I.Manager/pkg/clusterdb"
)

func Parse(lst []byte) ([]*clusterdb.ClusterDB, error) {

	tokens, err := tokenize(lst)
	if err != nil {
		return nil, err
	}

	dblist := make([]*clusterdb.ClusterDB, len(tokens))

	for i, v := range tokens {
		dblist[i] = tokenToClusterDB(v)
	}

	return dblist, nil

}

func tokenize(b []byte) ([]string, error) {
	re := regexp.MustCompile(
		`([0-9A-Fa-f]{8}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{12})` +
			`([,][^,]+){9}[,][^{]+[{][\d]([,][\d]{14}){2}([^}]+[}])[^}]+`)
	bs := re.FindAll(b, -1)
	if bs == nil {
		return nil, errors.New("cluser database items not found")
	}
	s := make([]string, len(bs))
	for i, val := range bs {
		s[i] = string(val)
		s[i] = strings.ReplaceAll(s[i], `"`, ``)
	}
	return s, nil
}

func tokenToClusterDB(token string) *clusterdb.ClusterDB {

	splittedToken := strings.Split(token, ",")

	db := &clusterdb.ClusterDB{}
	db.ID = splittedToken[0]
	db.Description = splittedToken[2]
	db.Ref, db.Server = parseConnectionString(splittedToken[8])
	db.Folder = db.Server
	db.Name = db.Ref

	if len(db.Description) > 0 {
		splittedDescription := strings.Split(db.Description, "/")
		db.Name = splittedDescription[len(splittedDescription)-1]
		db.Folder = strings.Join(splittedDescription[0:len(splittedDescription)-1], "/")
	}

	return db
}

func parseConnectionString(connectionString string) (ref, server string) {

	splittedConnectionString := strings.Split(connectionString, ";")
	for _, val := range splittedConnectionString {
		vLower := strings.ToLower(val)
		if strings.HasPrefix(vLower, "ref") {
			ref = strings.Split(val, "=")[1]
		} else if strings.HasPrefix(vLower, "srvr") {
			server = strings.Split(val, "=")[1]
		}
	}

	return

}
