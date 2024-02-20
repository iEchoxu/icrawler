package request

import "github.com/chromedp/chromedp"

var (
	Opts = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.IgnoreCertErrors,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		//chromedp.UserDataDir(dir),
		chromedp.WindowSize(1280, 1080),
	)
)
