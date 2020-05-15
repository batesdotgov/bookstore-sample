package purchasepersister

type ErrPurchaseInvalid struct {
	Message string
}

func (e ErrPurchaseInvalid) Error() string {
	return e.Message
}
