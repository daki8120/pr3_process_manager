// Файл process.go — структуры и типы для процессов
package main

import (
	"fmt"
	"strings"
	"time"
)

// ====================================================
// ПЕРЕЧИСЛЕНИЕ статусов процесса
// ====================================================

type ProcessStatus string

const (
	StatusRunning  ProcessStatus = "Running"
	StatusSleeping ProcessStatus = "Sleeping"
	StatusZombie   ProcessStatus = "Zombie"
	StatusStopped  ProcessStatus = "Stopped"
)

// ====================================================
// БАЗОВАЯ СТРУКТУРА: ресурсы (embedding)
// ====================================================

type ResourceUsage struct {
	CPUPercent float64
	MemoryMB   float64
}

func (r ResourceUsage) TotalLoad() float64 {
	return r.CPUPercent*0.6 + (r.MemoryMB/1024)*0.4*100
}

// ====================================================
// ОСНОВНАЯ СТРУКТУРА: процесс
// ====================================================

// Process представляет системный процесс
type Process struct {
	PID       int
	Name      string
	User      string
	Status    ProcessStatus
	Priority  int
	Threads   int
	StartTime time.Time // ЗАДАНИЕ 1: время запуска процесса

	ResourceUsage
}

// ====================================================
// МЕТОДЫ НА Process
// ====================================================

func (p Process) String() string {
	return fmt.Sprintf("Process{PID:%d, Name:%s, CPU:%.1f%%, Mem:%.1fMB, Status:%s}",
		p.PID, p.Name, p.CPUPercent, p.MemoryMB, p.Status)
}

func (p Process) IsAlive() bool {
	return p.Status == StatusRunning || p.Status == StatusSleeping
}

// Uptime — ЗАДАНИЕ 1: возвращает время работы процесса
func (p Process) Uptime() time.Duration {
	return time.Since(p.StartTime)
}

// formatUptime преобразует Duration в формат ЧЧ:ММ:СС
func formatUptime(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func (p *Process) Kill() {
	p.Status = StatusZombie
	p.CPUPercent = 0
	fmt.Printf("Процесс %d (%s) завершён\n", p.PID, p.Name)
}

func (p *Process) Suspend() {
	if p.Status == StatusRunning {
		p.Status = StatusStopped
		fmt.Printf("Процесс %d (%s) приостановлен\n", p.PID, p.Name)
	}
}

func (p *Process) Resume() {
	if p.Status == StatusStopped {
		p.Status = StatusRunning
		fmt.Printf("Процесс %d (%s) возобновлён\n", p.PID, p.Name)
	}
}

// FormatRow — теперь включает колонку Uptime
func (p Process) FormatRow() string {
	name := p.Name
	if len(name) > 20 {
		name = name[:17] + "..."
	}
	return fmt.Sprintf("  %-6d %-20s %-10s %-9s %5.1f%% %7.1f MB %3d %9s",
		p.PID, name, p.User, p.Status, p.CPUPercent, p.MemoryMB, p.Threads,
		formatUptime(p.Uptime()))
}

// ====================================================
// ИНТЕРФЕЙС ProcessManager
// ====================================================

type ProcessManager interface {
	Add(p Process)
	Remove(pid int) error
	FindByPID(pid int) (*Process, error)
	FindByName(name string) []*Process
	All() []*Process
	FilterByStatus(status ProcessStatus) []*Process
	TotalUsage() ResourceUsage
}

// ====================================================
// ЦВЕТА ANSI ДЛЯ СТАТУСОВ
// ====================================================

func statusColor(status ProcessStatus) string {
	switch status {
	case StatusRunning:
		return "\033[32m"
	case StatusSleeping:
		return "\033[36m"
	case StatusZombie:
		return "\033[31m"
	case StatusStopped:
		return "\033[33m"
	default:
		return "\033[0m"
	}
}

const colorReset = "\033[0m"

// FormatRowColored — с цветом и колонкой Uptime
func (p Process) FormatRowColored() string {
	name := p.Name
	if len(name) > 20 {
		name = name[:17] + "..."
	}
	color := statusColor(p.Status)
	return fmt.Sprintf("  %-6d %-20s %-10s %s%-9s%s %5.1f%% %7.1f MB %3d %9s",
		p.PID, name, p.User, color, p.Status, colorReset,
		p.CPUPercent, p.MemoryMB, p.Threads, formatUptime(p.Uptime()))
}

// printHeader — добавлена колонка "Uptime"
func printHeader() {
	line := strings.Repeat("─", 82)
	fmt.Println("  " + line)
	fmt.Printf("  %-6s %-20s %-10s %-9s %6s %9s %4s %9s\n",
		"PID", "Имя", "Пользователь", "Статус", "CPU%", "Память", "Поток", "Uptime")
	fmt.Println("  " + line)
}
