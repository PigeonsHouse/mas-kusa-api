package utils

import (
	"context"
	"sync"
	"time"

	"github.com/mattn/go-mastodon"
)

func CheckAccount(instance string, token string) (err error) {
	c := mastodon.NewClient(&mastodon.Config{
		Server:      instance,
		AccessToken: token,
	})
	_, err = c.GetAccountCurrentUser(context.Background())
	return
}

func GetUserName(instance string, token string) (acct string, err error) {
	var u *mastodon.Account
	c := mastodon.NewClient(&mastodon.Config{
		Server:      instance,
		AccessToken: token,
	})
	if u, err = c.GetAccountCurrentUser(context.Background()); err == nil {
		acct = u.Acct
	}
	return
}

func CountToot(instance string, token string, thisMonth bool) (baseDate time.Time, tootNumList []int, err error) {
	var (
		wg        sync.WaitGroup
		myAccount *mastodon.Account
		maxId     mastodon.ID
	)
	// mastodonの5分間のリクエスト上限
	const maxRequest = 300
	// 今が何月かの情報などを取得するため
	now := time.Now()

	// mastodonAPIを叩くためのインスタンス
	client := mastodon.NewClient(&mastodon.Config{
		Server:      instance,
		AccessToken: token,
	})

	// アカウントID取得に必要
	myAccount, err = client.GetAccountCurrentUser(context.Background())
	if err != nil {
		return
	}

	baseDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if !thisMonth {
		baseDate = baseDate.AddDate(0, -1, 0)
	}
	nextDate := baseDate.AddDate(0, 1, 0)
	// 調査する月の日数の保持
	MaxDay := nextDate.AddDate(0, 0, -1).Day()

	tootNumList = make([]int, MaxDay)
	// baseDate以降のトゥートをかき集める
	for i := 0; i < maxRequest-2; i++ {
		if toots, err := client.GetAccountStatuses(context.Background(), myAccount.ID, &mastodon.Pagination{MaxID: maxId, Limit: 40}); err == nil {
			last := toots[len(toots)-1]
			maxId = last.ID
			if last.CreatedAt.After(nextDate) {
				continue
			}
			wg.Add(1)
			go func(bd time.Time, nd time.Time, ts []*mastodon.Status, c *[]int) {
				defer wg.Done()
				routineCountUpToot(bd, nd, ts, c)
			}(baseDate, nextDate, toots, &tootNumList)
			if last.CreatedAt.Before(baseDate) {
				break
			}
		}
	}
	wg.Wait()

	return
}

func routineCountUpToot(baseDate time.Time, nextDate time.Time, toots []*mastodon.Status, counter *[]int) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	for _, toot := range toots {
		if toot.CreatedAt.After(baseDate) && toot.CreatedAt.Before(nextDate) {
			(*counter)[toot.CreatedAt.In(jst).Day()-1] += 1
		}
	}
}
