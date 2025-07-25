package services

import (
	"CodeQuizAI/dao"
	"context"
	"strconv"
)

// UserStatisticsResponse 用户统计响应结构
type UserStatisticsResponse struct {
	TotalQuestions int            `json:"total_questions"` // 总出题次数
	TotalPapers    int            `json:"total_papers"`    // 总试卷数量
	QuestionTypes  map[string]int `json:"question_types"`  // 题目类型分布
}

// ActiveDetail 活跃详情
type ActiveDetail struct {
	QuestionsCount int `json:"questions_count"` // 该时间段创建的题目数
	PapersCount    int `json:"papers_count"`    // 该时间段创建的试卷数
	Total          int `json:"total"`           // 总活跃数
}

// GetUserStatistics 获取用户统计信息
func GetUserStatistics(ctx context.Context, userID int64) (UserStatisticsResponse, error) {
	var result UserStatisticsResponse

	// 1. 统计总出题次数
	questionsCount, err := dao.Q.Question.WithContext(ctx).
		Where(dao.Question.UserID.Eq(userID)).
		Count()
	if err != nil {
		return result, err
	}
	result.TotalQuestions = int(questionsCount)

	// 2. 统计总试卷数量
	papersCount, err := dao.Q.Paper.WithContext(ctx).
		Where(dao.Paper.CreatorID.Eq(userID)).
		Count()
	if err != nil {
		return result, err
	}
	result.TotalPapers = int(papersCount)

	// 3. 统计题目类型分布
	questionTypes, err := getQuestionTypeDistribution(ctx, userID)
	if err != nil {
		return result, err
	}
	result.QuestionTypes = questionTypes

	return result, nil
}

// 获取题目类型分布
func getQuestionTypeDistribution(ctx context.Context, userID int64) (map[string]int, error) {
	type typeCount struct {
		Type  string
		Count int64
	}
	var results []typeCount

	err := dao.Q.Question.WithContext(ctx).
		Where(dao.Question.UserID.Eq(userID)).
		Select(dao.Question.QuestionType, dao.Question.QuestionType.Count().As("count")).
		Group(dao.Question.QuestionType).
		Scan(&results)

	if err != nil {
		return nil, err
	}

	distribution := make(map[string]int)
	for _, item := range results {
		distribution[item.Type] = int(item.Count)
	}
	return distribution, nil
}

// StatisticsOverview 整体统计数据结构
type StatisticsOverview struct {
	TotalUsers                int64            `json:"total_users"`
	TotalQuestions            int64            `json:"total_questions"`
	TotalPapers               int64            `json:"total_papers"`
	LanguageDistribution      map[string]int64 `json:"language_distribution"`
	AIModelUsage              map[string]int64 `json:"ai_model_usage"`
	PaperQuestionDistribution map[string]int64 `json:"paper_question_distribution"`
}

// GetStatisticsOverview 获取系统整体统计信息
func GetStatisticsOverview(ctx context.Context) (StatisticsOverview, error) {
	var overview StatisticsOverview
	var err error

	// 1. 统计总用户数
	overview.TotalUsers, err = dao.Q.User.WithContext(ctx).Count()
	if err != nil {
		return overview, err
	}

	// 2. 统计总题目数
	overview.TotalQuestions, err = dao.Q.Question.WithContext(ctx).Count()
	if err != nil {
		return overview, err
	}

	// 3. 统计总试卷数
	overview.TotalPapers, err = dao.Q.Paper.WithContext(ctx).Count()
	if err != nil {
		return overview, err
	}

	// 4. 各编程语言题目分布
	overview.LanguageDistribution, err = getLanguageDistribution(ctx)
	if err != nil {
		return overview, err
	}

	// 5. AI模型使用情况
	overview.AIModelUsage, err = getAIModelUsage(ctx)
	if err != nil {
		return overview, err
	}

	// 6. 试卷题目数量分布（使用字段表达式）
	overview.PaperQuestionDistribution, err = getPaperQuestionDistribution(ctx)
	if err != nil {
		return overview, err
	}

	return overview, nil
}

// 各编程语言题目分布
func getLanguageDistribution(ctx context.Context) (map[string]int64, error) {
	type langCount struct {
		Language string `json:"language"`
		Count    int64  `json:"count"`
	}

	var results []langCount
	err := dao.Q.Question.WithContext(ctx).
		Select(
			dao.Question.Language,
			dao.Question.Language.Count().As("count"),
		).
		Group(dao.Question.Language).
		Scan(&results)

	if err != nil {
		return nil, err
	}

	distribution := make(map[string]int64)
	for _, item := range results {
		distribution[item.Language] = item.Count
	}

	return distribution, nil
}

// AI模型使用情况
func getAIModelUsage(ctx context.Context) (map[string]int64, error) {
	type modelCount struct {
		AiModel string `json:"ai_model"`
		Count   int64  `json:"count"`
	}

	var results []modelCount
	// 使用字段表达式构建查询
	err := dao.Q.Question.WithContext(ctx).
		Select(
			dao.Question.AiModel,
			dao.Question.AiModel.Count().As("count"),
		).
		Group(dao.Question.AiModel).
		Scan(&results)

	if err != nil {
		return nil, err
	}

	usage := make(map[string]int64)
	for _, item := range results {
		usage[item.AiModel] = item.Count
	}

	return usage, nil
}

// 试卷题目数量分布
func getPaperQuestionDistribution(ctx context.Context) (map[string]int64, error) {
	type paperQuestionCount struct {
		PaperID int64 `json:"paper_id"`
		Count   int64 `json:"count"`
	}

	var results []paperQuestionCount
	// 统计每张试卷的题目数量
	err := dao.Q.PaperQuestion.WithContext(ctx).
		Select(
			dao.PaperQuestion.PaperID,
			dao.PaperQuestion.PaperID.Count().As("count"),
		).
		Group(dao.PaperQuestion.PaperID).
		Scan(&results)

	if err != nil {
		return nil, err
	}

	// 统计不同题目数量的试卷分布
	distribution := make(map[string]int64)
	for _, item := range results {
		key := strconv.FormatInt(item.Count, 10)
		distribution[key]++
	}

	return distribution, nil
}
