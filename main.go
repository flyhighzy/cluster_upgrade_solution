package main

import (
	"fmt"
	"test/solution"
)

func main() {
	nodeList := []solution.Node{{NodeName: "node1"}, {NodeName: "node2"}, {NodeName: "node3"}}

	// make app list
	app11 := solution.Application{NodeName: "node1", AppName: "app1"}
	app12 := solution.Application{NodeName: "node2", AppName: "app1"}

	app21 := solution.Application{NodeName: "node1", AppName: "app2"}
	app22 := solution.Application{NodeName: "node2", AppName: "app2"}

	app31 := solution.Application{NodeName: "node2", AppName: "app3"}
	app32 := solution.Application{NodeName: "node3", AppName: "app3"}

	appList := []solution.Application{app11, app12, app21, app22, app31, app32}

	// DisruptionBudget list
	db1 := solution.DisruptionBudget{DisruptionAllowed: 1, AppName: "app1"}
	db2 := solution.DisruptionBudget{DisruptionAllowed: 1, AppName: "app2"}
	db3 := solution.DisruptionBudget{DisruptionAllowed: 1, AppName: "app3"}
	disruptionBudgetList := []solution.DisruptionBudget{db1, db2, db3}

	result := solution.FindUpgradeSolution(nodeList, appList, disruptionBudgetList)
	fmt.Println(result)
}
