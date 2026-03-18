package rbac

import (
	"net/http"
	"strconv"
	"strings"

	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
	rbacsvc "go-admin-full/internal/services/rbac"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MenuController 定义MenuController相关结构。
type MenuController struct {
	svc         *rbacsvc.MenuService           // 菜单基础增删改查服务
	menuPermSvc *rbacsvc.MenuPermissionService // 菜单-权限绑定关系服务
}

// NewMenuController 创建MenuController实例。
func NewMenuController(db *gorm.DB) *MenuController {
	menuDao := rbacdao.NewMenuDAO(db)                                                 // 初始化菜单 DAO
	menuSvc := rbacsvc.NewMenuService(menuDao)                                        // 构造菜单服务
	menuPermSvc := rbacsvc.NewMenuPermissionService(rbacdao.NewMenuPermissionDao(db)) // 构造菜单-权限服务

	// 返回当前处理结果。
	return &MenuController{
		svc:         menuSvc,     // 绑定菜单服务
		menuPermSvc: menuPermSvc, // 绑定菜单权限服务
	}
}

// 定义当前类型结构。
type menuUpsertReq struct {
	// 处理当前语句逻辑。
	Title string `json:"title" binding:"required"`
	// 处理当前语句逻辑。
	Path string `json:"path" binding:"required"`
	// 处理当前语句逻辑。
	Icon string `json:"icon"`
	// 处理当前语句逻辑。
	ParentID uint `json:"parent_id"`
	// 处理当前语句逻辑。
	Component string `json:"component"`
	// 处理当前语句逻辑。
	OrderNum int `json:"order_num"`
}

// List 菜单列表查询。
func (mc *MenuController) List(c *gin.Context) {
	items, err := mc.svc.List(c.Request.Context()) // 拉取全量菜单列表
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // 返回服务端错误
		// 返回当前处理结果。
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items}) // 返回菜单数据
}

// Create 创建菜单。
func (mc *MenuController) Create(c *gin.Context) {
	// 声明当前变量。
	var req menuUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil { // 解析并校验请求体
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	m := models.Menu{
		Title:     strings.TrimSpace(req.Title),     // 菜单标题
		Path:      strings.TrimSpace(req.Path),      // 路由路径
		Icon:      strings.TrimSpace(req.Icon),      // 图标标识
		ParentID:  req.ParentID,                     // 父级菜单
		Component: strings.TrimSpace(req.Component), // 前端组件路径
		OrderNum:  req.OrderNum,                     // 排序权重
	}

	if m.Title == "" || m.Path == "" { // 基础字段校验
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "title/path不能为空"})
		// 返回当前处理结果。
		return
	}
	if m.ParentID > 0 { // 有父级时校验父级是否存在
		// 判断条件并进入对应分支逻辑。
		if _, err := mc.svc.Get(c.Request.Context(), m.ParentID); err != nil {
			// 调用c.JSON完成当前处理。
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent_id无效"})
			// 返回当前处理结果。
			return
		}
	}

	if err := mc.svc.Create(c.Request.Context(), &m); err != nil { // 入库创建菜单
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m}) // 返回创建结果
}

// Get 查询菜单详情。
func (mc *MenuController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) // 解析菜单 ID
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		// 返回当前处理结果。
		return
	}

	m, err := mc.svc.Get(c.Request.Context(), uint(id)) // 查询菜单
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m}) // 返回菜单详情
}

// Update 更新菜单。
func (mc *MenuController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) // 解析菜单 ID
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		// 返回当前处理结果。
		return
	}

	// 声明当前变量。
	var req menuUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil { // 解析更新字段
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	menu, err := mc.svc.Get(c.Request.Context(), uint(id)) // 获取旧数据
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	menu.Title = strings.TrimSpace(req.Title)         // 更新标题
	menu.Path = strings.TrimSpace(req.Path)           // 更新路径
	menu.Icon = strings.TrimSpace(req.Icon)           // 更新图标
	menu.ParentID = req.ParentID                      // 更新父级
	menu.Component = strings.TrimSpace(req.Component) // 更新组件路径
	menu.OrderNum = req.OrderNum                      // 更新排序

	if menu.Title == "" || menu.Path == "" { // 更新后校验必填字段
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "title/path不能为空"})
		// 返回当前处理结果。
		return
	}
	if menu.ParentID == uint(id) { // 不允许自引用
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "parent_id不能为自身"})
		// 返回当前处理结果。
		return
	}
	if menu.ParentID > 0 { // 校验父级有效性
		// 判断条件并进入对应分支逻辑。
		if _, err := mc.svc.Get(c.Request.Context(), menu.ParentID); err != nil {
			// 调用c.JSON完成当前处理。
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent_id无效"})
			// 返回当前处理结果。
			return
		}
	}

	if err := mc.svc.Update(c.Request.Context(), menu); err != nil { // 持久化更新
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menu}) // 返回更新结果
}

// Delete 删除菜单。
func (mc *MenuController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) // 解析菜单 ID
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		// 返回当前处理结果。
		return
	}

	if err := mc.svc.Delete(c.Request.Context(), uint(id)); err != nil { // 删除菜单
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "deleted"}) // 返回删除结果
}

// FrontendTree 用于前端动态路由，返回当前用户权限范围内的菜单树。
func (mc *MenuController) FrontendTree(c *gin.Context) {
	uid := c.GetUint("uid")                                    // 从 JWT 中获取用户 ID
	menus, err := mc.svc.ListForUser(c.Request.Context(), uid) // 查询用户可见菜单
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(500, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	tree := mc.svc.BuildMenuTreeFromList(menus) // 构建树形结构
	c.JSON(200, gin.H{"data": tree})            // 返回菜单树
}

// GetPermissions 查询菜单已绑定权限。
func (mc *MenuController) GetPermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) // 解析菜单 ID
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		// 返回当前处理结果。
		return
	}

	perms, err := mc.menuPermSvc.GetMenuPermissions(c.Request.Context(), uint(id)) // 查询绑定的权限点
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": perms}) // 返回权限列表
}

// BindPermissions 绑定菜单权限（全量替换）。
func (mc *MenuController) BindPermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) // 解析菜单 ID
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		// 返回当前处理结果。
		return
	}

	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { // 解析权限列表
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	if err := mc.menuPermSvc.BindPermissions(c.Request.Context(), uint(id), req.PermissionIDs); err != nil { // 执行绑定
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "menu permissions bound"}) // 返回绑定结果
}
