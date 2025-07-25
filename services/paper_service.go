package services

import (
	"CodeQuizAI/dao"
	"CodeQuizAI/models"
	"CodeQuizAI/utils"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// GetPapersRequest 试卷列表查询参数
type GetPapersRequest struct {
	Page     int    `form:"page"`      // 页码
	PageSize int    `form:"page_size"` // 每页条数
	Keyword  string `form:"keyword"`   // 标题关键词
	Sort     string `form:"sort"`      // 排序方式
}

// PaperListResponse 试卷列表响应
type PaperListResponse struct {
	List       []PaperItem `json:"list"`       // 试卷列表
	Pagination Pagination  `json:"pagination"` // 分页信息（复用之前的结构体）
}

// PaperItem 单份试卷信息（列表中展示的字段）
type PaperItem struct {
	ID          int64     `json:"id"`          // 试卷ID
	Title       string    `json:"title"`       // 试卷标题
	Description string    `json:"description"` // 试卷描述
	TotalScore  int       `json:"total_score"` // 总分
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
}

// GetUserPapers 查询用户创建的试卷列表
func GetUserPapers(
	ctx context.Context,
	creatorID int64,
	req GetPapersRequest,
) (PaperListResponse, error) {
	// 1. 构建查询条件：仅查询当前用户创建的试卷（未被软删除）
	query := dao.Q.Paper.WithContext(ctx).
		Where(dao.Q.Paper.CreatorID.Eq(creatorID))

	// 2. 关键词搜索（模糊匹配标题）
	if req.Keyword != "" {
		query = query.Where(dao.Q.Paper.Title.Like("%" + req.Keyword + "%"))
	}

	// 3. 排序
	switch req.Sort {
	case "created_at_asc":
		query = query.Order(dao.Q.Paper.CreatedAt.Asc())
	case "created_at_desc":
		query = query.Order(dao.Q.Paper.CreatedAt.Desc())
	default:
		return PaperListResponse{}, errors.New("无效的排序方式")
	}

	// 4. 分页计算
	offset := (req.Page - 1) * req.PageSize
	query = query.Limit(req.PageSize).Offset(offset)

	// 5. 执行查询
	papers, err := query.Find()
	if err != nil {
		return PaperListResponse{}, fmt.Errorf("查询失败：%w", err)
	}

	// 6. 查询总条数（用于分页信息）
	total, err := dao.Q.Paper.WithContext(ctx).
		Where(dao.Q.Paper.CreatorID.Eq(creatorID)).
		Count()
	if err != nil {
		return PaperListResponse{}, fmt.Errorf("统计总数失败：%w", err)
	}

	// 7. 转换响应格式
	var paperList []PaperItem
	for _, p := range papers {
		paperList = append(paperList, PaperItem{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			TotalScore:  p.TotalScore,
			CreatedAt:   p.CreatedAt,
		})
	}

	// 8. 计算总页数
	totalPages := (int(total) + req.PageSize - 1) / req.PageSize

	return PaperListResponse{
		List: paperList,
		Pagination: Pagination{
			Total:      total,
			Page:       int64(req.Page),
			PageSize:   int64(req.PageSize),
			TotalPages: int64(totalPages),
		},
	}, nil
}

// CreatePaperRequest 创建试卷的请求参数
type CreatePaperRequest struct {
	Title       string `json:"title" binding:"required"` // 试卷标题（必填）
	Description string `json:"description"`              // 试卷描述（可选）
	TotalScore  int    `json:"total_score"`              // 总分（可选，默认100）
}

// CreatePaperResponse 创建试卷的响应数据
type CreatePaperResponse struct {
	ID          int64     `json:"id"`          // 试卷ID
	Title       string    `json:"title"`       // 标题
	Description string    `json:"description"` // 描述
	TotalScore  int       `json:"total_score"` // 总分
	CreatorID   int64     `json:"creator_id"`  // 创建者ID
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
}

// CreatePaper 创建试卷
func CreatePaper(ctx context.Context, creatorID int64, req CreatePaperRequest) (*CreatePaperResponse, error) {
	// 1. 校验必填参数
	if req.Title == "" {
		return nil, errors.New("试卷标题不能为空")
	}

	// 2. 设置默认值（总分默认100）
	totalScore := req.TotalScore
	if totalScore <= 0 {
		totalScore = 100
	}

	// 3. 构造试卷模型
	paper := &models.Paper{
		Title:       req.Title,
		Description: req.Description,
		TotalScore:  totalScore,
		CreatorID:   creatorID,
	}

	// 4. 插入数据库
	if err := dao.Q.Paper.WithContext(ctx).Create(paper); err != nil {
		return nil, fmt.Errorf("数据库插入失败：%w", err)
	}

	// 5. 转换为响应格式并返回
	return &CreatePaperResponse{
		ID:          paper.ID,
		Title:       paper.Title,
		Description: paper.Description,
		TotalScore:  paper.TotalScore,
		CreatorID:   paper.CreatorID,
		CreatedAt:   paper.CreatedAt,
	}, nil
}

// PaperDetailResponse 试卷详情响应（包含关联题目）
type PaperDetailResponse struct {
	ID          int64                   `json:"id"`          // 试卷ID
	Title       string                  `json:"title"`       // 标题
	Description string                  `json:"description"` // 描述
	TotalScore  int                     `json:"total_score"` // 总分
	CreatedAt   time.Time               `json:"created_at"`  // 创建时间
	Questions   []PaperQuestionWithInfo `json:"questions"`   // 关联题目列表（带题目详情）
}

// PaperQuestionWithInfo 试卷中的题目信息（包含题目详情）
type PaperQuestionWithInfo struct {
	QuestionID    int64       `json:"question_id"`    // 题目ID
	QuestionOrder int         `json:"question_order"` // 题目顺序
	Score         int         `json:"score"`          // 题目分值
	QuestionInfo  QuestionDTO `json:"question_info"`  // 题目详情
}

// QuestionDTO 题目详情DTO
type QuestionDTO struct {
	ID           int64     `json:"id"`                    // 题目ID
	Title        string    `json:"title"`                 // 题目标题
	QuestionType string    `json:"question_type"`         // 题型（single/multiple）
	Options      string    `json:"options"`               // 选项（JSON格式字符串）
	Answer       string    `json:"answer"`                // 答案
	Explanation  string    `json:"explanation,omitempty"` // 解析（可选）
	Keywords     string    `json:"keywords,omitempty"`    // 关键词（可选）
	Language     string    `json:"language"`              // 编程语言
	AiModel      string    `json:"ai_model"`              // 使用的AI模型
	UserID       int64     `json:"user_id"`               // 创建者ID
	CreatedAt    time.Time `json:"created_at"`            // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`            // 更新时间
}

// GetPaperDetail 查询试卷详情（包含关联题目）
func GetPaperDetail(ctx context.Context, paperID, creatorID int64) (*PaperDetailResponse, error) {
	// 1. 查询试卷基本信息（验证存在性和权限）
	paper, err := dao.Q.Paper.WithContext(ctx).
		Where(
			dao.Paper.ID.Eq(paperID),
			dao.Paper.CreatorID.Eq(creatorID), // 仅允许查看自己创建的试卷
		).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrPaperNotFound
		}
		return nil, fmt.Errorf("查询试卷失败：%w", err)
	}

	// 2. 查询试卷与题目的关联关系
	paperQuestions, err := dao.Q.PaperQuestion.WithContext(ctx).
		Where(dao.PaperQuestion.PaperID.Eq(paperID)).
		Order(dao.PaperQuestion.QuestionOrder.Asc()). // 按顺序排序
		Find()
	if err != nil {
		return nil, fmt.Errorf("查询试卷题目关联失败：%w", err)
	}

	// 3. 提取题目ID列表，批量查询题目详情
	var questionIDs []int64
	for _, pq := range paperQuestions {
		questionIDs = append(questionIDs, pq.QuestionID)
	}
	if len(questionIDs) == 0 {
		// 试卷无题目时直接返回基本信息
		return &PaperDetailResponse{
			ID:          paper.ID,
			Title:       paper.Title,
			Description: paper.Description,
			TotalScore:  paper.TotalScore,
			CreatedAt:   paper.CreatedAt,
			Questions:   []PaperQuestionWithInfo{},
		}, nil
	}

	// 4. 批量查询题目详情
	questions, err := dao.Q.Question.WithContext(ctx).
		Where(
			dao.Question.ID.In(questionIDs...),
			dao.Question.DeletedAt.IsNull(), // 排除已删除题目
		).
		Find()
	if err != nil {
		return nil, fmt.Errorf("查询题目详情失败：%w", err)
	}

	// 5. 构建题目ID到详情的映射
	questionMap := make(map[int64]*models.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	// 6. 组装试卷题目列表
	var questionList []PaperQuestionWithInfo
	for _, pq := range paperQuestions {
		question, ok := questionMap[pq.QuestionID]
		if !ok {
			continue
		}

		questionList = append(questionList, PaperQuestionWithInfo{
			QuestionID:    pq.QuestionID,
			QuestionOrder: pq.QuestionOrder,
			Score:         pq.Score,
			QuestionInfo: QuestionDTO{
				ID:           question.ID,
				Title:        question.Title,
				QuestionType: question.QuestionType,
				Options:      question.Options,
				Answer:       question.Answer,
				Explanation:  question.Explanation,
				Keywords:     question.Keywords,
				Language:     question.Language,
				AiModel:      question.AiModel,
				UserID:       question.UserID,
				CreatedAt:    question.CreatedAt,
				UpdatedAt:    question.UpdatedAt,
			},
		})
	}

	// 7. 组装并返回试卷详情
	return &PaperDetailResponse{
		ID:          paper.ID,
		Title:       paper.Title,
		Description: paper.Description,
		TotalScore:  paper.TotalScore,
		CreatedAt:   paper.CreatedAt,
		Questions:   questionList,
	}, nil
}

