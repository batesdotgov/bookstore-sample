package purchaseshistory

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"text/template"

	"github.com/diegoholiveira/bookstore-sample/users"
)

type (
	PurchasesRepository struct {
		db *sql.DB
	}
)

var (
	queryBooksByPurchases = template.Must(template.New("query").Parse(`
		SELECT p.purchase_id, b.title, b.author, b.price, p.quantity
		FROM purchased_books AS p
		INNER JOIN books AS b ON p.book_id = b.id
		WHERE p.purchase_id IN ({{ . }})
	`))
)

func NewPurchasesRepository(db *sql.DB) PurchasesRepository {
	return PurchasesRepository{
		db: db,
	}
}

func (r PurchasesRepository) FindPurchasesByUser(ctx context.Context, user users.User) (users.Purchases, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, amount FROM purchases WHERE user_id = ?",
		user.ID,
	)
	if err != nil {
		return nil, err
	}

	mappedPurchases := make(map[uint64]users.Purchase)
	ids := make([]uint64, 0)

	for rows.Next() {
		p := users.Purchase{}
		p.Books = make(users.Books, 0)
		err := rows.Scan(
			&p.ID,
			&p.Amount,
		)
		if err != nil {
			return nil, err
		}
		mappedPurchases[p.ID] = p
		ids = append(ids, p.ID)
	}

	books, err := r.findPurchasedBooks(ctx, ids)
	if err != nil {
		return nil, err
	}

	for purchaseID, books := range books {
		purchase := mappedPurchases[purchaseID]
		purchase.Books = books
		mappedPurchases[purchaseID] = purchase
	}

	return toPurchasesSlice(mappedPurchases), nil
}

func stringify(ids []uint64) string {
	stringIDs := []string{}
	for _, id := range ids {
		stringIDs = append(stringIDs, strconv.Itoa(int(id)))
	}
	return strings.Join(stringIDs, ",")
}

func (r PurchasesRepository) findPurchasedBooks(ctx context.Context, ids []uint64) (map[uint64]users.Books, error) {
	var (
		buf strings.Builder
		err error
	)

	err = queryBooksByPurchases.Execute(&buf, stringify(ids))
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, buf.String())
	if err != nil {
		return nil, err
	}

	books := make(map[uint64]users.Books)

	for rows.Next() {
		var (
			book       users.Book
			purchaseID uint64
		)

		err := rows.Scan(
			&purchaseID,
			&book.Title,
			&book.Author,
			&book.Price,
			&book.Quantity,
		)
		if err != nil {
			return nil, err
		}

		books[purchaseID] = append(books[purchaseID], book)
	}

	return books, nil
}

func toPurchasesSlice(h map[uint64]users.Purchase) users.Purchases {
	purchases := make(users.Purchases, 0)
	for _, p := range h {
		purchases = append(purchases, p)
	}
	return purchases
}
