package dbbuilder

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/korableg/V8I.Manager/app/api/onecdb"
)

type (
	DBBuilder interface {
		Build(path string) ([]onecdb.DB, error)
	}

	Builder struct {
		regexp *regexp.Regexp
	}
)

func NewBuilder() (*Builder, error) {
	r, err := regexp.Compile(`([0-9A-Fa-f]{8}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{12})([,][^,]+){9}[,][^{]+[{][\d]([,][\d]{14}){2}([^}]+[}])[^}]+`)
	if err != nil {
		return nil, fmt.Errorf("compile regexp: %w", err)
	}

	return &Builder{regexp: r}, nil
}

func (b *Builder) Build(path string) ([]onecdb.DB, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	tokens := b.regexp.FindAll(data, -1)

	dbs := make([]onecdb.DB, 0, len(tokens))

	for _, token := range tokens {
		tokenPart := strings.Split(string(token), ",")

		uid, err := uuid.Parse(tokenPart[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse uuid: %w", err)
		}

		server, ref := parseConnection(tokenPart[8])

		db := onecdb.DB{
			UUID:    uid,
			Name:    ref,
			Connect: fmt.Sprintf(`Srvr="%s";Ref="%s"`, server, ref),
			Folder:  server,
		}

		dbs = append(dbs, db)
	}

	return dbs, nil
}

func parseConnection(connectionString string) (server, ref string) {
	splintedConnectionString := strings.Split(connectionString, ";")
	for _, val := range splintedConnectionString {
		vLower := strings.ToLower(val)
		if strings.HasPrefix(vLower, "ref") {
			ref = strings.Split(val, "=")[1]
		} else if strings.HasPrefix(vLower, "srvr") {
			server = strings.Split(val, "=")[1]
		}
	}

	return strings.Replace(server, `"`, "", -1), strings.Replace(ref, `"`, "", -1)
}
