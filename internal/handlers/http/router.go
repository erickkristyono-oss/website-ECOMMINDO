package httphandler

import (
	"ecommindo/internal/domain"
	"net/http"
)

type Router struct {
	Auth      *AuthHandler
	Service   *ServiceHandler
	Cart      *CartHandler
	Order     *OrderHandler
	JWTSecret []byte
}

func NewRouter(
	authUC domain.AuthUsecase,
	serviceUC domain.ServiceUsecase,
	cartUC domain.CartUsecase,
	orderUC domain.OrderUsecase,
	jwtSecret []byte,
	staticDir string,
) http.Handler {
	auth := NewAuthHandler(authUC)
	service := NewServiceHandler(serviceUC)
	cart := NewCartHandler(cartUC)
	order := NewOrderHandler(orderUC)

	requireAuth := AuthMiddleware(jwtSecret)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/register", auth.Register)
	mux.HandleFunc("POST /api/auth/login", auth.Login)
	mux.Handle("GET /api/auth/me", requireAuth(http.HandlerFunc(auth.Me)))

	mux.HandleFunc("GET /api/services", service.List)

	mux.Handle("GET /api/cart", requireAuth(http.HandlerFunc(cart.Get)))
	mux.Handle("POST /api/cart", requireAuth(http.HandlerFunc(cart.Add)))
	mux.Handle("DELETE /api/cart/{serviceID}", requireAuth(http.HandlerFunc(cart.Remove)))

	mux.Handle("POST /api/checkout", requireAuth(http.HandlerFunc(order.Checkout)))
	mux.Handle("GET /api/orders", requireAuth(http.HandlerFunc(order.List)))

	mux.Handle("/", http.FileServer(http.Dir(staticDir)))

	return CORSMiddleware(mux)
}
