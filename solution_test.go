package solution

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func generateNodes(n int) []Node {
	nodes := make([]Node, n)
	for idx := range nodes {
		name := fmt.Sprintf("node%d", idx+1)
		nodes[idx].NodeName = name
	}
	return nodes
}

func randomIntSlice(min int, max int, length int) []int {
	result := make([]int, length)
	for i := 0; i < length; i++ {
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(max-min) + min
		result[i] = num
	}
	return result
}

func generateRandomApps(appNum int, nodes *[]Node) ([]Application, []DisruptionBudget) {
	// each app have random int replicas between 10 and 200
	// and each replica have random node choosed
	result := make([]Application, 0)
	budgets := make([]DisruptionBudget, appNum)
	for i := 0; i < appNum; i++ {
		rand.Seed(time.Now().UnixNano())
		replicaNum := rand.Intn(190) + 10
		distributionIndice := randomIntSlice(0, len(*nodes), replicaNum)
		appName := fmt.Sprintf("app%d", i)
		for _, nodeIdx := range distributionIndice {
			result = append(result, Application{NodeName: (*nodes)[nodeIdx].NodeName, AppName: appName})
		}
		budgets[i].AppName = appName
		budgets[i].DisruptionAllowed = replicaNum / 5
	}
	return result, budgets
}

func checkResultValid(result *NodeGroupResult, nodeList []Node, appList []Application, disruptionBudgetList []DisruptionBudget) bool {
	nodeAppMap := make(map[string]NodeAppList, len(nodeList))
	for _, app := range appList {
		curApps, ok := nodeAppMap[app.NodeName]
		if ok {
			curApps.AddApp(&app)
		} else {
			newItem := NodeAppList{}
			newItem.AddApp(&app)
			nodeAppMap[app.NodeName] = newItem
		}
	}

	// appname : disruptionAllowed map
	budgetMap := make(map[string]int, len(appList))
	for _, budget := range disruptionBudgetList {
		budgetMap[budget.AppName] = budget.DisruptionAllowed
	}

	for _, ng := range result.NodeGroupList {
		curBudget := make(map[string]int, len(appList)) // appname: count
		for _, nodeName := range ng.NodeList {
			// calculate total disruption
			nodeAppList := nodeAppMap[nodeName]
			for _, item := range nodeAppList.NodeAppItems {
				curCount, ok := curBudget[item.AppName]
				if ok {
					curCount += item.Count
					if curCount > budgetMap[item.AppName] {
						return false
					}
				} else {
					curBudget[item.AppName] = item.Count
				}
			}
		}
	}
	return true
}

func TestNormalCase(t *testing.T) {
	nodeList := generateNodes(3)

	// make app list
	app11 := Application{NodeName: "node1", AppName: "app1"}
	app12 := Application{NodeName: "node2", AppName: "app1"}

	app21 := Application{NodeName: "node1", AppName: "app2"}
	app22 := Application{NodeName: "node2", AppName: "app2"}

	app31 := Application{NodeName: "node2", AppName: "app3"}
	app32 := Application{NodeName: "node3", AppName: "app3"}

	appList := []Application{app11, app12, app21, app22, app31, app32}

	// DisruptionBudget list
	db1 := DisruptionBudget{DisruptionAllowed: 1, AppName: "app1"}
	db2 := DisruptionBudget{DisruptionAllowed: 1, AppName: "app2"}
	db3 := DisruptionBudget{DisruptionAllowed: 1, AppName: "app3"}
	disruptionBudgetList := []DisruptionBudget{db1, db2, db3}

	result := FindUpgradeSolution(nodeList, appList, disruptionBudgetList)
	fmt.Println(result)
	assert.ElementsMatch(t, result.NodeGroupList[0].NodeList, []string{"node1", "node3"})
}

func TestAllSingleCase(t *testing.T) {
	nodeList := generateNodes(3)

	// make app list
	app11 := Application{NodeName: "node1", AppName: "app1"}
	app12 := Application{NodeName: "node2", AppName: "app1"}
	app13 := Application{NodeName: "node3", AppName: "app1"}

	app21 := Application{NodeName: "node1", AppName: "app2"}
	app22 := Application{NodeName: "node2", AppName: "app2"}

	app31 := Application{NodeName: "node2", AppName: "app3"}
	app32 := Application{NodeName: "node3", AppName: "app3"}

	appList := []Application{app11, app12, app13, app21, app22, app31, app32}

	// DisruptionBudget list
	db1 := DisruptionBudget{DisruptionAllowed: 1, AppName: "app1"}
	db2 := DisruptionBudget{DisruptionAllowed: 1, AppName: "app2"}
	db3 := DisruptionBudget{DisruptionAllowed: 1, AppName: "app3"}
	disruptionBudgetList := []DisruptionBudget{db1, db2, db3}

	result := FindUpgradeSolution(nodeList, appList, disruptionBudgetList)
	fmt.Println(result)
	// assert.ElementsMatch(t, result.NodeGroupList[0].NodeList, []string{"node1", "node3"})
}

func TestRandomLargeScaleCase(t *testing.T) {
	nodeList := generateNodes(5000)
	fmt.Println(nodeList[:10])

	appList, disruptionBudgetList := generateRandomApps(4000, &nodeList)

	// fmt.Printf("appList: %+v\n", appList[:10])
	// fmt.Printf("budgets: %+v\n", disruptionBudgetList[:10])

	start := time.Now()
	result := FindUpgradeSolution(nodeList, appList, disruptionBudgetList)
	end := time.Now()
	fmt.Println(end.Sub(start))
	// fmt.Println(result)
	totalNodes := 0
	for idx, nglist := range result.NodeGroupList {
		fmt.Printf("%d : %d\n", idx, len(nglist.NodeList))
		totalNodes += len(nglist.NodeList)
	}
	assert.Equal(t, 5000, totalNodes)

	ret := checkResultValid(&result, nodeList, appList, disruptionBudgetList)
	assert.Equal(t, ret, true)
}
