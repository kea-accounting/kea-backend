// Code generated by protoc-gen-go. DO NOT EDIT.
// source: invoice.proto

package datastructures

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Invoice_Status int32

const (
	Invoice_NEW     Invoice_Status = 0
	Invoice_SENT    Invoice_Status = 1
	Invoice_PAID    Invoice_Status = 2
	Invoice_DELETED Invoice_Status = 100
)

var Invoice_Status_name = map[int32]string{
	0:   "NEW",
	1:   "SENT",
	2:   "PAID",
	100: "DELETED",
}

var Invoice_Status_value = map[string]int32{
	"NEW":     0,
	"SENT":    1,
	"PAID":    2,
	"DELETED": 100,
}

func (x Invoice_Status) String() string {
	return proto.EnumName(Invoice_Status_name, int32(x))
}

func (Invoice_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3b1832ff34ba7c07, []int{0, 0}
}

// [START messages]
type Invoice struct {
	Id                   string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Created              *timestamp.Timestamp `protobuf:"bytes,2,opt,name=created,proto3" json:"created,omitempty"`
	LastUpdated          *timestamp.Timestamp `protobuf:"bytes,3,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	ClientId             string               `protobuf:"bytes,4,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Contact              string               `protobuf:"bytes,5,opt,name=contact,proto3" json:"contact,omitempty"`
	Project              string               `protobuf:"bytes,6,opt,name=project,proto3" json:"project,omitempty"`
	DateInvoice          string               `protobuf:"bytes,7,opt,name=date_invoice,json=dateInvoice,proto3" json:"date_invoice,omitempty"`
	PaymentTerms         int64                `protobuf:"varint,8,opt,name=payment_terms,json=paymentTerms,proto3" json:"payment_terms,omitempty"`
	DateDue              string               `protobuf:"bytes,9,opt,name=date_due,json=dateDue,proto3" json:"date_due,omitempty"`
	Status               Invoice_Status       `protobuf:"varint,10,opt,name=status,proto3,enum=datastructures.Invoice_Status" json:"status,omitempty"`
	Currency             string               `protobuf:"bytes,11,opt,name=currency,proto3" json:"currency,omitempty"`
	NetAmount            string               `protobuf:"bytes,12,opt,name=net_amount,json=netAmount,proto3" json:"net_amount,omitempty"`
	SalesTax             string               `protobuf:"bytes,13,opt,name=sales_tax,json=salesTax,proto3" json:"sales_tax,omitempty"`
	TotalAmount          string               `protobuf:"bytes,14,opt,name=total_amount,json=totalAmount,proto3" json:"total_amount,omitempty"`
	Lineitems            []*Invoice_LineItem  `protobuf:"bytes,15,rep,name=lineitems,proto3" json:"lineitems,omitempty"`
	UserId               string               `protobuf:"bytes,16,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	VatPeriod            string               `protobuf:"bytes,17,opt,name=vat_period,json=vatPeriod,proto3" json:"vat_period,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Invoice) Reset()         { *m = Invoice{} }
func (m *Invoice) String() string { return proto.CompactTextString(m) }
func (*Invoice) ProtoMessage()    {}
func (*Invoice) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b1832ff34ba7c07, []int{0}
}

func (m *Invoice) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Invoice.Unmarshal(m, b)
}
func (m *Invoice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Invoice.Marshal(b, m, deterministic)
}
func (m *Invoice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Invoice.Merge(m, src)
}
func (m *Invoice) XXX_Size() int {
	return xxx_messageInfo_Invoice.Size(m)
}
func (m *Invoice) XXX_DiscardUnknown() {
	xxx_messageInfo_Invoice.DiscardUnknown(m)
}

var xxx_messageInfo_Invoice proto.InternalMessageInfo

func (m *Invoice) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Invoice) GetCreated() *timestamp.Timestamp {
	if m != nil {
		return m.Created
	}
	return nil
}

func (m *Invoice) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

func (m *Invoice) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *Invoice) GetContact() string {
	if m != nil {
		return m.Contact
	}
	return ""
}

