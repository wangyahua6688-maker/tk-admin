package admin

// DrawRecordUpsertRequest 开奖区记录新增/编辑请求结构。
type DrawRecordUpsertRequest struct {
	SpecialLotteryID    *uint   `json:"special_lottery_id"`
	Issue               *string `json:"issue"`
	Year                *int    `json:"year"`
	DrawAt              *string `json:"draw_at"`
	NormalDrawResult    *string `json:"normal_draw_result"`
	SpecialDrawResult   *string `json:"special_draw_result"`
	DrawResult          *string `json:"draw_result"`
	DrawLabels          *string `json:"draw_labels"`
	PlaybackURL         *string `json:"playback_url"`
	SpecialSingleDouble *string `json:"special_single_double"`
	SpecialBigSmall     *string `json:"special_big_small"`
	SumSingleDouble     *string `json:"sum_single_double"`
	SumBigSmall         *string `json:"sum_big_small"`
	RecommendSix        *string `json:"recommend_six"`
	RecommendFour       *string `json:"recommend_four"`
	RecommendOne        *string `json:"recommend_one"`
	RecommendTen        *string `json:"recommend_ten"`
	SpecialCode         *string `json:"special_code"`
	NormalCode          *string `json:"normal_code"`
	Zheng1              *string `json:"zheng1"`
	Zheng2              *string `json:"zheng2"`
	Zheng3              *string `json:"zheng3"`
	Zheng4              *string `json:"zheng4"`
	Zheng5              *string `json:"zheng5"`
	Zheng6              *string `json:"zheng6"`
	IsCurrent           *int8   `json:"is_current"`
	Status              *int8   `json:"status"`
	Sort                *int    `json:"sort"`
}

// LotteryInfoUpsertRequest 图库内容新增/编辑请求结构。
type LotteryInfoUpsertRequest struct {
	SpecialLotteryID  *uint     `json:"special_lottery_id"`
	CategoryID        *uint     `json:"category_id"`
	CategoryTag       *string   `json:"category_tag"`
	Issue             *string   `json:"issue"`
	Year              *int      `json:"year"`
	Title             *string   `json:"title"`
	CoverImageURL     *string   `json:"cover_image_url"`
	DetailImageURL    *string   `json:"detail_image_url"`
	DrawCode          *string   `json:"draw_code"`
	NormalDrawResult  *string   `json:"normal_draw_result"`
	SpecialDrawResult *string   `json:"special_draw_result"`
	DrawResult        *string   `json:"draw_result"`
	DrawAt            *string   `json:"draw_at"`
	PlaybackURL       *string   `json:"playback_url"`
	LikesCount        *int64    `json:"likes_count"`
	CommentCount      *int64    `json:"comment_count"`
	FavoriteCount     *int64    `json:"favorite_count"`
	ReadCount         *int64    `json:"read_count"`
	PollEnabled       *int8     `json:"poll_enabled"`
	PollDefaultExpand *int8     `json:"poll_default_expand"`
	RecommendInfoIDs  *string   `json:"recommend_info_ids"`
	OptionNames       *[]string `json:"option_names"`
	IsCurrent         *int8     `json:"is_current"`
	Status            *int8     `json:"status"`
	Sort              *int      `json:"sort"`
}
