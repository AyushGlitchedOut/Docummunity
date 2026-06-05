package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/AyushGlitchedOut/Docummunity/consts"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

//We are Using a Token-Bucket Rate-Limiter. So, you can think of the RPS as the rate of filling of tokens in the
//Bucket and Burst as the max amount of tokens in the Bucket. Once the user spends all their tokens, their bucket is empty
//And they cant use anymore tokens and thus, have to wait for their tokens to regenerate

// Struct for a single Client, consisting of last request time and the limiter
type Client struct {
	Limiter     *rate.Limiter
	LastRequest time.Time
}

// NOTE TO SELF: Mutex is a feature of of Go which allows to put locks on global data, which here in this case is the ClientList, so that concurrent requests dont mess up the clientList's consistency
// Declaring the variables for List of All IP adresses, and the Mutex Lock
var (
	clientList = make(map[string]*Client)
	mu         sync.Mutex
)

func getRateLimiterForAddress(IP string) *rate.Limiter {

	//Lock and (for future) Unlock the Mutex
	mu.Lock()
	defer mu.Unlock()

	//Check if Client already exists in the List
	client, exists := clientList[IP]

	//If they dont exist, create a new RateLimiter for them, add them to the List and return the new Limiter
	if !exists {
		newLimiter := rate.NewLimiter(rate.Limit(consts.RequestsPerSecond), consts.Burst)

		clientList[IP] = &Client{
			Limiter:     newLimiter,
			LastRequest: time.Now(),
		}

		return newLimiter
	}

	//If they exist, update the LastRequest time and return the old limiter
	client.LastRequest = time.Now()
	return client.Limiter

}

//TODO: Put this somewhere else neatly

// Client List Cleanup service.
func StartClientListCleanupService() {
	go func() {
		for {

			//The service runs every minute
			time.Sleep(time.Minute)

			//Lock the Mutex
			mu.Lock()

			//Remove all IPs who have not done a single request in the last 10 minutes
			for IP, client := range clientList {
				if time.Since(client.LastRequest) > 10*time.Minute {
					delete(clientList, IP)
				}
			}

			//Unlock Mutex
			mu.Unlock()

		}
	}()

}

//TODO: This is a very simple Rate-Limiter. Consider Improvements in the future to make it safer and efficient

// Middleware for Simple Rate-Limiting
func RateLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//get Client's IP
		clientIP := ctx.ClientIP()

		//get Rate Limiter for obtained IP
		clientLimiter := getRateLimiterForAddress(clientIP)

		//429, If Client has exceeded their limits, return Error
		if !clientLimiter.Allow() {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate Limit exceeded",
			})
			ctx.Abort()
			return
		}

		//Proceed If No errors
		ctx.Next()
	}
}
