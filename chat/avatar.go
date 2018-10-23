package main

import (
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"strings"
)

var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得することができません。")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			//mはhash.Hash型でかつio.Writerを実装しているためWriteString()の第一引数になることができる
			io.WriteString(m, strings.ToLower(emailStr))
			//sumを使用するとその時点までに書着込まれた文字列を使用してハッシュ値が計算される
			return fmt.Sprintf("//wwww.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}
