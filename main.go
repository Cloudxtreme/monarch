package main

import (
	"flag"
	"fmt"
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
	version := os.Getenv("MONARCH_VERSION")

	if(version == "") {
		version = "0.1"
	}


	port := flag.Int("port", 8080, "Sets the port to listen on")
	depends_on_host := flag.String("depends_on_host", "", "Sets the host of another mesos-tester this mesos-tester depends on")
	depends_on_port := flag.Int("depends_on_port", 0, "Sets the port of another mesos-tester this mesos-tester depends on")


	flag.Parse()

	// grab environment variables, if any are set

	envPort := os.Getenv("MONARCH_PORT")
	envDepHost := os.Getenv("MONARCH_DEPENDS_ON_HOST")
	envDepPort := os.Getenv("MONARCH_DEPENDS_ON_PORT")

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

	fmt.Println("[magnetio-tester] --> Using dependency: " + dependency.Ip + ":" + strconv.Itoa(dependency.Port))

	// start the REST API

	r := gin.Default()
	// simple ping for health check
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong v2")
	})

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
			response.addHop(hostname, version)
			c.JSON(200, response)

		} else {

			c.Bind(&monarchs)
			statusCode, responseBody := dependency.call(monarchs, version)
			c.JSON(statusCode, responseBody)

		}

	})

	r.Run(":" + strconv.Itoa(*port))
}
