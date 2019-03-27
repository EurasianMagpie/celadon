package api

import "strings"

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"
import "github.com/EurasianMagpie/celadon/ipc"


func RegisterApiRoutes(r *gin.Engine) {
	apisubdomain := r.Group("/celadon")
	apisubdomain.GET("/rinfo", regionInfo)
	apisubdomain.GET("/ginfo", gameInfo)
	//apisubdomain.GET("/pinfo", priceInfo)
	apisubdomain.GET("/prank", priceRank)
	apisubdomain.GET("/gp", gamePrice)
	apisubdomain.GET("/sp", searchPrice)
	apisubdomain.GET("/recommend", queryRecommend)
	apisubdomain.GET("/plist", queryPriceList)
}

func regionInfo(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(200, formResult(301, string("invalid param id"), gin.H{}))
		return
	}

	r, err := db.QueryRegionInfo(id)
	if err == nil {
		d := gin.H{}
		if r != nil {
			d = gin.H{
				"id": r.Region_id,
				"name": r.Name,
				"cname": r.Cname,
			}
		}
		c.JSON(200, formResult(0, "", d))
	} else {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	}
}

func gameInfo(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(200, formResult(301, string("invalid param id"), gin.H{}))
		return
	}

	r, err := db.QueryGameInfo(id)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			d = formGameInfo(c, r)
			c.JSON(200, formResult(0, "", d))
		} else {
			c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
		}
	}
}

func priceInfo(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(200, formResult(301, string("invalid param id"), gin.H{}))
		return
	}

	r, err := db.QueryPriceInfo(id)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			d = gin.H {
				"id": r.Id,
				"discount": r.Discount,
				"price": r.Price,
				"lprice": r.LPrice,
				"lregion": r.LRegion,
				"hprice": r.HPrice,
				"hregion": r.HRegion,
			}
			c.JSON(200, formResult(0, "", d))
		} else {
			c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
		}
	}
}

func priceRank(c *gin.Context) {
	QueryPriceRank(c)
}

func gamePrice(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(200, formResult(301, string("invalid param id"), gin.H{}))
		return
	}

	r, err := db.QueryGamePrice(id)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			d = formGamePrice(c, *r)
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func searchPrice(c *gin.Context) {
	name := c.Query("name")
	if len(name) == 0 {
		c.JSON(200, formResult(301, string("invalid param name"), gin.H{}))
		return
	}
	r, err := db.QuerySearchGamePrice(name)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			var games []gin.H
			var ids []string
			for _, e := range *r {
				games = append(games, formGamePrice(c, e))
				ids = append(ids, e.Id)
			}
			if games != nil {
				d = gin.H {
					"games" : games,
				}
			}
			invokeIpcTask(ids)
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func queryRecommend(c *gin.Context) {
	r, err := db.QueryRecommendGames(20)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			var games []gin.H
			var ids []string
			for _, e := range *r {
				games = append(games, formGamePrice(c, e))
				ids = append(ids, e.Id)
			}
			if games != nil {
				d = gin.H {
					"games" : games,
				}
			}
			invokeIpcTask(ids)
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func queryPriceList(c *gin.Context) {
	ids := c.Query("ids")
	if len(ids) == 0 {
		c.JSON(200, formResult(301, string("invalid param ids"), gin.H{}))
		return
	}
	s := strings.Split(ids, ",")
	r, err := db.QueryPriceListByIds(s)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			var games []gin.H
			var ids []string
			for _, e := range *r {
				games = append(games, formGamePrice(c, e))
				ids = append(ids, e.Id)
			}
			if games != nil {
				d = gin.H {
					"games" : games,
				}
			}
			invokeIpcTask(ids)
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func invokeIpcTask(id []string) {
	go ipc.AddTask(id)
}