// DeletePaperResponse 删除相关响应结构体
type DeletePaperResponse struct {
	ID        int64     `json:"id"`         // 试卷ID
	DeletedAt time.Time `json:"deleted_at"` // 删除时间
}

// DeletePaper 软删除试卷及删除关联的题目关系
func DeletePaper(ctx context.Context, paperID, creatorID int64) (DeletePaperResponse, error) {
	// 1. 检查试卷是否存在且属于当前用户
	paper, err := dao.Q.Paper.WithContext(ctx).
		Unscoped().
		Where(
			dao.Q.Paper.ID.Eq(paperID),
			dao.Q.Paper.CreatorID.Eq(creatorID),
		).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			exists, checkErr := checkPaperExists(ctx, paperID)
			if checkErr != nil {
				return DeletePaperResponse{}, checkErr
			}
			if exists {
				return DeletePaperResponse{}, utils.ErrNoPermission
			}
			return DeletePaperResponse{}, utils.ErrPaperNotFound
		}
		return DeletePaperResponse{}, fmt.Errorf("查询试卷失败：%w", err)
	}

	// 2. 已删除则直接返回
	if !paper.DeletedAt.Time.IsZero() {
		return DeletePaperResponse{
			ID:        paperID,
			DeletedAt: paper.DeletedAt.Time,
		}, nil
	}

	// 3. 删除关联关系
	_, err = dao.Q.PaperQuestion.WithContext(ctx).
		Where(dao.Q.PaperQuestion.PaperID.Eq(paperID)).
		Delete()
	if err != nil {
		return DeletePaperResponse{}, fmt.Errorf("删除试卷题目关联失败：%w", err)
	}

	// 4. 软删除试卷（更新deleted_at字段）
	_, err = dao.Q.Paper.WithContext(ctx).
		Where(
			dao.Q.Paper.ID.Eq(paperID),
			dao.Q.Paper.CreatorID.Eq(creatorID),
		).
		Delete()
	if err != nil {
		return DeletePaperResponse{}, fmt.Errorf("软删除试卷失败：%w", err)
	}

	// 5. 查询删除时间并返回结果
	deletedPaper, err := dao.Q.Paper.WithContext(ctx).
		Unscoped().
		Where(dao.Q.Paper.ID.Eq(paperID)).
		First()
	if err != nil {
		return DeletePaperResponse{}, err
	}

	return DeletePaperResponse{
		ID:        paperID,
		DeletedAt: deletedPaper.DeletedAt.Time,
	}, nil
}

