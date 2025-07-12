package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thesphereonline/sphere/core/blockchain"
	"github.com/thesphereonline/sphere/core/consensus"
	"github.com/thesphereonline/sphere/core/txpool"
	"github.com/thesphereonline/sphere/core/types"
)

var (
	bc   = blockchain.NewBlockchain()
	pool = txpool.NewPool()
)

func StartServer() {
	r := gin.Default()

	r.GET("/chain", func(c *gin.Context) {
		c.JSON(http.StatusOK, bc.Chain)
	})

	r.POST("/tx", func(c *gin.Context) {
		var tx types.Transaction
		if err := c.ShouldBindJSON(&tx); err != nil {
			c.JSON(400, gin.H{"error": "invalid tx"})
			return
		}
		pool.Add(tx)
		c.JSON(200, gin.H{"status": "tx received"})
	})

	go autoMine()

	r.Run(":8080")
}

func autoMine() {
	for {
		time.Sleep(10 * time.Second)
		txs := pool.Flush()
		validator := consensus.SelectValidator()
		bc.AddBlock(txs, validator)
	}
}
