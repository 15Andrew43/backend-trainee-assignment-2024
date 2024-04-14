package myerrors

import "fmt"

type BannerExist struct {
	Feature_id int
	Tag_id     int
}

func (b BannerExist) Error() string {
	return fmt.Sprintf("баннер с feature_id = %v и tag_id = %v уже существует", b.Feature_id, b.Tag_id)
}