// 检查试卷是否存在
func checkPaperExists(ctx context.Context, paperID int64) (bool, error) {
	count, err := dao.Q.Paper.WithContext(ctx).
		Unscoped().
		Where(dao.Q.Paper.ID.Eq(paperID)).
		Count()
	return count > 0, err
}

// AddQuestionsToPaperRequest 添加题目到试卷的请求参数
type AddQuestionsToPaperRequest struct {
	Items []PaperQuestionItem `json:"items"` // 题目列表（包含顺序和分值）
}

// PaperQuestionItem 单题关联信息
type PaperQuestionItem struct {
	QuestionID    int64 `json:"question_id" binding:"required"`          // 题目ID
	QuestionOrder int   `json:"question_order" binding:"required,min=1"` // 题目顺序（从1开始）
	Score         int   `json:"score" binding:"required,min=1"`          // 题目分值（至少1分）
}

// AddQuestionsToPaperResponse 添加结果响应
type AddQuestionsToPaperResponse struct {
	PaperID        int64                 `json:"paper_id"`        // 试卷ID
	AddedCount     int                   `json:"added_count"`     // 成功添加的题目数量
	TotalQuestions int                   `json:"total_questions"` // 试卷当前总题目数
	TotalScore     int                   `json:"total_score"`     // 试卷当前总分值
	Items          []AddedQuestionResult `json:"items"`           // 各题添加结果
}

