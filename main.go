package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
)

var (
	hostname, _ = os.Hostname()
	dependency  DependentService
)

func main() {

	port := flag.Int("port", 8080, "Sets the port to listen on")
	depends_on_host := flag.String("depends_on_host", "", "Sets the host of another mesos-tester this mesos-tester depends on")
	depends_on_port := flag.Int("depends_on_port", 0, "Sets the port of another mesos-tester this mesos-tester depends on")
	version := flag.String("version", "0.2", "Set a fake version number for testing purposes")
	cookie_timeout := flag.Int("cookie_timeout", 30, "Sets the cookie timout of the /session endpoint in seconds")

	flag.Parse()

	// grab environment variables, if any are set

	envPort := os.Getenv("MONARCH_PORT")
	envDepHost := os.Getenv("MONARCH_DEPENDS_ON_HOST")
	envDepPort := os.Getenv("MONARCH_DEPENDS_ON_PORT")
	envVersion := os.Getenv("MONARCH_VERSION")
	envTimeout := os.Getenv("MONARCH_COOKIE_TIMEOUT")

	if envPort != "" {
		*port, _ = strconv.Atoi(envPort)
	}

	if envDepHost != "" {
		dependency.Ip = envDepHost
	} else {
		dependency.Ip = *depends_on_host
	}

	if envDepPort != "" {
		dependency.Port, _ = strconv.Atoi(envDepPort)
	} else {
		dependency.Port = *depends_on_port
	}

	if envVersion != "" {
		*version = envVersion
	}

	if envTimeout != "" {
		*cookie_timeout, _ = strconv.Atoi(envTimeout)
	}

	fmt.Println("[magnetio-tester] --> Using dependency: " + dependency.Ip + ":" + strconv.Itoa(dependency.Port))

	// start the REST API
	r := gin.Default()

	// create a session store
	store := sessions.NewCookieStore([]byte("store"))

	store.Options(sessions.Options{MaxAge: *cookie_timeout})
	r.Use(sessions.Sessions("MONARCH_SESSIONID", store))

	// simple ping for health check
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong v2")
	})

	// uses a simple session cookie to check session stickiness
	r.GET("/session", func(c *gin.Context) {
		session := sessions.Default(c)

		var count int
		v := session.Get("count")

		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count += 1
		}

		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})

	})

	// answers after a random wait time
	r.GET("/randomwait", func(c *gin.Context) {

		myrand := random(100, 500)

		time.Sleep(time.Duration(myrand) * time.Millisecond)
		c.String(200, "milliseconds waited: "+strconv.Itoa(myrand))
	})

	// simple host check to identify uniqueness
	r.GET("/host", func(c *gin.Context) {
		c.JSON(200, gin.H{"hostname": "" + hostname + ""})
	})

	// post the monarch.json to this endpoint to generate some work.
	// If the depends_on flag is set, we relay the work to another
	// instance of this application
	r.POST("/work", func(c *gin.Context) {

		var monarchs Monarch

		if dependency.Ip == "" {

			c.Bind(&monarchs)
			randomMonarchIndex := random(0, len(monarchs))

			var response TesterResponse
			response.MonarchCty = monarchs[randomMonarchIndex].Cty
			response.MonarchHse = monarchs[randomMonarchIndex].Hse
			response.MonarchNm = monarchs[randomMonarchIndex].Nm
			response.MonarchYrs = monarchs[randomMonarchIndex].Yrs
			response.setBackendTime()
			response.setEndTime()
			response.addHop(hostname, *version)
			c.JSON(200, response)

		} else {

			c.Bind(&monarchs)
			statusCode, responseBody := dependency.call(monarchs, *version)
			c.JSON(statusCode, responseBody)

		}

	})

	r.Run(":" + strconv.Itoa(*port))
}
