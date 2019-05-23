package router

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"log"
	"strings"
	"tech-park-db-hw/internal/pkg/controllers"
	"tech-park-db-hw/internal/pkg/middleware"
	"tech-park-db-hw/internal/pkg/middleware/logger"
)

type Route struct {
	Name    string
	Method  string
	Path    string
	Handler routing.Handler
}

type Routes []Route

func NewRouter() *routing.Router {
	router := routing.New()
	log.Println("Available routes:")
	for _, route := range r {
		log.Printf("%s: %v %s", route.Method, route.Path, route.Name)

		var handler routing.Handler
		handler = route.Handler

		router.To(
			route.Method,
			route.Path,
			logger.Logger(
				middleware.ApplyMiddlewares(
					handler,
					middleware.ContentTypeMiddleware),
				route.Name),
		)
	}

	return router
}

var r = Routes{
	Route{
		"Index",
		"GET",
		"/api/",
		controllers.Index,
	},

	Route{
		"Clear",
		strings.ToUpper("Post"),
		"/api/service/clear",
		controllers.Clear,
	},

	Route{
		"ForumCreate",
		strings.ToUpper("Post"),
		"/api/forum/create",
		controllers.ForumCreate,
	},

	Route{
		"ForumGetOne",
		strings.ToUpper("Get"),
		"/api/forum/<slug>/details",
		controllers.ForumGetOne,
	},

	Route{
		"ForumGetThreads",
		strings.ToUpper("Get"),
		"/api/forum/<slug>/threads",
		controllers.ForumGetThreads,
	},

	Route{
		"ForumGetUsers",
		strings.ToUpper("Get"),
		"/api/forum/<slug>/users",
		controllers.ForumGetUsers,
	},

	Route{
		"ThreadCreate",
		strings.ToUpper("Post"),
		"/api/forum/<slug>/create",
		controllers.ThreadCreate,
	},

	Route{
		"PostGetOne",
		strings.ToUpper("Get"),
		"/api/post/<id>/details",
		controllers.PostGetOne,
	},

	Route{
		"PostUpdate",
		strings.ToUpper("Post"),
		"/api/post/<id>/details",
		controllers.PostUpdate,
	},

	Route{
		"PostsCreate",
		strings.ToUpper("Post"),
		"/api/thread/<slug_or_id>/create",
		controllers.PostsCreate,
	},

	Route{
		"Status",
		strings.ToUpper("Get"),
		"/api/service/status",
		controllers.Status,
	},

	Route{
		"ThreadGetOne",
		strings.ToUpper("Get"),
		"/api/thread/<slug_or_id>/details",
		controllers.ThreadGetOne,
	},

	Route{
		"ThreadGetPosts",
		strings.ToUpper("Get"),
		"/api/thread/<slug_or_id>/posts",
		controllers.ThreadGetPosts,
	},

	Route{
		"ThreadUpdate",
		strings.ToUpper("Post"),
		"/api/thread/<slug_or_id>/details",
		controllers.ThreadUpdate,
	},

	Route{
		"ThreadVote",
		strings.ToUpper("Post"),
		"/api/thread/<slug_or_id>/vote",
		controllers.ThreadVote,
	},

	Route{
		"UserCreate",
		strings.ToUpper("Post"),
		"/api/user/<nickname>/create",
		controllers.UserCreate,
	},

	Route{
		"UserGetOne",
		strings.ToUpper("Get"),
		"/api/user/<nickname>/profile",
		controllers.UserGetOne,
	},

	Route{
		"UserUpdate",
		strings.ToUpper("Post"),
		"/api/user/<nickname>/profile",
		controllers.UserUpdate,
	},
}