// AddedQuestionResult 单题添加结果
type AddedQuestionResult struct {
	QuestionID int64  `json:"question_id"`
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
}

// AddQuestionsToPaper 向试卷添加题目
func AddQuestionsToPaper(ctx context.Context, paperID, creatorID int64, req AddQuestionsToPaperRequest) (AddQuestionsToPaperResponse, error) {
	// 1. 验证试卷是否存在且属于当前用户
	paper, err := dao.Q.Paper.WithContext(ctx).
		Where(
			dao.Q.Paper.ID.Eq(paperID),
			dao.Q.Paper.CreatorID.Eq(creatorID),
		).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AddQuestionsToPaperResponse{}, utils.ErrPaperNotFound
		}
		return AddQuestionsToPaperResponse{}, fmt.Errorf("查询试卷失败: %w", err)
	}

	// 2. 验证题目ID合法性（必须存在且属于当前用户）
	questionIDs := make([]int64, len(req.Items))
	for i, item := range req.Items {
		questionIDs[i] = item.QuestionID
	}

	// 2.1 批量查询题目是否存在
	questions, err := dao.Q.Question.WithContext(ctx).
		Where(
			dao.Q.Question.ID.In(questionIDs...),
			dao.Q.Question.UserID.Eq(creatorID), // 确保是用户自己的题目
		).
		Find()
	if err != nil {
		return AddQuestionsToPaperResponse{}, fmt.Errorf("查询题目失败: %w", err)
	}

	// 2.2 构建存在的题目ID映射
	existQuestionIDs := make(map[int64]bool)
	for _, q := range questions {
		existQuestionIDs[q.ID] = true
	}

	// 3. 检查试卷中已有的题目（避免重复添加）
	existingRelations, err := dao.Q.PaperQuestion.WithContext(ctx).
		Where(dao.Q.PaperQuestion.PaperID.Eq(paperID)).
		Find()
	if err != nil {
		return AddQuestionsToPaperResponse{}, fmt.Errorf("查询试卷题目关联失败: %w", err)
	}

	existingQIDs := make(map[int64]bool)
	existingScoreSum := 0
	for _, rel := range existingRelations {
		existingQIDs[rel.QuestionID] = true
		existingScoreSum += rel.Score
	}

	// 4. 检查题目顺序是否重复
	orderSet := make(map[int]bool)
	for _, item := range req.Items {
		if orderSet[item.QuestionOrder] {
			return AddQuestionsToPaperResponse{}, utils.ErrInvalidOrder
		}
		orderSet[item.QuestionOrder] = true
	}

	// 5. 计算新增题目的总分并校验是否超过试卷上限
	var newTotalScore int
	for _, item := range req.Items {
		newTotalScore += item.Score
	}

	totalAfterAdd := existingScoreSum + newTotalScore
	if totalAfterAdd > paper.TotalScore {
		return AddQuestionsToPaperResponse{}, fmt.Errorf(
			"%w（当前已有: %d, 新增: %d, 总计: %d, 试卷上限: %d）",
			utils.ErrScoreExceedTotal, existingScoreSum, newTotalScore, totalAfterAdd, paper.TotalScore,
		)
	}

	// 6. 批量创建关联记录
	var addedItems []*models.PaperQuestion
	var resultItems []AddedQuestionResult
	addedCount := 0

	for _, item := range req.Items {
		resItem := AddedQuestionResult{QuestionID: item.QuestionID}

		// 验证单题合法性
		if !existQuestionIDs[item.QuestionID] {
			resItem.Success = false
			resItem.Message = "题目不存在或无权限"
			resultItems = append(resultItems, resItem)
			continue
		}

		if existingQIDs[item.QuestionID] {
			resItem.Success = false
			resItem.Message = "题目已在试卷中"
			resultItems = append(resultItems, resItem)
			continue
		}

		// 合法题目加入待创建列表
		addedItems = append(addedItems, &models.PaperQuestion{
			PaperID:       paperID,
			QuestionID:    item.QuestionID,
			QuestionOrder: item.QuestionOrder,
			Score:         item.Score,
		})
		resItem.Success = true
		resultItems = append(resultItems, resItem)
		addedCount++
	}

	if len(resultItems) == 1 && !resultItems[0].Success {
		return AddQuestionsToPaperResponse{}, utils.ErrQuestionNotFound
	}

	// 7. 执行批量插入
	if len(addedItems) > 0 {
		err := dao.Q.PaperQuestion.WithContext(ctx).
			CreateInBatches(addedItems, 100)
		if err != nil {
			return AddQuestionsToPaperResponse{}, fmt.Errorf("添加题目失败: %w", err)
		}
	}

	// 8. 计算当前试卷总题目数和总分
	totalCount := len(existingRelations) + addedCount
	var truthTotalScore int
	for _, item := range addedItems {
		truthTotalScore += item.Score
	}
	currentTotalScore := existingScoreSum + truthTotalScore

	return AddQuestionsToPaperResponse{
		PaperID:        paperID,
		AddedCount:     addedCount,
		TotalQuestions: totalCount,
		TotalScore:     currentTotalScore,
		Items:          resultItems,
	}, nil
}

