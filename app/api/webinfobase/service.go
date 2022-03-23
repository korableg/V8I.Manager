package webinfobase

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"github.com/google/uuid"
	"text/template"
)

//go:embed webcommoninfobases.wsdl
var wsdlTmplt string

type (
	Service interface {
		WSDL() []byte
		CheckInfoBases(ctx context.Context, clientID uuid.UUID) (CheckInfoBaseResponse, error)
		GetInfoBases(ctx context.Context, clientID uuid.UUID) (GetInfoBasesResponse, error)
	}

	service struct {
		wsdl []byte
	}
)

func NewService(baseURL string) (*service, error) {
	w, err := template.New("wsdl").Parse(wsdlTmplt)
	if err != nil {
		return nil, fmt.Errorf("parse wsdl: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	if err = w.Execute(buf, baseURL); err != nil {
		return nil, fmt.Errorf("execute wsdl template: %w", err)
	}

	s := &service{
		wsdl: buf.Bytes(),
	}

	return s, nil
}

func (s *service) WSDL() []byte {
	return s.wsdl
}

func (s *service) CheckInfoBases(ctx context.Context, clientID uuid.UUID) (CheckInfoBaseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) GetInfoBases(ctx context.Context, clientID uuid.UUID) (GetInfoBasesResponse, error) {
	//TODO implement me
	panic("implement me")
}
