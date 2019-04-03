package main

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

/*
type TCache struct {
	*memcache.Client
}
*/

func main() {
	MemcachedAddresses := []string{"127.0.0.1:11211"}
	memcacheClient := memcache.New(MemcachedAddresses...)

	tc := &TCache{memcacheClient}

	mkey := "habrposts"
	tc.Delete(mkey)

	rebuild := func() (interface{}, []string, error) {
		habrPosts, err := GetHabrPosts()
		if err != nil {
			return nil, nil, err
		}
		return habrPosts, []string{"habrTag", "geektimes"}, nil
	}

	fmt.Println("\nTGet call #1")
	posts := RSS{}
	err := tc.TGet(mkey, 30, &posts, rebuild)
	fmt.Println("#1", len(posts.Items), "err:", err)

	fmt.Println("\nTGet call #2")
	posts = RSS{}
	err = tc.TGet(mkey, 30, &posts, rebuild)
	fmt.Println("#2", len(posts.Items), "err:", err)

	fmt.Println("\ninc tag habrTag")
	tc.Increment("habrTag", 1)

	go func() {
		// time.Sleep(time.Millisecond)
		fmt.Println("\nTGet call #async")
		posts = RSS{}
		err = tc.TGet(mkey, 30, &posts, rebuild)
		fmt.Println("#async", len(posts.Items), "err:", err)
	}()

	fmt.Println("\nTGet call #3")
	posts = RSS{}
	err = tc.TGet(mkey, 30, &posts, rebuild)
	fmt.Println("#3", len(posts.Items), "err:", err)
}
