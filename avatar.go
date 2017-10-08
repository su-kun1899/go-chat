package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
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

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}

	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}

	return "", ErrNoAvatarURL
}
