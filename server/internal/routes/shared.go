package routes

import "github.com/labstack/echo/v4"

func CurrentUser(c echo.Context) UserSession {
	u, ok := c.Get("user").(UserSession)
	if !ok {
		return UserSession{}
	}
	return u
}
