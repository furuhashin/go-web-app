package main

import "github.com/pkg/errors"

var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得することができません。")

type Avatar interface {
	GetAvatorURL(c *client) (string, error)
}
