package model

import "testing"

func TestFilterNotification(t *testing.T) {
	addrs := []string{
		"bchtest:qqv7g5uj8alewvuv92afngxj4t6j6sx4pczm5f8pkl",
		"bchtest:qpalmy832fp9ytdlx444sehajljnm554dulckcvjl5",
		"bchtest:qz04wg2jj75x34tge2v8w0l6r0repfcvcygv3t7sg5",
	}

	ret, err := FilterNotification(addrs,1541131796,1541131799)
	if err != nil {
		t.Error(err)
	}

	t.Log(ret)

}
