package webinfobase

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/korableg/V8I.Manager/app/internal/httplib"
	"io/ioutil"
	"net/http"
)

type (
	Handlers struct {
		validate *validator.Validate
	}
)

func NewHandlers(validate *validator.Validate) (*Handlers, error) {
	if validate == nil {
		return nil, errors.New("validator is nil")
	}

	h := &Handlers{validate: validate}

	return h, nil
}

func (h *Handlers) Register(r *mux.Router) *mux.Router {
	r.HandleFunc("/WebCommonInfoBases", h.WebCommonInfoBasesHead).Methods("HEAD")
	r.HandleFunc("/WebCommonInfoBases", h.WebCommonInfoBasesGet).Queries("wsdl", "").Methods("GET")
	r.HandleFunc("/WebCommonInfoBases/ws.cws", h.WebCommonInfoBasesPost).Methods("POST")

	return r
}

func (h *Handlers) WebCommonInfoBasesHead(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *Handlers) WebCommonInfoBasesGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(getWSDL())
}

func (h *Handlers) WebCommonInfoBasesPost(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httplib.WriteError(w, r.RequestURI, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = data
}

func getWSDL() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://schemas.xmlsoap.org/wsdl/"
		xmlns:soap12bind="http://schemas.xmlsoap.org/wsdl/soap12/"
		xmlns:soapbind="http://schemas.xmlsoap.org/wsdl/soap/"
		xmlns:tns="https://titovcode.com/WebCommonInfoBases"
		xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy"
		xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
		xmlns:xsd="http://www.w3.org/2001/XMLSchema"
		xmlns:xsd1="https://titovcode.com/WebCommonInfoBases"
		name="WebCommonInfoBases"
		targetNamespace="https://titovcode.com/WebCommonInfoBases">
	<types>
		<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"
				xmlns:xs1="https://titovcode.com/WebCommonInfoBases"
				targetNamespace="https://titovcode.com/WebCommonInfoBases"
				elementFormDefault="qualified">
			<xs:element name="CheckInfoBases">
				<xs:complexType>
					<xs:sequence>
						<xs:element name="ID"
								type="xs:string"
								nillable="true"/>
						<xs:element name="Code"
								type="xs:string"
								nillable="true"/>
					</xs:sequence>
				</xs:complexType>
			</xs:element>
			<xs:element name="CheckInfoBasesResponse">
				<xs:complexType>
					<xs:sequence>
						<xs:element name="return"
								type="xs:string"
								nillable="true"/>
						<xs:element name="Changed"
								type="xs:boolean"
								nillable="true"/>
						<xs:element name="URL"
								type="xs:string"
								nillable="true"/>
					</xs:sequence>
				</xs:complexType>
			</xs:element>
			<xs:element name="GetInfoBases">
				<xs:complexType>
					<xs:sequence>
						<xs:element name="ID"
								type="xs:string"
								nillable="true"/>
					</xs:sequence>
				</xs:complexType>
			</xs:element>
			<xs:element name="GetInfoBasesResponse">
				<xs:complexType>
					<xs:sequence>
						<xs:element name="return"
								type="xs:string"
								nillable="true"/>
						<xs:element name="ID"
								type="xs:string"
								nillable="true"/>
						<xs:element name="Code"
								type="xs:string"
								nillable="true"/>
						<xs:element name="Text"
								type="xs:string"
								nillable="true"/>
					</xs:sequence>
				</xs:complexType>
			</xs:element>
		</xs:schema>
	</types>
	<message name="CheckInfoBasesRequestMessage">
		<part name="parameters"
				element="tns:CheckInfoBases"/>
	</message>
	<message name="CheckInfoBasesResponseMessage">
		<part name="parameters"
				element="tns:CheckInfoBasesResponse"/>
	</message>
	<message name="GetInfoBasesRequestMessage">
		<part name="parameters"
				element="tns:GetInfoBases"/>
	</message>
	<message name="GetInfoBasesResponseMessage">
		<part name="parameters"
				element="tns:GetInfoBasesResponse"/>
	</message>
	<portType name="WebCommonInfoBasesPortType">
		<operation name="CheckInfoBases">
			<input message="tns:CheckInfoBasesRequestMessage"/>
			<output message="tns:CheckInfoBasesResponseMessage"/>
		</operation>
		<operation name="GetInfoBases">
			<input message="tns:GetInfoBasesRequestMessage"/>
			<output message="tns:GetInfoBasesResponseMessage"/>
		</operation>
	</portType>
	<binding name="WebCommonInfoBasesSoapBinding"
			type="tns:WebCommonInfoBasesPortType">
		<soapbind:binding style="document"
				transport="http://schemas.xmlsoap.org/soap/http"/>
		<operation name="CheckInfoBases">
			<soapbind:operation style="document"
					soapAction="https://titovcode.com/WebCommonInfoBases#WebCommonInfoBases:CheckInfoBases"/>
			<input>
				<soapbind:body use="literal"/>
			</input>
			<output>
				<soapbind:body use="literal"/>
			</output>
		</operation>
		<operation name="GetInfoBases">
			<soapbind:operation style="document"
					soapAction="https://titovcode.com/WebCommonInfoBases#WebCommonInfoBases:GetInfoBases"/>
			<input>
				<soapbind:body use="literal"/>
			</input>
			<output>
				<soapbind:body use="literal"/>
			</output>
		</operation>
	</binding>
	<binding name="WebCommonInfoBasesSoap12Binding"
			type="tns:WebCommonInfoBasesPortType">
		<soap12bind:binding style="document"
				transport="http://schemas.xmlsoap.org/soap/http"/>
		<operation name="CheckInfoBases">
			<soap12bind:operation style="document"
					soapAction="https://titovcode.com/WebCommonInfoBases#WebCommonInfoBases:CheckInfoBases"/>
			<input>
				<soap12bind:body use="literal"/>
			</input>
			<output>
				<soap12bind:body use="literal"/>
			</output>
		</operation>
		<operation name="GetInfoBases">
			<soap12bind:operation style="document"
					soapAction="https://titovcode.com/WebCommonInfoBases#WebCommonInfoBases:GetInfoBases"/>
			<input>
				<soap12bind:body use="literal"/>
			</input>
			<output>
				<soap12bind:body use="literal"/>
			</output>
		</operation>
	</binding>
	<service name="WebCommonInfoBases">
		<port name="WebCommonInfoBasesSoap"
				binding="tns:WebCommonInfoBasesSoapBinding">
			<documentation> 
				<wsi:Claim xmlns:wsi="http://ws-i.org/schemas/conformanceClaim/"
						conformsTo="http://ws-i.org/profiles/basic/1.1"/>
			</documentation>
			<soapbind:address location="http://192.168.2.2:8080/WebCommonInfoBases/ws.cws"/>
		</port>
		<port name="WebCommonInfoBasesSoap12"
				binding="tns:WebCommonInfoBasesSoap12Binding">
			<soap12bind:address location="http://192.168.2.2:8080/WebCommonInfoBases/ws.cws"/>
		</port>
	</service>
</definitions>`)
}
