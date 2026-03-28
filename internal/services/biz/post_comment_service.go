package biz

import (
	"context"
	"errors"
	"strings"

	bizdao "go-admin/internal/dao/biz"
	"go-admin/internal/models"
	"gorm.io/gorm"
)

// PostCommentService 帖子评论业务服务。
type PostCommentService struct {
	dao            *bizdao.PostCommentDAO
	postArticleDAO *bizdao.PostArticleDAO
	lookupSvc      *UserOpsLookupService
}

// NewPostCommentService 创建帖子评论服务。
func NewPostCommentService(db *gorm.DB) *PostCommentService {
	return &PostCommentService{
		dao:            bizdao.NewPostCommentDAO(db),
		postArticleDAO: bizdao.NewPostArticleDAO(db),
		lookupSvc:      NewUserOpsLookupService(db),
	}
}

// ListPostComments 查询帖子评论并补充用户信息。
func (s *PostCommentService) ListPostComments(ctx context.Context, postID uint) ([]map[string]interface{}, error) {
	// 拉取评论列表。
	rows, err := s.dao.ListPostComments(ctx, postID, 500)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, err
	}

	// 收集用户ID。
	userIDs := make([]uint, 0)
	// 定义并初始化当前变量。
	seen := make(map[uint]struct{})
	// 循环处理当前数据集合。
	for _, row := range rows {
		// 判断条件并进入对应分支逻辑。
		if _, ok := seen[row.UserID]; ok {
			// 处理当前语句逻辑。
			continue
		}
		// 更新当前变量或字段值。
		seen[row.UserID] = struct{}{}
		// 更新当前变量或字段值。
		userIDs = append(userIDs, row.UserID)
	}

	// 批量查询用户信息。
	userMap := make(map[uint]models.WUser)
	// 判断条件并进入对应分支逻辑。
	if len(userIDs) > 0 {
		// 定义并初始化当前变量。
		users, uerr := s.lookupSvc.GetUsersByIDs(ctx, userIDs)
		// 判断条件并进入对应分支逻辑。
		if uerr != nil {
			// 返回当前处理结果。
			return nil, uerr
		}
		// 循环处理当前数据集合。
		for _, u := range users {
			// 更新当前变量或字段值。
			userMap[u.ID] = u
		}
	}

	// 组装返回结构。
	items := make([]map[string]interface{}, 0, len(rows))
	// 循环处理当前数据集合。
	for _, row := range rows {
		// 定义并初始化当前变量。
		u := userMap[row.UserID]
		// 更新当前变量或字段值。
		items = append(items, map[string]interface{}{
			// 处理当前语句逻辑。
			"id": row.ID,
			// 处理当前语句逻辑。
			"post_id": row.PostID,
			// 处理当前语句逻辑。
			"user_id": row.UserID,
			// 处理当前语句逻辑。
			"parent_id": row.ParentID,
			// 处理当前语句逻辑。
			"content": row.Content,
			// 处理当前语句逻辑。
			"likes": row.Likes,
			// 处理当前语句逻辑。
			"status": row.Status,
			// 处理当前语句逻辑。
			"created_at": row.CreatedAt,
			// 处理当前语句逻辑。
			"username": u.Username,
			// 处理当前语句逻辑。
			"nickname": u.Nickname,
			// 处理当前语句逻辑。
			"user_type": u.UserType,
		})
	}

	// 返回当前处理结果。
	return items, nil
}

// CreatePostComment 新增帖子评论（机器人/官方）。
func (s *PostCommentService) CreatePostComment(ctx context.Context, postID uint, userID uint, parentID uint, content string, status *int8) (*models.WComment, error) {
	// 参数校验。
	if postID == 0 {
		// 返回当前处理结果。
		return nil, errors.New("post_id required")
	}
	// 判断条件并进入对应分支逻辑。
	if userID == 0 {
		// 返回当前处理结果。
		return nil, errors.New("user_id required")
	}
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(content) == "" {
		// 返回当前处理结果。
		return nil, errors.New("content required")
	}

	// 校验用户类型。
	if !s.lookupSvc.IsUserTypes(ctx, userID, "robot", "official") {
		// 返回当前处理结果。
		return nil, errors.New("user must be robot or official account")
	}

	// 校验帖子存在。
	post, err := s.postArticleDAO.GetPostArticleByID(ctx, postID)
	if err != nil {
		// 返回当前处理结果。
		return nil, errors.New("post not found")
	}
	if post == nil {
		// 返回当前处理结果。
		return nil, errors.New("post not found")
	}

	// 组装评论模型。
	item := models.WComment{
		// 处理当前语句逻辑。
		PostID: postID,
		// 处理当前语句逻辑。
		UserID: userID,
		// 处理当前语句逻辑。
		ParentID: parentID,
		// 调用strings.TrimSpace完成当前处理。
		Content: strings.TrimSpace(content),
		// 处理当前语句逻辑。
		Likes: 0,
		// 处理当前语句逻辑。
		Status: 1,
	}
	// 判断条件并进入对应分支逻辑。
	if status != nil {
		// 更新当前变量或字段值。
		item.Status = *status
	}

	// 落库。
	if err := s.dao.CreatePostComment(ctx, &item); err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &item, nil
}