func (m *Invoice) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

func (m *Invoice) GetDateInvoice() string {
	if m != nil {
		return m.DateInvoice
	}
	return ""
}

func (m *Invoice) GetPaymentTerms() int64 {
	if m != nil {
		return m.PaymentTerms
	}
	return 0
}

func (m *Invoice) GetDateDue() string {
	if m != nil {
		return m.DateDue
	}
	return ""
}

func (m *Invoice) GetStatus() Invoice_Status {
	if m != nil {
		return m.Status
	}
	return Invoice_NEW
}

func (m *Invoice) GetCurrency() string {
	if m != nil {
		return m.Currency
	}
	return ""
}

func (m *Invoice) GetNetAmount() string {
	if m != nil {
		return m.NetAmount
	}
	return ""
}

func (m *Invoice) GetSalesTax() string {
	if m != nil {
		return m.SalesTax
	}
	return ""
}

func (m *Invoice) GetTotalAmount() string {
	if m != nil {
		return m.TotalAmount
	}
	return ""
}

func (m *Invoice) GetLineitems() []*Invoice_LineItem {
	if m != nil {
		return m.Lineitems
	}
	return nil
}

func (m *Invoice) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *Invoice) GetVatPeriod() string {
	if m != nil {
		return m.VatPeriod
	}
	return ""
}

type Invoice_LineItem struct {
	Quantity             int64    `protobuf:"varint,1,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Description          string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Price                string   `protobuf:"bytes,3,opt,name=price,proto3" json:"price,omitempty"`
	SalesTaxRate         int64    `protobuf:"varint,4,opt,name=sales_tax_rate,json=salesTaxRate,proto3" json:"sales_tax_rate,omitempty"`
	NetAmount            string   `protobuf:"bytes,5,opt,name=net_amount,json=netAmount,proto3" json:"net_amount,omitempty"`
	TaxAmount            string   `protobuf:"bytes,6,opt,name=tax_amount,json=taxAmount,proto3" json:"tax_amount,omitempty"`
	TotalAmount          string   `protobuf:"bytes,7,opt,name=total_amount,json=totalAmount,proto3" json:"total_amount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Invoice_LineItem) Reset()         { *m = Invoice_LineItem{} }
func (m *Invoice_LineItem) String() string { return proto.CompactTextString(m) }
func (*Invoice_LineItem) ProtoMessage()    {}
func (*Invoice_LineItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b1832ff34ba7c07, []int{0, 0}
}

func (m *Invoice_LineItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Invoice_LineItem.Unmarshal(m, b)
}
func (m *Invoice_LineItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Invoice_LineItem.Marshal(b, m, deterministic)
}
func (m *Invoice_LineItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Invoice_LineItem.Merge(m, src)
}
func (m *Invoice_LineItem) XXX_Size() int {
	return xxx_messageInfo_Invoice_LineItem.Size(m)
}
func (m *Invoice_LineItem) XXX_DiscardUnknown() {
	xxx_messageInfo_Invoice_LineItem.DiscardUnknown(m)
}

var xxx_messageInfo_Invoice_LineItem proto.InternalMessageInfo

func (m *Invoice_LineItem) GetQuantity() int64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *Invoice_LineItem) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Invoice_LineItem) GetPrice() string {
	if m != nil {
		return m.Price
	}
	return ""
}

func (m *Invoice_LineItem) GetSalesTaxRate() int64 {
	if m != nil {
		return m.SalesTaxRate
	}
	return 0
}

func (m *Invoice_LineItem) GetNetAmount() string {
	if m != nil {
		return m.NetAmount
	}
	return ""
}

func (m *Invoice_LineItem) GetTaxAmount() string {
	if m != nil {
		return m.TaxAmount
	}
	return ""
}

func (m *Invoice_LineItem) GetTotalAmount() string {
	if m != nil {
		return m.TotalAmount
	}
	return ""
}

