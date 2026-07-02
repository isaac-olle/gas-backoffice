package admin_assertion_keys

type AdminAssertionKeyPrimitives struct {
	Id        string
	AdminId   string
	Algorithm string
	PublicKey []byte
}

type OwnershipPrimitives struct {
	AdminId string
}

func (o Ownership) ToPrimitives() OwnershipPrimitives {
	return OwnershipPrimitives{AdminId: o.AdminId.String()}
}

func (aak *AdminAssertionKey) ToPrimitives() AdminAssertionKeyPrimitives {
	return AdminAssertionKeyPrimitives{
		Id:        aak.id.String(),
		AdminId:   aak.ownership.AdminId.String(),
		Algorithm: aak.algorithm.String(),
		PublicKey: aak.key.Bytes(),
	}
}
