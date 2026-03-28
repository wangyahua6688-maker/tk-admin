package biz

import (
	"context"
	"errors"

	"go-admin/internal/models"
	"gorm.io/gorm"
)

// PostArticleDAO 发帖数据访问层。
type PostArticleDAO struct {
	db *gorm.DB
}

// NewPostArticleDAO 创建发帖 DAO。
func NewPostArticleDAO(db *gorm.DB) *PostArticleDAO {
	return &PostArticleDAO{db: db}
}

// ListPostArticles 查询帖子列表（按官方/网友区分）。
func (d *PostArticleDAO) ListPostArticles(ctx context.Context, isOfficial int8, limit int) ([]models.WPostArticle, error) {
	// 定义并初始化当前变量。
	query := d.db.WithContext(ctx).Model(&models.WPostArticle{}).Where("is_official = ?", isOfficial).Order("id DESC")
	// 判断条件并进入对应分支逻辑。
	if limit > 0 {
		// 更新当前变量或字段值。
		query = query.Limit(limit)
	}
	// 定义并初始化当前变量。
	items := make([]models.WPostArticle, 0)
	// 判断条件并进入对应分支逻辑。
	if err := query.Find(&items).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return items, nil
}

// CreatePostArticle 新增帖子。
func (d *PostArticleDAO) CreatePostArticle(ctx context.Context, item *models.WPostArticle) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdatePostArticle 更新帖子。
func (d *PostArticleDAO) UpdatePostArticle(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WPostArticle{}).Where("id = ?", id).Updates(updates).Error
}

// DeletePostArticle 删除帖子。
func (d *PostArticleDAO) DeletePostArticle(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.WPostArticle{}, id).Error
}

// GetPostArticleByID 查询帖子是否存在。
func (d *PostArticleDAO) GetPostArticleByID(ctx context.Context, id uint) (*models.WPostArticle, error) {
	// 声明当前变量。
	var post models.WPostArticle
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).Select("id,user_id,is_official,status").Where("id = ?", id).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &post, nil
}
