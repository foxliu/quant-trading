package account

import (
	"context"
	"errors"
	"sync"
)

/*
Scheduler
=========

账户调度器负责:
- 多账户资源分配
- 账户间优先级管理
- 账户级风险约束
- 账户状态监控

设计原则:
- 账户是资源,需要统一调度
- 支持账户间隔离
- 支持账户优先级
*/
type Scheduler struct {
	mu       sync.RWMutex
	accounts map[string]*Context // accountID -> Context
	priority map[string]int      // accountID -> priority (数字越小优先级越高)

	ctx    context.Context
	cancel context.CancelFunc
}

// NewScheduler 创建账户调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		accounts: make(map[string]*Context),
		priority: make(map[string]int),
	}
}

// Start 启动调度器
func (s *Scheduler) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

// Register 注册账户
func (s *Scheduler) Register(accountID string, ctx *Context, priority int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.accounts[accountID]; exists {
		return errors.New("account already registered: " + accountID)
	}

	s.accounts[accountID] = ctx
	s.priority[accountID] = priority
	return nil
}

// Unregister 注销账户
func (s *Scheduler) Unregister(accountID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.accounts[accountID]; !exists {
		return errors.New("account not found: " + accountID)
	}

	delete(s.accounts, accountID)
	delete(s.priority, accountID)
	return nil
}

// Get 获取账户上下文
func (s *Scheduler) Get(accountID string) (*Context, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ctx, exists := s.accounts[accountID]
	if !exists {
		return nil, errors.New("account not found: " + accountID)
	}

	return ctx, nil
}

// List 列出所有账户ID
func (s *Scheduler) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]string, 0, len(s.accounts))
	for id := range s.accounts {
		ids = append(ids, id)
	}
	return ids
}

// GetPriority 获取账户优先级
func (s *Scheduler) GetPriority(accountID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	priority, exists := s.priority[accountID]
	if !exists {
		return 0, errors.New("account not found: " + accountID)
	}

	return priority, nil
}

// SetPriority 设置账户优先级
func (s *Scheduler) SetPriority(accountID string, priority int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.accounts[accountID]; !exists {
		return errors.New("account not found: " + accountID)
	}

	s.priority[accountID] = priority
	return nil
}

// GetAllSnapshots 获取所有账户快照
func (s *Scheduler) GetAllSnapshots() []Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshots := make([]Snapshot, 0, len(s.accounts))
	for _, ctx := range s.accounts {
		ctx.mu.Lock()
		snapshot := Snapshot{
			Balance: ctx.balance,
			At:      ctx.updateAt,
		}
		ctx.mu.Unlock()
		snapshots = append(snapshots, snapshot)
	}

	return snapshots
}
