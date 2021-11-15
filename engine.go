package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func Run(config Config) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("window-size", fmt.Sprintf("%d,%d", 1920, 1080)),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	if config.Login != nil {
		runScript(ctx, config.Login)
	}

	if config.Playback != nil {
		runScript(ctx, config.Playback)
	}
}

func runScript(ctx context.Context, config *ScriptConfig) {
	script := LoadScript(config)

	if err := chromedp.Run(ctx, chromedp.Navigate(script.StartUrl)); err != nil {
		panic(err)
	}
	time.Sleep(time.Duration(config.ActionDelayTime) * time.Second)

	for _, action := range script.Actions {
		targetCtx := ctx
		targets, err := chromedp.Targets(ctx)
		if err != nil {
			panic(err)
		}
		for _, target := range targets {
			if target.URL == action.TabUrl {
				targetCtx, _ = chromedp.NewContext(ctx, chromedp.WithTargetID(target.TargetID))
				break
			}
		}

		switch action.Type {
		case "click":
			if err := chromedp.Run(targetCtx, chromedp.Click(action.XPath, chromedp.NodeReady, chromedp.BySearch)); err != nil {
				panic(err)
			}
		case "text":
			if err := chromedp.Run(targetCtx, chromedp.SendKeys(action.XPath, action.Value, chromedp.NodeReady, chromedp.BySearch)); err != nil {
				panic(err)
			}
		}

		time.Sleep(time.Duration(config.ActionDelayTime) * time.Second)
	}
}
