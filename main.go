// Файл main.go — основная логика программы
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// ====================================================
// ГЕНЕРАТОР ТЕСТОВЫХ ДАННЫХ
// ====================================================

func generateTestProcesses() []Process {
	processNames := []string{
		"System", "svchost.exe", "explorer.exe", "chrome.exe",
		"code.exe", "python.exe", "node.exe", "postgres.exe",
		"nginx", "redis-server", "docker", "go.exe",
		"winlogon.exe", "lsass.exe", "csrss.exe", "wininit.exe",
	}

	users := []string{"SYSTEM", "Student", "Administrator", "Network Service"}
	statuses := []ProcessStatus{
		StatusRunning, StatusRunning, StatusRunning,
		StatusSleeping, StatusSleeping,
		StatusZombie,
		StatusStopped,
	}

	var processes []Process

	for i, name := range processNames {
		pid := 1000 + i*100 + rand.Intn(99)

		status := statuses[rand.Intn(len(statuses))]
		var cpu, mem float64

		switch name {
		case "chrome.exe":
			cpu = 5.0 + rand.Float64()*15.0
			mem = 200 + rand.Float64()*300
		case "System":
			cpu = 0.1 + rand.Float64()*0.5
			mem = 50 + rand.Float64()*50
		case "postgres.exe":
			cpu = 2.0 + rand.Float64()*8.0
			mem = 100 + rand.Float64()*200
		default:
			cpu = rand.Float64() * 5.0
			mem = 10 + rand.Float64()*100
		}

		if status != StatusRunning {
			cpu = 0
		}

		// ЗАДАНИЕ 1: случайное время запуска (от 1 минуты до 10 часов назад)
		secondsAgo := rand.Intn(36000) + 60
		startTime := time.Now().Add(-time.Duration(secondsAgo) * time.Second)

		processes = append(processes, Process{
			PID:       pid,
			Name:      name,
			User:      users[rand.Intn(len(users))],
			Status:    status,
			Priority:  rand.Intn(40) - 20,
			Threads:   1 + rand.Intn(32),
			StartTime: startTime,
			ResourceUsage: ResourceUsage{
				CPUPercent: cpu,
				MemoryMB:   mem,
			},
		})
	}

	return processes
}

// ====================================================
// ГЛАВНАЯ ФУНКЦИЯ
// ====================================================

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║         Менеджер процессов — Go Process Manager          ║")
	fmt.Println("║              МДК.01.04 Практическая работа № 3           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	manager := NewProcessManager()

	for _, p := range generateTestProcesses() {
		manager.Add(p)
	}

	// ── Демонстрация 1: вывод всех процессов (с Uptime) ────
	fmt.Println("=== Все процессы (сортировка по PID) ===")
	manager.SortBy(SortByPID)
	manager.Print()
	fmt.Println()

	// ── Демонстрация 2: сортировка по памяти ──────────────
	fmt.Println("=== Топ процессов по использованию памяти ===")
	manager.SortBy(SortByMemory)
	manager.Print()
	fmt.Println()

	// ── Демонстрация 3: сортировка по CPU ──────────────────
	fmt.Println("=== Топ процессов по загрузке CPU ===")
	manager.SortBy(SortByCPU)
	manager.Print()
	fmt.Println()

	// ── Демонстрация 4: фильтрация Running процессов ───────
	fmt.Println("=== Только Running процессы ===")
	running := manager.FilterByStatus(StatusRunning)
	printHeader()
	for _, p := range running {
		fmt.Println(p.FormatRowColored())
	}
	fmt.Printf("  Найдено: %d Running процессов\n", len(running))
	fmt.Println()

	// ── Демонстрация 5: поиск по имени ────────────────────
	fmt.Println("=== Поиск процессов по имени 'chrome' ===")
	found := manager.FindByName("chrome")
	if len(found) > 0 {
		printHeader()
		for _, p := range found {
			fmt.Println(p.FormatRowColored())
		}
	} else {
		fmt.Println("  Процессы не найдены")
	}
	fmt.Println()

	// ── Демонстрация 6: поиск по PID и Kill ───────────────
	allProcs := manager.All()
	if len(allProcs) > 0 {
		targetPID := allProcs[0].PID
		fmt.Printf("=== Поиск процесса PID=%d и его завершение ===\n", targetPID)

		proc, err := manager.FindByPID(targetPID)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
		} else {
			fmt.Printf("Найден: %s | Uptime: %s\n", proc, formatUptime(proc.Uptime()))
			proc.Kill()
			fmt.Printf("После Kill: %s\n", proc)
		}
	}
	fmt.Println()

	// ── Демонстрация 7: статистика по статусам ─────────────
	fmt.Println("=== Статистика по статусам ===")
	manager.Stats()

	// ── Демонстрация 8: ЗАДАНИЕ 2 — группировка по пользователю
	fmt.Println()
	fmt.Println("=== Группировка по пользователям ===")
	manager.PrintUserStats()

	// ── Демонстрация 9: удаление процесса ─────────────────
	fmt.Println()
	fmt.Println("=== Удаление процесса ===")

	err := manager.Remove(99999)
	if err != nil {
		fmt.Printf("Ожидаемая ошибка: %v\n", err)
	}

	if len(manager.All()) > 0 {
		pidToRemove := manager.All()[0].PID
		err = manager.Remove(pidToRemove)
		if err != nil {
			fmt.Printf("Ошибка удаления: %v\n", err)
		} else {
			fmt.Printf("Процесс PID=%d успешно удалён\n", pidToRemove)
			fmt.Printf("Осталось процессов: %d\n", len(manager.All()))
		}
	}

	// ── Демонстрация 10: ЗАДАНИЕ 3 (доп.) — полиморфизм ────
	// Один и тот же интерфейс ProcessManager реализуют ДВА разных типа.
	// Это демонстрация принципа "утиной типизации" в Go.
	fmt.Println()
	fmt.Println("=== Демонстрация полиморфизма через интерфейс ProcessManager ===")

	var pm ProcessManager // переменная интерфейсного типа

	pm = NewProcessManager() // можно присвоить SimpleProcessManager
	pm.Add(Process{PID: 1, Name: "test1", User: "demo",
		Status: StatusRunning, StartTime: time.Now()})
	fmt.Printf("SimpleProcessManager: %d процессов\n", len(pm.All()))

	pm = NewPriorityProcessManager() // и PriorityProcessManager — тоже подходит!
	pm.Add(Process{PID: 2, Name: "test2", Priority: 10, User: "demo",
		Status: StatusRunning, StartTime: time.Now()})
	pm.Add(Process{PID: 3, Name: "test3", Priority: -5, User: "demo",
		Status: StatusRunning, StartTime: time.Now()})
	fmt.Printf("PriorityProcessManager: %d процессов (автосортировка по приоритету)\n", len(pm.All()))
	for _, p := range pm.All() {
		fmt.Printf("  PID=%d Priority=%d\n", p.PID, p.Priority)
	}

	// ── Режим командной строки: --kill PID ─────────────────
	if len(os.Args) > 1 && os.Args[1] == "--kill" && len(os.Args) > 2 {
		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Неверный PID: %s\n", os.Args[2])
			os.Exit(1)
		}
		proc, err := manager.FindByPID(pid)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			os.Exit(1)
		}
		proc.Kill()
	}
}
