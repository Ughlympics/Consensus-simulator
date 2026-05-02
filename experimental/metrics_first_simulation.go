package experimental

import (
	"fmt"
)

func RunSimulation(
	N int,
	M int,
	K int,
	R int,
	E int,
	gamma float64,
	alpha float64,
	pOff float64,
	pInv float64,
	tau float64,
	useWhales bool,
) GlobalMetrics {

	nodes := GenerateNodes(N, gamma)

	allBlocks := []Block{}
	totalSlots := 0

	captureCount := 0
	concentrationSum := 0.0
	epochs := 0

	roundsDone := 0

	for roundsDone < R {

		// 1. candidate selection (random or with whales)
		candidates := SelectCandidates(nodes, M)

		// 2. voting (with or without whales)
		if useWhales {
			candidates = VoteWhaleForOneCandidate(nodes, candidates, K)
		} else {
			VoteRandom(nodes, candidates)
		}

		// 3. obvious delegates
		delegates := SelectDelegates(candidates, K)

		// 4. determining how many rounds to simulate in this epoch (E or remaining rounds)
		roundsLeft := R - roundsDone
		currentE := E
		if roundsLeft < E {
			currentE = roundsLeft
		}

		// 5. gen blocks
		blocks := GenerateBlocks(delegates, currentE, alpha, pOff, pInv, tau)

		allBlocks = append(allBlocks, blocks...)
		totalSlots += currentE

		// 6. local metrics for the current epoch
		epochMetrics := ComputeMetrics(blocks, currentE, delegates)

		if epochMetrics.Capture {
			captureCount++
		}

		concentrationSum += epochMetrics.Concentration
		epochs++

		roundsDone += currentE
	}

	// 7. gloval metrics for the whole simulation
	global := ComputeMetrics(allBlocks, totalSlots, []Delegate{})

	// 8. add capture and concentration to global metrics
	if epochs > 0 {
		global.Concentration = concentrationSum / float64(epochs)
		global.Capture = float64(captureCount)/float64(epochs) > 0.5
	}

	return global
}

func PrintMetrics(m GlobalMetrics) {
	fmt.Println("===== Результати симуляції DPoS =====")

	fmt.Printf("Загальна кількість слотів: %d\n", m.TotalSlots)
	fmt.Printf("Створено блоків: %d\n", m.ProducedBlocks)

	fmt.Printf("Валідні блоки: %d\n", m.ValidBlocks)
	fmt.Printf("Невалідні блоки: %d\n", m.InvalidBlocks)

	fmt.Printf("Доступність мережі (Availability): %.2f%%\n", m.Availability*100)
	fmt.Printf("Частка невалідних блоків (Invalid Share): %.2f%%\n", m.InvalidShare*100)
	fmt.Printf("Частота форків (Fork Rate): %.2f%%\n", m.ForkRate*100)

	fmt.Printf("Середня затримка (Latency): %.2f слотів\n", m.Latency)

	if m.Capture {
		fmt.Println("⚠️ Мережа була захоплена (більшість делегатів зловмисні)")
	} else {
		fmt.Println("✅ Мережа не була захоплена")
	}

	fmt.Printf("Рівень централізації (Concentration): %.2f%%\n", m.Concentration*100)

	fmt.Println("=====================================")
}

/////////////////////////////////////////////////////////
////////////////
/////////////////////////////////////////////////////////
//RRRRRRRRROOOOOOOOAAAAAAARRRRRRR

// func RunSimulation(
// 	N int,
// 	M int,
// 	K int,
// 	R int,
// 	E int,
// 	gamma float64,
// 	alpha float64,
// 	pOff float64,
// 	pInv float64,
// 	tau float64,
// 	useWhales bool,
// ) GlobalMetrics {

// 	nodes := GenerateNodes(N, gamma)

// 	allBlocks := []Block{}

// 	roundsDone := 0

// 	for roundsDone < R {

// 		candidates := SelectCandidates(nodes, M)

// 		if useWhales {
// 			candidates = VoteWhaleForOneCandidate(nodes, candidates, K)
// 		} else {
// 			VoteRandom(nodes, candidates)
// 		}

// 		delegates := SelectDelegates(candidates, K)

// 		roundsLeft := R - roundsDone
// 		currentE := E
// 		if roundsLeft < E {
// 			currentE = roundsLeft
// 		}

// 		blocks := GenerateBlocks(delegates, currentE, alpha, pOff, pInv, tau)

// 		allBlocks = append(allBlocks, blocks...)

// 		roundsDone += currentE
// 	}

// 	return ComputeMetrics(allBlocks)
// }
