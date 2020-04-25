package implementation

import "github.com/mikelsr/bspl/proto"

func testProtocol() proto.Protocol {
	buyer := proto.Role("Buyer")
	seller := proto.Role("Seller")
	p := proto.Protocol{
		Name:  "ProtoName",
		Roles: []proto.Role{buyer, seller},
		Params: []proto.Parameter{
			{Name: "ID", Key: true, Io: proto.Out},
			{Name: "item", Io: proto.Out},
			{Name: "price", Io: proto.Out},
		},
		Actions: []proto.Action{
			{Name: "Offer", From: buyer, To: seller, Params: []proto.Parameter{
				{Name: "ID", Key: true, Io: proto.In},
				{Name: "item", Io: proto.In},
				{Name: "price", Io: proto.Out},
			}},
			{Name: "Request", From: buyer, To: seller, Params: []proto.Parameter{
				{Name: "ID", Key: true, Io: proto.Out},
				{Name: "item", Io: proto.Out},
			}},
		},
	}
	return p
}

func testInstance() *Instance {
	p := testProtocol()
	roles := Roles{
		proto.Role("Buyer"):  "B",
		proto.Role("Seller"): "S",
	}
	i := NewInstance(p, roles)
	values := make(Values)
	for _, param := range p.Parameters() {
		values[param.String()] = "X"
	}
	i.values = values
	return i
}