func init() {
	proto.RegisterEnum("datastructures.Invoice_Status", Invoice_Status_name, Invoice_Status_value)
	proto.RegisterType((*Invoice)(nil), "datastructures.Invoice")
	proto.RegisterType((*Invoice_LineItem)(nil), "datastructures.Invoice.LineItem")
}

func init() { proto.RegisterFile("invoice.proto", fileDescriptor_3b1832ff34ba7c07) }

var fileDescriptor_3b1832ff34ba7c07 = []byte{
	// 550 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0x5f, 0x6f, 0xd3, 0x3e,
	0x14, 0xfd, 0xa5, 0xd9, 0x9a, 0xf6, 0xa6, 0xeb, 0xaf, 0x58, 0x48, 0x98, 0xa2, 0x41, 0x18, 0x3c,
	0xf4, 0x29, 0x93, 0x0a, 0xe2, 0x0d, 0xa4, 0x49, 0xed, 0x43, 0xa5, 0x69, 0x9a, 0xb2, 0x22, 0x1e,
	0x23, 0x2f, 0xbe, 0x4c, 0x46, 0xcd, 0x1f, 0xec, 0x9b, 0x69, 0xfb, 0x08, 0x7c, 0x50, 0xbe, 0x07,
	0xb2, 0x9d, 0x6c, 0xa2, 0x08, 0xf1, 0x96, 0x7b, 0xce, 0x3d, 0xc7, 0xbe, 0xe7, 0x3a, 0x70, 0xa4,
	0xaa, 0xdb, 0x5a, 0x15, 0x98, 0x36, 0xba, 0xa6, 0x9a, 0x4d, 0xa5, 0x20, 0x61, 0x48, 0xb7, 0x05,
	0xb5, 0x1a, 0xcd, 0xfc, 0xd5, 0x4d, 0x5d, 0xdf, 0xec, 0xf0, 0xd4, 0xb1, 0xd7, 0xed, 0xd7, 0x53,
	0x52, 0x25, 0x1a, 0x12, 0x65, 0xe3, 0x05, 0x27, 0x3f, 0x22, 0x88, 0x36, 0xde, 0x82, 0x4d, 0x61,
	0xa0, 0x24, 0x0f, 0x92, 0x60, 0x31, 0xce, 0x06, 0x4a, 0xb2, 0xf7, 0x10, 0x15, 0x1a, 0x05, 0xa1,
	0xe4, 0x83, 0x24, 0x58, 0xc4, 0xcb, 0x79, 0xea, 0xed, 0xd2, 0xde, 0x2e, 0xdd, 0xf6, 0x76, 0x59,
	0xdf, 0xca, 0x3e, 0xc2, 0x64, 0x27, 0x0c, 0xe5, 0x6d, 0x23, 0x9d, 0x34, 0xfc, 0xa7, 0x34, 0xb6,
	0xfd, 0x9f, 0x7d, 0x3b, 0x7b, 0x01, 0xe3, 0x62, 0xa7, 0xb0, 0xa2, 0x5c, 0x49, 0x7e, 0xe0, 0xee,
	0x32, 0xf2, 0xc0, 0x46, 0x32, 0x0e, 0x51, 0x51, 0x57, 0x24, 0x0a, 0xe2, 0x87, 0x8e, 0xea, 0x4b,
	0xcb, 0x34, 0xba, 0xfe, 0x86, 0x05, 0xf1, 0xa1, 0x67, 0xba, 0x92, 0xbd, 0x86, 0x89, 0x75, 0xce,
	0xbb, 0xa0, 0x78, 0xe4, 0xe8, 0xd8, 0x62, 0xfd, 0xe0, 0x6f, 0xe0, 0xa8, 0x11, 0xf7, 0xa5, 0x3d,
	0x94, 0x50, 0x97, 0x86, 0x8f, 0x92, 0x60, 0x11, 0x66, 0x93, 0x0e, 0xdc, 0x5a, 0x8c, 0x3d, 0x87,
	0x91, 0xf3, 0x91, 0x2d, 0xf2, 0xb1, 0x3f, 0xc2, 0xd6, 0xab, 0x16, 0xd9, 0x07, 0x18, 0x1a, 0x12,
	0xd4, 0x1a, 0x0e, 0x49, 0xb0, 0x98, 0x2e, 0x5f, 0xa6, 0xbf, 0xaf, 0x21, 0xed, 0x0e, 0x4a, 0xaf,
	0x5c, 0x57, 0xd6, 0x75, 0xb3, 0x39, 0x8c, 0x8a, 0x56, 0x6b, 0xac, 0x8a, 0x7b, 0x1e, 0x77, 0xa3,
	0x76, 0x35, 0x3b, 0x06, 0xa8, 0x90, 0x72, 0x51, 0xd6, 0x6d, 0x45, 0x7c, 0xe2, 0xd8, 0x71, 0x85,
	0x74, 0xe6, 0x00, 0x1b, 0x93, 0x11, 0x3b, 0x34, 0x39, 0x89, 0x3b, 0x7e, 0xe4, 0xb5, 0x0e, 0xd8,
	0x8a, 0x3b, 0x3b, 0x32, 0xd5, 0x24, 0x76, 0xbd, 0x7a, 0xea, 0x47, 0x76, 0x58, 0xa7, 0xff, 0x04,
	0xe3, 0x9d, 0xaa, 0x50, 0x11, 0x96, 0x86, 0xff, 0x9f, 0x84, 0x8b, 0x78, 0x99, 0xfc, 0xed, 0xd6,
	0xe7, 0xaa, 0xc2, 0x0d, 0x61, 0x99, 0x3d, 0x4a, 0xd8, 0x33, 0x88, 0x5a, 0x83, 0xda, 0x2e, 0x69,
	0xe6, 0xdc, 0x87, 0xb6, 0xdc, 0x48, 0x7b, 0xef, 0x5b, 0x41, 0x79, 0x83, 0x5a, 0xd5, 0x92, 0x3f,
	0xf1, 0xf7, 0xbe, 0x15, 0x74, 0xe9, 0x80, 0xf9, 0xcf, 0x00, 0x46, 0xbd, 0x9f, 0x9d, 0xff, 0x7b,
	0x2b, 0x2a, 0x52, 0x74, 0xef, 0x9e, 0x5d, 0x98, 0x3d, 0xd4, 0x2c, 0x81, 0x58, 0xa2, 0x29, 0xb4,
	0x6a, 0x48, 0xd5, 0x95, 0x7b, 0x80, 0x76, 0x6b, 0x8f, 0x10, 0x7b, 0x0a, 0x87, 0x8d, 0xb6, 0x1b,
	0x0d, 0x1d, 0xe7, 0x0b, 0xf6, 0x16, 0xa6, 0x0f, 0xc1, 0xe4, 0x5a, 0x10, 0xba, 0x47, 0x14, 0x66,
	0x93, 0x3e, 0x9d, 0x4c, 0x10, 0xee, 0xa5, 0x7b, 0xb8, 0x9f, 0xee, 0x31, 0x80, 0x95, 0x77, 0xb4,
	0x7f, 0x50, 0x63, 0x12, 0x77, 0x1d, 0xbd, 0x9f, 0x6f, 0xf4, 0x47, 0xbe, 0x27, 0x4b, 0x18, 0xfa,
	0x65, 0xb3, 0x08, 0xc2, 0x8b, 0xf5, 0x97, 0xd9, 0x7f, 0x6c, 0x04, 0x07, 0x57, 0xeb, 0x8b, 0xed,
	0x2c, 0xb0, 0x5f, 0x97, 0x67, 0x9b, 0xd5, 0x6c, 0xc0, 0x62, 0x88, 0x56, 0xeb, 0xf3, 0xf5, 0x76,
	0xbd, 0x9a, 0xc9, 0xeb, 0xa1, 0xfb, 0x37, 0xde, 0xfd, 0x0a, 0x00, 0x00, 0xff, 0xff, 0xc5, 0x27,
	0xb2, 0xda, 0xd4, 0x03, 0x00, 0x00,
}
