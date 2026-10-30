// Файл manager.go — реализация интерфейса ProcessManager
package main

import (
	"fmt"
	"sort"
	"strings"
)

// ====================================================
// РЕАЛИЗАЦИЯ № 1: SimpleProcessManager (хранит как есть)
// ====================================================

type SimpleProcessManager struct {
	processes []*Process
}

func NewProcessManager() *SimpleProcessManager {
	return &SimpleProcessManager{
		processes: make([]*Process, 0),
	}
}

func (m *SimpleProcessManager) Add(p Process) {
	pCopy := p
	m.processes = append(m.processes, &pCopy)
}

func (m *SimpleProcessManager) Remove(pid int) error {
	for i, p := range m.processes {
		if p.PID == pid {
			m.processes = append(m.processes[:i], m.processes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("процесс с PID %d не найден", pid)
}

func (m *SimpleProcessManager) FindByPID(pid int) (*Process, error) {
	for _, p := range m.processes {
		if p.PID == pid {
			return p, nil
		}
	}
	return nil, fmt.Errorf("процесс с PID %d не найден", pid)
}

func (m *SimpleProcessManager) FindByName(name string) []*Process {
	var result []*Process
	nameLower := strings.ToLower(name)
	for _, p := range m.processes {
		if strings.Contains(strings.ToLower(p.Name), nameLower) {
			result = append(result, p)
		}
	}
	return result
}

func (m *SimpleProcessManager) All() []*Process {
	return m.processes
}

func (m *SimpleProcessManager) FilterByStatus(status ProcessStatus) []*Process {
	var result []*Process
	for _, p := range m.processes {
		if p.Status == status {
			result = append(result, p)
		}
	}
	return result
}

func (m *SimpleProcessManager) TotalUsage() ResourceUsage {
	var total ResourceUsage
	for _, p := range m.processes {
		total.CPUPercent += p.CPUPercent
		total.MemoryMB += p.MemoryMB
	}
	return total
}

// GroupByUser — ЗАДАНИЕ 2: группирует процессы по пользователю
func (m *SimpleProcessManager) GroupByUser() map[string][]*Process {
	groups := make(map[string][]*Process)
	for _, p := range m.processes {
		groups[p.User] = append(groups[p.User], p)
	}
	return groups
}

// PrintUserStats — выводит статистику ресурсов по каждому пользователю
func (m *SimpleProcessManager) PrintUserStats() {
	groups := m.GroupByUser()
	fmt.Println("\n  ── Статистика по пользователям ───────────────")
	for user, procs := range groups {
		var totalCPU, totalMem float64
		for _, p := range procs {
			totalCPU += p.CPUPercent
			totalMem += p.MemoryMB
		}
		fmt.Printf("  %-20s: %d процессов | CPU: %5.1f%% | Память: %7.1f МБ\n",
			user, len(procs), totalCPU, totalMem)
	}
}

// ── Сортировка ──────────────────────────────────────

type SortField int

const (
	SortByPID SortField = iota
	SortByCPU
	SortByMemory
	SortByName
	SortByStatus
)

func (m *SimpleProcessManager) SortBy(field SortField) {
	sort.Slice(m.processes, func(i, j int) bool {
		pi, pj := m.processes[i], m.processes[j]
		switch field {
		case SortByCPU:
			return pi.CPUPercent > pj.CPUPercent
		case SortByMemory:
			return pi.MemoryMB > pj.MemoryMB
		case SortByName:
			return pi.Name < pj.Name
		case SortByStatus:
			return string(pi.Status) < string(pj.Status)
		default:
			return pi.PID < pj.PID
		}
	})
}

func (m *SimpleProcessManager) Print() {
	printHeader()
	for _, p := range m.processes {
		fmt.Println(p.FormatRowColored())
	}
	fmt.Printf("  %s\n", strings.Repeat("─", 82))

	total := m.TotalUsage()
	fmt.Printf("  Итого: %d процессов | CPU: %.1f%% | Память: %.1f МБ\n",
		len(m.processes), total.CPUPercent, total.MemoryMB)
}

func (m *SimpleProcessManager) Stats() {
	statuses := []ProcessStatus{StatusRunning, StatusSleeping, StatusZombie, StatusStopped}
	fmt.Println("\n  ── Статистика по статусам ───────────────")
	for _, s := range statuses {
		count := len(m.FilterByStatus(s))
		if count > 0 {
			fmt.Printf("  %s%-9s%s: %d\n", statusColor(s), s, colorReset, count)
		}
	}
}

// ====================================================
// РЕАЛИЗАЦИЯ № 2: PriorityProcessManager (ЗАДАНИЕ 3, доп.)
// Всегда хранит процессы отсортированными по приоритету
// ====================================================

type PriorityProcessManager struct {
	processes []*Process
}

func NewPriorityProcessManager() *PriorityProcessManager {
	return &PriorityProcessManager{processes: make([]*Process, 0)}
}

// sortByPriority — внутренний метод поддержания порядка
func (m *PriorityProcessManager) sortByPriority() {
	sort.Slice(m.processes, func(i, j int) bool {
		return m.processes[i].Priority < m.processes[j].Priority
	})
}

func (m *PriorityProcessManager) Add(p Process) {
	pCopy := p
	m.processes = append(m.processes, &pCopy)
	m.sortByPriority() // поддерживаем инвариант сортировки
}

func (m *PriorityProcessManager) Remove(pid int) error {
	for i, p := range m.processes {
		if p.PID == pid {
			m.processes = append(m.processes[:i], m.processes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("процесс с PID %d не найден", pid)
}

func (m *PriorityProcessManager) FindByPID(pid int) (*Process, error) {
	for _, p := range m.processes {
		if p.PID == pid {
			return p, nil
		}
	}
	return nil, fmt.Errorf("процесс с PID %d не найден", pid)
}

func (m *PriorityProcessManager) FindByName(name string) []*Process {
	var result []*Process
	nameLower := strings.ToLower(name)
	for _, p := range m.processes {
		if strings.Contains(strings.ToLower(p.Name), nameLower) {
			result = append(result, p)
		}
	}
	return result
}

func (m *PriorityProcessManager) All() []*Process {
	return m.processes
}

func (m *PriorityProcessManager) FilterByStatus(status ProcessStatus) []*Process {
	var result []*Process
	for _, p := range m.processes {
		if p.Status == status {
			result = append(result, p)
		}
	}
	return result
}

func (m *PriorityProcessManager) TotalUsage() ResourceUsage {
	var total ResourceUsage
	for _, p := range m.processes {
		total.CPUPercent += p.CPUPercent
		total.MemoryMB += p.MemoryMB
	}
	return total
}

func (m *PriorityProcessManager) Print() {
	printHeader()
	for _, p := range m.processes {
		fmt.Println(p.FormatRowColored())
	}
}
