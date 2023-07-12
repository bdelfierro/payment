package postgreSQL

const (
	CheckActiveCart = `select cart_id from carts where status = 'ACTIVE' and user_id = $1`

	CreateCart = `INSERT INTO carts (user_id, createtimestamp) VALUES ($1,CURRENT_TIMESTAMP)`

	Upsert = `INSERT INTO cart_items (cart_id, product_id, quantity, createtimestamp) VALUES ($1,$2,$3,CURRENT_TIMESTAMP) ON CONFLICT (cart_id, product_id) DO UPDATE SET quantity = cart_items.quantity + excluded.quantity`

	DeleteEmptyItems = `delete from cart_items where cart_id = $1 and quantity <= 0`

	DeleteProductFromCart = `delete from cart_items where cart_id = $1 and product_id = $2`

	UpdateCartStatus = `UPDATE carts SET status = $1 WHERE cart_id = $2`

	GetCartDetails = `select c.user_id, c.cart_id, c.status, ci.product_id, ci.quantity, p.name, p.price, p.currency, p.description, p.image_url, ci.quantity * p.price as total_item_price
from carts c inner join cart_items ci on c.cart_id = ci.cart_id inner join products p on p.product_id = ci.product_id where c.cart_id = $1`
)
