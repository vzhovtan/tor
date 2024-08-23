package tor

import (
	"io"
	"log"
	"os"
	"text/template"
)

type Intf struct {
	Int_type         string
	Int_id           string
	Int_desc         string
	Rfc_address      string
	Ipv6_rfc_address string
}

type Base struct {
	Mgmt_lo1_id string
	Mgmt_lo2_id string
	Mgmt_lo3_id string
	Mgmt_iso_id string
	Asn_id      string
	Rd_id       string
}

type L3info struct {
	Neighbor_id   string
	Neighbor_desc string
}

type Device struct {
	Device_id        string
	Management_ipv6  string
	Management_ip    string
	Network_password string
	Tacacs_key       string
	Bootstrap        bool
	Snmp_community   string
	Base             Base
	Interfaces       []Intf
	L3information    []L3info
}

func Main() {
	dev := BuildDevice()
	err := BuildConfig(os.Stdout, dev)
	if err != nil {
		log.Fatal(err)
	}
}

func BuildDevice() *Device {
	intf := Intf{
		Int_type:         "uplink",
		Int_id:           "111",
		Int_desc:         "TEST_DESCRIPTION",
		Rfc_address:      "192.168.1.1",
		Ipv6_rfc_address: "2001:11:11",
	}
	base := Base{
		Mgmt_lo1_id: "1.1.1.1",
		Mgmt_lo2_id: "2.2.2.2",
		Mgmt_lo3_id: "3.3.3.3",
		Mgmt_iso_id: "ISO-ID_TEST",
		Asn_id:      "6530_TEST",
		Rd_id:       "123:456",
	}
	l3 := L3info{
		Neighbor_id:   "TEST_NE_ID",
		Neighbor_desc: "TEST_NE_DESCR",
	}
	return &Device{
		Device_id:        "TEST_DEVICE_ID",
		Management_ipv6:  "2001:22:22:22",
		Management_ip:    "10.10.10.10",
		Network_password: "PW",
		Tacacs_key:       "TACACS_KEY",
		Snmp_community:   "TEST_COMMUNITY",
		Bootstrap:        false,
		Base:             base,
		Interfaces: []Intf{
			intf,
		},
		L3information: []L3info{
			l3,
		},
	}
}

func BuildConfig(out io.Writer, dev *Device) error {
	t, err := template.ParseFiles("conf.tmpl")
	if err != nil {
		return err
	}
	err = t.Execute(out, dev)
	if err != nil {
		return err
	}
	return nil
}
