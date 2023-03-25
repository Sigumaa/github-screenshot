package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	url := flag.String("url", "", "GitHub Code URL")
	out := flag.String("out", "out.png", "Output file")
	flag.Parse()

	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*url, *out); err != nil {
		log.Fatal(err)
	}
}

func run(url, out string) error {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.DisableGPU,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	selector := "readme-toc.content"
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(url, selector, &buf)); err != nil {
		return err
	}

	if err := os.WriteFile(out, buf, 0o644); err != nil {
		return err
	}

	return nil
}
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(sel),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
