package main

import (
	"log/slog"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type DataJob struct {
	id int
}

type JobInfo struct {
	id       int
	duration uint64
}

func main() {
	timeStart := time.Now()

	poolSize := runtime.GOMAXPROCS(0)
	slog.Info("Pool size", "size", poolSize)

	jobCh := make(chan DataJob, poolSize)
	resCh := make(chan JobInfo, poolSize)

	var total uint64 = 0

	var wg sync.WaitGroup

	go func() {
		for i := 0; i < poolSize; i++ {
			jobCh <- DataJob{id: i}
		}
		close(jobCh)
	}()

	for job := range jobCh {
		wg.Add(1)

		go job.doDataJob(resCh, &wg)
	}

	for i := 0; i < poolSize; i++ {
		go func() {
			result := <-resCh
			atomic.AddUint64(&total, result.duration)
		}()
	}

	wg.Wait()
	duration := getDuration(timeStart)
	close(resCh)

	slog.Info("Total jobs done duration", "seconds", total)
	slog.Info("Main goroutine work duration", "seconds", duration)
}

func (j *DataJob) doDataJob(ch chan JobInfo, wg *sync.WaitGroup) {
	defer sendJobInfo(ch, j.id, time.Now())
	defer wg.Done()

	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
}

func sendJobInfo(ch chan JobInfo, jobId int, timeStart time.Time) {
	duration := getDuration(timeStart)

	ch <- JobInfo{id: jobId, duration: duration}

	slog.Info("Done datajob", "id", jobId, "duration", duration)
}

func getDuration(timeStart time.Time) uint64 {
	timeEnd := time.Now()
	return uint64(timeEnd.Sub(timeStart).Seconds())
}

// Задача 3: Пул рабочих горутин с подсчётом времени и сбором результатов
// Описание задачи:
// Вы разрабатываете систему для обработки заданий с использованием пула рабочих горутин.
// Каждое задание представляет собой задачу, которую нужно выполнить, например,
// обработать данные или выполнить вычисления.
// Задания поступают в канал, и несколько рабочих горутин берут задания из канала и
// обрабатывают их параллельно.

// В этой задаче вам нужно реализовать пул горутин, которые:

// Будут обрабатывать поступающие задания.
// Для каждого задания будут считать, сколько времени прошло с момента начала его обработки,
// и добавлять это время к общему счётчику времени.
// Будут собирать информацию о каждом выполненном задании
// (например, ID задания и время его выполнения) в отдельный канал.
// После завершения всех заданий программа должна вывести:

// Результаты выполнения каждого задания (ID задания и время выполнения).
// Общее время, которое все горутины потратили на обработку заданий.
// Требования:
// Пул горутин: Необходимо создать пул рабочих горутин,
// которые будут параллельно обрабатывать задания.
// Подсчёт времени:
// Каждая рабочая горутина должна засекать время начала и окончания работы над заданием,
// и атомарно добавлять это время к общему счётчику времени работы всех горутин.
// Сбор результатов:
// После выполнения задания горутина должна отправлять информацию о выполнении
// (ID задания и время выполнения) в канал для сбора результатов.
// Завершение программы:
// После завершения всех заданий программа должна вывести список результатов и общее время,
// которое все горутины потратили на выполнение заданий.
