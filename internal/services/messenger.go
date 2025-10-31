package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"signal/main/internal/models"
	"signal/main/internal/utils"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var ctx context.Context

func ForwardMessengerMessages() {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium-browser"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("start-maximized", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		startUrl := "https://static.xx.fbcdn.net/rsrc.php/ye/r/GI6432-g72t.ico" // messenger icon
		// targetUrl := "https://static.xx.fbcdn.net/rsrc.php/yy/r/XFhtdTsftOC.ogg" // audio message
		// targetUrl := "https://www.messenger.com/api/graphql/" // new message request

		switch e := ev.(type) {
		case *network.EventRequestWillBeSent:
			req := e.Request
			if req.URL == startUrl {
				go onNewMessage()
			}
		}
	})

	Login()
	log.Printf("Login successful, waiting for the messages...")
	SendSignalMessage("Login successful, waiting for the messages...", []string{os.Getenv("MY_NUMBER")})
	select {}
}

func Login() {
	password, pin := utils.GetSecrets()

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.messenger.com/`),
		chromedp.WaitVisible(`[title='Decline optional cookies']`, chromedp.ByQuery),
		chromedp.Click(`[title='Decline optional cookies']`, chromedp.ByQuery),
		chromedp.WaitVisible(`#email`, chromedp.ByQuery),
		chromedp.SendKeys(`#email`, os.Getenv("EMAIL"), chromedp.ByID),
		chromedp.SendKeys(`#pass`, password, chromedp.ByID),
		chromedp.WaitVisible(`#loginbutton`, chromedp.ByQuery),
		chromedp.Evaluate(`document.querySelector('#loginbutton').scrollIntoView({behavior: 'auto', block: 'center'})`, nil),
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`#loginbutton`, chromedp.ByQuery),
		chromedp.WaitVisible(`[aria-label='PIN']`, chromedp.ByQuery),
		chromedp.SendKeys(`[aria-label='PIN']`, pin, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatal("Error logging in: ", err)
	}
}

func onNewMessage() {
	urlToMessageMap := ReadUnreadMessages()

	msg := GetLastUnreadMessage(urlToMessageMap)
	SendSignalMessage(msg, []string{os.Getenv("MY_NUMBER")})
}

func ReadUnreadMessages() map[string]string {
	var link map[string]string
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`
			links = {};
			unreadMessageLinks = Array.from(
			document.querySelectorAll('a[href^="/t"], a[href^="/e2ee"]')
			).filter(
			(a) =>
				a.querySelector('span[data-visualcompletion="ignore"]') !== null &&
				a.ariaLabel == null
			);
			unreadMessageLinks.forEach((link) => {
				let splitted = link.innerText.split("\n");
				let text =
					!link.href.contains("e2ee") && splitted[1] && splitted[1].trim() !== ""
					? ' (' + splitted[2].split(":")[0] + ') - '
					: " - ";
				links[link.href] = splitted[0] + text;
			});
			links;
		`, &link),
	)

	if err != nil {
		models.SaveErrorLog(nil, debug.Stack(), "Error reading unread messages", err, 0)
	}

	return link
}

func GetLastUnreadMessage(urlToMessageMap map[string]string) string {
	var messages []string

	for url, msgStart := range urlToMessageMap {
		var msgEnd string
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			// chromedp.WaitNotVisible(`img[draggable="false"]`, chromedp.ByQuery),
			chromedp.Sleep(1*time.Second),
			chromedp.Evaluate(`
			lastMessage =  Array.from(
				document.querySelectorAll('div[dir="auto"], img[alt="Open photo"]')
			).pop();

			message = '';
			if (!lastMessage) {
				message = '';
			} else if(lastMessage.tagName.toLowerCase() === 'div') {
				let isLink = lastMessage.innerHTML.includes('<a');
				let textRegEx = /<span\b[^>]*>.*?alt="([^"]*)".*?<\/span>/g;
				let linkRegEx = /<span.+ href="([^"]+)" .+>/g;
				message = lastMessage.innerHTML.replaceAll(
					isLink ? linkRegEx : textRegEx,
					'$1'
				);
			} else {
				message = 'Sent a photo';
			}
			message;
		`, &msgEnd),
		)

		if err != nil {
			models.SaveErrorLog([]byte(url), debug.Stack(), "Error getting unread message", err, 0)
			continue
		}

		messages = append(messages, msgStart+msgEnd)
	}

	return strings.Join(messages, "\n")
}

func CheckChromedpHealth() (string, error) {
	if ctx == nil {
		return "Unhealthy", fmt.Errorf("context is nil")
	}

	if err := ctx.Err(); err != nil {
		return "Unhealthy", fmt.Errorf("context is no longer active: %w", err)
	}

	healthCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := chromedp.Run(healthCtx, chromedp.Navigate("https://www.messenger.com/")); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "Unhealthy", fmt.Errorf("browser did not respond in time: %w", err)
		}
		return "Unhealthy", fmt.Errorf("chromedp error: %w", err)
	}

	return "Healthy", nil
}