// RemoveQuestionFromPaperRequest 移除试卷中题目的请求参数
type RemoveQuestionFromPaperRequest struct {
	PaperID    int64 // 试卷ID（由路径参数获取）
	QuestionID int64 // 题目ID（由路径参数获取）
}

// RemoveQuestionFromPaperResponse 移除试卷中题目的响应
type RemoveQuestionFromPaperResponse struct {
	PaperID           int64 `json:"paper_id"`            // 试卷ID
	RemovedQuestionID int64 `json:"removed_question_id"` // 被移除的题目ID
	RemainingCount    int   `json:"remaining_count"`     // 剩余题目数量
	RemainingScore    int   `json:"remaining_score"`     // 剩余题目总分
}

// RemoveQuestionFromPaper 从试卷中移除指定题目
func RemoveQuestionFromPaper(ctx context.Context, paperID, questionID, creatorID int64) (RemoveQuestionFromPaperResponse, error) {
	// 1. 验证试卷是否存在且属于当前用户
	_, err := dao.Q.Paper.WithContext(ctx).
		Where(
			dao.Q.Paper.ID.Eq(paperID),
			dao.Q.Paper.CreatorID.Eq(creatorID),
		).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return RemoveQuestionFromPaperResponse{}, utils.ErrPaperNotFound
		}
		return RemoveQuestionFromPaperResponse{}, fmt.Errorf("查询试卷失败: %w", err)
	}

	// 2. 验证题目是否存在于该试卷中
	relation, err := dao.Q.PaperQuestion.WithContext(ctx).
		Where(
			dao.Q.PaperQuestion.PaperID.Eq(paperID),
			dao.Q.PaperQuestion.QuestionID.Eq(questionID),
		).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return RemoveQuestionFromPaperResponse{}, errors.New("题目不存在于该试卷中")
		}
		return RemoveQuestionFromPaperResponse{}, fmt.Errorf("查询题目关联失败: %w", err)
	}

	// 3. 删除试卷与题目的关联记录
	if _, err := dao.Q.PaperQuestion.WithContext(ctx).
		Where(
			dao.Q.PaperQuestion.PaperID.Eq(paperID),
			dao.Q.PaperQuestion.QuestionID.Eq(questionID),
		).
		Delete(relation); err != nil {
		return RemoveQuestionFromPaperResponse{}, fmt.Errorf("移除题目失败: %w", err)
	}

	// 4. 查询剩余题目总分
	remainingRelations, err := dao.Q.PaperQuestion.WithContext(ctx).
		Where(dao.Q.PaperQuestion.PaperID.Eq(paperID)).
		Find()
	if err != nil {
		return RemoveQuestionFromPaperResponse{}, fmt.Errorf("查询剩余题目失败: %w", err)
	}

	remainingScore := 0
	for _, rel := range remainingRelations {
		remainingScore += rel.Score
	}

	// 5. 构建响应
	return RemoveQuestionFromPaperResponse{
		PaperID:           paperID,
		RemovedQuestionID: questionID,
		RemainingCount:    len(remainingRelations),
		RemainingScore:    remainingScore,
	}, nil
}

