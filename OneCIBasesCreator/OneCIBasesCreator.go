package OneCIBasesCreator

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"regexp"
	"strings"
)

type DB struct {
	ID string
	Server string
	Ref string
	Description string
	Name string
	Folder string
}

func Create(inPathLst, outPathiBases []string) error {

	DBs, err := getDBs(inPathLst)
	if err != nil {
		return err
	}

	iBases, err := dbsToIBases(DBs)
	if err != nil {
		return err
	}

	for _, path := range outPathiBases {
		err = ioutil.WriteFile(path, iBases, 0644)
		if err != nil {
			return err
		}
	}

	return nil

}

func getDBs(inPathLst []string) ([]*DB, error) {

	dirtyBytes, err := readLst(inPathLst)
	if err != nil {
		return nil, err
	}

	rawDBs, err := findDBs(dirtyBytes)
	if err != nil {
		return nil, err
	}

	DBs := make([]*DB, 0)

	for _, rawDB := range rawDBs {
		DBs = append(DBs, convertRawDBToDB(rawDB))
	}

	return DBs, nil

}

func readLst(pathLst []string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, p := range pathLst {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func findDBs(b []byte) ([]string, error) {
	re := regexp.MustCompile(
		`([0-9A-Fa-f]{8}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{12})` +
			`([,][^,]+){9}[,][^{]+[{][\d]([,][\d]{14}){2}([^}]+[}])[^}]+`)
	bs := re.FindAll(b, -1)
	if bs == nil {
		return nil, errors.New("DB not found in the bytes slice")
	}
	s := make([]string, len(bs))
	for i, val := range bs {
		s[i] = string(val)
	}
	return s, nil
}

func convertRawDBToDB(rawDB string) *DB {
	db := &DB{}

	fields := strings.Split(rawDB, ",")

	connStr := strings.Trim(fields[8], "\"")
	connStrFields := strings.Split(connStr, ";")
	for _, val := range connStrFields {
		vLower := strings.ToLower(val)
		if strings.HasPrefix(vLower, "ref") {
			db.Ref = strings.Split(val, "=")[1]
		} else if strings.HasPrefix(vLower, "srvr") {
			db.Server = strings.Split(val, "=")[1]
		}
	}

	db.ID = strings.Trim(fields[0], "\"")
	db.Description = strings.Trim(fields[2], "\"")

	if len(db.Description) > 0 {
		descrFields := strings.Split(db.Description, "/")
		db.Name = descrFields[len(descrFields)-1]
		db.Folder = strings.Join(descrFields[0:len(descrFields)-1], "/")
	} else {
		db.Folder = db.Server
		db.Name = db.Ref
	}

	return db
}

func dbsToIBases(dbs []*DB) ([]byte, error) {
	const iBaseItem = "[{{.Name}}]\r\nConnect=Srvr=\"{{.Server}}\";Ref=\"{{.Ref}}\";\r\nID={{.ID}}\r\nFolder={{.Folder}}\r\n"

	buf := bytes.NewBuffer(nil)
	t := template.Must(template.New("iBaseItem").Parse(iBaseItem))
	for _, db := range dbs {
		err := t.Execute(buf, db)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}