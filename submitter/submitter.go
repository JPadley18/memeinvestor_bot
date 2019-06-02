package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/thecsw/mira"
)

const (
	StickyComment = `
**INVESTMENTS GO HERE - ONLY DIRECT REPLIES TO ME WILL BE PROCESSED**

To prevent thread spam and other natural disasters, I only respond to direct replies. Other commands will be ignored and may be penalized. Let's keep our marketplace clean!

---
- Visit [meme.market](https://meme.market) for help, market statistics, and investor profiles.
- Visit /r/MemeInvestor_bot for questions or suggestions about me.
- Support the project via our [patreon](https://www.patreon.com/memeinvestor_bot)
- New user? Lost or confused? Reply '!help' to this message, or visit the [Wiki](https://www.reddit.com/r/MemeEconomy/wiki/index) for a more in depth explanation.
`

	AttachSource = `
----
Psst, u/%NAME%, you can invoke '!template https://imgur.com/...' command to publicly post your template!
`
)

func main() {
	// Authenticate
	r, _ := mira.Init(mira.ReadCredsFromEnv())
	// Get handler to listen
	c, stop := r.StreamNewPosts("memeinvestor_test")
	var status error = nil
	fmt.Println("source | [time] | thing_id | author | submitter | time elapsed | status")
	go func() {
		for {
			post := <-c
			start := time.Now()
			to_post := StickyComment
			to_post += strings.ReplaceAll(AttachSource, "%NAME%", post.GetAuthor())
			comment, err := r.Comment(post.GetId(), to_post)
			status = err
			status = r.Distinguish(comment.GetId(), "yes", true)
			finish := time.Now()
			// Output the worker log
			fmt.Printf("%v [%v] %v %v %v %v \"%v\"\n",
				"submitter",
				start.Format(time.RFC1123),
				post.GetId(),
				post.GetAuthor(),
				post.GetSubreddit(),
				finish.Sub(start),
				status,
			)
		}
	}()
	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt)
	// Block until received
	<-sigint
	fmt.Printf("Shutting down in %v\n", r.Stream.PostListInterval*time.Second)
	// Stop the streaming
	stop <- true
	// Just to be sure, sleep for a while
	time.Sleep(r.Stream.PostListInterval)
	os.Exit(0)
}
