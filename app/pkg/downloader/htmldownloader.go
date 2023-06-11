package downloader

import (
	"PokemonCrawler/app/pkg/helpers"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"golang.org/x/net/html"
)

type PageScraper struct {
	pw      *playwright.Playwright
	browser playwright.Browser
}

func NewWebPageScraper() (*PageScraper, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	return &PageScraper{
		pw:      pw,
		browser: browser,
	}, nil
}

func (ps *PageScraper) ScrapeMoveHtml(urls []string, version string) error {
	for i, url := range urls {
		page, err := ps.browser.NewPage()
		if err != nil {
			return err
		}

		if _, err = page.Goto(url, playwright.PageGotoOptions{}); err != nil {
			return err
		}

		content, err := page.Content()
		if err != nil {
			return err
		}

		fileName := helpers.SavePath("html\\move", version, "move.html")
		if len(urls) > 1 {
			dir, file := filepath.Split(fileName)
			fileName = filepath.Join(dir, strconv.Itoa(i+1), file)
		}

		err = helpers.EnsureDirExists(fileName)
		if err != nil {
			return err
		}

		utf8Content, _ := ConvertCharsetToUTF8(content)

		err = ioutil.WriteFile(fileName, []byte(utf8Content), 0644)
		if err != nil {
			return err
		}

		fmt.Printf("HTML saved to %s\n", fileName)

		if err = page.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (ps *PageScraper) ScrapePokemonHtml(urls []string, version string, versionName, searchNo string, index int) error {
	for _, url := range urls {
		fileName, err := ps.scrapeAndSavePage(url, version, versionName, searchNo)
		if err != nil {
			return err
		}

		// 保存したhtmlを読み込む
		fileContent, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatalf("could not read file: %v", err)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(fileContent)))

		// 他のフォルム
		forms := getForms(doc, index)

		for _, form := range forms {
			searchNo := "n" + form
			println("処理中...:" + searchNo + "(" + versionName + ")")

			anotherFormSearchNo := ""
			// シルヴァディのタイプ別のページはないため、シルヴァディ（ノーマル）のページを読み込む
			if index == 773 {
				anotherFormSearchNo = searchNo[:4]
			} else {
				anotherFormSearchNo = searchNo
			}

			url := helpers.CreateSearchURL(versionName, anotherFormSearchNo)

			_, err := ps.scrapeAndSavePage(url, version, versionName, anotherFormSearchNo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ps *PageScraper) scrapeAndSavePage(url, version string, versionName string, searchNo string) (string, error) {
	const (
		maxRetry = 3
	)

	fileName := helpers.SavePath("html\\pokemon", version, searchNo+".html")

	for i := 0; i < maxRetry; i++ {
		page, err := ps.browser.NewPage()
		if err != nil {
			return "", err
		}

		// リダイレクト時の処理
		page.On("Request.RedirectedFrom", func(request playwright.Request) {
			if versionName == "swsh" {
				redirectedURL := request.URL()
				re := regexp.MustCompile(`/(legends_arceus|bdsp)/zukan/(.+)`)
				if re.MatchString(redirectedURL) {
					newURL := re.ReplaceAllString(redirectedURL, "/sm/zukan/$2")
					page.Goto(newURL, playwright.PageGotoOptions{Timeout: playwright.Float(30000)})
				}
			}
		})

		page.On("Request.RedirectedTo", func(request playwright.Request) {
			if versionName == "swsh" {
				redirectedURL := request.URL()
				re := regexp.MustCompile(`/(legends_arceus|bdsp)/zukan/(.+)`)
				if re.MatchString(redirectedURL) {
					newURL := re.ReplaceAllString(redirectedURL, "/sm/zukan/$2")
					page.Goto(newURL, playwright.PageGotoOptions{Timeout: playwright.Float(30000)})
				}
			}
		})

		page.On("Request.Finished", func(request playwright.Request) {
			fmt.Println("Request finished: ", request.URL())
		})

		res, err := page.Goto(url, playwright.PageGotoOptions{Timeout: playwright.Float(30000)})
		print(res)
		if err != nil {
			fmt.Printf("Failed to load page, attempt %d: %v\n", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if versionName == "swsh" {
			redirectedURL := res.URL()
			re := regexp.MustCompile(`/(legends_arceus|bdsp)/zukan/(.+)`)
			if re.MatchString(redirectedURL) {
				newURL := re.ReplaceAllString(redirectedURL, "/sm/zukan/$2")
				page.Goto(newURL, playwright.PageGotoOptions{Timeout: playwright.Float(30000)})
			}
		}

		if versionName == "sv" {
			redirectedURL := res.URL()
			re := regexp.MustCompile(`/(legends_arceus|bdsp)/zukan/(.+)`)
			if re.MatchString(redirectedURL) {
				newURL := re.ReplaceAllString(redirectedURL, "/sm/zukan/$2")
				page.Goto(newURL, playwright.PageGotoOptions{Timeout: playwright.Float(30000)})
			}
		}

		content, err := page.Content()
		if err != nil {
			return "", err
		}

		if strings.Contains(content, "このページは存在しません。") {
			return "", fmt.Errorf("page does not exist")
		}

		err = helpers.EnsureDirExists(fileName)
		if err != nil {
			return "", err
		}

		utf8Content, _ := ConvertCharsetToUTF8(content)

		err = ioutil.WriteFile(fileName, []byte(utf8Content), 0644)
		if err != nil {
			return "", err
		}

		fmt.Printf("HTML saved to %s\n", fileName)

		if err = page.Close(); err != nil {
			return "", err
		}
		break
	}

	return fileName, nil
}

func (ps *PageScraper) Close() error {
	if err := ps.browser.Close(); err != nil {
		return err
	}

	if err := ps.pw.Stop(); err != nil {
		return err
	}

	return nil
}

func getForms(page *goquery.Document, index int) []string {
	if index == 773 {
		return []string{
			"773a", // かくとう
			"773b", // ひこう
			"773c", // どく
			"773d", // じめん
			"773e", // いわ
			"773f", // むし
			"773g", // ゴースト
			"773h", // はがね
			"773i", // ほのお
			"773j", // みず
			"773k", // くさ
			"773l", // でんき
			"773m", // エスパー
			"773n", // こおり
			"773o", // ドラゴン
			"773p", // あく
			"773q", // フェアリー
		}
	}

	var forms []string
	// 他のフォルムがあるか
	page.Find(".select_list:not(.gen_list):first-child").Each(func(index int, s *goquery.Selection) {
		s.Find("li > a").Each(func(index int, s1 *goquery.Selection) {
			text, _ := s1.Attr("href")
			slice := strings.Split(text, "/")
			form := strings.Replace(slice[3], "n", "", 1)
			forms = append(forms, form)
		})
	})
	return forms
}

// htmlのメタタグのcharsetをUTF-8に書き換え、書き換えたhtmlを返します。
func ConvertCharsetToUTF8(content string) (string, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("Error parsing HTML: %w", err)
	}

	traverseHTMLNodes(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			for i, a := range n.Attr {
				if a.Key == "charset" {
					n.Attr[i].Val = "UTF-8"
				}
			}
		}
	})

	// 書き換えたhtmlを文字列に変換
	var b strings.Builder
	err = html.Render(&b, doc)
	if err != nil {
		return "", fmt.Errorf("Error rendering HTML: %w", err)
	}

	return b.String(), nil
}

// htmlのノードにアクセスし、各ノードに対してfuncterを実行します。
func traverseHTMLNodes(n *html.Node, functer func(*html.Node)) {
	if n == nil {
		return
	}
	functer(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseHTMLNodes(c, functer)
	}
}
