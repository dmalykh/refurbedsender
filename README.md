# Sender

## Usage

``` go
package main  
  
import (  
   "context"  
   "github.com/dmalykh/refurbedsender" 
   "github.com/dmalykh/refurbedsender/gate/http"
   "github.com/dmalykh/refurbedsender/queue/list"
   "log"
   )

func main() {
	// Build notifier  
	var queue = list.NewListQueue() 
	gate, err := http.NewHTTPGate(&http.Config{  
	   URL:     `<url for queries>`,    
	})  
	if err != nil {  
		log.Fatal(err)
	}  
	var sender = refurbedsender.NewSender(queue, gate, false)

	// Print errors from channel
	go func() {  
	   for err := range s.Errors() {  
		   log.Println(err.Error())
	   }  
	}()
	
	// Run service with throttling middleware
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := sender.Run(ctx, middleware.WithThrottlingMiddleware(10, 5, 1*time.Second)); err != nil {
			log.Fatal(err)
		}
	}()
	
	// Finally, we can send messages!
	sender.Send(ctx, sender.NewMessage(`Wow! My text here!`)
	sender.Send(ctx, sender.NewMessage([]byte(`And bytes accepted too!`))

	cancel()
}
```