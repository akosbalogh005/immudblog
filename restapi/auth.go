package restapi

import (
	"fmt"
	"immudblog/config"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

const ROLE_READ = "read"
const ROLE_READWRITE = "readwrite"
const ROLE_WRITE = "write"

type userData struct {
	Password string
	Role     string
}

var users map[string]userData

func initUsers() {
	users = make(map[string]userData)
	items := strings.Split(config.ServerFlags.AuthUsers, ",")
	for _, item := range items {
		p := strings.Split(item, ":")
		if len(p) != 3 {
			log.Fatalf("Invalid Auth config: %s  (%s)", config.ServerFlags.AuthUsers, item)
		}
		users[p[0]] = userData{p[1], p[2]}
	}
	log.Debugf("users map Initialized")
}

func checkAuthorization(perms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		basicAuth(c, perms...)
	}
}

func basicAuth(c *gin.Context, perms ...string) {
	// Get the Basic Authentication credentials
	if len(users) == 0 {
		initUsers()
	}

	user, password, hasAuth := c.Request.BasicAuth()
	userData, hasUser := users[user]
	roleOk := false
	if hasUser {
		for _, perm := range perms {
			if perm == userData.Role {
				roleOk = true
			}
		}
	}
	if hasAuth && hasUser && roleOk && userData.Password == password {
		log.Debugf("User %s is autenticated [%s]", user, userData.Role)
		c.Set("username", user)
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("not authorized"))
		return
	}
}
