package biz

import (
	"context"

	bizdao "go-admin/internal/dao/biz"
	"go-admin/internal/models"
	"gorm.io/gorm"
)

// PostArticleService 发帖业务服务。
type PostArticleService struct {
	dao *bizdao.PostArticleDAO
}

// NewPostArticleService 创建发帖服务。
func NewPostArticleService(db *gorm.DB) *PostArticleService {
	return &PostArticleService{dao: bizdao.NewPostArticleDAO(db)}
}

// ListPostArticles 查询帖子列表（按官方/网友区分）。
func (s *PostArticleService) ListPostArticles(ctx context.Context, isOfficial bool, limit int) ([]models.WPostArticle, error) {
	// 定义并初始化当前变量。
	flag := int8(0)
	// 判断条件并进入对应分支逻辑。
	if isOfficial {
		// 更新当前变量或字段值。
		flag = 1
	}
	// 返回当前处理结果。
	return s.dao.ListPostArticles(ctx, flag, limit)
}

// CreatePostArticle 新增帖子。
func (s *PostArticleService) CreatePostArticle(ctx context.Context, item *models.WPostArticle) error {
	// 返回当前处理结果。
	return s.dao.CreatePostArticle(ctx, item)
}

// UpdatePostArticle 更新帖子。
func (s *PostArticleService) UpdatePostArticle(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return s.dao.UpdatePostArticle(ctx, id, updates)
}

// DeletePostArticle 删除帖子。
func (s *PostArticleService) DeletePostArticle(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.DeletePostArticle(ctx, id)
}