// UpdateQuestionOrderRequest 调整题目顺序的请求参数
type UpdateQuestionOrderRequest struct {
	QuestionOrders []QuestionOrderItem `json:"question_orders"` // 题目ID与新顺序的列表
}

// QuestionOrderItem 单个题目的顺序信息
type QuestionOrderItem struct {
	QuestionID    int64 `json:"question_id"`    // 题目ID
	QuestionOrder int   `json:"question_order"` // 新顺序值
}

// UpdateQuestionOrderResponse 调整题目顺序的响应
type UpdateQuestionOrderResponse struct {
	PaperID      int64 `json:"paper_id"`      // 试卷ID
	UpdatedCount int   `json:"updated_count"` // 成功更新的题目数量
}

// UpdateQuestionOrder 调整试卷中题目的顺序（确保所有题目顺序唯一，包括未修改的题目）
func UpdateQuestionOrder(ctx context.Context, paperID, creatorID int64, req UpdateQuestionOrderRequest) (UpdateQuestionOrderResponse, error) {
	// 1. 验证试卷是否存在且属于当前用户
	_, err := dao.Q.Paper.WithContext(ctx).
		Where(
			dao.Q.Paper.ID.Eq(paperID),
			dao.Q.Paper.CreatorID.Eq(creatorID),
		).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UpdateQuestionOrderResponse{}, utils.ErrPaperNotFound
		}
		return UpdateQuestionOrderResponse{}, fmt.Errorf("查询试卷失败: %w", err)
	}

	// 2. 查询该试卷中所有题目的现有顺序（包括未被本次修改的题目）
	var existingOrders []int // 声明用于接收顺序的切片
	err = dao.Q.PaperQuestion.WithContext(ctx).
		Where(dao.Q.PaperQuestion.PaperID.Eq(paperID)).
		Pluck(dao.Q.PaperQuestion.QuestionOrder, &existingOrders)
	if err != nil {
		return UpdateQuestionOrderResponse{}, fmt.Errorf("查询现有题目顺序失败: %w", err)
	}

	// 3. 提取本次请求中要更新的新顺序，并检查重复性
	newOrders := make([]int, 0, len(req.QuestionOrders))
	newOrderMap := make(map[int]bool)
	for _, item := range req.QuestionOrders {
		// 检查新顺序内部是否重复
		if newOrderMap[item.QuestionOrder] {
			return UpdateQuestionOrderResponse{}, fmt.Errorf("新顺序 %d 重复，请确保所有新顺序唯一", item.QuestionOrder)
		}
		newOrderMap[item.QuestionOrder] = true
		newOrders = append(newOrders, item.QuestionOrder)
	}

	// 4. 检查新顺序是否与现有未修改的题目顺序重复
	existingOrderMap := make(map[int]bool)
	for _, order := range existingOrders {
		isBeingUpdated := false
		for _, item := range req.QuestionOrders {
			// 查询该题目原来的顺序
			oldRel, err := dao.Q.PaperQuestion.WithContext(ctx).
				Where(
					dao.Q.PaperQuestion.PaperID.Eq(paperID),
					dao.Q.PaperQuestion.QuestionID.Eq(item.QuestionID),
				).
				First()
			if err == nil && oldRel.QuestionOrder == order {
				isBeingUpdated = true
				break
			}
		}
		if !isBeingUpdated {
			existingOrderMap[order] = true
		}
	}

	// 校验新顺序是否与未修改的顺序冲突
	for _, newOrder := range newOrders {
		if existingOrderMap[newOrder] {
			return UpdateQuestionOrderResponse{}, fmt.Errorf("新顺序 %d 已被其他题目使用，请更换", newOrder)
		}
	}

	// 5. 提取所有题目ID用于验证是否属于该试卷
	var questionIDs []int64
	for _, item := range req.QuestionOrders {
		questionIDs = append(questionIDs, item.QuestionID)
	}

	// 6. 验证这些题目是否都属于该试卷
	existingRelations, err := dao.Q.PaperQuestion.WithContext(ctx).
		Where(
			dao.Q.PaperQuestion.PaperID.Eq(paperID),
			dao.Q.PaperQuestion.QuestionID.In(questionIDs...),
		).
		Find()
	if err != nil {
		return UpdateQuestionOrderResponse{}, fmt.Errorf("查询题目关联失败: %w", err)
	}

	existingMap := make(map[int64]bool)
	for _, rel := range existingRelations {
		existingMap[rel.QuestionID] = true
	}
	for _, item := range req.QuestionOrders {
		if !existingMap[item.QuestionID] {
			return UpdateQuestionOrderResponse{}, fmt.Errorf("题目ID %d 不存在于该试卷中", item.QuestionID)
		}
	}

	// 7. 批量更新题目顺序
	updatedCount := 0
	for _, item := range req.QuestionOrders {
		_, err := dao.Q.PaperQuestion.WithContext(ctx).
			Where(
				dao.Q.PaperQuestion.PaperID.Eq(paperID),
				dao.Q.PaperQuestion.QuestionID.Eq(item.QuestionID),
			).
			Update(dao.Q.PaperQuestion.QuestionOrder, item.QuestionOrder)
		if err == nil {
			updatedCount++
		}
	}

	return UpdateQuestionOrderResponse{
		PaperID:      paperID,
		UpdatedCount: updatedCount,
	}, nil
}

