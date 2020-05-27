// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package fiber

import (
	"io/ioutil"
	"net"
	"net/http/httptest"
	"testing"
	"time"

	utils "github.com/gofiber/utils"
	fasthttp "github.com/valyala/fasthttp"
)

func testStatus200(t *testing.T, app *App, url string, method string) {
	req := httptest.NewRequest(method, url, nil)

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}

// func Test_App_Methods(t *testing.T) {

// }

func Test_App_Nested_Params(t *testing.T) {
	app := New()

	app.Get("/test", func(c *Ctx) {
		c.Status(400).Send("Should move on")
	})
	app.Get("/test/:param", func(c *Ctx) {
		c.Status(400).Send("Should move on")
	})
	app.Get("/test/:param/test", func(c *Ctx) {
		c.Status(400).Send("Should move on")
	})
	app.Get("/test/:param/test/:param2", func(c *Ctx) {
		c.Status(200).Send("Good job")
	})

	req := httptest.NewRequest("GET", "/test/john/test/doe", nil)
	resp, err := app.Test(req)

	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}

func Test_App_Use_Params(t *testing.T) {
	app := New()

	app.Use("/prefix/:param", func(c *Ctx) {
		utils.AssertEqual(t, "john", c.Params("param"))
	})

	app.Use("/:param/*", func(c *Ctx) {
		utils.AssertEqual(t, "john", c.Params("param"))
		utils.AssertEqual(t, "doe", c.Params("*"))
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/prefix/john", nil))
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")

	resp, err = app.Test(httptest.NewRequest("GET", "/john/doe", nil))
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}

func Test_App_Use_Params_Group(t *testing.T) {
	app := New()

	group := app.Group("/prefix/:param/*")
	group.Use("/", func(c *Ctx) {
		c.Next()
	})
	group.Get("/test", func(c *Ctx) {
		utils.AssertEqual(t, "john", c.Params("param"))
		utils.AssertEqual(t, "doe", c.Params("*"))
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/prefix/john/doe/test", nil))
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}

func Test_App_Order(t *testing.T) {
	app := New()

	app.Get("/test", func(c *Ctx) {
		c.Write("1")
		c.Next()
	})

	app.All("/test", func(c *Ctx) {
		c.Write("2")
		c.Next()
	})

	app.Use(func(c *Ctx) {
		c.Write("3")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")

	body, err := ioutil.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, "123", string(body))
}
func Test_App_Methods(t *testing.T) {
	var dummyHandler = func(c *Ctx) {}

	app := New()

	app.Connect("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "CONNECT")

	app.Put("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "PUT")

	app.Post("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "POST")

	app.Delete("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "DELETE")

	app.Head("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "HEAD")

	app.Patch("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "PATCH")

	app.Options("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "OPTIONS")

	app.Trace("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "TRACE")

	app.Get("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "GET")

	app.All("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "POST")

	app.Use("/:john?/:doe?", dummyHandler)
	testStatus200(t, app, "/john/doe", "GET")

}

func Test_App_New(t *testing.T) {
	app := New()
	app.Get("/", func(*Ctx) {

	})

	appConfig := New(&Settings{
		Immutable: true,
	})
	appConfig.Get("/", func(*Ctx) {

	})
}

func Test_App_Shutdown(t *testing.T) {
	app := New(&Settings{
		DisableStartupMessage: true,
	})
	_ = app.Shutdown()
}

func Test_App_Static(t *testing.T) {
	app := New()

	grp := app.Group("/v1")

	grp.Static("/v2", ".github/auth_assign.yml")
	app.Static("/*", ".github/FUNDING.yml")
	app.Static("/john", "./.github")

	req := httptest.NewRequest("GET", "/john/stale.yml", nil)
	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	utils.AssertEqual(t, false, resp.Header.Get("Content-Length") == "")

	req = httptest.NewRequest("GET", "/yesyes/john/doe", nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	utils.AssertEqual(t, false, resp.Header.Get("Content-Length") == "")

	req = httptest.NewRequest("GET", "/john/stale.yml", nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	utils.AssertEqual(t, false, resp.Header.Get("Content-Length") == "")

	req = httptest.NewRequest("GET", "/v1/v2", nil)
	resp, err = app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	utils.AssertEqual(t, false, resp.Header.Get("Content-Length") == "")
}

func Test_App_Group(t *testing.T) {
	var dummyHandler = func(c *Ctx) {}

	app := New()

	grp := app.Group("/test")
	grp.Get("/", dummyHandler)
	testStatus200(t, app, "/test", "GET")

	grp.Get("/:demo?", dummyHandler)
	testStatus200(t, app, "/test/john", "GET")

	grp.Connect("/CONNECT", dummyHandler)
	testStatus200(t, app, "/test/CONNECT", "CONNECT")

	grp.Put("/PUT", dummyHandler)
	testStatus200(t, app, "/test/PUT", "PUT")

	grp.Post("/POST", dummyHandler)
	testStatus200(t, app, "/test/POST", "POST")

	grp.Delete("/DELETE", dummyHandler)
	testStatus200(t, app, "/test/DELETE", "DELETE")

	grp.Head("/HEAD", dummyHandler)
	testStatus200(t, app, "/test/HEAD", "HEAD")

	grp.Patch("/PATCH", dummyHandler)
	testStatus200(t, app, "/test/PATCH", "PATCH")

	grp.Options("/OPTIONS", dummyHandler)
	testStatus200(t, app, "/test/OPTIONS", "OPTIONS")

	grp.Trace("/TRACE", dummyHandler)
	testStatus200(t, app, "/test/TRACE", "TRACE")

	grp.All("/ALL", dummyHandler)
	testStatus200(t, app, "/test/ALL", "POST")

	grp.Use("/USE", dummyHandler)
	testStatus200(t, app, "/test/USE/oke", "GET")

	api := grp.Group("/v1")
	api.Post("/", dummyHandler)

	resp, err := app.Test(httptest.NewRequest("POST", "/test/v1/", nil))
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	//utils.AssertEqual(t, "/test/v1", resp.Header.Get("Location"), "Location")

	api.Get("/users", dummyHandler)
	resp, err = app.Test(httptest.NewRequest("GET", "/test/v1/UsErS", nil))
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	//utils.AssertEqual(t, "/test/v1/users", resp.Header.Get("Location"), "Location")
}

func Test_App_Listen(t *testing.T) {
	app := New(&Settings{
		DisableStartupMessage: true,
	})
	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Listen(4003))

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Listen("4010"))
}

func Test_App_Serve(t *testing.T) {
	app := New(&Settings{
		DisableStartupMessage: true,
		Prefork:               true,
	})
	ln, err := net.Listen("tcp4", ":4020")
	utils.AssertEqual(t, nil, err)

	go func() {
		time.Sleep(1000 * time.Millisecond)
		utils.AssertEqual(t, nil, app.Shutdown())
	}()

	utils.AssertEqual(t, nil, app.Serve(ln))
}

// go test -v -run=^$ -bench=Benchmark_App_ETag -benchmem -count=4
func Benchmark_App_ETag(b *testing.B) {
	app := New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)
	c.Send("Hello, World!")
	for n := 0; n < b.N; n++ {
		setETag(c, false)
	}
	utils.AssertEqual(b, `"13-1831710635"`, string(c.Fasthttp.Response.Header.Peek(HeaderETag)))
}

// go test -v -run=^$ -bench=Benchmark_App_ETag_Weak -benchmem -count=4
func Benchmark_App_ETag_Weak(b *testing.B) {
	app := New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)
	c.Send("Hello, World!")
	for n := 0; n < b.N; n++ {
		setETag(c, true)
	}
	utils.AssertEqual(b, `W/"13-1831710635"`, string(c.Fasthttp.Response.Header.Peek(HeaderETag)))
}