package request

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"icrawler/crawlers/configs"
	"log"
	"os"
	"strings"
	"time"
)

type Browser struct {
	targetID         target.ID
	browserContextID cdp.BrowserContextID
	BrowserCtx       context.Context
}

func NewRemoteAllocator(ctx context.Context, url string) *Browser {
	browser := &Browser{}

	allocatorContext, _ := chromedp.NewRemoteAllocator(ctx, url, chromedp.NoModifyURL)

	opts := []chromedp.ContextOption{chromedp.WithNewBrowserContext()}
	browser.BrowserCtx, _ = chromedp.NewContext(allocatorContext, opts...)

	// Start the browser....
	if err := chromedp.Run(browser.BrowserCtx, chromedp.Navigate("about:blank")); err != nil {
		log.Fatalln("failed to connect to chrome.")
	}

	browser.targetID = chromedp.FromContext(browser.BrowserCtx).Target.TargetID
	browser.browserContextID = chromedp.FromContext(browser.BrowserCtx).BrowserContextID

	return browser
}

func NewExecAllocator(ctx context.Context, opts []chromedp.ExecAllocatorOption) *Browser {
	browser := &Browser{}
	ctx, _ = chromedp.NewExecAllocator(context.Background(), opts...)
	browser.BrowserCtx, _ = chromedp.NewContext(ctx) //chromedp.WithDebugf(log.Printf),
	// Start the browser....
	if err := chromedp.Run(browser.BrowserCtx, chromedp.Navigate("about:blank")); err != nil {
		log.Fatalln("failed to connect to chrome.")
	}

	browser.targetID = chromedp.FromContext(browser.BrowserCtx).Target.TargetID
	browser.browserContextID = chromedp.FromContext(browser.BrowserCtx).BrowserContextID

	return browser
}

func (b *Browser) NewTab(ctx context.Context) (context.Context, context.CancelFunc, error) {
	newCtx, cancel := chromedp.NewContext(ctx)

	//attach the tab to the browser.
	if err := chromedp.Run(newCtx); err != nil {
		log.Fatalln("unable to spin up tab.")
	}

	return newCtx, cancel, nil
}

func Response(url, rootXpath string, wb *Browser) string {
	var resultHTML string

	tabCtx, cancel, err := wb.NewTab(wb.BrowserCtx)
	if err != nil {
		log.Println(err)
	}
	defer cancel()

	jsText := GetElementValueByJS(rootXpath)

	if err := chromedp.Run(tabCtx,
		chromedp.Navigate(url),
		chromedp.EvaluateAsDevTools(jsText, &resultHTML),
		chromedp.Sleep(configs.Load().TimeoutCount*time.Second),
	); err != nil {
		fmt.Printf("unable to navigate to page: %v", err)
	}

	//fmt.Println(resultHTML)
	return resultHTML
}

func Navigate(url, rootXpath string, wb *Browser) string {
	var resultHTML string

	tabCtx, cancel, err := wb.NewTab(wb.BrowserCtx)
	if err != nil {
		log.Println(err)
	}
	defer cancel()

	if err := chromedp.Run(tabCtx,
		chromedp.Navigate(url),
		chromedp.OuterHTML(rootXpath, &resultHTML, chromedp.ByJSPath),
		chromedp.Sleep(configs.Load().TimeoutCount*time.Second),
	); err != nil {
		fmt.Printf("unable to navigate to page: %v", err)
	}

	//fmt.Println(resultHTML)
	return resultHTML
}

func GetElementValueByJS(sel string) (js string) {
	const funcJS = `if (sources.length > 0) {
                sources[0].outerHTML;
            } else {
                '未找到标签';
            }`

	invokeFuncJS := `var sources =` + sel + `;`
	return strings.Join([]string{invokeFuncJS, funcJS}, " ")
}

func (b *Browser) CloseTab(ctx context.Context) error {
	// 关闭选项卡
	targets, err := chromedp.Targets(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range targets {
		if t.Title == "about:blank" {
			tabCtx, _ := chromedp.NewContext(ctx, chromedp.WithTargetID(t.TargetID))
			if err := chromedp.Run(tabCtx, page.Close()); err != nil {
				log.Fatal(err)
			}
		}
	}

	// 另外一种关闭选项卡的方法
	//tabCtx, cancel = chromedp.NewContext(ctx, chromedp.WithTargetID(t.TargetID))
	//if err := chromedp.Run(tabCtx); err != nil {
	//	log.Fatal(err)
	//}
	//// call cancel() to close the tab
	//cancel()

	return nil
}

func (b *Browser) CloseBrowser(ctx context.Context) error {
	browserCtx := chromedp.FromContext(ctx).Browser
	if err := browserCtx.Process().Signal(os.Kill); err != nil {
		return err
	}
	return nil
}

func (b *Browser) CleanTemp(ctx context.Context) error {
	//tempDir := chromedp.FromContext(wb.BrowserCtx).Browser.userDataDir
	//if _, err := os.Lstat(tempDir); !os.IsNotExist(err) {
	//	log.Fatalf("temporary user data dir %q not deleted", tempDir)
	//}

	return nil
}
