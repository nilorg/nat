package model

import (
	"gorm.io/gorm"
)

const (
	_TCP    = "TCP"
	_UDP    = "UDP"
	_IPv4   = "IPv4"
	_IPv6   = "IPv6"
	_Domain = "Domain"
)

type PortForwardType string

const PortForwardTypeTCP PortForwardType = _TCP
const PortForwardTypeUDP PortForwardType = _UDP

type LocalAddressType string

const LocalAddressTypeIPv4 LocalAddressType = _IPv4
const LocalAddressTypeIPv6 LocalAddressType = _IPv6

type RemoteAddressType string

const RemoteAddressTypeIPv4 RemoteAddressType = _IPv4
const RemoteAddressTypeIPv6 RemoteAddressType = _IPv6
const RemoteAddressTypeDomain RemoteAddressType = _Domain

type PortForward struct {
	gorm.Model
	Type             PortForwardType   `gorm:"column:type;type:char(3);not null;default:'';comment:'类型'" json:"type"`
	LocalPort        int               `gorm:"column:local_port;type:int(5);not null;default:0;comment:'本地端口'" json:"local_port"`
	LocalAddress     string            `gorm:"column:local_address;type:varchar(255);not null;default:'';comment:'本地地址'" json:"local_address"`
	LOcalAddressType LocalAddressType  `gorm:"column:local_address_type;type:char(6);not null;default:'';comment:'本地地址类型'" json:"local_address_type"`
	RemoteAddress    string            `gorm:"column:remote_address;type:varchar(255);not null;default:'';comment:'远程地址'" json:"remote_address"`
	RemoteAddressTyp RemoteAddressType `gorm:"column:remote_address_type;type:char(6);not null;default:'';comment:'远程地址类型'" json:"remote_address_type"`
	RemotePort       int               `gorm:"column:remote_port;type:int(5);not null;default:0;comment:'远程端口'" json:"remote_port"`
	Timing           int64             `gorm:"column:timing;type:bigint;not null;default:30000;comment:'定时时间'" json:"timing"`
}

type Domain struct {
	gorm.Model
	Domain string `gorm:"column:domain;type:varchar(255);not null;default:'';comment:'域名'" json:"domain"`
	IP     string `gorm:"column:ip;type:varchar(255);not null;default:'';comment:'IP'" json:"ip"`
}

type DomainHistory struct {
	gorm.Model
	DomainID uint   `gorm:"column:domain_id;type:int(10) unsigned;not null;default:0;comment:'域名ID'" json:"domain_id"`
	OldIP    string `gorm:"column:old_ip;type:varchar(255);not null;default:'';comment:'旧IP'" json:"old_ip"`
	NewIP    string `gorm:"column:new_ip;type:varchar(255);not null;default:'';comment:'新IP'" json:"new_ip"`
}
