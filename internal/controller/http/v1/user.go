package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/s4lat/gosavingsbot/internal/domain"
	"github.com/s4lat/gosavingsbot/internal/usecase"
	"net/http"
)

type UserRoutes struct {
	uc *usecase.UseCase
}

func NewUserRoutes(uc *usecase.UseCase) *UserRoutes {
	return &UserRoutes{
		uc: uc,
	}
}

// Return html page with script that send user init data to server
func (ur *UserRoutes) Index(c echo.Context) error {
	body := `<html>
    <head>
        <script src='https://telegram.org/js/telegram-web-app.js?56'></script>
        <script>
            document.addEventListener('DOMContentLoaded', function () {
                // Create buttons
                const sendButton = document.createElement('button');
                sendButton.innerText = 'Get me';
                document.body.appendChild(sendButton);

                const refreshButton = document.createElement('button');
                refreshButton.innerText = 'Refresh';
                document.body.appendChild(refreshButton);

                const newButton = document.createElement('button');
                newButton.innerText = 'New';
                document.body.appendChild(newButton);

                const responseContainer = document.createElement('div');
                responseContainer.id = 'response-container';
                document.body.appendChild(responseContainer);

				const el = document.createElement('p')
				el.style = "word-wrap: break-word;"
				el.innerText = ` + "`${btoa(Telegram.WebApp.initData)}`\n" + `document.body.appendChild(el)

                // Helper function for making requests
                function makeRequest(method, url) {
                    if (Telegram.WebApp && Telegram.WebApp.initData) {
                        const initData = Telegram.WebApp.initData;
                        const encodedInitData = btoa(initData);

                        fetch(url, {
                            method: method,
                            headers: {
                                'Authorization': encodedInitData
                            }
                        })
                        .then(response => {
                            if (!response.ok) {
                                throw new Error(` + "`Error: ${response.statusText}`" + `);
                            }
                            return response.json();
                        })
                        .then(data => {
                            responseContainer.innerHTML = ` + "`Response: ${JSON.stringify(data)}`" + `;
                        })
                        .catch(error => {
                            responseContainer.innerHTML = ` + "`Error: ${error.message}`" + `;
                        });
                    } else {
                        responseContainer.innerHTML = 'Error: Telegram WebApp is not available.';
                    }
                }

                // Event listeners for buttons
                sendButton.addEventListener('click', function () {
                    makeRequest('GET', '/api/v1/users/me');
                });

                refreshButton.addEventListener('click', function () {
                    makeRequest('GET', '/api/v1/users/me?update=true');
                });

                newButton.addEventListener('click', function () {
                    makeRequest('POST', '/api/v1/users');
                });
            });
        </script>
    </head>
    <body>
        <!-- Content is dynamically created and managed by JavaScript -->
    </body>
</html>`

	return c.HTML(http.StatusOK, body)
}

func (ur *UserRoutes) GetMe(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		l            = getLoggerFromEchoContext(c)
		user         = getServiceUser(c)
		tgWebAppUser = getTgWebAppUser(c)
	)

	update := c.QueryParam("update")
	if update != "" {
		var err error
		user, err = ur.uc.User.UpdateUser(ctx, domain.User{
			Id:       user.Id,
			Username: tgWebAppUser.Username,
		})
		if err != nil {
			l.Errorf("can't update user: %s", err)
			return c.JSON(http.StatusInternalServerError, newErrorResponse("can't update user"))
		}
	}

	return c.JSON(http.StatusOK, user)
}

func (ur *UserRoutes) CreateNewUser(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		l            = getLoggerFromEchoContext(c)
		tgWebAppUser = getTgWebAppUser(c)
	)

	user, err := ur.uc.User.CreateUser(ctx, domain.User{
		Id:       tgWebAppUser.Id,
		Username: tgWebAppUser.Username,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return c.JSON(http.StatusOK, domain.User{Id: tgWebAppUser.Id, Username: tgWebAppUser.Username})
		}

		l.Errorf("can't create user: %s", err)
		return c.JSON(http.StatusInternalServerError, newErrorResponse("can't create user"))
	}

	return c.JSON(http.StatusOK, user)
}
