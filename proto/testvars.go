package proto

func testProtocol() Protocol {
	buyer := Role("Buyer")
	seller := Role("Seller")
	p := Protocol{
		Name:  "Noncircular",
		Roles: []Role{buyer, seller},
		Params: []Parameter{
			{Name: "ID", Key: true, Io: Out},
			{Name: "item", Io: Out},
			{Name: "price", Io: Out},
		},
		Actions: []Action{
			{Name: "Request", From: buyer, To: seller, Params: []Parameter{
				{Name: "ID", Key: true, Io: Out},
				{Name: "item", Io: Out},
			}},
			{Name: "Offer", From: buyer, To: seller, Params: []Parameter{
				{Name: "ID", Key: true, Io: In},
				{Name: "item", Io: In},
				{Name: "price", Io: Out},
			}},
		},
	}
	return p
}