package solution

import (
	"fmt"
	"sort"
	"strings"
)

type DisruptionBudget struct {
	DisruptionAllowed int
	AppName           string
}

type Application struct {
	NodeName string
	AppName  string
}

type Node struct {
	NodeName string
}

type NodeAppItem struct {
	AppName string
	Count   int
}

type NodeWeight struct {
	NodeName string
	Weight   float32
}

type NodeAppList struct {
	NodeAppItems []*NodeAppItem
}

func (nal *NodeAppList) AddApp(app *Application) {
	find := false
	for _, item := range nal.NodeAppItems {
		if item.AppName == app.AppName {
			item.Count++
			find = true
			break
		}
	}
	if !find {
		nal.NodeAppItems = append(nal.NodeAppItems, &NodeAppItem{AppName: app.AppName, Count: 1})
	}
}

type NodeGroup struct {
	NodeList []string       // nodeName list
	AppCount map[string]int // appName : count, as a cache
}

// nodeAppMap: {nodeName : NodeAppList}
// budgetMap: {appName : budgetAllowed(int)}
// if can add to NodeGroup, return true, else false
func (ng *NodeGroup) TryToAdd(nodeName string, nodeAppList *NodeAppList, budgetMap *map[string]int) bool {
	// always new node here
	newAppCount := ng.AppCount
	for _, nodeAppItem := range nodeAppList.NodeAppItems {
		curCount, ok := ng.AppCount[nodeAppItem.AppName]
		if ok {
			// check if exceed budget
			if curCount+nodeAppItem.Count > (*budgetMap)[nodeAppItem.AppName] {
				return false
			}
			newAppCount[nodeAppItem.AppName] = curCount + nodeAppItem.Count
		} else if nodeAppItem.Count > (*budgetMap)[nodeAppItem.AppName] {
			return false
		} else {
			newAppCount[nodeAppItem.AppName] = nodeAppItem.Count
		}
	}
	// after check, it's ok to add
	ng.NodeList = append(ng.NodeList, nodeName)
	ng.AppCount = newAppCount
	return true
}

type NodeGroupResult struct {
	NodeGroupList []*NodeGroup
}

func (ngr *NodeGroupResult) AddNode(nodeName string, nodeAppList *NodeAppList, budgetMap *map[string]int) {
	added := false
	for _, ng := range ngr.NodeGroupList {
		ok := ng.TryToAdd(nodeName, nodeAppList, budgetMap)
		if ok {
			added = true
			break
		}
	}
	if !added {
		newNg := NodeGroup{NodeList: []string{}, AppCount: make(map[string]int)}
		newNg.TryToAdd(nodeName, nodeAppList, budgetMap) // should be true
		ngr.NodeGroupList = append(ngr.NodeGroupList, &newNg)
	}
	compareFunc := func(i int, j int) bool {
		return len(ngr.NodeGroupList[i].NodeList) > len(ngr.NodeGroupList[j].NodeList)
	}
	sort.Slice(ngr.NodeGroupList, compareFunc)
}

func (ngr NodeGroupResult) String() string {
	ret := []string{}
	for idx, ng := range ngr.NodeGroupList {
		ret = append(ret, fmt.Sprintf("%d: [%s]", idx, strings.Join(ng.NodeList, ",")))
	}
	return strings.Join(ret, "\n")
}

// return sorted node weight list by weight
func getNodeWeight(nodeAppMap *map[string]NodeAppList, budgetMap *map[string]int) []NodeWeight {
	nodeWeightList := make([]NodeWeight, 0)
	for nodeName, appList := range *nodeAppMap {
		weight := float32(0.0)
		for _, app := range appList.NodeAppItems {
			count := app.Count
			budgetAllowed := (*budgetMap)[app.AppName]
			tmpWeight := float32(count) / float32(budgetAllowed)
			weight += tmpWeight
		}
		nodeWeightList = append(nodeWeightList, NodeWeight{NodeName: nodeName, Weight: weight})
	}
	sort.Slice(nodeWeightList, func(i, j int) bool {
		return nodeWeightList[i].Weight < nodeWeightList[j].Weight
	})
	return nodeWeightList
}

// find quickest upgrade solution and return node names
func FindUpgradeSolution(nodeList []Node, appList []Application, disruptionBudgetList []DisruptionBudget) NodeGroupResult {
	// make map:  {nodeName: NodeAppList}
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

	// get node weight map
	nodeWeightList := getNodeWeight(&nodeAppMap, &budgetMap)

	result := NodeGroupResult{}
	// 后续优化：使用堆，现在先遍历
	for _, nw := range nodeWeightList {
		nodeName := nw.NodeName
		naList := nodeAppMap[nodeName]
		result.AddNode(nodeName, &naList, &budgetMap)
	}
	return result
}
