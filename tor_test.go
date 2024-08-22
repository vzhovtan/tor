package tor_test

import (
	"github/vzhovtan/tor"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var TestBase = tor.Base{
	Mgmt_lo1_id: "1.1.1.1",
	Mgmt_lo2_id: "2.2.2.2",
	Mgmt_lo3_id: "3.3.3.3",
	Mgmt_iso_id: "ISO-ID_TEST",
	Asn_id:      "6530_TEST",
	Rd_id:       "123:456",
}

var TestInterface = tor.Intf{
	Int_type:         "uplink",
	Int_id:           "111",
	Int_desc:         "TEST_DESCRIPTION",
	Rfc_address:      "192.168.1.1",
	Ipv6_rfc_address: "2001:11:11",
}

var TestL3 = tor.L3info{
	Neighbor_id:   "TEST_NE_ID",
	Neighbor_desc: "TEST_NE_DESCR",
}

var testDeviceResult = tor.Device{
	Device_id:        "TEST_DEVICE_ID",
	Management_ipv6:  "2001:22:22:22",
	Management_ip:    "10.10.10.10",
	Network_password: "PW TACACS_KEY",
	Tacacs_key:       "PW TACACS_KEY",
	Bootstrap:        false,
	Snmp_community:   "TEST_COMMUNITY",
	Base:             TestBase,
	Interfaces: []tor.Intf{
		TestInterface,
	},
	L3information: []tor.L3info{
		TestL3,
	},
}

var testFullResult = `conf t
  no ipv6 access-list custom-copp-acl-tacacs6
  ipv6 access-list custom-copp-acl-tacacs6
    permit tcp any any eq 2832
    permit tcp any eq 2832 any
  class-map type control-plane match-any custom-copp-class-management
   no match access-group name custom-copp-acl-tacacs
   no match access-group name custom-copp-acl-radius
   no match access-group name custom-copp-acl-radius6
   match access-group name custom-copp-acl-tacacs6
 exit
 control-plane
     service-policy input custom-copp-policy-strict

hostname TEST_DEVICE_ID
cli alias name wr copy running-config startup-config
system vlan long-name
vdc TEST_DEVICE_ID id 1
  limit-resource vlan minimum 16 maximum 4094
  limit-resource vrf minimum 2 maximum 4096
  limit-resource port-channel minimum 0 maximum 511
  limit-resource m4route-mem minimum 58 maximum 58
  limit-resource m6route-mem minimum 8 maximum 8

feature bash-shell
feature scp-server
feature tacacs+
cfs eth distribute
nv overlay evpn
feature bgp
feature isis
feature pbr
feature interface-vlan
feature vn-segment-vlan-based
feature lacp
feature dhcp
feature lldp
feature bfd
feature nv overlay
feature ngoam

clock timezone PST -8 0
clock summer-time PDT 2 Sun Mar 02:00 1 Sun Nov 02:00 60

archive
  path bootflash:startup-config
  time-period 0
  maximum 14

ngoam install acl

ip domain-lookup
ip domain-name net.google.com
ip name-server 2001:AAAA:BBBB::8888 2001:AAAA:BBBB::8844


system default switchport
system jumbomtu 9000
ip ping source-interface loopback0
ip traceroute source-interface loopback0
ip dns source-interface loopback0
ip ssh source-interface loopback0

rmon event 1 description FATAL(1) owner PMON@FATAL
rmon event 2 description CRITICAL(2) owner PMON@CRITICAL
rmon event 3 description ERROR(3) owner PMON@ERROR
rmon event 4 description WARNING(4) owner PMON@WARNING
rmon event 5 description INFORMATION(5) owner PMON@INFO
snmp-server context default vrf default
snmp-server counter cache timeout 10
snmp-server packetsize 1400
snmp-server globalEnforcePriv
snmp-server user netops network-operator auth sha TEST_COMMUNITY priv aes-128 TEST_COMMUNITY


ntp server AAAA:BBBB:CCCC::1 use-vrf default
ntp server 2001:AAAA:BBBB::4 use-vrf default
ntp server 2001:BBBB:CCCC::8 use-vrf default
ntp server 2001:CCCC:DDDD::C use-vrf default

ntp source-interface loopback0

fabric forwarding anycast-gateway-mac 0000.5100.5300
ip igmp snooping vxlan
system vlan nve-overlay id 1100-1999

username admin password 0 PW role network-admin
username netops password 0 PW role network-admin

tacacs-server key 0 TACACS_KEY
ip tacacs source-interface loopback0
tacacs-server timeout 4

tacacs-server host AAAA:BBBB:CCCC::1 port 2832
tacacs-server host AAAA:BBBB:CCCC::2 port 2832
tacacs-server host AAAA:BBBB:CCCC::3 port 2832
aaa group server tacacs+ TacServer
    server AAAA:BBBB:CCCC::1
    server AAAA:BBBB:CCCC::2
    server AAAA:BBBB:CCCC::3

    source-interface loopback0

aaa authentication login default group TacServer local
aaa authentication login console local
aaa authorization config-commands default group TacServer local

aaa authorization commands default group TacServer

aaa accounting default group TacServer
aaa authentication login invalid-username-log
aaa authentication login error-enable

route-map RM-PERMIT permit 10
route-map RM-REDIST-DIRECT permit 10
  set community 65203:5 additive
route-map RM-REDIST-STATIC permit 10
  set community 65203:5 additive

service dhcp
ip dhcp relay
ipv6 dhcp relay

vrf context management
system nve peer-vni-counter



interface Ethernet111
description TEST_DESCRIPTION
no switchport
mtu 9216
no ip redirects
ip address 192.168.1.1
ipv6 address 2001:11:11
no ipv6 redirects
isis network point-to-point
ip router isis underlay
ipv6 router isis underlay
no isis passive-interface level-2
no shutdown



interface nve1
description NVE
no shutdown
host-reachability protocol bgp
source-interface loopback2
global ingress-replication protocol bgp

interface loopback0
description management
ip address 10.10.10.10
ip router isis underlay

ipv6 address 2001:22:22:22
ipv6 router isis underlay


interface loopback1
description IGP/BGP
ip address 1.1.1.1
ip router isis underlay

interface loopback2
description VTEP
ip address 2.2.2.2
ip router isis underlay

interface loopback3
description MCAST
ip address 3.3.3.3
ip router isis underlay

line console
  exec-timeout 180
line vty

router isis underlay
  net ISO-ID_TEST
  is-type level-2
  address-family ipv4 unicast
    bfd
  address-family ipv6 unicast
    bfd
    multi-topology
  passive-interface default level-1-2

router bgp 6530_TEST
router-id 123:456
log-neighbor-changes
  address-family l2vpn evpn
    maximum-paths ibgp 8
  template peer FABRIC-SPINE
    bfd
    remote-as 6530_TEST
    update-source loopback1
    address-family l2vpn evpn
      send-community
      send-community extended

  neighbor TEST_NE_ID
    description TEST_NE_DESCR
    inherit peer FABRIC-SPINE`

func TestBuildDevice(t *testing.T) {
	dev := tor.BuildDevice()
	if !cmp.Equal(dev, testDeviceResult) {
		t.Errorf("test BuildDevice Failed - error")
	}
}

// func TestBuildConfig(t *testing.T) {
// 	out := new(bytes.Buffer)
// 	deviceResult := tor.BuildConfig(out, dev)
// 	if err != nil {
// 		t.Errorf("test TreeDir Failed - error")
// 	}
// 	result := out.String()
// 	if result != testDirResult {
// 		t.Errorf("test TreeDir Failed - results not match\nGot:\n%v\nExpected:\n%v\n", result, testDirResult)
// 	}
// }
