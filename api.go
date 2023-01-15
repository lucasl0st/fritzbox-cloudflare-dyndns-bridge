package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Api struct {
	users map[string]User
	ctx   context.Context
	cache map[string]string
}

func (api *Api) RegisterEndpoints(r *gin.Engine) {
	r.GET("/dyn", api.fritzDynDnsUpdate)
}

func (api *Api) fritzDynDnsUpdate(c *gin.Context) {
	fritzParams := FritzParams{
		Domain:       c.Query("domain"),
		DualStack:    c.Query("dualstack") == "1",
		Ip6Addr:      c.Query("ip6addr"),
		Ip6LanPrefix: c.Query("ip6lanprefix"),
		IpAddr:       c.Query("ipaddr"),
		Username:     c.Query("username"),
		Password:     c.Query("password"),
	}

	user, found := api.users[fritzParams.Username]

	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "username or password wrong"})
		return
	}

	if user.Password != fritzParams.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "username or password wrong"})
		return
	}

	for _, allowedDomain := range user.AllowedDomains {
		if allowedDomain == fritzParams.Domain {
			noChange, err, v4error, v6error := api.updateDomain(api.ctx, fritzParams, user)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "error while updating record", "error": err.Error()})
				fmt.Printf("could not update dns record %s, error: %s \n", fritzParams.Domain, err.Error())
				return
			}

			if v4error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "error while updating v4 record", "error": v4error.Error()})
				fmt.Printf("could not update dns record %s with ipv4 %s, error: %s \n", fritzParams.Domain, fritzParams.IpAddr, v4error.Error())
				return
			}

			if v6error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "error while updating v6 record", "error": v6error.Error()})
				fmt.Printf("could not update dns record %s with ipv6 %s, error: %s \n", fritzParams.Domain, fritzParams.Ip6Addr, v6error.Error())
				return
			}

			fmt.Printf("updted dns record %s with ipv4 %s and ipv6 %s \n", fritzParams.Domain, fritzParams.IpAddr, fritzParams.Ip6Addr)

			if noChange {
				c.JSON(http.StatusAccepted, gin.H{"message": "domain not updated because no change required"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "updated domain"})
			return
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"message": "user does not have permission for this domain"})
}

func (api *Api) updateDomain(ctx context.Context, params FritzParams, user User) (bool, error, error, error) {
	f := ParseFQDN(params.Domain)

	zone := f.ZoneName + "." + f.TLD
	name := ""

	if api.cache[name+zone+"A"] == params.IpAddr && api.cache[name+zone+"AAAA"] == params.Ip6Addr {
		return true, nil, nil, nil
	}

	zoneId, err := getZoneId(user, zone)

	if err != nil {
		return true, err, nil, nil
	}

	for i := 0; i < len(f.Subdomains); i++ {
		name += f.Subdomains[i] + "."
	}

	var v4error error
	var v6error error

	if len(params.IpAddr) > 0 {
		v4error = updateDomainRecord(ctx, user, zoneId, params.IpAddr, name, zone, "A")
	} else {
		v4error = errors.New("no v4 address provided")
	}

	if v4error == nil {
		api.cache[name+zone+"A"] = params.IpAddr
	}

	if len(params.Ip6Addr) > 0 {
		v6error = updateDomainRecord(ctx, user, zoneId, params.Ip6Addr, name, zone, "AAAA")
	} else {
		v6error = errors.New("no v6 address provided")
	}

	if v6error == nil {
		api.cache[name+zone+"AAAA"] = params.Ip6Addr
	}

	return false, nil, v4error, v6error
}
