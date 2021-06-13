package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLはErrNoAvatarURLを返すべき")
	}
	testUrl := "http://url-to-avatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("値が存在する場合、AuthAvatar.GetAvatarURLはErrNoAvatarURLを返すべきではない")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURLは正しいURLを返すべき")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Errorf(fmt.Sprintf("%s", err))
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではない")
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURLが%sという誤った値を返した", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	var permAllUsersCanAccess = 0777
	ioutil.WriteFile(filename, []byte{}, fs.FileMode(permAllUsersCanAccess))
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURLはエラーを返すべきではない")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURLが%sという誤った値を返した", url)
	}
}
