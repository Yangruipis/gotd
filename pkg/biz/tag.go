package biz

import "github.com/Yangruipis/gotd/pkg/core"

func (b *Biz) DeleteTag(tagId uint32) error {
	return b.tagDao.Delete(b.ctx, tagId)
}

func (b *Biz) CreateTaskTag(tagName string, taskId uint32) (*core.Tag, error) {
	return b.tagDao.CreateTaskTag(b.ctx, &core.Tag{
		TagName: tagName,
	}, taskId)
}

func (b *Biz) CreateTag(tagName string) (*core.Tag, error) {
	return b.tagDao.Create(b.ctx, &core.Tag{
		TagName: tagName,
	})
}

func (b *Biz) ListAllTag() ([]*core.Tag, error) {
	return b.tagDao.List(b.ctx, core.TagFilterParam{})
}
