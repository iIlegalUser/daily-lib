package test

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Result string

func find(query string) (Result, error) {
	// select {} //模拟故障
	return Result(fmt.Sprintf("result for %q", query)), nil
}

func TestDo(t *testing.T) {
	var g singleflight.Group
	const n = 5
	waited := int32(n)
	done := make(chan struct{})
	key := "https://weibo.com/1227368500/H3GIgngon"
	for i := 0; i < n; i++ {
		go func(j int) {
			v, _, shared := g.Do(key, func() (interface{}, error) {
				ret, err := find(key)
				return ret, err
			})
			if atomic.AddInt32(&waited, -1) == 0 {
				close(done)
			}
			fmt.Printf("index: %d, val: %v, shared: %v\n", j, v, shared)
		}(i)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		fmt.Println("Do hangs")
	}
}

func TestDoChan(t *testing.T) {
	var g singleflight.Group
	var wg sync.WaitGroup
	const n = 5
	key := "https://weibo.com/1227368500/H3GIgngon"
	for i := 0; i < n; i++ {
		go func(j int) {
			wg.Add(1)
			defer wg.Done()
			ch := g.DoChan(key, func() (interface{}, error) {
				ret, err := find(key)
				return ret, err
			})
			timeout := time.After(500 * time.Millisecond)

			var ret singleflight.Result
			select {
			case ret = <-ch:
				fmt.Printf("index: %d, val: %v, shared: %v\n", j, ret.Val, ret.Shared)
			case <-timeout:
				fmt.Printf("%d: timeout\n", j)
				return
			}
		}(i)
	}
	wg.Wait()
}
