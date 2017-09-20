package main

import (
	"errors"
)

// ErrNoAvatarURL はAvatarインスタンスがアバターのErrNoAvatarURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

// Avatar はユーザーのプロフィール画像を表す型
type Avatar interface {
	// GetAvatarURLは指定されたクライアントのアバターのURLを返す
	// 問題が発生した場合はエラーを返す
	// URLが取得できなかった場合には ErrNoAvatarURL を返します
	GetAvatarURL(c *client) (string, error)
}
