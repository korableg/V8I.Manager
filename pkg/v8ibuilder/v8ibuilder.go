package v8ibuilder

import (
	"bytes"
	"text/template"

	"github.com/korableg/V8I.Manager/pkg/clusterdb"
)

func Build(clusterDBS ...*clusterdb.ClusterDB) ([]byte, error) {

	iBaseTemplate := getiBaseTemplate()

	buf := bytes.NewBuffer(nil)

	for _, db := range clusterDBS {
		err := iBaseTemplate.Execute(buf, db)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func getiBaseTemplate() *template.Template {
	const iBaseTemplate = "[{{.Name}}]\r\nConnect=Srvr=\"{{.Server}}\";Ref=\"{{.Ref}}\";\r\nID={{.ID}}\r\nFolder={{.Folder}}\r\n"

	return template.Must(template.New("iBaseItem").Parse(iBaseTemplate))

}
