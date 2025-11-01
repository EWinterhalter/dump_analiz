package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		os.Exit(1)
	}
	defer file.Close()

	reSource := regexp.MustCompile(`\w+\s+\d+\s+\d+:\d+:\d+\s+([\w\-.]+)\s`)

	reASA := regexp.MustCompile(`%ASA-\d+-(\d+):`)
	reIP := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

	sources := map[string]bool{}
	messageTypes := map[string]bool{}
	ipSet := map[string]bool{}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	matchedLines := 0

	fmt.Println("=== АНАЛИЗ ФОРМАТА ЛОГА ===")

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if lineCount <= 3 {
			fmt.Printf("Строка %d: %s\n", lineCount, line)
		}

		if m := reSource.FindStringSubmatch(line); m != nil {
			source := m[1]
			sources[source] = true
			matchedLines++
			if matchedLines <= 5 {
				fmt.Printf("Найден источник '%s' в строке: %s\n", source, line[:100])
			}
		}

		if m := reASA.FindStringSubmatch(line); m != nil {
			messageTypes[m[1]] = true
		}

		for _, ip := range reIP.FindAllString(line, -1) {
			ipSet[ip] = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		os.Exit(1)
	}

	fmt.Println("\n==== РЕЗУЛЬТАТЫ ====")
	fmt.Printf("Обработано строк: %d\n", lineCount)
	fmt.Printf("Количество различных источников сообщений: %d\n", len(sources))
	printSet("Источники сообщений", sources)

	fmt.Printf("\nКоличество различных типов сообщений Cisco ASA: %d\n", len(messageTypes))
	printSet("Типы сообщений", messageTypes)

	fmt.Printf("\nКоличество различных IP-адресов: %d\n", len(ipSet))

	if len(ipSet) > 0 {
		ips := make([]string, 0, len(ipSet))
		for ip := range ipSet {
			ips = append(ips, ip)
		}
		sort.Strings(ips)
		for i := 0; i < len(ips) && i <= len(ips); i++ {
			fmt.Printf("  %s\n", ips[i])
		}
	}
}

func printSet(title string, m map[string]bool) {
	if len(m) == 0 {
		fmt.Printf("%s: не найдено\n", title)
		return
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Printf("%s (%d):\n", title, len(keys))
	for _, k := range keys {
		fmt.Printf("  %s\n", k)
	}
}