// UpdatePaperRequest 更新试卷信息的请求参数
type UpdatePaperRequest struct {
	Title       string `json:"title"`       // 试卷标题（可选）
	Description string `json:"description"` // 试卷描述（可选）
	TotalScore  int    `json:"total_score"` // 试卷总分（可选）
}

// UpdatePaperResponse 更新试卷信息的响应
type UpdatePaperResponse struct {
	ID        int64     `json:"id"`         // 试卷ID
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// UpdatePaper 更新试卷信息
func UpdatePaper(ctx context.Context, paperID, creatorID int64, req UpdatePaperRequest) (UpdatePaperResponse, error) {
	// 1. 验证试卷是否存在且属于当前用户
	paper, err := dao.Q.Paper.WithContext(ctx).
		Where(
			dao.Q.Paper.ID.Eq(paperID),
			dao.Q.Paper.CreatorID.Eq(creatorID),
		).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 检查试卷是否存在
			exists, checkErr := checkPaperExists(ctx, paperID)
			if checkErr != nil {
				return UpdatePaperResponse{}, checkErr
			}
			if exists {
				return UpdatePaperResponse{}, utils.ErrNoPermission
			}
			return UpdatePaperResponse{}, utils.ErrPaperNotFound
		}
		return UpdatePaperResponse{}, fmt.Errorf("查询试卷失败: %w", err)
	}

	// 2. 构建更新字段（只更新非空/有效字段）
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	// 3. 总分校验：如果传递了新总分，需要检查是否不低于题目总分
	if req.TotalScore > 0 {
		// 3.1 查询该试卷下所有题目的分值总和
		questionTotalScore, err := getPaperQuestionsTotalScore(ctx, paperID)
		if err != nil {
			return UpdatePaperResponse{}, fmt.Errorf("计算题目总分失败: %w", err)
		}

		// 3.2 验证新总分是否不低于题目总分
		if req.TotalScore < questionTotalScore {
			return UpdatePaperResponse{}, fmt.Errorf("试卷总分不能低于题目总分（当前题目总分：%d）", questionTotalScore)
		}

		updates["total_score"] = req.TotalScore
	}

	// 4. 若没有需要更新的字段，直接返回
	if len(updates) == 0 {
		return UpdatePaperResponse{
			ID:        paperID,
			UpdatedAt: paper.UpdatedAt,
		}, fmt.Errorf("至少提供一个需要更新的字段")
	}

	// 5. 执行更新操作
	_, err = dao.Q.Paper.WithContext(ctx).
		Where(
			dao.Q.Paper.ID.Eq(paperID),
			dao.Q.Paper.CreatorID.Eq(creatorID),
		).
		Updates(updates)
	if err != nil {
		return UpdatePaperResponse{}, fmt.Errorf("更新试卷失败: %w", err)
	}

	// 6. 查询更新后的试卷信息（获取最新的updated_at）
	updatedPaper, err := dao.Q.Paper.WithContext(ctx).
		Where(dao.Q.Paper.ID.Eq(paperID)).
		First()
	if err != nil {
		return UpdatePaperResponse{}, fmt.Errorf("获取更新后试卷信息失败: %w", err)
	}

	return UpdatePaperResponse{
		ID:        paperID,
		UpdatedAt: updatedPaper.UpdatedAt,
	}, nil
}

// 计算试卷下所有题目的总分
func getPaperQuestionsTotalScore(ctx context.Context, paperID int64) (int, error) {
	// 查询试卷题目关联表中该试卷的所有记录
	relations, err := dao.Q.PaperQuestion.WithContext(ctx).
		Where(dao.Q.PaperQuestion.PaperID.Eq(paperID)).
		Find()
	if err != nil {
		return 0, err
	}

	// 累加所有题目的分值
	total := 0
	for _, rel := range relations {
		total += rel.Score
	}

	return total, nil
}
